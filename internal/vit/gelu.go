package vit

import "math"

// invSqrt2 is 1/√2 used in the exact-erf GELU.
const invSqrt2 = 0.7071067811865475

// gelu applies the exact erf-based GELU activation in place. This matches the
// default behavior of torch.nn.GELU (approximate="none"), which DINOv2 uses.
func gelu(x []float32) {
	for i, v := range x {
		f := float64(v)
		x[i] = float32(0.5 * f * (1 + math.Erf(f*invSqrt2)))
	}
}
