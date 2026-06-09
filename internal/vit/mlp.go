package vit

// MLP is the per-block feed-forward network: fc1 -> GELU -> fc2.
// Linear weights are stored in PyTorch layout [out_features, in_features].
type MLP struct {
	FC1Weight []float32 // [Hidden, Dim]
	FC1Bias   []float32 // [Hidden]
	FC2Weight []float32 // [Dim, Hidden]
	FC2Bias   []float32 // [Dim]
	Dim       int
	Hidden    int
}

// Forward returns the MLP applied row-wise to x ([seqLen, Dim]).
// The returned slice owns its own backing array.
func (m *MLP) Forward(x []float32, seqLen int) []float32 {
	h := make([]float32, seqLen*m.Hidden)
	gemmNT(seqLen, m.Hidden, m.Dim, 1, x, m.FC1Weight, 0, h)
	addRowBias(h, m.FC1Bias, seqLen, m.Hidden)
	gelu(h)
	out := make([]float32, seqLen*m.Dim)
	gemmNT(seqLen, m.Dim, m.Hidden, 1, h, m.FC2Weight, 0, out)
	addRowBias(out, m.FC2Bias, seqLen, m.Dim)
	return out
}
