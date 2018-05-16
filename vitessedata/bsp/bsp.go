package bsp

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/golang/protobuf/proto"
	"github.com/vitesse-ftian/dggo/vitessedata/proto/xdrive"
)

type stepMsg struct {
	step int32
	msg  *xdrive.XMsg
}

type Bsp struct {
	ncli       int64
	me         int64
	globalStep int64
	msgCnt     int64
	sess       string
	db         fdb.Database
	dir        directory.DirectorySubspace
}

//
// Public APIs
//
func BspInit(cf string, sess string, ncli int64, me int64) *Bsp {
	var bsp Bsp
	bsp.ncli = ncli
	bsp.me = me
	bsp.globalStep = 0
	bsp.msgCnt = 0
	bsp.sess = sess

	// Open database.
	fdb.MustAPIVersion(510)
	if cf != "" {
		bsp.db = fdb.MustOpen(cf, []byte("DB"))
	} else {
		bsp.db = fdb.MustOpenDefault()
	}

	path := []string{"__Deepgreen_sys", "bsp", bsp.sess}

	var err error
	bsp.dir, err = directory.CreateOrOpen(bsp.db, path, nil)
	if err != nil {
		panic(err)
	}
	return &bsp
}

func (bsp *Bsp) Deinit() {
	path := []string{"__Deepgreen_sys", "bsp", bsp.sess}
	_, err := directory.Root().Remove(bsp.db, path)
	if err != nil {
		panic(err)
	}
}

func (bsp *Bsp) GlobalStep() int64 {
	return bsp.globalStep
}

func (bsp *Bsp) Send(route int64, msg *xdrive.XMsg) {
	msgba, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	_, err = bsp.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		ktup := tuple.Tuple{bsp.globalStep + 1, route, bsp.me, bsp.msgCnt}
		k := bsp.dir.Pack(ktup)
		tr.Set(k, msgba)
		return nil, nil
	})

	bsp.msgCnt += 1
}

func (bsp *Bsp) Recv() *xdrive.XMsg {
	msg, err := bsp.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		ktup := tuple.Tuple{bsp.globalStep, bsp.me}
		k := bsp.dir.Pack(ktup)
		kr, _ := fdb.PrefixRange(k)
		krr := tr.GetRange(kr, fdb.RangeOptions{})
		kri := krr.Iterator()
		if kri.Advance() {
			kv := kri.MustGet()
			tr.Clear(kv.Key)
			return kv.Value, nil
		} else {
			return nil, nil
		}
	})

	if err != nil {
		panic(err)
	}

	var xmsg xdrive.XMsg
	err = proto.Unmarshal(msg.([]byte), &xmsg)
	if err != nil {
		panic(err)
	}
	return &xmsg
}

type syncData struct {
	nend  int
	ncont int
	w     fdb.FutureNil
}

func (bsp *Bsp) Sync(end bool) bool {
	ktup := tuple.Tuple{bsp.globalStep}
	k := bsp.dir.Pack(ktup)
	// First, set my flag.
	bsp.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		v := tr.Get(k).MustGet()
		if len(v) == 0 {
			v = make([]byte, bsp.ncli)
		}

		v[bsp.me] = 1
		if end {
			v[bsp.me] = 2
		}
		tr.Set(k, v)
		return nil, nil
	})

	var sd syncData
	for {
		bsp.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			sd.nend = 0
			sd.ncont = 0
			v := tr.Get(k).MustGet()
			for _, st := range v {
				if st == 1 {
					sd.ncont += 1
				} else if st == 2 {
					sd.nend += 1
				}
			}

			if sd.ncont+sd.nend != int(bsp.ncli) {
				sd.w = tr.Watch(k)
			}
			return nil, nil
		})

		if sd.nend+sd.ncont == int(bsp.ncli) {
			bsp.globalStep += 1
			return sd.nend == int(bsp.ncli)
		} else {
			sd.w.BlockUntilReady()
		}
	}
}
