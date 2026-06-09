# Zernike

> Float64 perceptual hash using rotation-invariant Zernike moment magnitudes

## Overview

Zernike is a perceptual hash algorithm that uses Zernike moments to create a rotation-invariant image descriptor. It computes complex Zernike moments up to a specified degree and stores their magnitudes, making the hash invariant to image rotation.

This algorithm is ideal for:

* Matching images regardless of rotation angle
* Finding similar shapes and structures
* Applications requiring rotation invariance
* Pattern recognition and image retrieval

> [!NOTE]
> Zernike returns a **Float64** hash type (not Binary). Use L2 (Euclidean) distance for comparison.

## How It Works

1. **Resize**: Image is resized to the specified dimensions (default: 64×64)
2. **Grayscale**: Converts to grayscale
3. **Zernike Computation**: Computes Zernike moments for all valid orders up to maximum degree
4. **Magnitude Extraction**: Stores |A(n,m)| values (magnitude only, discarding phase)
5. **Normalization**: Normalizes by DC component for intensity invariance

The magnitude-based representation makes the hash rotation-invariant while remaining sensitive to shape and structure.

## Constructor

```go
func NewZernike(opts ...ZernikeOption) (Zernike, error)
```

Creates a new Zernike hasher with optional configuration.

### Available Options

* `WithSize(width, height uint)` - Set resize dimensions (default: 64×64)
* `WithInterpolation(interp Interpolation)` - Set interpolation method (default: Bilinear)
* `WithDegree(degree int)` - Set maximum Zernike degree (default: 8)
* `WithDistance(fn DistanceFunc)` - Override default L2 distance function

## Usage Example

```go
package main

import (
    "fmt"
    "image"
    _ "image/png"
    "os"

    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create Zernike hasher with higher degree for more detail
    hasher, err := imghash.NewZernike(
        imghash.WithDegree(12),
        imghash.WithSize(128, 128),
    )
    if err != nil {
        panic(err)
    }

    // Load images (potentially rotated versions)
    img1, _ := loadImage("original.png")
    img2, _ := loadImage("rotated.png")

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
    
    // Rotated images should have small distance
    if distance < 0.1 {
        fmt.Println("Images match (rotation-invariant)")
    }
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

- **`width`** `uint` (default: `64`) — Image resize width
- **`height`** `uint` (default: `64`) — Image resize height
- **`interpolation`** `Interpolation` (default: `Bilinear`) — Resize interpolation method
- **`degree`** `int` (default: `8`) — Maximum Zernike degree (must be > 0). Higher values capture more detail but increase hash size.
- **`distanceFunc`** `DistanceFunc` (default: `similarity.L2`) — Distance comparison function

## Hash Type

Returns `hashtype.Float64` - a slice of float64 values containing normalized Zernike moment magnitudes.

Hash size depends on the degree parameter:

* Degree 8: 20 values (excluding A(0,0))
* Degree 12: 42 values
* Generally: approximately (degree+1)×(degree+2)/4 values

## Distance Metric

Default comparison uses **L2 (Euclidean) distance**. Lower values indicate more similar images.

You can override with:

* `similarity.L1` - Manhattan distance
* `similarity.Cosine` - Cosine distance
* Custom distance function

## Degree Parameter

The `degree` parameter controls the maximum order n:

* **Lower degrees (4-8)**: Faster computation, captures coarse structure
* **Higher degrees (12-16)**: More detailed representation, larger hash size
* Trade-off between discriminative power and computational cost

## Reference

Based on: A. Khotanzad and Y. H. Hong, "Invariant Image Recognition by Zernike Moments" (1990).
