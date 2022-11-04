package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const v_max = 25;		// m/s
const d = 75; 			// m

func f(x float64) float64 {
	if x <= 0 {
		return 0
	} else if x < d {
		return x * v_max / d
	}
	return v_max
}

func main() {
	p := plot.New()
	
	p.Title.Text = "Graph of f(x)"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	pts := make(plotter.XYs, 100)
	for i := - 10; i < pts.Len() - 10; i++ {
		pts[i+10].X = float64(i)
		pts[i+10].Y = float64(f(float64(i)))
	}

	plotutil.AddLinePoints(p, "f(x)", pts)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "problem1.png"); err != nil {
		panic(err)
	}
}
