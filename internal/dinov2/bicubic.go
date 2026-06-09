package dinov2

import (
	"math"

	"github.com/ajdnik/imghash/v2/internal/vit"
)

// Positional-embedding interpolation constants. DINOv2 ViT-S/14+reg was
// pretrained at 518x518 (37x37 patch grid, 1370 positions including CLS). At
// inference time the model interpolates this grid down to the actual patch
// count (16x16 for our locked 224x224 input). We reproduce the upstream
// behavior (facebookresearch/dinov2 vision_transformer.py
// interpolate_pos_encoding) once at weight-load time so Forward never pays
// the cost.
const (
	PosEmbedNativeLen = 1370 // 1 + 37*37
	posEmbedNativeM   = 37
	posEmbedTargetM   = vit.ImageSize / vit.PatchSize // 16
	cubicA            = -0.75
)

// InterpolatePosEmbed bicubically resamples the patch portion of the native
// positional embedding from a 37x37 grid to a 16x16 grid and prepends the
// CLS positional embedding unchanged. Output shape: [vit.PosEmbedLen,
// vit.Dim] = [257, 384].
func InterpolatePosEmbed(raw []float32) []float32 {
	dim := vit.Dim
	out := make([]float32, vit.PosEmbedLen*dim)
	copy(out[:dim], raw[:dim])

	patchRaw := raw[dim:]
	patchOut := out[dim:]
	// ONNX Resize with coordinate_transformation_mode="pytorch_half_pixel"
	// computes source positions as src = (dst + 0.5) * in_size / out_size - 0.5.
	// proteus-photos' DINOv2 export bakes a Resize node with explicit
	// `sizes=[1, dim, 16, 16]` so PyTorch's interpolate_offset (the 0.1 kludge
	// used when the scale_factor path is taken) has no effect here.
	invScale := float64(posEmbedNativeM) / float64(posEmbedTargetM)

	type axisSample struct {
		idx [4]int
		w   [4]float64
	}
	samples := make([]axisSample, posEmbedTargetM)
	for t := 0; t < posEmbedTargetM; t++ {
		src := (float64(t)+0.5)*invScale - 0.5
		base := int(math.Floor(src)) - 1
		frac := src - float64(base+1)
		var s axisSample
		s.w[0] = cubicKernel(1 + frac)
		s.w[1] = cubicKernel(frac)
		s.w[2] = cubicKernel(1 - frac)
		s.w[3] = cubicKernel(2 - frac)
		for k := 0; k < 4; k++ {
			s.idx[k] = clampInt(base+k, 0, posEmbedNativeM-1)
		}
		samples[t] = s
	}

	for ty := 0; ty < posEmbedTargetM; ty++ {
		ys := samples[ty]
		for tx := 0; tx < posEmbedTargetM; tx++ {
			xs := samples[tx]
			dstRow := patchOut[(ty*posEmbedTargetM+tx)*dim : (ty*posEmbedTargetM+tx+1)*dim]
			for d := 0; d < dim; d++ {
				var v float64
				for ky := 0; ky < 4; ky++ {
					srcRow := patchRaw[(ys.idx[ky]*posEmbedNativeM)*dim : (ys.idx[ky]*posEmbedNativeM+posEmbedNativeM)*dim]
					wy := ys.w[ky]
					var rowAcc float64
					for kx := 0; kx < 4; kx++ {
						rowAcc += xs.w[kx] * float64(srcRow[xs.idx[kx]*dim+d])
					}
					v += wy * rowAcc
				}
				dstRow[d] = float32(v)
			}
		}
	}
	return out
}

// cubicKernel is the a=-0.75 cubic convolution kernel used by PyTorch's
// F.interpolate(mode="bicubic").
func cubicKernel(t float64) float64 {
	t = math.Abs(t)
	switch {
	case t < 1:
		return (cubicA+2)*t*t*t - (cubicA+3)*t*t + 1
	case t < 2:
		return cubicA*t*t*t - 5*cubicA*t*t + 8*cubicA*t - 4*cubicA
	}
	return 0
}

func clampInt(v, lo, hi int) int {
	switch {
	case v < lo:
		return lo
	case v > hi:
		return hi
	}
	return v
}
