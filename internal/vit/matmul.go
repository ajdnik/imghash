// Package vit implements a hand-rolled forward pass for the DINOv2 ViT-S/14
// (with register tokens) Vision Transformer used by the DINOHash perceptual
// hash. The package is internal: only the dinohash package may import it.
package vit

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas32"
)

// gemm computes C = alpha*A*B + beta*C where A is [m, k], B is [k, n] and
// C is [m, n]. The matrices are stored in row-major float32 with contiguous
// strides equal to the column count.
func gemm(m, n, k int, alpha float32, a, b []float32, beta float32, c []float32) {
	A := blas32.General{Rows: m, Cols: k, Stride: k, Data: a}
	B := blas32.General{Rows: k, Cols: n, Stride: n, Data: b}
	C := blas32.General{Rows: m, Cols: n, Stride: n, Data: c}
	blas32.Gemm(blas.NoTrans, blas.NoTrans, alpha, A, B, beta, C)
}

// gemmNT computes C = A * Bᵀ + beta*C where A is [m, k], B is [n, k] and
// C is [m, n]. Used for the attention QKᵀ product.
func gemmNT(m, n, k int, alpha float32, a, b []float32, beta float32, c []float32) {
	A := blas32.General{Rows: m, Cols: k, Stride: k, Data: a}
	B := blas32.General{Rows: n, Cols: k, Stride: k, Data: b}
	C := blas32.General{Rows: m, Cols: n, Stride: n, Data: c}
	blas32.Gemm(blas.NoTrans, blas.Trans, alpha, A, B, beta, C)
}

// addRowBias adds the per-column bias vector b (length cols) to every row of
// the row-major matrix x (rows*cols entries).
func addRowBias(x []float32, b []float32, rows, cols int) {
	for r := 0; r < rows; r++ {
		row := x[r*cols : r*cols+cols]
		for i, v := range b {
			row[i] += v
		}
	}
}

// softmaxRow applies a numerically stable softmax to each row of x in place.
// x has shape [rows, cols] in row-major order.
func softmaxRow(x []float32, rows, cols int) {
	for r := 0; r < rows; r++ {
		row := x[r*cols : r*cols+cols]
		max := row[0]
		for _, v := range row[1:] {
			if v > max {
				max = v
			}
		}
		var sum float32
		for i, v := range row {
			e := float32(math.Exp(float64(v - max)))
			row[i] = e
			sum += e
		}
		inv := 1 / sum
		for i := range row {
			row[i] *= inv
		}
	}
}

// scaleRows multiplies every column of every row of x ([rows, cols]) by the
// per-column factor g (length cols). Used for LayerScale.
func scaleRows(x []float32, g []float32, rows, cols int) {
	for r := 0; r < rows; r++ {
		row := x[r*cols : r*cols+cols]
		for i, v := range g {
			row[i] *= v
		}
	}
}

// addInPlace computes a += b elementwise; both have length n.
func addInPlace(a, b []float32) {
	for i, v := range b {
		a[i] += v
	}
}
