package vit

import (
	"math"
	"testing"
)

func TestModelForward_TokenAssemblyAndHead(t *testing.T) {
	m := &Model{
		PatchEmbedW: make([]float32, Dim*Channels*PatchSize*PatchSize),
		PatchEmbedB: make([]float32, Dim),
		CLSToken:    make([]float32, Dim), RegTokens: make([]float32, NumRegisters*Dim),
		PosEmbed:  make([]float32, PosEmbedLen*Dim),
		NormGamma: make([]float32, Dim), NormBeta: make([]float32, Dim),
		HeadW: make([]float32, HashBits*Dim), HeadB: make([]float32, HashBits),
	}
	b := &Block{
		Dim: Dim, Eps: 0,
		Norm1Gamma: make([]float32, Dim), Norm1Beta: make([]float32, Dim),
		Norm2Gamma: make([]float32, Dim), Norm2Beta: make([]float32, Dim),
		Attn: &Attention{Dim: Dim, NumHeads: NumHeads, HeadDim: HeadDim,
			QKVWeight: make([]float32, 3*Dim*Dim), QKVBias: make([]float32, 3*Dim),
			ProjWeight: make([]float32, Dim*Dim), ProjBias: make([]float32, Dim)},
		MLP: &MLP{Dim: Dim, Hidden: 1,
			FC1Weight: make([]float32, Dim), FC1Bias: make([]float32, 1),
			FC2Weight: make([]float32, Dim), FC2Bias: make([]float32, Dim)},
	}
	m.Blocks = []*Block{b}
	// Break channel symmetry so final normalization retains differences in
	// token assembly and attention, rather than reducing every channel to +/-1.
	b.Attn.ProjBias[0] = 1
	// Alternating channels have mean zero. CLS = +/-4, registers = +/-1,
	// patches = -/+1 after positional embeddings. Uniform attention averages
	// their normalized values, so its CLS contribution is +/- (5-256)/261.
	for i := 0; i < Dim; i++ {
		sign := float32(1)
		if i%2 != 0 {
			sign = -1
		}
		m.CLSToken[i] = 3 * sign
		m.PosEmbed[i] = sign
		m.PatchEmbedB[i] = 2 * sign
		for r := 0; r < NumRegisters; r++ {
			m.RegTokens[r*Dim+i] = sign
		}
		for p := 0; p < NumPatches; p++ {
			m.PosEmbed[(p+1)*Dim+i] = -3 * sign
		}
		b.Norm1Gamma[i] = 1
		b.Norm2Gamma[i] = 1
		b.Attn.QKVWeight[(2*Dim+i)*Dim+i] = 1 // V is identity; Q and K are zero.
		b.Attn.ProjWeight[i*Dim+i] = 1
		m.NormGamma[i] = 2
		m.NormBeta[i] = 0.25
	}
	for i := 0; i < HashBits; i++ {
		m.HeadW[i*Dim+i] = 3
		m.HeadB[i] = float32(i) / 10
	}
	out := m.Forward(make([]float32, Channels*ImageSize*ImageSize))
	cls := 4 + float64(1+NumRegisters-NumPatches)/float64(SeqLen)
	mean := 1.0 / Dim
	variance := cls*cls + 2*cls/Dim + 1.0/Dim - mean*mean
	for i, got := range out {
		sign := 1.0
		if i%2 != 0 {
			sign = -1
		}
		value := sign * cls
		if i == 0 {
			value++
		}
		normalized := (value - mean) / math.Sqrt(variance+float64(Eps))
		want := 3*(2*normalized+0.25) + float64(m.HeadB[i])
		if math.Abs(float64(got)-want) > 1e-5 {
			t.Errorf("output[%d] = %v, want %v", i, got, want)
		}
	}
}
