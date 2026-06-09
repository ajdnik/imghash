package vit

// Architecture constants for DINOv2 ViT-S/14 with 4 register tokens.
// These match the upstream configuration in facebookresearch/dinov2.
const (
	ImageSize    = 224
	PatchSize    = 14
	Channels     = 3
	Dim          = 384
	NumHeads     = 6
	HeadDim      = 64 // Dim / NumHeads
	MLPHidden    = 1536
	NumBlocks    = 12
	NumPatches   = (ImageSize / PatchSize) * (ImageSize / PatchSize) // 256
	NumRegisters = 4
	SeqLen       = 1 + NumRegisters + NumPatches // 261
	PosEmbedLen  = 1 + NumPatches                // 257 (CLS + patches; registers get no pos embed)
	HashBits     = 96
	Eps          = float32(1e-6)
)

// Model holds all weights for the ViT-S/14+reg backbone fused with the 96-bit
// PCA-aligned hash head. All slices are owned by the model; callers must not
// mutate them.
type Model struct {
	PatchEmbedW []float32 // [Dim, Channels*PatchSize*PatchSize]
	PatchEmbedB []float32 // [Dim]
	CLSToken    []float32 // [Dim]
	RegTokens   []float32 // [NumRegisters * Dim]
	PosEmbed    []float32 // [PosEmbedLen * Dim]
	Blocks      []*Block
	NormGamma   []float32 // [Dim]
	NormBeta    []float32 // [Dim]
	HeadW       []float32 // [HashBits * Dim]
	HeadB       []float32 // [HashBits]
}

// Forward runs the full DINOHash pipeline on a single image already
// preprocessed into NCHW float32 with shape [Channels, ImageSize, ImageSize].
// It returns the 96-element pre-binarization vector; the caller is responsible
// for sign-thresholding into a 12-byte hash.
func (m *Model) Forward(input []float32) [HashBits]float32 {
	// 1. Patch embed -> [NumPatches, Dim].
	patches := patchEmbed(input, m.PatchEmbedW, m.PatchEmbedB, ImageSize, PatchSize, Channels, Dim)

	// 2. Build token sequence [CLS, REG_0..REG_3, PATCH_0..PATCH_255] of shape [SeqLen, Dim].
	//    DINOv2's prepare_tokens_with_masks order: cls + pos[0:1], then registers
	//    (no pos), then patches + pos[1:].
	x := make([]float32, SeqLen*Dim)
	// Row 0: CLS + pos[0]
	for i := 0; i < Dim; i++ {
		x[i] = m.CLSToken[i] + m.PosEmbed[i]
	}
	// Rows 1..1+NumRegisters: register tokens, no pos embed
	for r := 0; r < NumRegisters; r++ {
		dst := x[(1+r)*Dim : (1+r+1)*Dim]
		src := m.RegTokens[r*Dim : (r+1)*Dim]
		copy(dst, src)
	}
	// Rows 1+NumRegisters..: patches + pos[1:NumPatches+1]
	for p := 0; p < NumPatches; p++ {
		row := x[(1+NumRegisters+p)*Dim : (1+NumRegisters+p+1)*Dim]
		patch := patches[p*Dim : (p+1)*Dim]
		pos := m.PosEmbed[(1+p)*Dim : (1+p+1)*Dim]
		for i := 0; i < Dim; i++ {
			row[i] = patch[i] + pos[i]
		}
	}

	// 3. 12 transformer blocks.
	for _, b := range m.Blocks {
		b.Forward(x, SeqLen)
	}

	// 4. Final LN, then take CLS (row 0).
	layerNorm(x, m.NormGamma, m.NormBeta, SeqLen, Dim, Eps)
	cls := x[:Dim]

	// 5. Head linear: [1, Dim] @ [HashBits, Dim]ᵀ + bias -> [HashBits].
	var out [HashBits]float32
	gemmNT(1, HashBits, Dim, 1, cls, m.HeadW, 0, out[:])
	for i, b := range m.HeadB {
		out[i] += b
	}
	return out
}
