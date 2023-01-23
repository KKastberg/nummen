package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

func print_matrix(matrix *mat.Dense) {
	var rows, _ = matrix.Dims()

	for i := 0; i < rows; i++ {
		for j := 0; j < rows+1; j++ {
			fmt.Printf("%f ", matrix.At(i, j))
		}
		fmt.Println("\n")
	}

}

func p2() {
	a := mat.NewDense(7, 1, []float64{0, 0.5, 1, 1.5, 2, 2.99, 3})
	b := mat.NewDense(7, 1, []float64{0, 0.52, 1.09, 1.75, 2.45, 3.5, 4})

	// Least square
	//least_square(a, b)

	// polynomial interpolation
	polynomial_interpolation(a, b)

}

func polynomial_interpolation(a *mat.Dense, b *mat.Dense) {
	var rows, _ = a.Dims()
	vander_matrix := mat.NewDense(rows, rows+1, nil)

	for row := 0; row < rows; row++ {
		//A.Set(row, 0, 1) // First column in every row needs to be equal to 1
		for col := 0; col < rows+1; col++ {
			vander_matrix.Set(row, col, math.Pow(a.At(row, 0), float64(col)))
		}
	}

	print_matrix(vander_matrix)

	 := mat.NewDense(rows+1, 1, nil)


}

func least_square(a *mat.Dense, b *mat.Dense) *mat.Dense {
	A := create_A(a)

	// Create the matrix (A^T*A)^-1
	A_T_A := mat.NewDense(2, 2, nil)
	A_T_A.Mul(A.T(), A)
	A_T_A_inv := mat.NewDense(2, 2, nil)
	A_T_A_inv.Inverse(A_T_A)
	fmt.Printf("A_T_A: %v\n", A_T_A)
	fmt.Printf("A_T_A_inv: %v\n", A_T_A_inv)

	// Create the matrix A^T*b
	A_T_b := mat.NewDense(2, 1, nil)
	A_T_b.Mul(A.T(), b)
	fmt.Printf("A_T_b: %v\n", A_T_b)

	// Calculate the solution x
	x := mat.NewDense(2, 1, nil)
	x.Mul(A_T_A_inv, A_T_b)

	fmt.Printf("a: %v\n", a)
	fmt.Printf("A: %v\n", A)
	fmt.Printf("b: %v\n", b)
	fmt.Printf("x: %v\n", x)

	return A
}

func create_A(a *mat.Dense) *mat.Dense {
	A := mat.NewDense(7, 2, nil)
	for i := 0; i < 7; i++ {
		A.Set(i, 0, a.At(i, 0))
		A.Set(i, 1, a.At(i, 0)*a.At(i, 0))
	}
	return A
}

// Entry point
func main() {
	p2()
}
