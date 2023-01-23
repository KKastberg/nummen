package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func p3() {
	x := mat.NewDense(6, 1, []float64{150, 200, 300, 500, 1000, 2000});
	y := mat.NewDense(6, 1, []float64{2, 3, 4, 5, 6, 7});

	p3a5(x, y);
}

func p3a5(x *mat.Dense, y *mat.Dense) {
	// 1/U(x) - 1/8 = a/x
	A := mat.NewDense(6,1, nil)
	for i := 0; i < 6; i++ {
		A.Set(i, 0, 1/x.At(i,0));
	}

	b := mat.NewDense(6, 1, nil);
	for i := 0; i < 6; i++ {
		b.Set(i, 0, 1/y.At(i,0) - 1.0/8.0);
	}

	least_square(A, b);
}

func p3a6(x *mat.Dense, y *mat.Dense) {
	// 1/U(x) - 1/8 = a/x
	A := mat.NewDense(6,1, nil)
	for i := 0; i < 6; i++ {
		A.Set(i, 0, 1/x.At(i,0));
	}

	b := mat.NewDense(6, 1, nil);
	for i := 0; i < 6; i++ {
		b.Set(i, 0, 1/y.At(i,0) - 1.0/8.0);
	}

	least_square(A, b);
}

func p2() {
	a := mat.NewDense(7, 1, []float64{0, 0.5, 1, 1.5, 2, 2.99, 3});
	b := mat.NewDense(7, 1, []float64{0, 0.52, 1.09, 1.75, 2.45, 3.5, 4});
	
	// Least square
	A := mat.NewDense(7,2, nil)
	for i := 0; i < 7; i++ {
		A.Set(i, 0, a.At(i,0));
		A.Set(i, 1, a.At(i,0)*a.At(i,0));
	}
	least_square(A, b);

	// polynomial interpolation
	polynomial_interpolation(a, b);

}

func polynomial_interpolation(a *mat.Dense, b *mat.Dense) {
	// Create the interpolator
	print("Polynomial interpolation: \n");
}

func least_square(A *mat.Dense, b *mat.Dense) *mat.Dense {
	_, A_cols := A.Dims();
	_, b_cols := b.Dims();
	
	// Create the matrix (A^T*A)^-1
	A_T_A := mat.NewDense(A_cols, A_cols, nil)
	A_T_A.Mul(A.T(), A)
	A_T_A_inv := mat.NewDense(A_cols, A_cols, nil)
	A_T_A_inv.Inverse(A_T_A)
	fmt.Printf("A_T_A: %v\n", A_T_A)
	fmt.Printf("A_T_A_inv: %v\n", A_T_A_inv)

	// Create the matrix A^T*b
	A_T_b := mat.NewDense(A_cols, b_cols, nil)
	A_T_b.Mul(A.T(), b)
	fmt.Printf("A_T_b: %v\n", A_T_b)

	// Calculate the solution x
	x := mat.NewDense(A_cols, b_cols, nil)
	x.Mul(A_T_A_inv, A_T_b)

	fmt.Printf("A: %v\n", A)
	fmt.Printf("b: %v\n", b)
	fmt.Printf("x: %v\n", x)

	return x;
}

// Entry point
func main() {
	p3()
}
