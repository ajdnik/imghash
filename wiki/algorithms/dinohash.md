# DINOHash

> Deep-learning perceptual hash from "DINOHash: Robust Deep Perceptual Hashing using DINOv2 Features" (Khare et al., 2025).

## Overview

DINOHash is a deep-learning perceptual hash backed by a frozen DINOv2 ViT-S/14 (with register tokens) backbone fused with a 96-bit PCA-aligned linear head. Compares using Hamming distance.

Unlike the classical algorithms in this package, DINOHash captures high-level semantic features that survive crops, recolors, lossy re-encoding, and adversarial edits designed to defeat hashes like PDQ or PHash.

> [!NOTE]
> DINOHash is the only algorithm in imghash that ships its model weights in a separate Go module so importers do not pay the ~85 MB embed cost unless they opt in.

## When to Use

Use DINOHash when:

- You need robustness to crops, color shifts, and aggressive re-encoding
- Adversarial robustness matters (e.g. content moderation, near-duplicate detection)
- You can afford around one second of CPU per image
- Pre-filtering with a faster hash (PDQ) for high-throughput pipelines is acceptable

Avoid DINOHash when:

- You need sub-millisecond latency
- Your input set is restricted to bit-identical or trivially-edited duplicates (classical hashes are faster and sufficient)
- You cannot ship the ~85 MB weights blob with your binary

## Weights

This package ships no model weights. The canonical blob lives in a sibling module `github.com/ajdnik/imghash/v2/dinoweights`:

```go
import (
  "github.com/ajdnik/imghash/v2"
  "github.com/ajdnik/imghash/v2/dinoweights"
)

d, err := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
```

The dinoweights module exports a single `Blob []byte` variable embedded via `go:embed`. The blob is a safetensors fp32 file extracted from the fused 96-bit ONNX checkpoint published by [proteus-photos/dinohash-perceptual-hash](https://github.com/proteus-photos/dinohash-perceptual-hash).

## Constructor

```go
func NewDINOHash(opts ...DINOHashOption) (*DINOHash, error)
```

The constructor returns immediately; weights are parsed lazily on the first `Calculate` call. `Calculate` returns `ErrNoWeights` when neither `WithSafetensorsBlob` nor `WithDINOWeights` was supplied.

## Available Options

- `WithSafetensorsBlob(b []byte)` — required unless `WithDINOWeights` is supplied
- `WithDINOWeights(p WeightsProvider)` — advanced: supply a custom `WeightsProvider`
- `WithDistance(fn DistanceFunc)` — override the default Hamming metric

> [!TIP]
> Input size, ImageNet normalization, hash length, and bicubic positional-embedding interpolation are intentionally not exposed as options. Exposing them would let callers silently produce hashes that diverge from the published reference.

## Usage Example

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash/v2"
  "github.com/ajdnik/imghash/v2/dinoweights"
)

func main() {
  d, err := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
  if err != nil {
    panic(err)
  }

  h1, err := imghash.HashFile(d, "image1.jpg")
  if err != nil {
    panic(err)
  }

  h2, err := imghash.HashFile(d, "image2.jpg")
  if err != nil {
    panic(err)
  }

  dist, err := d.Compare(h1, h2)
  if err != nil {
    panic(err)
  }

  fmt.Printf("Distance: %v / 96 bits\n", dist)
}
```

## Default Settings

- **Hash type:** `Binary` (12 bytes, 96 bits)
- **Default metric:** Hamming distance
- **Backbone:** DINOv2 ViT-S/14 with register tokens (~21 M parameters)
- **Input:** 224x224 RGB, ImageNet-normalized
- **Resize:** `BilinearExact`
- **Positional embedding:** bicubic interpolation from native 37x37 grid

## How It Works

1. Resize the input image to 224x224 with bilinear interpolation
2. Normalize each pixel using ImageNet mean and standard deviation
3. Extract 256 patch embeddings via a 14x14 stride-14 conv2d
4. Prepend the CLS token and 4 register tokens; add bicubically-interpolated positional embeddings to the patch slice
5. Run 12 transformer blocks (pre-LN attention + GELU MLP with LayerScale)
6. Apply final LayerNorm and take the CLS row
7. Project to 96 dimensions via the fused PCA-aligned linear head
8. Pack each sign bit (MSB-first) into a 12-byte `Binary` hash

## Custom Weights Provider

For checkpoint variants, networked weights, or alternate storage formats, implement `WeightsProvider` and pass via `WithDINOWeights`:

```go
type WeightsProvider interface {
    Tensors() (map[string]Tensor, error)
}
```

The interface returns the canonical tensor name → shape/data map. Required tensor names match the safetensors layout shipped in `dinoweights`. The exported `imghash.ParseSafetensors` helper can decode any compatible blob.

## Performance Characteristics

> [!TIP]
> Single-image latency is around one second on a modern CPU. For high-throughput dedup pipelines pair DINOHash with a faster pre-filter (PDQ) and invoke DINOHash only on candidates within a PDQ Hamming threshold.

Inference runs in pure Go on top of `gonum` BLAS. There is no CGO dependency, no ONNX runtime, and no GPU acceleration. The first `Calculate` call parses the ~85 MB safetensors blob and constructs the ViT model; subsequent calls reuse it.

## DINOHash vs. Other Algorithms

| Algorithm | Hash size | Robust to crops | Robust to recolor | Adversarial-resistant | Throughput |
|-----------|-----------|-----------------|-------------------|-----------------------|------------|
| PHash | 64 bits | low | medium | low | ~ms |
| PDQ | 256 bits | medium | medium | low | ~ms |
| DINOHash | 96 bits | high | high | high | ~1s |

## References

- Paper: [DINOHash: Robust Deep Perceptual Hashing using DINOv2 Features](https://arxiv.org/abs/2503.11195) (Khare et al., 2025)
- Reference implementation: [proteus-photos/dinohash-perceptual-hash](https://github.com/proteus-photos/dinohash-perceptual-hash)
- Backbone: [facebookresearch/dinov2](https://github.com/facebookresearch/dinov2)
