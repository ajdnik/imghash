package imghash

import (
	"image"
	"sync"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/dinov2"
	"github.com/ajdnik/imghash/v2/internal/vit"
	"github.com/ajdnik/imghash/v2/similarity"
)

// Tensor is a single named weight loaded from a safetensors blob. Shape is
// the logical multi-dimensional shape (PyTorch row-major); Data is the
// contiguous float32 buffer. Callers must not mutate either field after
// returning from a WeightsProvider — DINOHash retains references to the
// slice.
type Tensor = dinov2.Tensor

// WeightsProvider supplies the named float32 tensors that constitute a
// DINOHash model. The canonical source is the byte slice in the sibling
// module github.com/ajdnik/imghash/dinoweights, supplied via
// WithSafetensorsBlob. Implement this interface to load weights from disk,
// network, or a custom format.
type WeightsProvider interface {
	// Tensors returns the model's weight map keyed by canonical name.
	// Implementations may lazily decode on first call. The returned map
	// should not be mutated by the caller.
	Tensors() (map[string]Tensor, error)
}

// ParseSafetensors decodes a safetensors fp32 blob into a name->Tensor map.
// Only the "F32" dtype is accepted; fp32 is preserved on purpose so the
// returned values are bit-identical to the source initializers.
//
// Exported so external WeightsProvider implementations can share the same
// decoder; most callers should use WithSafetensorsBlob instead.
func ParseSafetensors(blob []byte) (map[string]Tensor, error) {
	return dinov2.ParseSafetensors(blob)
}

// safetensorsProvider wraps a raw safetensors blob as a WeightsProvider.
type safetensorsProvider []byte

// Tensors implements WeightsProvider.
func (s safetensorsProvider) Tensors() (map[string]Tensor, error) {
	return ParseSafetensors(s)
}

// DINOHash is a perceptual hash backed by a frozen DINOv2 ViT-S/14 (with
// register tokens) backbone fused with a 96-bit PCA-aligned head. It follows
// the method in "DINOHash: Robust Deep Perceptual Hashing using DINOv2
// Features" (Khare et al., 2025, https://arxiv.org/abs/2503.11195) and reuses
// the fused 96-bit ONNX checkpoint published by
// https://github.com/proteus-photos/dinohash-perceptual-hash.
//
// Unlike the other perceptual hashes in this package, DINOHash runs a
// hand-rolled ViT forward pass in pure Go (no CGO, no ONNX runtime). Inference
// is therefore substantially slower (around one second per image on a modern
// laptop) and the API is intentionally minimal — image preprocessing and hash
// length are anchored to the published reference so hashes remain comparable
// across implementations.
//
// # Weights
//
// This package ships no model weights. Supply them via WithSafetensorsBlob
// (most convenient) or WithDINOWeights (for custom providers). The default
// blob lives in the sibling module github.com/ajdnik/imghash/dinoweights:
//
//	import (
//	    "github.com/ajdnik/imghash/v2"
//	    "github.com/ajdnik/imghash/dinoweights"
//	)
//
//	d, err := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
//
// The zero value of DINOHash is not usable; create one with NewDINOHash and
// pass it by pointer. The model holds ~30 MB of float32 weights once loaded.
type DINOHash struct {
	weights  WeightsProvider
	distFunc DistanceFunc

	once    sync.Once
	model   *vit.Model
	initErr error
}

// NewDINOHash returns a DINOHash configured with the given options. Model
// weights are not parsed until the first call to Calculate; constructor cost
// is therefore negligible. Calculate returns ErrNoWeights if no weights
// source was supplied.
func NewDINOHash(opts ...DINOHashOption) (*DINOHash, error) {
	d := &DINOHash{}
	for _, o := range opts {
		o.applyDINOHash(d)
	}
	return d, nil
}

// Calculate returns the 96-bit DINOHash of img as a 12-byte hashtype.Binary.
// The first call asks the configured weights source for the model tensors
// and constructs the ViT; subsequent calls reuse the same model.
func (d *DINOHash) Calculate(img image.Image) (hashtype.Hash, error) {
	d.once.Do(func() {
		if d.weights == nil {
			d.initErr = ErrNoWeights
			return
		}
		tensors, err := d.weights.Tensors()
		if err != nil {
			d.initErr = err
			return
		}
		m, err := dinov2.LoadModel(tensors)
		if err != nil {
			d.initErr = err
			return
		}
		d.model = m
	})
	if d.initErr != nil {
		return nil, d.initErr
	}
	input := dinov2.ImageToTensor(img)
	pre := d.model.Forward(input)
	// Bit packing matches numpy.packbits with bitorder="big" (MSB-first),
	// which the Python reference uses to serialize the 96 sign bits.
	hash := make(hashtype.Binary, vit.HashBits/8)
	for i, v := range pre {
		if v >= 0 {
			_ = hash.SetReverse(uint(i))
		}
	}
	return hash, nil
}

// Compare returns the distance between two DINOHashes. Without WithDistance,
// it uses similarity.Hamming over the 12-byte binary representation. Compare
// does not require a weights source.
func (d *DINOHash) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateBinaryCompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if d.distFunc != nil {
		return d.distFunc(h1, h2)
	}
	return similarity.Hamming(h1, h2)
}
