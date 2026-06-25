# EHD (Edge Histogram Descriptor)

> MPEG-7 Edge Histogram Descriptor perceptual hash

## Overview

EHD (Edge Histogram Descriptor) is an MPEG-7 standard perceptual hash algorithm that captures the spatial distribution of edges in an image. It describes edge orientations across different regions of the image, making it effective for texture and shape-based similarity.

## How It Works

The EHD algorithm processes images through the following steps:

1. **Preprocessing**: The image is resized and converted to grayscale
2. **Grid division**: The image is divided into a 4×4 grid of cells
3. **Local edge detection**: For each cell, 2×2 pixel blocks are analyzed for edge patterns
4. **Edge classification**: Each edge is classified into one of 5 types:
   * Vertical
   * Horizontal
   * 45° diagonal
   * 135° diagonal
   * Non-directional
5. **Histogram computation**: A 5-bin histogram is built for each grid cell
6. **Quantization**: Each histogram bin is quantized to 3 bits using MPEG-7 quantization tables

The result is an 80-element descriptor (4×4 grid × 5 edge types).

> [!NOTE]
> EHD returns a **UInt8** hash type (80 bytes), not Binary or Float64.

## When to Use EHD

EHD is particularly effective for:

* Texture-based image retrieval
* Shape and structure matching
* Images with strong edge patterns
* Object recognition based on contours
* Document or logo matching

It is less suitable for:

* Images with very subtle textures
* Smooth gradients without clear edges
* Color-based matching (uses grayscale only)

## Constructor

```go
func NewEHD(opts ...EHDOption) (EHD, error)
```

### Options

- **`WithSize`** `func(width, height uint) SizeOption`
  Sets the resize dimensions before hash computation. The image is resized to this size before processing.
  **Default**: 256×256 pixels

- **`WithInterpolation`** `func(interp Interpolation) InterpolationOption`
  Sets the interpolation method used when resizing the image.
  **Default**: `Bilinear`
  Available options: `NearestNeighbor`, `Bilinear`, `Bicubic`, `Lanczos`

- **`WithDistance`** `func(fn DistanceFunc) DistanceOption`
  Overrides the default distance function used by the `Compare` method.
  **Default**: `similarity.L1` (Manhattan distance)
  Available functions: `Hamming`, `L1`, `L2`, `Cosine`, `ChiSquare`, `PCC`, `Jaccard`

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
	// Create EHD hasher with default settings
	hasher, err := imghash.NewEHD()
	if err != nil {
		panic(err)
	}

	// Or with custom options
	customHasher, err := imghash.NewEHD(
		imghash.WithSize(512, 512),
		imghash.WithInterpolation(imghash.Bicubic),
	)
	if err != nil {
		panic(err)
	}

	// Load images
	img1, _ := loadImage("image1.jpg")
	img2, _ := loadImage("image2.jpg")

	// Calculate hashes
	hash1, err := hasher.Calculate(img1)
	if err != nil {
		panic(err)
	}

	hash2, err := hasher.Calculate(img2)
	if err != nil {
		panic(err)
	}

	// Compare hashes
	distance, err := hasher.Compare(hash1, hash2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("EHD distance: %.2f\n", distance)
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

| Parameter         | Default Value  |
| ----------------- | -------------- |
| Width             | 256 pixels     |
| Height            | 256 pixels     |
| Interpolation     | Bilinear       |
| Distance Function | L1 (Manhattan) |
| Hash Length       | 80 bytes       |
| Grid Size         | 4×4            |
| Edge Types        | 5              |

## Technical Details

* **Hash Type**: `hashtype.UInt8`
* **Hash Length**: 80 bytes (4×4 grid cells × 5 edge types)
* **Default Comparison**: L1 (Manhattan) distance
* **Color Processing**: Converted to grayscale
* **Edge Threshold**: 11.0
* **Quantization**: 3 bits per bin (8 levels) using MPEG-7 tables
* **Block Size**: 2×2 pixels for edge detection

## Edge Types

The five edge types detected by EHD:

1. **Vertical**: Edges running vertically
2. **Horizontal**: Edges running horizontally
3. **45° Diagonal**: Edges from bottom-left to top-right
4. **135° Diagonal**: Edges from top-left to bottom-right
5. **Non-directional**: Edges without clear orientation

## References

* MPEG-7 Edge Histogram Descriptor standard
* Uses standardized quantization tables for 3-bit histogram encoding
