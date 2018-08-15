package xplot

type DataCol interface {
	Name() string
	Color() int
	Symbol() string
	Data() []float64
	Text() []string
}

type DataTable interface {
	NCol() int
	Col(i int) DataCol
	Refresh()
}

type Function struct {
	Name   string
	F      func(float64) float64
	Color  int
	Symbol string
}

type dataRangeCol struct {
	name  string
	begin float64
	end   float64
	step  float64
	data  []float64
}

func (dr *dataRangeCol) init(begin, end, step float64) {
	dr.name = "x"
	dr.begin = begin
	dr.end = end
	dr.step = step
	for v := begin; v <= end; v += step {
		dr.data = append(dr.data, v)
	}
}

func (dr *dataRangeCol) Name() string {
	return dr.name
}
func (dr *dataRangeCol) Color() int {
	return 0
}
func (dr *dataRangeCol) Symbol() string {
	return ""
}
func (dr *dataRangeCol) Data() []float64 {
	return dr.data
}
func (dr *dataRangeCol) Text() []string {
	return nil
}

type funcCol struct {
	dr   *dataRangeCol
	f    Function
	data []float64
}

func newFuncCol(dr *dataRangeCol, f Function) *funcCol {
	return &funcCol{dr: dr, f: f}
}

func (fc *funcCol) Name() string {
	return fc.f.Name
}

func (fc *funcCol) Color() int {
	return fc.f.Color
}

func (fc *funcCol) Symbol() string {
	return fc.f.Symbol
}

func (fc *funcCol) Data() []float64 {
	if fc.data == nil {
		for _, x := range fc.dr.Data() {
			fc.data = append(fc.data, fc.f.F(x))
		}
	}
	return fc.data
}

func (fc *funcCol) Text() []string {
	return nil
}

type FunctionDataTable struct {
	x  dataRangeCol
	ys []*funcCol
}

func (f *FunctionDataTable) NCol() int {
	return len(f.ys) + 1
}

func (f *FunctionDataTable) Col(i int) DataCol {
	if i == 0 {
		return &f.x
	}
	return f.ys[i-1]
}

func (f *FunctionDataTable) Refresh() {
	for _, y := range f.ys {
		y.data = nil
	}
}

func NewFuncDataTable(xbegin, xend, xstep float64) *FunctionDataTable {
	var fdt FunctionDataTable
	fdt.x.init(xbegin, xend, xstep)
	return &fdt
}

func (f *FunctionDataTable) AddFunc(ff Function) *FunctionDataTable {
	f.ys = append(f.ys, newFuncCol(&f.x, ff))
	return f
}
