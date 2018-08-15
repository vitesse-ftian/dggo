package xplot

import (
	"fmt"
	tm "github.com/buger/goterm"
	"math"
	"testing"
)

func TestFuncData(t *testing.T) {
	f1 := Function{
		Name: "f1",
		F: func(x float64) float64 {
			return 1.0
		},
	}

	fx := Function{
		Name: "x",
		F: func(x float64) float64 {
			return x
		},
	}

	fxplus := Function{
		Name: "x+1",
		F: func(x float64) float64 {
			return x + 1.0
		},
		Color:  tm.RED,
		Symbol: "x",
	}

	t.Run("f", func(t *testing.T) {
		fdt := NewFuncDataTable(0, 10, 0.1)
		fdt.AddFunc(f1).AddFunc(fx).AddFunc(fxplus)

		fmt.Printf("Functions Datatable.\n", fdt.NCol())

		for c := 0; c < fdt.NCol(); c++ {
			col := fdt.Col(c)
			fmt.Printf("%d-th Function, name %s, color %d, symbol %s.\n", c, col.Name(), col.Color(), col.Symbol())
			data := col.Data()
			fmt.Printf("first 3 values %v, last 3 values %v.\n", data[:3], data[len(data)-3:])
		}
	})

	t.Run("goterm", func(t *testing.T) {
		chart := tm.NewLineChart(200, 20)
		data := new(tm.DataTable)
		data.AddColumn("Time")
		data.AddColumn("Sin(x)")
		data.AddColumn("Cos(x+1)")
		for i := 0.1; i < 10; i += 0.1 {
			data.AddRow(i, math.Sin(i), math.Cos(i+1))
		}
		tm.Println(chart.Draw(data))
		tm.Flush()
		tm.MoveCursorUp(22)
		tm.Println(chart.Draw(data))
		tm.Flush()
	})

	t.Run("canvas", func(t *testing.T) {
		canvas, err := NewCanvasSize(200, 20)
		if err != nil {
			t.Error(err)
		}

		dt := NewFuncDataTable(0, 10, 0.1)
		dt.AddFunc(f1).AddFunc(fx).AddFunc(fxplus)
		canvas.drawAxis(dt)

		tm.Println(canvas.Draw())
		tm.Flush()
	})
}
