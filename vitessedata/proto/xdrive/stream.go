package xdrive

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
)

func ProtostreamRead(f io.Reader, pb proto.Message) error {
	var magic int32
	err := binary.Read(f, binary.LittleEndian, &magic)
	if err != nil {
		return err
	}

	if magic != 0x20aa30bb {
		return fmt.Errorf("Protostream read bad magic.")
	}

	var msgsz int32
	err = binary.Read(f, binary.LittleEndian, &msgsz)
	if err != nil {
		return err
	}

	buf := make([]byte, msgsz)
	rsz := int32(0)
	for rsz < msgsz {
		looprsz, err := f.Read(buf[rsz:])
		if err != nil {
			return err
		}
		rsz += int32(looprsz)
	}

	err = proto.Unmarshal(buf, pb)
	return err
}

func ProtostreamWrite(f io.Writer, pb proto.Message) error {
	msg, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	var magic, msgsz int32

	magic = 0x20aa30bb
	msgsz = int32(len(msg))

	err = binary.Write(f, binary.LittleEndian, &magic)
	if err != nil {
		return err
	}

	err = binary.Write(f, binary.LittleEndian, &msgsz)
	if err != nil {
		return err
	}

	if msgsz > 0 {
		wsz, err := f.Write(msg)
		if err != nil || int32(wsz) != msgsz {
			return fmt.Errorf("delim write short write msg")
		}
	}

	return nil
}

func ReadXMsg(f io.Reader) (*XMsg, error) {
	var msg XMsg
	err := ProtostreamRead(f, &msg)
	return &msg, err
}

func WriteXMsg(f io.Writer, msg *XMsg) error {
	return ProtostreamWrite(f, msg)
}
