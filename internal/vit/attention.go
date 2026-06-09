package vit

import "math"

// Attention is a multi-head self-attention layer matching the DINOv2 ViT
// implementation. Weights are stored in PyTorch layout: linear weights have
// shape [out_features, in_features].
type Attention struct {
	QKVWeight  []float32 // [3*Dim, Dim]
	QKVBias    []float32 // [3*Dim]
	ProjWeight []float32 // [Dim, Dim]
	ProjBias   []float32 // [Dim]
	Dim        int
	NumHeads   int
	HeadDim    int // Dim / NumHeads
}

// Forward returns scaled-dot-product multi-head attention of x.
// x has shape [seqLen, Dim]; the returned slice owns its own backing array.
func (a *Attention) Forward(x []float32, seqLen int) []float32 {
	d := a.Dim
	qkv := make([]float32, seqLen*3*d)
	gemmNT(seqLen, 3*d, d, 1, x, a.QKVWeight, 0, qkv)
	addRowBias(qkv, a.QKVBias, seqLen, 3*d)

	out := make([]float32, seqLen*d)
	hd := a.HeadDim
	scale := float32(1.0 / math.Sqrt(float64(hd)))

	packQ := make([]float32, seqLen*hd)
	packK := make([]float32, seqLen*hd)
	packV := make([]float32, seqLen*hd)
	scores := make([]float32, seqLen*seqLen)
	headOut := make([]float32, seqLen*hd)

	for h := 0; h < a.NumHeads; h++ {
		qOff := h * hd
		kOff := d + h*hd
		vOff := 2*d + h*hd
		for r := 0; r < seqLen; r++ {
			row := qkv[r*3*d:]
			copy(packQ[r*hd:r*hd+hd], row[qOff:qOff+hd])
			copy(packK[r*hd:r*hd+hd], row[kOff:kOff+hd])
			copy(packV[r*hd:r*hd+hd], row[vOff:vOff+hd])
		}
		gemmNT(seqLen, seqLen, hd, scale, packQ, packK, 0, scores)
		softmaxRow(scores, seqLen, seqLen)
		gemm(seqLen, hd, seqLen, 1, scores, packV, 0, headOut)
		for r := 0; r < seqLen; r++ {
			copy(out[r*d+h*hd:r*d+h*hd+hd], headOut[r*hd:r*hd+hd])
		}
	}

	result := make([]float32, seqLen*d)
	gemmNT(seqLen, d, d, 1, out, a.ProjWeight, 0, result)
	addRowBias(result, a.ProjBias, seqLen, d)
	return result
}
