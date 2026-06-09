package vit

import (
	"math"
	"testing"
)

func TestMLP_IdentityNoOpBias(t *testing.T) {
	// dim=2, hidden=2, identity weights, zero biases.
	// FC1: x -> x; GELU(x); FC2: GELU(x) -> GELU(x).
	m := &MLP{
		Dim: 2, Hidden: 2,
		FC1Weight: []float32{1, 0, 0, 1},
		FC1Bias:   []float32{0, 0},
		FC2Weight: []float32{1, 0, 0, 1},
		FC2Bias:   []float32{0, 0},
	}
	x := []float32{1, -1, 2, 3}
	got := m.Forward(x, 2)
	for i, v := range []float32{1, -1, 2, 3} {
		f := float64(v)
		want := float32(0.5 * f * (1 + math.Erf(f*invSqrt2)))
		if math.Abs(float64(got[i]-want)) > 1e-5 {
			t.Errorf("got[%d] = %v, want %v", i, got[i], want)
		}
	}
}

func TestMLP_BiasApplied(t *testing.T) {
	m := &MLP{
		Dim: 2, Hidden: 2,
		FC1Weight: []float32{1, 0, 0, 1},
		FC1Bias:   []float32{0, 0},
		FC2Weight: []float32{1, 0, 0, 1},
		FC2Bias:   []float32{10, -20},
	}
	x := []float32{0, 0}
	got := m.Forward(x, 1)
	// GELU(0) = 0; fc2 output = bias.
	want := []float32{10, -20}
	for i, v := range want {
		if math.Abs(float64(got[i]-v)) > 1e-5 {
			t.Errorf("got[%d] = %v, want %v", i, got[i], v)
		}
	}
}
