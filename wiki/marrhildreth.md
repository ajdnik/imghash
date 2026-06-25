# Marr-Hildreth Hash

> Edge-based perceptual hash using Laplacian of Gaussian

## Overview

Marr-Hildreth Hash is a perceptual hashing algorithm based on the method described in "Implementation and Benchmarking of Perceptual Image Hash Functions" by Zauner et al. It uses the Marr-Hildreth edge detector (Laplacian of Gaussian) to extract structural features.

The algorithm applies Gaussian blur, histogram equalization, and a Marr-Hildreth kernel to detect edges, then computes block sums and thresholds them to create the hash.

## When to Use

Use Marr-Hildreth Hash when you need:

* **Edge-based features** for robust matching
* **High robustness** to various transformations
* **Academic rigor** with benchmarking results
* **Larger hash sizes** for better discrimination
* **Scale-space analysis** through kernel parameters

> [!NOTE]
> Marr-Hildreth Hash produces a 72-byte (576-bit) hash by default, making it suitable for applications requiring high precision.

## Constructor

```go
func NewMarrHildreth(opts ...MarrHildrethOption) (MarrHildreth, error)
```

### Available Options

* `WithSize(width, height uint)` - Sets the resize dimensions
* `WithInterpolation(interp Interpolation)` - Sets the resize interpolation method
* `WithScale(scale float64)` - Sets the scale parameter for Marr-Hildreth kernel
* `WithAlpha(alpha float64)` - Sets the alpha parameter for Marr-Hildreth kernel
* `WithKernelSize(size int)` - Sets the Gaussian kernel size for pre-blur
* `WithSigma(sigma float64)` - Sets the Gaussian sigma for pre-blur (0 = auto)
* `WithDistance(fn DistanceFunc)` - Overrides the default Hamming distance function

### Supported Interpolation Methods

* `NearestNeighbor`
* `Bilinear`
* `Bicubic` (default)
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
    // Create Marr-Hildreth hasher with default settings
    mh, err := imghash.NewMarrHildreth()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(mh, "photo.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Marr-Hildreth hash: %v\n", hash)
}
```

### With Custom Options

```go
// Create Marr-Hildreth hasher with custom parameters
mh, err := imghash.NewMarrHildreth(
    imghash.WithScale(1.5),
    imghash.WithAlpha(2.5),
    imghash.WithSize(512, 512),
    imghash.WithInterpolation(imghash.Bicubic),
    imghash.WithKernelSize(9),
    imghash.WithSigma(1.5),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(mh, "photo.jpg")
```

## Default Settings

* **Hash size**: 576 bits (72 bytes)
* **Resize dimensions**: 512×512 pixels
* **Interpolation**: Bicubic
* **Scale**: 1.0
* **Alpha**: 2.0
* **Gaussian kernel size**: 7
* **Gaussian sigma**: 0 (auto-computed)
* **Block size**: 16×16 pixels
* **Number of blocks**: 31×31
* **Sub-block size**: 3×3
* **Sub-block stride**: 4
* **Distance metric**: Hamming distance

## How It Works

The Marr-Hildreth Hash algorithm:

1. Converts image to grayscale
2. Applies Gaussian blur with specified kernel size and sigma
3. Resizes to 512×512 pixels
4. Applies histogram equalization for contrast enhancement
5. Applies 2D Marr-Hildreth (Laplacian of Gaussian) filter
6. Divides filtered image into 16×16 pixel blocks (31×31 blocks total)
7. Computes sum of each block
8. Creates hash by:
   * Stepping through blocks with 3×3 sub-blocks at stride 4
   * Computing average of each 3×3 sub-block
   * Comparing each value to the sub-block average
   * Setting bit to 1 if value > average, 0 otherwise
9. Produces a 576-bit binary hash

> [!TIP]
> The Marr-Hildreth kernel is a Laplacian of Gaussian (LoG) operator that detects edges at multiple scales, making the hash robust to various transformations.

## Scale and Alpha Parameters

The Marr-Hildreth kernel is computed using the formula:

```
LoG(x, y) = (2 - (x² + y²)/σ²) * exp(-(x² + y²)/(2σ²))
```

Where:

* `scale` controls the decomposition level
* `alpha` determines the kernel sigma: `σ = 4 * alpha^level`

```go
// Smaller scale for finer details
mh1, _ := imghash.NewMarrHildreth(
    imghash.WithScale(0.5),
    imghash.WithAlpha(1.5),
)

// Larger scale for coarser features
mh2, _ := imghash.NewMarrHildreth(
    imghash.WithScale(2.0),
    imghash.WithAlpha(3.0),
)
```

## Comparison

Marr-Hildreth hashes are compared using Hamming distance:

```go
mh, _ := imghash.NewMarrHildreth()

h1, _ := imghash.HashFile(mh, "original.jpg")
h2, _ := imghash.HashFile(mh, "modified.jpg")

dist, err := mh.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 50 {
    fmt.Println("Images are similar")
}
```

> [!NOTE]
> Due to the larger hash size (576 bits), threshold values for similarity will be proportionally higher than for 64-bit hashes.

## Performance Characteristics

* **Speed**: Slow (Gaussian blur, histogram equalization, LoG filter, block processing)
* **Memory**: Moderate (72-byte hash, 512×512 intermediate image)
* **Robustness**: Very high for various transformations
* **Rotation**: Not rotation-invariant
* **Scaling**: Very robust to scaling

## Hash Size Calculation

```
gridSteps = ((31 - 3) / 4) + 1 = 8
hashBits = 8 × 8 × 3 × 3 = 576 bits = 72 bytes
```

## References

* [Implementation and Benchmarking of Perceptual Image Hash Functions - Zauner et al.](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking)

> [!WARNING]
> Marr-Hildreth Hash is computationally expensive and requires significant processing. Use it when robustness is more important than speed.
