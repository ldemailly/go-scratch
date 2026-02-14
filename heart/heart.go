package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type heart struct {
	n int
}

func (h heart) Dims() (c, r int) { return h.n, h.n }
func (h heart) X(c int) float64  { return -1.5 + 3.0*float64(c)/float64(h.n-1) }
func (h heart) Y(r int) float64  { return -1.5 + 3.0*float64(r)/float64(h.n-1) }
func (h heart) Z(c, r int) float64 {
	x := h.X(c)
	y := h.Y(r)
	a := x*x + y*y - 1
	return a*a*a - x*x*y*y*y
}

func main() {
	p := plot.New()
	p.Title.Text = "Happy Valentine [(x^2+y^2-1)^3 = x^2 y^3]"
	grid := heart{n: 400}
	pal := palette.Heat(2, 1)
	contour := plotter.NewContour(grid, []float64{0}, pal)
	p.Add(contour)
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "heart.png"); err != nil {
		panic(err)
	}
}
