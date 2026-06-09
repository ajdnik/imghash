# GIST

> Float64 holistic scene descriptor using oriented Gabor filters on a spatial grid

## Overview

GIST is a holistic scene descriptor that captures the overall spatial structure and dominant orientations of an image. It applies a bank of oriented Gabor filters at multiple scales and pools the responses over a coarse spatial grid to create a global image representation.

This algorithm is ideal for:

* Scene recognition and classification
* Finding images with similar spatial layouts
* Matching images with similar overall structure
* Content-based image retrieval for scenes

> [!NOTE]
> GIST returns a **Float64** hash type (not Binary). Use Cosine distance for comparison.

## How It Works

1. **Resize**: Image is resized to the specified dimensions (default: 64×64)
2. **Grayscale**: Converts to grayscale
3. **Normalization**: Normalizes pixel values (zero mean, unit variance)
4. **Gabor Filtering**: Applies oriented Gabor filters at 3 scales with multiple orientations:
   * Scale 1 (wavelength 4): 8 orientations
   * Scale 2 (wavelength 8): 8 orientations
   * Scale 3 (wavelength 12): 4 orientations
5. **Spatial Pooling**: Divides image into grid cells and averages filter responses
6. **L2 Normalization**: Normalizes the final descriptor vector

The resulting descriptor captures both the spatial envelope and textural properties of the scene.

## Constructor

```go
func NewGIST(opts ...GISTOption) (GIST, error)
```

Creates a new GIST hasher with optional configuration.

### Available Options

* `WithSize(width, height uint)` - Set resize dimensions (default: 64×64)
* `WithInterpolation(interp Interpolation)` - Set interpolation method (default: Bilinear)
* `WithGridSize(x, y uint)` - Set spatial grid dimensions (default: 4×4)
* `WithDistance(fn DistanceFunc)` - Override default Cosine distance function

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
    // Create GIST hasher with finer spatial grid
    hasher, err := imghash.NewGIST(
        imghash.WithGridSize(8, 8),
        imghash.WithSize(128, 128),
    )
    if err != nil {
        panic(err)
    }

    // Load scene images
    scene1, _ := loadImage("beach.jpg")
    scene2, _ := loadImage("beach_similar.jpg")

    // Calculate hashes
    hash1, err := hasher.Calculate(scene1)
    if err != nil {
        panic(err)
    }

    hash2, err := hasher.Calculate(scene2)
    if err != nil {
        panic(err)
    }

    // Compare using Cosine distance
    distance, err := hasher.Compare(hash1, hash2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Scene distance: %.4f\n", distance)
    
    // Similar scenes should have low cosine distance
    if distance < 0.2 {
        fmt.Println("Scenes are similar")
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
- **`gridX`** `uint` (default: `4`) — Number of grid cells in X direction (must be > 0)
- **`gridY`** `uint` (default: `4`) — Number of grid cells in Y direction (must be > 0)
- **`distanceFunc`** `DistanceFunc` (default: `similarity.Cosine`) — Distance comparison function

## Hash Type

Returns `hashtype.Float64` - a normalized descriptor vector.

Hash size = `gridX × gridY × totalOrientations`

* With default 4×4 grid: 320 values (4 × 4 × 20 orientations)
* With 8×8 grid: 1280 values

## Distance Metric

Default comparison uses **Cosine distance**. Lower values indicate more similar scene structure.

You can override with:

* `similarity.L2` - Euclidean distance
* `similarity.L1` - Manhattan distance
* Custom distance function

## Grid Size Parameter

The grid size controls spatial granularity:

* **Smaller grids (4×4)**: Captures global scene structure, smaller hash
* **Larger grids (8×8, 16×16)**: More spatial detail, larger hash, more discriminative
* Trade-off between descriptor size and spatial precision

## Gabor Filter Bank

Fixed configuration (3 scales, 20 total orientations):

* **Wavelength 4**: 8 orientations (fine details)
* **Wavelength 8**: 8 orientations (medium scale)
* **Wavelength 12**: 4 orientations (coarse structure)

## Reference

Based on: A. Oliva and A. Torralba, "Modeling the Shape of the Scene: A Holistic Representation of the Spatial Envelope" (2001).
