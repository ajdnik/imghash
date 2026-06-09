# PHash (Perceptual Hash)

> DCT-based perceptual hash with weighted Hamming distance

## Overview

PHash (Perceptual Hash) is a DCT-based perceptual hashing algorithm described in "Implementation and Benchmarking of Perceptual Image Hash Functions" by Zauner et al. It produces a 64-bit hash by analyzing low-frequency DCT coefficients.

This algorithm is more robust than simple hashing methods and supports weighted Hamming distance for improved accuracy.

## When to Use

Use PHash when you need:

* **DCT-based robustness** to minor transformations
* **64-bit compact hash** for efficient storage
* **Weighted distance metrics** for better discrimination
* **Research-backed algorithm** with benchmarking results
* **Balance** between robustness and hash size

> [!NOTE]
> PHash is widely used in academic research and commercial applications for perceptual image matching.

## Constructor

```go
func NewPHash(opts ...PHashOption) (PHash, error)
```

### Available Options

* `WithSize(width, height uint)` - Sets the resize dimensions (DCT input size)
* `WithInterpolation(interp Interpolation)` - Sets the resize interpolation method
* `WithWeights(weights []float64)` - Sets per-byte weights for weighted Hamming distance
* `WithDistance(fn DistanceFunc)` - Overrides the default weighted Hamming distance function

### Supported Interpolation Methods

* `NearestNeighbor`
* `Bilinear`
* `BilinearExact` (default)
* `Bicubic`
* `MitchellNetravali`
* `Lanczos2`
* `Lanczos3`

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create PHash with default settings
    phash, err := imghash.NewPHash()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(phash, "photo.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("PHash: %v\n", hash)
}
```

### With Custom Options

```go
// Create PHash with custom DCT size and interpolation
phash, err := imghash.NewPHash(
    imghash.WithSize(64, 64),  // Larger DCT input
    imghash.WithInterpolation(imghash.Lanczos3),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(phash, "photo.jpg")
```

### With Weighted Distance

```go
// Create PHash with custom per-byte weights
weights := []float64{1.0, 1.2, 1.5, 1.8, 2.0, 2.2, 2.5, 3.0}

phash, err := imghash.NewPHash(
    imghash.WithWeights(weights),
)
if err != nil {
    panic(err)
}

// Compare uses weighted Hamming distance
h1, _ := imghash.HashFile(phash, "image1.jpg")
h2, _ := imghash.HashFile(phash, "image2.jpg")

dist, err := phash.Compare(h1, h2)
```

## Default Settings

* **Hash size**: 64 bits (8 bytes)
* **Resize dimensions**: 32×32 pixels
* **Interpolation**: BilinearExact
* **DCT block**: Top-left 8×8 coefficients (excluding DC)
* **Distance metric**: Weighted Hamming distance
* **Default weights**: All 1.0 (uniform weighting)

## How It Works

The PHash algorithm:

1. Resizes the image to 32×32 pixels
2. Converts to grayscale
3. Computes 2D Discrete Cosine Transform (DCT)
4. Extracts top-left 8×8 DCT coefficients
5. Removes the DC component (strongest frequency)
6. Computes the mean of remaining coefficients
7. Thresholds coefficients against the mean:
   * Sets bit to 1 if coefficient > mean
   * Sets bit to 0 if coefficient ≤ mean
8. Produces a 64-bit binary hash

> [!TIP]
> The DCT captures frequency information, making PHash robust to minor image modifications while remaining sensitive to structural changes.

## Weighted Hamming Distance

PHash supports weighted Hamming distance, where each byte can have a different weight:

```go
// Higher weights for more significant bytes
weights := []float64{
    1.0,  // Byte 0 (lowest frequencies)
    1.5,  // Byte 1
    2.0,  // Byte 2
    2.5,  // Byte 3
    3.0,  // Byte 4
    3.5,  // Byte 5
    4.0,  // Byte 6
    4.5,  // Byte 7 (highest frequencies)
}

phash, _ := imghash.NewPHash(imghash.WithWeights(weights))
```

This allows you to emphasize or de-emphasize certain frequency ranges based on your use case.

## Comparison

PHash uses weighted Hamming distance by default:

```go
phash, _ := imghash.NewPHash()

h1, _ := imghash.HashFile(phash, "original.jpg")
h2, _ := imghash.HashFile(phash, "modified.jpg")

dist, err := phash.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 15 {
    fmt.Println("Images are perceptually similar")
}
```

## Performance Characteristics

* **Speed**: Moderate (DCT computation required)
* **Memory**: Low (64-bit hash)
* **Robustness**: High for minor edits, compression
* **Rotation**: Not rotation-invariant
* **Scaling**: Robust to scaling

## PHash vs. Other Algorithms

| Algorithm    | Hash Size | Speed     | Robustness |
| ------------ | --------- | --------- | ---------- |
| **PHash**    | 64 bits   | Medium    | High       |
| Average Hash | 64 bits   | Very Fast | Low        |
| PDQ          | 256 bits  | Medium    | Very High  |
| Block Mean   | Variable  | Slow      | Very High  |

## References

* [Implementation and Benchmarking of Perceptual Image Hash Functions - Zauner et al.](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking)

> [!WARNING]
> PHash is not rotation-invariant. For rotation-robust hashing, consider RASH or Block Mean with rotation variants.
