# RASH (Rotation Aware Spatial Hash)

> Rotation-invariant perceptual hash using ring sampling

## Overview

RASH (Rotation Aware Spatial Hash) is a perceptual hash designed to be robust against image rotation. It combines concentric ring sampling for rotation invariance, a 1-D DCT for frequency compaction, and median thresholding for binarization.

Because ring-mean features are inherently rotation-invariant (rotating the image only permutes pixels within a ring, leaving its mean unchanged), the resulting hash stays stable under arbitrary rotations.

## When to Use

Use RASH when you need:

* **Rotation-invariant hashing** without multiple hash variants
* **Compact hash** (64-bit default) with rotation robustness
* **Ring-based spatial features** for circular symmetry
* **Efficient rotation handling** compared to Block Mean rotation variants
* **DCT-based frequency analysis** of radial features

> [!NOTE]
> RASH provides rotation invariance in a compact 64-bit hash, unlike Block Mean rotation variants that require 6,144+ bits.

## Constructor

```go
func NewRASH(opts ...RASHOption) (RASH, error)
```

### Available Options

* `WithSize(width, height uint)` - Sets the resize dimensions
* `WithRings(rings int)` - Sets the number of concentric rings
* `WithSigma(sigma float64)` - Sets the Gaussian blur sigma
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
    // Create RASH hasher with default settings
    rash, err := imghash.NewRASH()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(rash, "photo.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("RASH hash: %v\n", hash)
}
```

### With Custom Options

```go
// Create RASH hasher with custom parameters
rash, err := imghash.NewRASH(
    imghash.WithRings(240),           // More rings for finer detail
    imghash.WithSigma(1.5),           // Stronger blur
    imghash.WithSize(512, 512),       // Larger input image
    imghash.WithInterpolation(imghash.Bicubic),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(rash, "photo.jpg")
```

## Default Settings

* **Hash size**: 64 bits (8 bytes)
* **Resize dimensions**: 256×256 pixels
* **Number of rings**: 180
* **Gaussian sigma**: 1.0
* **Interpolation**: Bilinear
* **Distance metric**: Hamming distance

## How It Works

The RASH algorithm:

1. Resizes the image to specified dimensions (default 256×256)
2. Converts to grayscale
3. Applies Gaussian blur with specified sigma
4. Divides the image into concentric rings around the center
5. Computes the mean pixel intensity for each ring:
   * Ring width = max\_radius / number\_of\_rings
   * Each pixel assigned to a ring based on distance from center
6. Applies 1-D DCT to the ring mean vector
7. Keeps the first 64 low-frequency coefficients (skipping DC component)
8. Computes the median of these coefficients
9. Thresholds coefficients against the median:
   * Sets bit to 1 if coefficient > median
   * Sets bit to 0 if coefficient ≤ median
10. Produces a 64-bit binary hash

> [!TIP]
> The key insight: rotating the image only shuffles pixels within each ring, leaving the ring means unchanged. This makes RASH inherently rotation-invariant.

## Ring Sampling Illustration

```
        Ring 0 (center)
       /
      /  Ring 1
     /  /
    /  /  Ring 2
   /  /  /
  ┌──────────┐
  │    ∘     │  ← Image center
  │   ∘∘∘    │
  │  ∘∘∘∘∘   │  ← Concentric rings
  │ ∘∘∘∘∘∘∘  │
  │  ∘∘∘∘∘   │
  │   ∘∘∘    │
  │    ∘     │
  └──────────┘
```

Each ring contains pixels at a similar distance from the center.

## Number of Rings

The number of rings affects hash quality and rotation invariance:

```go
// Fewer rings: Coarser features, faster
rash_coarse, _ := imghash.NewRASH(
    imghash.WithRings(64),  // Only 64 rings
)

// More rings: Finer features, more detail
rash_fine, _ := imghash.NewRASH(
    imghash.WithRings(360),  // 360 rings (1° resolution)
)
```

> [!NOTE]
> With 180 rings (default), the DCT extracts 64 coefficients from 180 ring means, providing good frequency compaction.

## Gaussian Blur

The sigma parameter controls pre-blur strength:

```go
// No blur
rash_sharp, _ := imghash.NewRASH(
    imghash.WithSigma(0),
)

// Strong blur for noise reduction
rash_smooth, _ := imghash.NewRASH(
    imghash.WithSigma(2.5),
)
```

Blur reduces noise and high-frequency details, improving robustness.

## Comparison

RASH hashes are compared using Hamming distance:

```go
rash, _ := imghash.NewRASH()

h1, _ := imghash.HashFile(rash, "original.jpg")
h2, _ := imghash.HashFile(rash, "rotated.jpg")  // Same image, rotated

dist, err := rash.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 10 {
    fmt.Println("Images are the same (rotation-invariant match)")
}
```

## RASH vs. Other Rotation-Robust Methods

| Algorithm             | Hash Size  | Speed     | Rotation Robustness |
| --------------------- | ---------- | --------- | ------------------- |
| **RASH**              | 64 bits    | Fast      | Excellent           |
| Block Mean (Rotation) | 6,144 bits | Very Slow | Excellent           |
| PHash                 | 64 bits    | Medium    | None                |
| PDQ                   | 256 bits   | Medium    | None                |

## Performance Characteristics

* **Speed**: Fast (ring sampling, 1-D DCT)
* **Memory**: Very low (64-bit hash)
* **Rotation**: Invariant to arbitrary rotations
* **Scaling**: Robust to scaling
* **Cropping**: Not robust to cropping (changes ring structure)

## Use Cases

* **Logo detection** with unknown rotation
* **Duplicate detection** for user-uploaded photos (arbitrary orientations)
* **Satellite/aerial imagery** matching
* **Artwork comparison** with varying orientations
* **Medical imaging** (scans at different angles)

## Limitations

> [!WARNING]
> RASH is not robust to cropping or significant perspective changes, as these alter the concentric ring structure. For crop-robust hashing, consider other methods.

## Technical Inspiration

RASH is inspired by ring-partition hashing literature:

* Tang et al., "Robust Image Hashing with Ring Partition and Invariant Vector Distance" (2016)
* De Roover et al., "Robust image hashing based on radial variance of pixels" (2005)

> [!TIP]
> For maximum rotation robustness with minimal hash size, RASH is the best choice. For rotation plus other transformations, consider combining multiple hash types.
