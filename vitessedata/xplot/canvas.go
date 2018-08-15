/*
XPlot
*/

package xplot

import (
	"fmt"
	"github.com/buger/goterm"
	"math"
	"strings"
)

func rangeToLabels(xbegin, xend float64, nlabel int) []string {
	step := (xend - xbegin) / float64(nlabel-1)
	var ret []string
	for i := 0; i < nlabel; i++ {
		ret = append(ret, fmt.Sprintf("%.1f", xbegin+step*float64(i)))
	}
	return ret
}

func dataRange(data []float64) (minY, maxY float64) {
	minY = math.Inf(1)
	maxY = math.Inf(-1)
	for _, v := range data {
		minY = math.Min(minY, v)
		maxY = math.Max(maxY, v)
	}
	return
}

func dtXLabels(dt DataTable, n int) []string {
	col := dt.Col(0)
	if col.Text() != nil {
		return col.Text()
	} else {
		minv, maxv := dataRange(col.Data())
		return rangeToLabels(minv, maxv, n)
	}
}

func dtYRange(dt DataTable) (minY, maxY float64) {
	minY = math.Inf(1)
	maxY = math.Inf(-1)
	for i := 1; i < dt.NCol(); i++ {
		col := dt.Col(i)
		if col.Data() != nil {
			low, high := dataRange(col.Data())
			minY = math.Min(minY, low)
			maxY = math.Max(maxY, high)
		}
	}
	return
}

func min(x, y int) int {
	if x >= y {
		return y
	}
	return x
}

func max(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

type Canvas struct {
	width  int
	height int
	buf    [][]string

	paddingLeft  int
	paddingRight int
	paddingY     int
}

func NewCanvasSize(ww, hh int) (*Canvas, error) {
	h := goterm.Height()
	w := goterm.Width()
	if h < 20 || w < 60 {
		return nil, fmt.Errorf("minimal console size is 60x20.")
	}

	if h > w/3 {
		h = w / 3
	} else {
		h = h - 2
	}

	var c Canvas
	c.width = min(w, ww)
	c.height = min(h, hh)
	c.paddingY = 2

	c.createBuf()
	return &c, nil
}

func NewCanvas() (*Canvas, error) {
	return NewCanvasSize(1000000, 1000000)
}

func (c *Canvas) createBuf() {
	c.buf = make([][]string, c.height)
	for i := 0; i < c.height; i++ {
		c.buf[i] = make([]string, c.width)
		for j := 0; j < c.width; j++ {
			c.buf[i][j] = " "
		}
	}
}

func (c *Canvas) plotArea() [][]string {
	area := make([][]string, c.height-c.paddingY)
	for i := 0; i < c.height-c.paddingY; i++ {
		area[i] = c.buf[i+c.paddingY][c.paddingLeft : c.width-c.paddingRight]
	}
	return area
}

func (c *Canvas) writeText(x, y int, s string) {
	for i, ss := range strings.Split(s, "") {
		c.buf[x][y+i] = ss
	}
}

func (c *Canvas) writeTextLeft(x, y int, s string) {
	ss := strings.Split(s, "")
	ll := len(ss)
	for i, sss := range ss {
		c.buf[x][y-ll+i] = sss
	}
}

func (c *Canvas) drawAxis(dt DataTable) {
	xl := dtXLabels(dt, 5)
	miny, maxy := dtYRange(dt)
	yl := rangeToLabels(miny, maxy, c.height-c.paddingY)

	c.paddingLeft = len(xl[0])/2 + 1
	c.paddingRight = len(xl[len(xl)-1])/2 + 1

	for _, ss := range yl {
		c.paddingLeft = max(c.paddingLeft, len(ss)+1)
	}
	c.paddingRight += (c.width - c.paddingLeft - c.paddingRight) % (len(xl) - 1)

	// draw x axis.
	xstep := (c.width - c.paddingLeft - c.paddingRight) / (len(xl) - 1)
	for idx, ss := range xl {
		c.writeText(0, c.paddingLeft+xstep*idx-len(ss)/2, ss)
	}
	for i := c.paddingLeft; i < c.width-c.paddingRight; i++ {
		c.buf[1][i] = "-"
	}

	for y := 2; y < c.height; y++ {
		if (y-2)%5 == 0 {
			c.writeTextLeft(y, c.paddingLeft, fmt.Sprintf("%s|", yl[y-2]))
		} else {
			c.writeTextLeft(y, c.paddingLeft, "|")
		}
	}
}

func (c *Canvas) Draw() (out string) {
	for row := c.height - 1; row >= 0; row-- {
		out += strings.Join(c.buf[row], "") + "\n"
	}
	return
}

type LineChart struct {
}

type ScatterPlot struct {
}

type BarChart struct {
}

type Histogram struct {
}
