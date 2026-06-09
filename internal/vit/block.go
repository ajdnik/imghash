package vit

// Block is a single pre-LN transformer block:
//
//	x = x + ls1 ⊙ attn(LN(x))
//	x = x + ls2 ⊙ mlp(LN(x))
//
// LayerScale (ls1, ls2) is applied only if the gamma slice is non-nil. DINOv2
// ViT-S/14+reg uses LayerScale with init_value=1e-5.
type Block struct {
	Norm1Gamma []float32 // [Dim]
	Norm1Beta  []float32 // [Dim]
	Attn       *Attention
	LS1Gamma   []float32 // [Dim] or nil if not present
	Norm2Gamma []float32 // [Dim]
	Norm2Beta  []float32 // [Dim]
	MLP        *MLP
	LS2Gamma   []float32 // [Dim] or nil if not present
	Dim        int
	Eps        float32
}

// Forward applies the block to x ([seqLen, Dim]) in place.
func (b *Block) Forward(x []float32, seqLen int) {
	h := make([]float32, seqLen*b.Dim)
	copy(h, x)
	layerNorm(h, b.Norm1Gamma, b.Norm1Beta, seqLen, b.Dim, b.Eps)
	h = b.Attn.Forward(h, seqLen)
	if b.LS1Gamma != nil {
		scaleRows(h, b.LS1Gamma, seqLen, b.Dim)
	}
	addInPlace(x, h)

	h2 := make([]float32, seqLen*b.Dim)
	copy(h2, x)
	layerNorm(h2, b.Norm2Gamma, b.Norm2Beta, seqLen, b.Dim, b.Eps)
	h2 = b.MLP.Forward(h2, seqLen)
	if b.LS2Gamma != nil {
		scaleRows(h2, b.LS2Gamma, seqLen, b.Dim)
	}
	addInPlace(x, h2)
}
