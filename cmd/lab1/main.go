package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const v_max float64 = 25 // Max speed [m/s]
const d float64 = 75     // Max distance between cars [m]
const M int = 10         // Count of cars [cars]
const h float64 = 0.1    // Step size [step]
const t float64 = 40     // Time instant [s]
const g float64 = 5      // Speed of first car [m/s]
const di float64 = d     // Initial distances between cars [m]
const iterMax int = 4    // Max iterations for fixedpoints method

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
	for i := -10; i < pts.Len()-10; i++ {
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
	var n int = int(math.Round(t / h)) // Count of steps
	allCarPos := make([][]float64, n)  // List of all car positions

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {
		carPos[i-1] = float64(i) * di
	}
	for i := 0; i < n; i++ {
		nextCarPos, _ := euler_step(carPos)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	plot_position_graphs("problem 2", allCarPos)
	generate_car_plots("problem 2", allCarPos)

}

func p7() {
	var n int = int(math.Round(t / h)) // Count of steps
	allCarPos := make([][]float64, n)  // List of all car positions

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {
		carPos[i-1] = float64(i) * di
	}

	for i := 0; i < n; i++ {
		nextCarPos := next_time_step(carPos, iterMax)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	plot_position_graphs("Uppgift7", allCarPos)

}

func next_time_step(carPos []float64, _iterMax int) []float64 {
	nextCarPos := make([]float64, M)
	nextCarPos[M-1] = carPos[M-1] + h*g
	for i := M - 2; i >= 0; i-- {
		nextCarPos[i] = next_car_time_step(carPos[i], nextCarPos[i+1], _iterMax)
	}
	// fmt.Println(nextCarPos)
	return nextCarPos
}

func next_car_time_step(carPos float64, nextCarPos float64, _iterMax int) float64 {
	guess := carPos
	for i := 0; i < _iterMax; i++ {
		//fmt.Println(guess)
		guess = fixedpoint_iteration(guess, carPos, nextCarPos)
	}
	return guess
}

func fixedpoint_iteration(guess float64, carPos float64, nextCarPos float64) float64 {
	return carPos + h*f(nextCarPos-guess)
}

func p8() {
	var n int = int(math.Round(t / h)) // Count of steps
	allCarPos := make([][]float64, n)  // List of all car positions

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {
		carPos[i-1] = float64(i) * di
	}
	// fmt.Println(carPos)

	for i := 0; i < n; i++ {
		nextCarPos := back_euler_step(carPos)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	plot_position_graphs("Uppgift8", allCarPos)
}

func p8v2() {
	var n int = int(math.Round(t / h)) // Count of steps
	max_iter := 50
	abs_error := make([]float64, max_iter)

	// Generate list of initial car positions
	initialCarPos := make([]float64, M)
	for i := 1; i <= M; i++ {
		initialCarPos[i-1] = float64(i) * di
	}

	// Calculate final car positions with back euler x_1(40)
	carPos := initialCarPos[:]
	for i := 0; i < n; i++ {
		nextCarPos := back_euler_step(carPos)
		carPos = nextCarPos
	}
	ground_truth := carPos[:]

	// Calculate final car positions with fixed point euler x_1(40) for 1-50 fixed point iterations
	for i := 0; i < max_iter; i++ {

		carPos := initialCarPos[:]

		for j := 0; j < n; j++ {
			nextCarPos := next_time_step(carPos, i)
			carPos = nextCarPos
		}

		abs_error[i] = math.Abs(carPos[0] - ground_truth[0])
	}

	fmt.Println(abs_error)
	plot_fix_point_error(abs_error)
}

func back_euler_step(carPos []float64) []float64 {
	nextCarPos := make([]float64, M)
	nextCarPos[M-1] = carPos[M-1] + h*g
	for i := M - 2; i >= 0; i-- {
		nextCarPos[i] = car_back_euler_step(carPos[i], nextCarPos[i+1])
	}
	//fmt.Println(nextCarPos)
	return nextCarPos
}

func car_back_euler_step(carPos float64, nextCarPos float64) float64 {
	return (d*carPos + h*v_max*nextCarPos) / (d + h*v_max)
}

func get_car_speed(carPos float64, nextCarPos float64) float64 {
	return f(nextCarPos - carPos) // Calculate speed of car i
}

func get_all_car_speeds(carPos []float64) []float64 {
	carSpeed := make([]float64, M) // List of car speeds
	for i := 0; i < M-1; i++ {
		carSpeed[i] = get_car_speed(carPos[i], carPos[i+1]) // Calculate speed of car i
	}
	carSpeed[M-1] = g // Speed of car M is constant g
	return carSpeed
}

func get_new_car_pos(carPos []float64, carSpeed []float64) []float64 {
	newCarPos := make([]float64, M) // List of car positions
	for i := 0; i < M; i++ {
		newCarPos[i] = carPos[i] + carSpeed[i]*h
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
	fmt.Println("Running problem 7 & 8")
	p1()
	p2()
	p7()
	p8()
	p8v2()
}

func main() {
	run_all()
}
