package vit

import (
	"slices"
	"testing"
)

func TestBlockForward_ResidualAndLayerScale(t *testing.T) {
	for _, tc := range []struct {
		name     string
		ls1, ls2 []float32
		want     []float32
	}{
		{"without_scale", nil, nil, []float32{6, 9, 8, 11}},
		{"with_scale", []float32{2, 3}, []float32{4, 5}, []float32{19, 33, 21, 35}},
		{"zero_scale", []float32{0, 0}, []float32{0, 0}, []float32{1, 2, 3, 4}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Zero projection weights isolate the known attention and MLP biases.
			b := &Block{
				Dim: 2, Eps: Eps,
				Norm1Gamma: []float32{1, 1}, Norm1Beta: make([]float32, 2),
				Norm2Gamma: []float32{1, 1}, Norm2Beta: make([]float32, 2),
				LS1Gamma: tc.ls1, LS2Gamma: tc.ls2,
				Attn: &Attention{Dim: 2, NumHeads: 1, HeadDim: 2,
					QKVWeight: make([]float32, 12), QKVBias: make([]float32, 6),
					ProjWeight: make([]float32, 4), ProjBias: []float32{1, 2}},
				MLP: &MLP{Dim: 2, Hidden: 1,
					FC1Weight: make([]float32, 2), FC1Bias: make([]float32, 1),
					FC2Weight: make([]float32, 2), FC2Bias: []float32{4, 5}},
			}
			x := []float32{1, 2, 3, 4}
			b.Forward(x, 2)
			if !slices.Equal(x, tc.want) {
				t.Errorf("Forward = %v, want %v", x, tc.want)
			}
		})
	}
}
