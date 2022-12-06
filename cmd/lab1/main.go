package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// -------------------------
// Run this script with:
// cd cmd/lab1 && go run .
// -------------------------

// Parameters
const v_max float64 = 25 // Max speed [m/s]
const d float64 = 75     // Max distance between cars [m]
const M int = 10         // Count of cars [cars]
const h float64 = 0.1    // Step size [step]
const t float64 = 40     // Time instant [s]
const g float64 = 5      // Speed of first car [m/s]
const di float64 = d     // Initial distances between cars [m]
const iterMax int = 10   // Max iterations for fixedpoints method

func f(x float64) float64 {
	// Function f(x) as given in the assignment
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
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "uppgift1.png"); err != nil {
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

	// Calculate final car positions with forward euler
	for i := 0; i < n; i++ {
		nextCarPos, _ := euler_step(carPos)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	plot_position_graphs("Uppgift2", allCarPos)
}

func p3() {
	var n int = int(math.Round(t / h)) // Count of steps
	allCarPos := make([][]float64, n)  // List of all car positions

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {
		carPos[i-1] = float64(i) * di
	}

	// Calculate final car positions with forward euler
	for i := 0; i < n; i++ {
		nextCarPos, _ := euler_step(carPos)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	generate_car_plots("uppgift3", allCarPos)
}

func p7() {
	var n int = int(math.Round(t / h)) // Count of steps
	allCarPos := make([][]float64, n)  // List of all car positions

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {
		carPos[i-1] = float64(i) * di
	}

	// Calculate final car positions with fixedpoint iteration
	for i := 0; i < n; i++ {
		nextCarPos := fpi_step_all_cars(carPos, iterMax)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	plot_position_graphs("Uppgift7", allCarPos)
}

func fpi_step_all_cars(carPos []float64, _iterMax int) []float64 {
	// Calculate next car positions for all cars with fixedpoint iteration
	nextCarPos := make([]float64, M)
	nextCarPos[M-1] = carPos[M-1] + h*g
	for i := M - 2; i >= 0; i-- {
		nextCarPos[i] = fpi_step_single_car(carPos[i], nextCarPos[i+1], _iterMax)
	}
	return nextCarPos
}

func fpi_step_single_car(carPos float64, nextCarPos float64, _iterMax int) float64 {
	// Calculate next car position for a single car with fixedpoint iteration
	guess := carPos
	for i := 0; i < _iterMax; i++ {
		guess = fixedpoint_iteration(guess, carPos, nextCarPos)
	}
	return guess
}

func fixedpoint_iteration(guess float64, carPos float64, nextCarPos float64) float64 {
	// Calculate next fixedpoint guess for one car
	return carPos + h*f(nextCarPos-guess)
}

func p8A() {
	var n int = int(math.Round(t / h)) // Count of steps
	allCarPos := make([][]float64, n)  // List of all car positions

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {
		carPos[i-1] = float64(i) * di
	}

	// Calculate final car positions with back euler
	for i := 0; i < n; i++ {
		nextCarPos := back_euler_step(carPos)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	plot_position_graphs("Uppgift8", allCarPos)
}

func p8B() {
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
			nextCarPos := fpi_step_all_cars(carPos, i)
			carPos = nextCarPos
		}
		// Caclulate absolute error between ground truth (back euler) and fixed point iteration
		abs_error[i] = math.Abs(carPos[0] - ground_truth[0])
	}

	plot_fix_point_error(abs_error)
}

func back_euler_step(carPos []float64) []float64 {
	// Calculate next position of all cars with back euler
	nextCarPos := make([]float64, M)
	nextCarPos[M-1] = carPos[M-1] + h*g
	for i := M - 2; i >= 0; i-- {
		nextCarPos[i] = car_back_euler_step(carPos[i], nextCarPos[i+1])
	}
	return nextCarPos
}

func car_back_euler_step(carPos float64, nextCarPos float64) float64 {
	// Calculate next position of one specific car with back euler
	return (d*carPos + h*v_max*nextCarPos) / (d + h*v_max)
}

func get_car_speed(carPos float64, nextCarPos float64) float64 {
	// Calculate speed of one specific car
	return f(nextCarPos - carPos) // Calculate speed of car i
}

func get_all_car_speeds(carPos []float64) []float64 {
	// Calculate speed of all cars
	carSpeeds := make([]float64, M) // List of car speeds
	for i := 0; i < M-1; i++ {
		carSpeeds[i] = get_car_speed(carPos[i], carPos[i+1]) // Calculate speed of car i
	}
	carSpeeds[M-1] = g // Speed of car M is constant g
	return carSpeeds
}

func get_new_car_pos(carPos []float64, carSpeed []float64) []float64 {
	// Calculate new car positions for all cars using euler forward
	newCarPos := make([]float64, M) // List of car positions
	for i := 0; i < M; i++ {
		newCarPos[i] = carPos[i] + carSpeed[i]*h
	}
	return newCarPos
}

func euler_step(carPos []float64) ([]float64, []float64) {
	// Calculate one step with forward euler for all cars
	carSpeed := get_all_car_speeds(carPos)
	newCarPos := get_new_car_pos(carPos, carSpeed)
	return newCarPos, carSpeed
}


func run_all() {
	// Run all problems
	p1()
	fmt.Println("Problem 1 done")
	p2()
	fmt.Println("Problem 2 done")
	p3()	// Runs slow
	fmt.Println("Problem 3 done")
	p7()
	fmt.Println("Problem 7 done")
	p8A()  // Use back euler to calculate car positions
	fmt.Println("Problem 8A done")
	p8B()  // Compare back euler with fixed point euler
	fmt.Println("Problem 8B done")
}

// Entry point
func main() {
	run_all()
}
