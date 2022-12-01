package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func plot_position_graphs(name string, allCarPos [][]float64) {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Car positions (%s)", name)
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
		line.Color = color.RGBA{B: uint8(255 * i / M), A: 255, R: uint8(255 - 255*i/M), G: uint8(255/2 + 255*i/M/2)}
		p.Add(line)
	}

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, fmt.Sprintf("%s.png", name)); err != nil {
		panic(err)
	}
}

func plot_fix_point_error(errors []float64) {
	//p := plot.New()
	//p.Title.Text = "Errors"
	//p.X.Label.Text = "Iterations"
	//p.Y.Label.Text = "Error"
	//p.BackgroundColor = color.Transparent
	//p.Y.Scale = plot.LogScale{}

	//// Generate data points with f(x)
	//for i := 0; i < len(errors); i++ {

	//	pts := make(plotter.XYs, len(errors))
	//	pts[i].X = float64(i)
	//	pts[i].Y = float64(errors[i])

	//	line, err := plotter.NewLine(pts)
	//	if err != nil {
	//		panic(err)
	//	}
	//	line.Color = color.RGBA{B: uint8(255 * i / M), A: 255, R: uint8(255 - 255*i/M), G: uint8(255/2 + 255*i/M/2)}
	//	p.Add(line)
	//}

	//boxPlot(values)
	p := plot.New()
	p.Title.Text = "bar plot"
	// p.Y.Scale = plot.LogScale{}

	

	var values plotter.Values
	for i := 0; i < len(errors); i++ {
			values = append(values, math.Log10(errors[i]))
	}

	bar, err := plotter.NewBarChart(values, 15)
	if err != nil {
		panic(err)
	}
	p.Add(bar)


	if err := p.Save(5*vg.Inch, 5*vg.Inch, "bar.png"); err != nil {
		panic(err)
	}
}

func generate_car_plots(name string, allCarPos [][]float64) {
	for i := 0; i < len(allCarPos); i++ {
		p := plot.New()
		p.Title.Text = fmt.Sprintf("Car positions (%s)", name)
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Car index"
		p.BackgroundColor = color.Transparent
		p.X.Min = 0
		p.X.Max = allCarPos[len(allCarPos)-1][len(allCarPos[0])-1] + 100

		// Generate data points with f(x)
		pts := make(plotter.XYs, len(allCarPos))
		for j := 0; j < len(allCarPos[0]); j++ {
			pts[j].Y = float64(j)
			pts[j].X = allCarPos[i][j]
		}
		plotutil.AddScatters(p, pts)

		// Save the plot to a PNG file.
		if err := p.Save(4*vg.Inch, 4*vg.Inch, fmt.Sprintf("%s_%d.png", name, i)); err != nil {
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
