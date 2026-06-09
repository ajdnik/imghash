package vit

import (
	"math"
	"testing"
)

func TestLayerNorm_ZeroMeanUnitVar(t *testing.T) {
	// One row of 4 values; gamma=1, beta=0; output should have ~0 mean and ~1 var.
	x := []float32{1, 2, 3, 4}
	gamma := []float32{1, 1, 1, 1}
	beta := []float32{0, 0, 0, 0}
	layerNorm(x, gamma, beta, 1, 4, 1e-6)
	var sum float64
	for _, v := range x {
		sum += float64(v)
	}
	if math.Abs(sum/4) > 1e-3 {
		t.Errorf("normalized mean = %v, want ~0", sum/4)
	}
	var ss float64
	for _, v := range x {
		ss += float64(v) * float64(v)
	}
	// Biased variance, since LN divides by N not N-1.
	v := ss / 4
	if math.Abs(v-1) > 1e-3 {
		t.Errorf("normalized variance = %v, want ~1", v)
	}
}

func TestLayerNorm_GammaBeta(t *testing.T) {
	x := []float32{1, 2, 3, 4}
	gamma := []float32{2, 2, 2, 2}
	beta := []float32{5, 5, 5, 5}
	layerNorm(x, gamma, beta, 1, 4, 1e-6)
	// After unit-variance normalization, applying gamma=2 and beta=5 shifts the
	// distribution to mean 5 and variance 4.
	var sum float64
	for _, v := range x {
		sum += float64(v)
	}
	mean := sum / 4
	if math.Abs(mean-5) > 1e-3 {
		t.Errorf("mean = %v, want 5", mean)
	}
	var ss float64
	for _, v := range x {
		d := float64(v) - mean
		ss += d * d
	}
	if math.Abs(ss/4-4) > 1e-3 {
		t.Errorf("variance = %v, want 4", ss/4)
	}
}

func TestLayerNorm_KnownValues(t *testing.T) {
	// Two rows of 4 values; verify per-row independence.
	x := []float32{1, 2, 3, 4, 10, 20, 30, 40}
	gamma := []float32{1, 1, 1, 1}
	beta := []float32{0, 0, 0, 0}
	layerNorm(x, gamma, beta, 2, 4, 1e-6)
	// Per-row biased std = sqrt(((-1.5)² + (-0.5)² + 0.5² + 1.5²)/4) = sqrt(1.25).
	std := math.Sqrt(1.25)
	wantRow0 := []float32{
		float32(-1.5 / std), float32(-0.5 / std),
		float32(0.5 / std), float32(1.5 / std),
	}
	wantRow1 := []float32{
		float32(-15.0 / (10 * std)), float32(-5.0 / (10 * std)),
		float32(5.0 / (10 * std)), float32(15.0 / (10 * std)),
	}
	for i, v := range wantRow0 {
		if math.Abs(float64(x[i]-v)) > 1e-3 {
			t.Errorf("row0[%d] = %v, want %v", i, x[i], v)
		}
	}
	for i, v := range wantRow1 {
		if math.Abs(float64(x[4+i]-v)) > 1e-3 {
			t.Errorf("row1[%d] = %v, want %v", i, x[4+i], v)
		}
	}
}
