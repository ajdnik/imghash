# PDQ Hash

> Facebook's Perceptual Image Hash for robust deduplication

## Overview

PDQ (Photo-DNA Quality) is a perceptual hash developed by Facebook (now Meta) that produces a 256-bit hash robust to JPEG compression, rescaling, and minor edits while remaining fast enough for large-scale deduplication.

The algorithm applies a Discrete Cosine Transform (DCT) to a 64×64 grayscale image, extracts a 16×16 coefficient block, and thresholds it against the median to produce a binary hash.

## When to Use

Use PDQ when you need:

* **High robustness** to JPEG compression and minor edits
* **Large-scale deduplication** of images (designed for production use)
* **256-bit hashes** for better discrimination than 64-bit alternatives
* **Industry-standard** hash used by content moderation systems

> [!NOTE]
> PDQ is particularly effective for detecting duplicate photos, memes, and content variations at scale.

## Constructor

```go
func NewPDQ(opts ...PDQOption) (PDQ, error)
```

### Available Options

* `WithInterpolation(interp Interpolation)` - Sets the resize interpolation method
* `WithDistance(fn DistanceFunc)` - Overrides the default Hamming distance function

### Supported Interpolation Methods

* `NearestNeighbor`
* `Bilinear` (default)
* `Bicubic`
* `MitchellNetravali`
* `Lanczos2`
* `Lanczos3`
* `BilinearExact`

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create PDQ hasher with default settings
    pdq, err := imghash.NewPDQ()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(pdq, "photo.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("PDQ hash: %v\n", hash)
}
```

### With Custom Options

```go
// Create PDQ hasher with custom interpolation
pdq, err := imghash.NewPDQ(
    imghash.WithInterpolation(imghash.Bicubic),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(pdq, "photo.jpg")
```

## Default Settings

* **Hash size**: 256 bits (32 bytes)
* **Resize dimensions**: 64×64 pixels (fixed)
* **Interpolation**: Bilinear
* **Distance metric**: Hamming distance
* **Processing**: DCT-based with Jarosz filter

## Technical Details

The PDQ algorithm:

1. Resizes the image to 64×64 pixels
2. Converts to grayscale
3. Applies Jarosz box filter (window=2, reps=2)
4. Computes DCT on the filtered image
5. Extracts top-left 16×16 DCT coefficients
6. Thresholds coefficients against their median
7. Produces a 256-bit binary hash

## Comparison

PDQ hashes are compared using Hamming distance by default. Lower distances indicate more similar images.

```go
dist, err := pdq.Compare(hash1, hash2)
if err != nil {
    panic(err)
}

if dist < 20 {
    fmt.Println("Images are very similar")
}
```

## References

* [PDQ and TMK+PDQF by Facebook](https://github.com/facebook/ThreatExchange/tree/main/pdq)

> [!TIP]
> PDQ is optimized for real-world photo matching and is used by major platforms for content moderation and copyright detection.
