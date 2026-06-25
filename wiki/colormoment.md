# ColorMoment

> Float64 perceptual hash using color invariant moments from YCrCb and HSV color spaces

## Overview

ColorMoment is a perceptual hash algorithm that uses color invariant moments to create a robust image fingerprint. It operates by computing Hu moments from both YCrCb and HSV color spaces, combining them into a single Float64 hash descriptor.

This algorithm is particularly effective for:

* Color image matching where color information is important
* Finding images with similar color distributions
* Detecting images that have been color-adjusted or filtered
* Applications requiring rotation-invariant features

> [!NOTE]
> ColorMoment returns a **Float64** hash type (not Binary). Use L2 (Euclidean) distance for comparison.

## How It Works

1. **Resize**: Image is resized to the specified dimensions (default: 512×512)
2. **Gaussian Blur**: Applies Gaussian smoothing to reduce noise
3. **Color Space Conversion**: Converts to both YCrCb and HSV color spaces
4. **Moment Extraction**: Computes image moments from each color space
5. **Hu Moments**: Calculates Hu invariant moments for rotation invariance
6. **Descriptor**: Combines Hu moments from both color spaces into a single vector

Based on the method described in "Perceptual Hashing for Color Images Using Invariant Moments" by Tang et al.

## Constructor

```go
func NewColorMoment(opts ...ColorMomentOption) (ColorMoment, error)
```

Creates a new ColorMoment hasher with optional configuration.

### Available Options

* `WithSize(width, height uint)` - Set resize dimensions (default: 512×512)
* `WithInterpolation(interp Interpolation)` - Set interpolation method (default: Bicubic)
* `WithKernelSize(size int)` - Set Gaussian kernel size (default: 3)
* `WithSigma(sigma float64)` - Set Gaussian kernel sigma (default: 0)
* `WithDistance(fn DistanceFunc)` - Override default L2 distance function

## Usage Example

```go
package main

import (
    "fmt"
    "image"
    _ "image/jpeg"
    "os"

    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create ColorMoment hasher with custom settings
    hasher, err := imghash.NewColorMoment(
        imghash.WithSize(512, 512),
        imghash.WithKernelSize(5),
        imghash.WithSigma(1.0),
    )
    if err != nil {
        panic(err)
    }

    // Load images
    img1, _ := loadImage("photo1.jpg")
    img2, _ := loadImage("photo2.jpg")

    // Calculate hashes
    hash1, err := hasher.Calculate(img1)
    if err != nil {
        panic(err)
    }

    hash2, err := hasher.Calculate(img2)
    if err != nil {
        panic(err)
    }

    // Compare using L2 distance
    distance, err := hasher.Compare(hash1, hash2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Distance: %.4f\n", distance)
}

func loadImage(path string) (image.Image, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    img, _, err := image.Decode(f)
    return img, err
}
```

## Default Settings

- **`width`** `uint` (default: `512`) — Image resize width
- **`height`** `uint` (default: `512`) — Image resize height
- **`interpolation`** `Interpolation` (default: `Bicubic`) — Resize interpolation method
- **`kernel`** `int` (default: `3`) — Gaussian kernel size (must be > 0)
- **`sigma`** `float64` (default: `0`) — Gaussian kernel standard deviation (0 = auto-calculate)
- **`distanceFunc`** `DistanceFunc` (default: `similarity.L2`) — Distance comparison function

## Hash Type

Returns `hashtype.Float64` - a slice of float64 values containing the combined Hu moments from HSV and YCrCb color spaces (typically 14 values: 7 from each color space).

## Distance Metric

Default comparison uses **L2 (Euclidean) distance**. Lower values indicate more similar images.

You can override with:

* `similarity.L1` - Manhattan distance
* `similarity.Cosine` - Cosine distance
* `similarity.ChiSquare` - Chi-square distance
* Custom distance function

## Reference

Based on: [Perceptual Hashing for Color Images Using Invariant Moments](https://www.researchgate.net/publication/286870507_Perceptual_hashing_for_color_images_using_invariant_moments) - Tang et al.
