package main

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const v_max float64 = 25;		// Max speed [m/s]
const d float64 = 75; 			// Max distance between cars [m]
const M int = 10				    // Count of cars [cars]
const h float64 = 0.5				// Step size [step]
const t float64 = 40;		    // Time instant [s]
const g float64 = 5.0;		  // Speed of first car [m/s]
const di float64 = d;		    // Initial distances between cars [m]

func f(x float64) float64 {
	if x <= 0 {
		return 0
	} else if x < d {
		return x * v_max / d
	}
	return v_max
}

func p1() {
	// Create a plot
	p := plot.New()
	p.Title.Text = "Graph of f(x)"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	// Generate data points with f(x)
	pts := make(plotter.XYs, 100)
	for i := - 10; i < pts.Len() - 10; i++ {
		pts[i+10].X = float64(i)
		pts[i+10].Y = float64(f(float64(i)))
	}

	plotutil.AddLinePoints(p, "f(x)", pts)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "problem1.svg"); err != nil {
		panic(err)
	}
}

func p2() {
	const n int = int(t/h)  // Count of steps

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {			
		carPos[i - 1] = float64(i) * di
	}
	fmt.Println(carPos)
	for i := 0; i < n; i++ {
		nextCarPos, carSpeed := euler_step(carPos)
		fmt.Println(carSpeed)
		carPos = nextCarPos
	}

}

func get_car_speed(carPos float64, nextCarPos float64) float64 {
	return f(nextCarPos - carPos)		// Calculate speed of car i
}

func get_all_car_speeds(carPos []float64) []float64 {
	carSpeed := make([]float64, M)		// List of car speeds
	for i := 0; i < M - 1; i++ {
		carSpeed[i] = get_car_speed(carPos[i], carPos[i + 1]) 	// Calculate speed of car i
	}
	carSpeed[M - 1] = g		// Speed of car M is constant g
	return carSpeed
}

func get_new_car_pos(carPos []float64, carSpeed []float64) []float64 {
	newCarPos := make([]float64, M)		// List of car positions
	for i := 0; i < M; i++ {
		newCarPos[i] = carPos[i] + carSpeed[i] * h		
	}
	return newCarPos
}

func euler_step(carPos []float64) ([]float64, []float64) {
	carSpeed := get_all_car_speeds(carPos)
	newCarPos := get_new_car_pos(carPos, carSpeed)
	return newCarPos, carSpeed
}

func run_all() {
	// fmt.Println("Running problem 1"); p1()
	fmt.Println("Running problem 2"); p2()
}

func main() {
	run_all()
}
