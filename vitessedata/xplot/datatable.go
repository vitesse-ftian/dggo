package xplot

import (
	"fmt"
	"math"
)

type DataTable interface {
	NRow() int
	NCol() int
	ColumnName(col int) string
	ColumnType(col int) string
	ColData(col int) []float64
	XLabels() []string

	Refresh()
}

func RangeToLabels(xbegin, xend float64, nlabel int) []string {
	step := (xend - xbegin) / float64(nlabel-1)
	var ret []string
	for i := 0; i < nlabel; i++ {
		ret = append(ret, fmt.Sprintf("%.1f", xbegin+step*float64(i)))
	}
	return ret
}

func DataTableYRange(dt DataTable) (minY, maxY float64) {
	minY = math.Inf(1)
	maxY = math.Inf(-1)
	for c := 0; c < dt.NCol(); c++ {
		data := dt.ColData(c)
		for r := 0; r < dt.NRow(); r++ {
			minY = math.Min(minY, data[r])
			maxY = math.Max(maxY, data[r])
		}
	}
	return
}

type Function struct {
	Name string
	F    func(float64) float64
}

type FunctionDataTable struct {
	xbegin  float64
	xend    float64
	xstep   float64
	nxlabel int
	cols    []Function
	data    [][]float64
}

func (f *FunctionDataTable) NRow() int {
	return int((f.xend-f.xbegin)/f.xstep) + 1
}

func (f *FunctionDataTable) NCol() int {
	return len(f.cols)
}

func (f *FunctionDataTable) ColumnName(i int) string {
	return f.cols[i].Name
}

func (f *FunctionDataTable) ColumnType(i int) string {
	return "float64"
}

func (f *FunctionDataTable) ColData(col int) []float64 {
	// Compute and cache.   F could be very expensive.
	nrow := f.NRow()
	if len(f.data[col]) != nrow {
		f.data[col] = make([]float64, nrow)
		for i := 0; i < nrow; i++ {
			f.data[col][i] = f.cols[col].F(f.xbegin + f.xstep*float64(i))
		}
	}
	return f.data[col]
}

func (f *FunctionDataTable) XLabels() []string {
	return RangeToLabels(f.xbegin, f.xend, f.nxlabel)
}

func (f *FunctionDataTable) Refresh() {
	for i := 0; i < f.NCol(); i++ {
		f.data[i] = make([]float64, 0)
	}
}

func (f *FunctionDataTable) SetNXLabel(nx int) *FunctionDataTable {
	f.nxlabel = nx
	return f
}

func (f *FunctionDataTable) AddFunc(ff Function) *FunctionDataTable {
	f.cols = append(f.cols, ff)
	f.data = append(f.data, make([]float64, 0))
	return f
}

func NewFuncDataTable(xbegin, xend, xstep float64) *FunctionDataTable {
	var fdt FunctionDataTable
	fdt.xbegin = xbegin
	fdt.xend = xend
	fdt.xstep = xstep
	fdt.nxlabel = 5
	return &fdt
}
