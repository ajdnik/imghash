# WHash (Wavelet Hash)

> Haar wavelet transform-based perceptual hash

## Overview

WHash is a perceptual hash based on the Haar wavelet transform. It applies a multi-level 2-D discrete wavelet transform (DWT) to the image and thresholds the low-frequency (LL) coefficients against their median.

The wavelet transform decomposes the image into different frequency bands, with the LL subband containing the low-frequency components that capture the image's overall structure.

## When to Use

Use WHash when you need:

* **Wavelet-based features** instead of DCT
* **Multi-resolution analysis** through decomposition levels
* **Compact representation** of low-frequency content
* **Alternative to DCT-based methods** like PHash
* **Configurable hash size** through decomposition levels

> [!NOTE]
> WHash is particularly effective for images where wavelet decomposition better captures salient features than frequency-domain transforms.

## Constructor

```go
func NewWHash(opts ...WHashOption) (WHash, error)
```

### Available Options

* `WithSize(width, height uint)` - Sets the final LL subband dimensions (hash size = width × height bits)
* `WithLevel(level int)` - Sets the number of Haar DWT decomposition levels
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
    // Create WHash with default settings
    whash, err := imghash.NewWHash()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(whash, "photo.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("WHash: %v\n", hash)
}
```

### With Custom Options

```go
// Create WHash with custom decomposition level and size
whash, err := imghash.NewWHash(
    imghash.WithSize(16, 16),     // 256-bit hash
    imghash.WithLevel(4),          // 4 decomposition levels
    imghash.WithInterpolation(imghash.Bicubic),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(whash, "photo.jpg")
```

## Default Settings

* **Hash size**: 64 bits (8 bytes)
* **LL subband size**: 8×8 pixels
* **Decomposition levels**: 3
* **Resize dimensions**: Computed as width × 2^level × height × 2^level
* **Interpolation**: Bilinear
* **Distance metric**: Hamming distance

## How It Works

The WHash algorithm:

1. Computes resize dimensions: `(width × 2^level, height × 2^level)`
   * For default 8×8 with level=3: resizes to 64×64
2. Resizes the image to computed dimensions
3. Converts to grayscale
4. Applies Haar DWT for the specified number of levels
5. Extracts the LL (low-low) subband of size width×height
6. Computes the median of LL coefficients
7. Thresholds LL against the median:
   * Sets bit to 1 if coefficient > median
   * Sets bit to 0 if coefficient ≤ median
8. Produces a binary hash (width × height bits)

> [!TIP]
> Each decomposition level halves the dimensions, so the final LL subband after `level` passes is exactly width×height pixels.

## Understanding Decomposition Levels

The level parameter controls how many times the Haar DWT is applied:

```go
// Level 1: One decomposition pass
whash1, _ := imghash.NewWHash(
    imghash.WithSize(8, 8),
    imghash.WithLevel(1),
)
// Resizes to 16×16, applies 1 DWT, extracts 8×8 LL subband

// Level 3 (default): Three decomposition passes
whash3, _ := imghash.NewWHash(
    imghash.WithSize(8, 8),
    imghash.WithLevel(3),
)
// Resizes to 64×64, applies 3 DWT passes, extracts 8×8 LL subband

// Level 5: Five decomposition passes
whash5, _ := imghash.NewWHash(
    imghash.WithSize(8, 8),
    imghash.WithLevel(5),
)
// Resizes to 256×256, applies 5 DWT passes, extracts 8×8 LL subband
```

Higher levels capture coarser features but require more computation.

## Hash Size vs. Decomposition

| Width × Height | Levels | Resize To | Hash Bits |
| -------------- | ------ | --------- | --------- |
| 8×8            | 3      | 64×64     | 64        |
| 16×16          | 3      | 128×128   | 256       |
| 8×8            | 4      | 128×128   | 64        |
| 16×16          | 4      | 256×256   | 256       |

## Comparison

WHash uses Hamming distance by default:

```go
whash, _ := imghash.NewWHash()

h1, _ := imghash.HashFile(whash, "original.jpg")
h2, _ := imghash.HashFile(whash, "modified.jpg")

dist, err := whash.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 10 {
    fmt.Println("Images are similar")
}
```

## Wavelet vs. DCT

| Aspect           | WHash (Wavelet)       | PHash (DCT)       |
| ---------------- | --------------------- | ----------------- |
| **Transform**    | Haar DWT              | DCT               |
| **Localization** | Time-frequency        | Frequency         |
| **Levels**       | Configurable          | Fixed             |
| **Best for**     | Multi-scale features  | Frequency content |
| **Speed**        | Fast (Haar is simple) | Moderate          |

## Performance Characteristics

* **Speed**: Fast (Haar DWT is computationally efficient)
* **Memory**: Low to moderate (depends on hash size)
* **Robustness**: Good for minor modifications
* **Rotation**: Not rotation-invariant
* **Scaling**: Robust to scaling

## References

* [Wavelet Image Hash in Python](https://fullstackml.com/wavelet-image-hash-in-python-3504571f3b08)

> [!WARNING]
> The level parameter must be positive. Higher levels require larger input images and more computation but may capture coarser features.
