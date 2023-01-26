package main

import (
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func print_matrix(matrix *mat.Dense) {
	var rows, columns = matrix.Dims()
	fmt.Println("Matrix (", rows, "x", columns, "): ")
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			fmt.Printf("%f ", matrix.At(i, j))
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func print_vector(vector mat.VecDense) {
	var rows, _ = vector.Dims()

	for i := 0; i < rows; i++ {
		fmt.Printf("%f ", vector.At(i, 0))

	}
	fmt.Println("")
}

func p3() {
	x := mat.NewDense(6, 1, []float64{150, 200, 300, 500, 1000, 2000})
	y := mat.NewDense(6, 1, []float64{2, 3, 4, 5, 6, 7})

	p3a6(x, y)
}

func p3a5(x *mat.Dense, y *mat.Dense) {
	// 1/U(x) - 1/8 = a/x
	samples, _ := x.Dims()

	A := mat.NewDense(samples, 1, nil)
	for i := 0; i < samples; i++ {
		A.Set(i, 0, 1/x.At(i, 0))
	}

	b := mat.NewDense(samples, 1, nil)
	for i := 0; i < samples; i++ {
		b.Set(i, 0, 1/y.At(i, 0)-1.0/8.0)
	}

	a := least_square(A, b).At(0, 0)
	fmt.Println("a=", a)
}

func p3a6(x *mat.Dense, y *mat.Dense) {
	// ln(8-y) = ln(a) + cln(x); d := ln(c)
	samples, _ := x.Dims()

	A := mat.NewDense(samples, 2, nil)
	for i := 0; i < samples; i++ {
		A.Set(i, 0, math.Exp(1))
		A.Set(i, 1, math.Log(x.At(i, 0)))
	}
	print_matrix(A)

	b := mat.NewDense(samples, 1, nil)
	for i := 0; i < samples; i++ {
		b.Set(i, 0, math.Log(8-y.At(i, 0)))
	}
	print_matrix(b)

	X := least_square(A, b)
	a := X.At(0, 0)
	c := X.At(1, 0)
	fmt.Println("a=", a, "c=", c)
}

func p2() {
	a := mat.NewDense(7, 1, []float64{0, 0.5, 1, 1.5, 2, 2.99, 3})
	b := mat.NewDense(7, 1, []float64{0, 0.52, 1.09, 1.75, 2.45, 3.5, 4})

	// Least square
	A := mat.NewDense(7, 2, nil)
	for i := 0; i < 7; i++ {
		A.Set(i, 0, a.At(i, 0))
		A.Set(i, 1, a.At(i, 0)*a.At(i, 0))
	}
	least_vec := least_square(A, b)

	// polynomial interpolation
	coeff_vec := polynomial_interpolation(a, b)

	p := plot.New()
	p.Title.Text = "Uppgift 2"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"
	p.BackgroundColor = color.White

	// Create scatter graph
	data_points := make(plotter.XYs, 7)
	for i := range a.RawMatrix().Data {
		data_points[i].X = a.At(i, 0)
		data_points[i].Y = b.At(i, 0)
	}
	scatter, _ := plotter.NewScatter(data_points)

	// Create polynomial graph
	polynomial := plotter.NewFunction(func(x float64) float64 {
		var result float64

		x_1 := coeff_vec.At(1, 0)
		x_2 := coeff_vec.At(2, 0)
		x_3 := coeff_vec.At(3, 0)
		x_4 := coeff_vec.At(4, 0)
		x_5 := coeff_vec.At(5, 0)
		x_6 := coeff_vec.At(6, 0)

		result = x_1*x + x_2*math.Pow(x, 2) + x_3*math.Pow(x, 3) + x_4*math.Pow(x, 4) + x_5*math.Pow(x, 5) + x_6*math.Pow(x, 6)

		return result

	})
	polynomial.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}

	// Create least-squares graph
	least_squares_function := plotter.NewFunction(func(x float64) float64 {
		var result float64

		a := least_vec.At(0, 0)
		b := least_vec.At(1, 0)

		result = b*math.Pow(x, 2) + a*x

		return result

	})
	least_squares_function.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	p.Add(scatter, polynomial, least_squares_function, plotter.NewGrid())
	p.X.Min = 0
	p.X.Max = 3
	p.Y.Min = -3
	p.Y.Max = 4

	p.Legend.Add("Scatter", scatter)
	p.Legend.Add("Polynomial interpolation", polynomial)
	p.Legend.Add("Least-squares method", least_squares_function)

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "problem2.png"); err != nil {
		panic(err)
	}

}

func polynomial_interpolation(a *mat.Dense, b *mat.Dense) mat.VecDense {
	var rows, _ = a.Dims()
	vander_matrix := mat.NewDense(rows, rows, nil)

	for row := 0; row < rows; row++ {
		//A.Set(row, 0, 1) // First column in every row needs to be equal to 1
		for col := 0; col < rows; col++ {
			vander_matrix.Set(row, col, math.Pow(a.At(row, 0), float64(col)))
		}
	}

	fmt.Printf("Vander matrix dimensions: %dx%d\n", rows, rows)
	print_matrix(vander_matrix)

	y := mat.NewVecDense(rows, nil)
	_ = y.CopyVec(b.ColView(0))

	var c mat.VecDense
	c.SolveVec(vander_matrix, y)

	print_vector(c)

	return c

}

func least_square(A *mat.Dense, b *mat.Dense) *mat.Dense {
	_, A_cols := A.Dims()
	_, b_cols := b.Dims()

	// Create the matrix (A^T*A)^-1
	A_T_A := mat.NewDense(A_cols, A_cols, nil)
	A_T_A.Mul(A.T(), A)
	A_T_A_inv := mat.NewDense(A_cols, A_cols, nil)
	A_T_A_inv.Inverse(A_T_A)

	// Create the matrix A^T*b
	A_T_b := mat.NewDense(A_cols, b_cols, nil)
	A_T_b.Mul(A.T(), b)

	// Calculate the solution x
	x := mat.NewDense(A_cols, b_cols, nil)
	x.Mul(A_T_A_inv, A_T_b)

	return x
}

// Entry point
func main() {
	p3()
}
