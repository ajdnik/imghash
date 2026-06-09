package vit

import (
	"math"
	"testing"
)

func TestAttention_IdentityQKVProj(t *testing.T) {
	// seq=2, dim=2, heads=1, headDim=2.
	// QKV weight is the [6x2] identity-block matrix so that q/k/v each pass x through.
	// Proj weight is the [2x2] identity.
	a := &Attention{
		Dim: 2, NumHeads: 1, HeadDim: 2,
		QKVWeight: []float32{
			1, 0, // q row 0
			0, 1, // q row 1
			1, 0, // k row 0
			0, 1, // k row 1
			1, 0, // v row 0
			0, 1, // v row 1
		},
		QKVBias: []float32{0, 0, 0, 0, 0, 0},
		ProjWeight: []float32{
			1, 0,
			0, 1,
		},
		ProjBias: []float32{0, 0},
	}
	x := []float32{1, 0, 0, 1}
	got := a.Forward(x, 2)
	// Hand-computed expected:
	// scores = QKᵀ/√2 -> [[1/√2, 0], [0, 1/√2]]
	// softmax row -> p = e^(1/√2)/(e^(1/√2)+1)
	p := math.Exp(math.Sqrt2/2) / (math.Exp(math.Sqrt2/2) + 1)
	// out = softmax · V; V = [[1,0],[0,1]] so:
	// out[0] = [p, 1-p]
	// out[1] = [1-p, p]
	// proj is identity, so final = out.
	want := []float32{float32(p), float32(1 - p), float32(1 - p), float32(p)}
	for i, v := range want {
		if math.Abs(float64(got[i]-v)) > 1e-5 {
			t.Errorf("got[%d] = %v, want %v", i, got[i], v)
		}
	}
}

func TestAttention_BiasApplied(t *testing.T) {
	// Same as above but with a non-zero ProjBias so we know it was added.
	a := &Attention{
		Dim: 2, NumHeads: 1, HeadDim: 2,
		QKVWeight: []float32{
			1, 0, 0, 1,
			1, 0, 0, 1,
			1, 0, 0, 1,
		},
		QKVBias:    []float32{0, 0, 0, 0, 0, 0},
		ProjWeight: []float32{1, 0, 0, 1},
		ProjBias:   []float32{7, -3},
	}
	x := []float32{1, 0, 0, 1}
	got := a.Forward(x, 2)
	p := math.Exp(math.Sqrt2/2) / (math.Exp(math.Sqrt2/2) + 1)
	want := []float32{
		float32(p) + 7, float32(1-p) - 3,
		float32(1-p) + 7, float32(p) - 3,
	}
	for i, v := range want {
		if math.Abs(float64(got[i]-v)) > 1e-5 {
			t.Errorf("got[%d] = %v, want %v", i, got[i], v)
		}
	}
}
