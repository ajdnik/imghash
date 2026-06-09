package vit

import (
	"math"
	"testing"
)

func TestGemm(t *testing.T) {
	// A (2x3) * B (3x2) = C (2x2).
	a := []float32{1, 2, 3, 4, 5, 6}
	b := []float32{7, 8, 9, 10, 11, 12}
	c := make([]float32, 4)
	gemm(2, 2, 3, 1, a, b, 0, c)
	want := []float32{58, 64, 139, 154}
	for i, v := range want {
		if c[i] != v {
			t.Errorf("c[%d] = %v, want %v", i, c[i], v)
		}
	}
}

func TestGemmNT(t *testing.T) {
	// A (2x3), B (2x3) -> C (2x2) = A * Bᵀ.
	a := []float32{1, 2, 3, 4, 5, 6}
	b := []float32{7, 8, 9, 10, 11, 12}
	c := make([]float32, 4)
	gemmNT(2, 2, 3, 1, a, b, 0, c)
	// Row 0 of A · col 0 of Bᵀ = (1,2,3)·(7,8,9) = 50.
	// Row 0 of A · col 1 of Bᵀ = (1,2,3)·(10,11,12) = 68.
	// Row 1 of A · col 0 of Bᵀ = (4,5,6)·(7,8,9) = 122.
	// Row 1 of A · col 1 of Bᵀ = (4,5,6)·(10,11,12) = 167.
	want := []float32{50, 68, 122, 167}
	for i, v := range want {
		if c[i] != v {
			t.Errorf("c[%d] = %v, want %v", i, c[i], v)
		}
	}
}

func TestAddRowBias(t *testing.T) {
	x := []float32{1, 2, 3, 4, 5, 6}
	b := []float32{10, 20, 30}
	addRowBias(x, b, 2, 3)
	want := []float32{11, 22, 33, 14, 25, 36}
	for i, v := range want {
		if x[i] != v {
			t.Errorf("x[%d] = %v, want %v", i, x[i], v)
		}
	}
}

func TestSoftmaxRow_BasicSum(t *testing.T) {
	x := []float32{1, 2, 3, 4, 5, 6}
	softmaxRow(x, 2, 3)
	for r := 0; r < 2; r++ {
		var sum float32
		for c := 0; c < 3; c++ {
			sum += x[r*3+c]
		}
		if math.Abs(float64(sum-1)) > 1e-5 {
			t.Errorf("row %d sum = %v, want 1", r, sum)
		}
	}
}

func TestSoftmaxRow_NumericalStability(t *testing.T) {
	// Without max-subtraction this row would overflow to NaN/Inf.
	x := []float32{1000, 1000.001, -1000}
	softmaxRow(x, 1, 3)
	for i, v := range x {
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			t.Fatalf("x[%d] = %v: softmax must be stable for large logits", i, v)
		}
	}
	var sum float32
	for _, v := range x {
		sum += v
	}
	if math.Abs(float64(sum-1)) > 1e-5 {
		t.Errorf("sum = %v, want 1", sum)
	}
	if x[2] > 1e-6 {
		t.Errorf("expected x[2] ≈ 0 for very small logit, got %v", x[2])
	}
}

func TestScaleRows(t *testing.T) {
	x := []float32{1, 2, 3, 4, 5, 6}
	g := []float32{10, 100, 1000}
	scaleRows(x, g, 2, 3)
	want := []float32{10, 200, 3000, 40, 500, 6000}
	for i, v := range want {
		if x[i] != v {
			t.Errorf("x[%d] = %v, want %v", i, x[i], v)
		}
	}
}

func TestAddInPlace(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{10, 20, 30}
	addInPlace(a, b)
	want := []float32{11, 22, 33}
	for i, v := range want {
		if a[i] != v {
			t.Errorf("a[%d] = %v, want %v", i, a[i], v)
		}
	}
}
