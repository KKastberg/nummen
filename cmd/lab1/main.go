package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"os"

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
	allCarPos := make([][]float64, n)		// List of all car positions

	// Generate list of initial car positions
	carPos := make([]float64, M)
	for i := 1; i <= M; i++ {			
		carPos[i - 1] = float64(i) * di
	}
	for i := 0; i < n; i++ {
		nextCarPos, _ := euler_step(carPos)
		allCarPos[i] = carPos
		carPos = nextCarPos
	}

	plot_position_graphs(allCarPos)
	generate_car_plots(allCarPos)

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

func plot_position_graphs(allCarPos [][]float64) {
	p := plot.New()
	p.Title.Text = "Car Positions"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "x"
	p.BackgroundColor = color.Transparent

	// Generate data points with f(x)
	for i := 0; i < M; i++ {
		pts := make(plotter.XYs, len(allCarPos))
		for j := 0; j < len(allCarPos); j++ {
			pts[j].X = float64(j) * h
			pts[j].Y = allCarPos[j][i]
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			panic(err)
		}
		line.Color = color.RGBA{B: uint8(255*i/M), A: 255, R: uint8(255-255*i/M), G: uint8(255/2 + 255*i/M/2)}
		p.Add(line)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "problem1.png"); err != nil {
		panic(err)
	}
}

func generate_car_plots(allCarPos [][]float64) {
	for i := 0; i < len(allCarPos); i++ {
		p := plot.New()
		p.Title.Text = "Car Positions"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Car index"
		p.BackgroundColor = color.Transparent
		p.X.Min = 0
		p.X.Max = allCarPos[len(allCarPos)-1][len(allCarPos[0]) - 1] + 100

		// Generate data points with f(x)
		pts := make(plotter.XYs, len(allCarPos))
		for j := 0; j < len(allCarPos[0]); j++ {
			pts[j].Y = float64(j)
			pts[j].X = allCarPos[i][j]
		}
		plotutil.AddScatters(p, pts)

		// Save the plot to a PNG file.
		if err := p.Save(4*vg.Inch, 4*vg.Inch, fmt.Sprintf("problem3_%d.png", i)); err != nil {
			panic(err)
		}

	}
}

func generate_gif() {
	// Create a new animated GIF.
	anim := gif.GIF{LoopCount: 0}
	for i := 0; i < 80; i++ {
		// Open the image file.
		f, err := os.Open(fmt.Sprintf("problem3_%d.png", i))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// Decode the image.
		img, err := png.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)

		// Append the image to the GIF.
		anim.Delay = append(anim.Delay, 0)
		anim.Image = append(anim.Image, palettedImage) 
	}

	// Write the animated GIF to a file.
	f, err := os.OpenFile("problem3.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gif.EncodeAll(f, &anim)
}

func run_all() {
	// fmt.Println("Running problem 1"); p1()
	fmt.Println("Running problem 2"); p2()
}

func main() {
	run_all()
}
