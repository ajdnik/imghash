package vit

// patchEmbed runs the non-overlapping 14x14 patch convolution on an NCHW
// [3, 224, 224] image. It is equivalent to:
//
//	Conv2d(3 -> Dim, kernel=PatchSize, stride=PatchSize) without padding.
//
// The implementation flattens each spatial patch into a 3*PatchSize*PatchSize
// vector (im2col) and then runs a single GEMM against the convolution kernel
// stored in PyTorch layout [Dim, 3, PatchSize, PatchSize] (flattened to
// [Dim, 3*PatchSize*PatchSize]).
func patchEmbed(input, weight, bias []float32, imageSize, patchSize, channels, dim int) []float32 {
	patchSpan := channels * patchSize * patchSize
	gridSide := imageSize / patchSize
	numPatches := gridSide * gridSide
	patches := make([]float32, numPatches*patchSpan)
	pixelCount := imageSize * imageSize
	for py := 0; py < gridSide; py++ {
		for px := 0; px < gridSide; px++ {
			patchIdx := py*gridSide + px
			dstBase := patchIdx * patchSpan
			for kc := 0; kc < channels; kc++ {
				channelBase := kc * pixelCount
				dstCBase := dstBase + kc*patchSize*patchSize
				for ky := 0; ky < patchSize; ky++ {
					srcRowStart := channelBase + (py*patchSize+ky)*imageSize + px*patchSize
					copy(patches[dstCBase+ky*patchSize:dstCBase+ky*patchSize+patchSize],
						input[srcRowStart:srcRowStart+patchSize])
				}
			}
		}
	}
	out := make([]float32, numPatches*dim)
	gemmNT(numPatches, dim, patchSpan, 1, patches, weight, 0, out)
	addRowBias(out, bias, numPatches, dim)
	return out
}
