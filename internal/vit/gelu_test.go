package vit

import (
	"math"
	"testing"
)

func TestGELU_KnownValues(t *testing.T) {
	// Hand-computed reference values using 0.5*x*(1+erf(x/√2)).
	tests := []struct {
		in, want float32
	}{
		{0, 0},
		{1, float32(0.5 * (1 + math.Erf(invSqrt2)))},
		{-1, float32(-0.5 * (1 + math.Erf(-invSqrt2)))},
		{2, float32(0.5 * 2 * (1 + math.Erf(2*invSqrt2)))},
		{-3, float32(0.5 * -3 * (1 + math.Erf(-3*invSqrt2)))},
	}
	for _, tt := range tests {
		x := []float32{tt.in}
		gelu(x)
		if math.Abs(float64(x[0]-tt.want)) > 1e-5 {
			t.Errorf("gelu(%v) = %v, want %v", tt.in, x[0], tt.want)
		}
	}
}

func TestGELU_AsymptoteHigh(t *testing.T) {
	// For large positive x, GELU(x) ≈ x.
	x := []float32{10}
	gelu(x)
	if math.Abs(float64(x[0]-10)) > 1e-4 {
		t.Errorf("gelu(10) = %v, want ≈ 10", x[0])
	}
}

func TestGELU_AsymptoteLow(t *testing.T) {
	// For large negative x, GELU(x) ≈ 0.
	x := []float32{-10}
	gelu(x)
	if math.Abs(float64(x[0])) > 1e-4 {
		t.Errorf("gelu(-10) = %v, want ≈ 0", x[0])
	}
}
