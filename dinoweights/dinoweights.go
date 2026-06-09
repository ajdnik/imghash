// Package dinoweights ships the default model weights used by
// github.com/ajdnik/imghash/v2/dinohash. The blob is the fused 96-bit ONNX
// checkpoint dinov2_vits14_reg_96bit.onnx published by
// https://github.com/proteus-photos/dinohash-perceptual-hash, converted to
// safetensors fp32. fp32 is preserved on purpose so the embedded values are
// bit-identical to the ONNX initializers.
//
// The package is intentionally trivial: it exports the raw bytes only. This
// keeps dinoweights decoupled from dinohash's Go API so a release of one
// module does not force a release of the other. dinohash supplies the
// parser; dinoweights supplies the bytes.
//
// Usage:
//
//	import (
//	    "github.com/ajdnik/imghash/v2/dinohash"
//	    "github.com/ajdnik/imghash/v2/dinoweights"
//	)
//
//	d, err := dinohash.NewDINOHash(dinohash.WithSafetensorsBlob(dinoweights.Blob))
package dinoweights

import _ "embed"

// Blob is the embedded safetensors fp32 weight blob (~85 MB). The bytes are
// in the standard safetensors v0.4 layout: an 8-byte little-endian header
// length followed by a JSON header and concatenated raw tensor data. Pass
// this slice to dinohash.WithSafetensorsBlob; do not mutate it.
//
//go:embed assets/dinov2_vits14_reg_96bit.safetensors
var Blob []byte
