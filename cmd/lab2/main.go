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
	var rows, cols = vector.Dims()
	fmt.Println("Matrix (", rows, "x", cols, "): ")
	for i := 0; i < rows; i++ {
		fmt.Printf("%f ", vector.At(i, 0))
	}
	fmt.Println("")
}

func p3() {
	x := mat.NewDense(6, 1, []float64{150, 200, 300, 500, 1000, 2000})
	y := mat.NewDense(6, 1, []float64{2, 3, 4, 5, 6, 7})

	//p3a5(x, y)
	//p3a6a(x, y)
	p3a6b(x, y)
	// coeff_vec := polynomial_interpolation(x, y)
	// fmt.Println("coeff_vec=", coeff_vec)

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

	sol, residual_vec := least_square(A, b)
	a := sol.At(0, 0)

	// mse := mse()

	fmt.Println("a=", a, "rmse=", rmse(residual_vec))
}

func p3a6a(x *mat.Dense, y *mat.Dense) {
	// ln(8-y) = ln(a) + bln(x); d := ln(a)
	samples, _ := x.Dims()

	A := mat.NewDense(samples, 2, nil)
	for i := 0; i < samples; i++ {
		A.Set(i, 0, 1)
		A.Set(i, 1, math.Log(x.At(i, 0)))
	}

	b_ := mat.NewDense(samples, 1, nil)
	for i := 0; i < samples; i++ {
		b_.Set(i, 0, math.Log(8-y.At(i, 0)))
	}

	X, residual_vec:= least_square(A, b_)
	d := X.At(0, 0)
	b := X.At(1, 0)
	a := math.Exp(d)
	fmt.Println("a=", a, "b=", b, "rmse", rmse(residual_vec))
}

func p3a6b(x *mat.Dense, y *mat.Dense) {
	// y = 8 - ax^b
	samples, _ := x.Dims()

	precision := 1e-9999
	max_iterations := 1000000
	a0 := 5.0
	b0 := 0.1

	for n := 0; n < max_iterations; n++ {
		A := mat.NewDense(samples, 2, nil)
		for i := 0; i < samples; i++ {
			J0, J1 := J(x.At(i, 0), a0, b0)
			A.Set(i, 0, J0)
			A.Set(i, 1, J1)
		}
		// print_matrix(A)

		b_ := mat.NewDense(samples, 1, nil)
		for i := 0; i < samples; i++ {
			b_.Set(i, 0, y.At(i, 0) - U(x.At(i, 0), a0, b0))
		}
		// print_matrix(b_)

		X, _ := least_square(A, b_)
		// fmt.Println(mat.Norm(X, 1))

		if mat.Norm(X, 1) < precision {
			break
		}

		a0 = a0 + X.At(0, 0)
		b0 = b0 + X.At(1, 0)
	}
	// fmt.Println(n)
	fmt.Println("a=", a0, "b=", b0)
}

func U(x float64, a float64, b float64) float64 {
	return 8 - a*math.Pow(x, b)
}

func J(x float64, a float64, b float64) (float64, float64) {
	return -math.Pow(x, b), -a*math.Log(x)*math.Pow(x, b)
}





// func gaussNewton(x0 *mat.Dense, U func(*mat.Dense), J func(*mat.Dense), precision float64, max_iterations int) {
// 	x := x0
// 	for i := 0; i < max_iterations; i++ {
// 		x1 = x - Jx^-1 * Ux
// 		if norm(x1-x) < precision {
// 			break
// 		}
// 		x = x1
// 	}
// 	return x
// }

func p2() {
	a := mat.NewDense(7, 1, []float64{0, 0.5, 1, 1.5, 2, 2.99, 3})
	b := mat.NewDense(7, 1, []float64{0, 0.52, 1.09, 1.75, 2.45, 3.5, 4})

	// Least square
	A := mat.NewDense(7, 2, nil)
	for i := 0; i < 7; i++ {
		A.Set(i, 0, a.At(i, 0))
		A.Set(i, 1, a.At(i, 0)*a.At(i, 0))
	}
	least_vec, residual_vec := least_square(A, b)

	// polynomial interpolation
	coeff_vec := polynomial_interpolation(a, b)
	print_vector(coeff_vec)

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
	rmse := rmse(residual_vec)

	fmt.Printf("The RMSE for least-squares method in problem 2 is: %f\n", rmse)

}

func rmse(residual_vec *mat.Dense) float64 {
	// Calculate the errors
	// Finding least-square errors:
	// 1. Find the magnitude/norm of the residual vector r, ||r||
	residual_norm := residual_vec.Norm(2)

	// 2. Find the Squared Error SE = ||r||^2
	squared_error := math.Pow(residual_norm, 2)

	// 3. Root Mean Squared Error: RMSE = sqrt(SE/n)
	n, _ := residual_vec.Dims()
	rmse := math.Sqrt(squared_error / float64(n))

	return rmse
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

	// fmt.Printf("Vander matrix dimensions: %dx%d\n", rows, rows)
	// print_matrix(vander_matrix)

	y := mat.NewVecDense(rows, nil)
	_ = y.CopyVec(b.ColView(0))

	var c mat.VecDense
	c.SolveVec(vander_matrix, y)

	//print_vector(c)

	return c

}

func least_square(A *mat.Dense, b *mat.Dense) (*mat.Dense, *mat.Dense) {
	A_rows, A_cols := A.Dims()
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

	_, x_cols := x.Dims()

	//print_matrix(A)
	//print_matrix(x)

	Ax := mat.NewDense(A_rows, x_cols, nil)
	Ax.Mul(A, x)

	residual_vec := mat.NewDense(A_rows, x_cols, nil)
	residual_vec.Sub(b, Ax)

	return x, residual_vec
}

func mse(predictions, actual []float64) float64 {
	var sumSquaredError float64
	for i := 0; i < len(predictions); i++ {
			sumSquaredError += math.Pow(predictions[i]-actual[i], 2)
	}
	return sumSquaredError / float64(len(predictions))
}

// Entry point
func main() {
	p2()
}
