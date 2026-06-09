// Package dinov2 implements the non-public machinery behind the DINOHash
// perceptual hash exposed at the root of the imghash module: image
// preprocessing, safetensors blob parsing, bicubic positional-embedding
// interpolation, and weight-map → ViT model loading.
//
// The package is internal: only imghash and its tests may import it. The
// canonical Tensor type lives here and is aliased back to the root package
// so the public WeightsProvider interface in imghash can reference it
// without creating an import cycle.
package dinov2

// Tensor is a single named weight loaded from a safetensors blob. Shape is
// the logical multi-dimensional shape (PyTorch row-major); Data is the
// contiguous float32 buffer.
type Tensor struct {
	Shape []int
	Data  []float32
}
