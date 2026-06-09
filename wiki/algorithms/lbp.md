# LBP (Local Binary Patterns)

> Local Binary Patterns perceptual hash for texture analysis

## Overview

LBP (Local Binary Patterns) is a perceptual hash algorithm that excels at texture description and classification. It computes local binary patterns for each pixel by comparing it with its neighbors, then builds histograms of these patterns across grid cells.

## How It Works

The LBP algorithm processes images through the following steps:

1. **Preprocessing**: The image is resized and converted to grayscale
2. **LBP code computation**: For each pixel, compare its value with its 8 neighbors (clockwise from right):
   * If neighbor ≥ center: set bit to 1
   * If neighbor \< center: set bit to 0
   * Creates an 8-bit code (0-255) for each pixel
3. **Grid division**: The LBP code image is divided into grid cells
4. **Histogram computation**: A 256-bin histogram is built for each cell
5. **Normalization**: Each histogram is normalized (0-255 range)
6. **Concatenation**: All cell histograms are concatenated into a single descriptor

> [!NOTE]
> LBP returns a **UInt8** hash type with size = gridX × gridY × 256 bytes.

## When to Use LBP

LBP is particularly effective for:

* Texture classification and recognition
* Face recognition and detection
* Material identification
* Surface analysis
* Pattern matching in grayscale images
* Rotation-invariant texture matching (with rotation-invariant variants)

It is less suitable for:

* Color-based matching (uses grayscale only)
* Images where color is more important than texture
* Smooth images without texture variation

## Constructor

```go
func NewLBP(opts ...LBPOption) (LBP, error)
```

### Options

- **`WithSize`** `func(width, height uint) SizeOption` — Sets the resize dimensions before hash computation. The image is resized to this size before processing.

  **Default**: 256×256 pixels

- **`WithInterpolation`** `func(interp Interpolation) InterpolationOption` — Sets the interpolation method used when resizing the image.

  **Default**: `Bilinear`

  Available options: `NearestNeighbor`, `Bilinear`, `Bicubic`, `Lanczos`

- **`WithGridSize`** `func(x, y uint) GridSizeOption` — Sets the number of grid cells used to divide the image for spatial histogram computation.

  **Default**: 1×1 (single histogram for entire image)

  Larger grids capture more spatial information but produce longer hashes.

- **`WithDistance`** `func(fn DistanceFunc) DistanceOption` — Overrides the default distance function used by the `Compare` method.

  **Default**: `similarity.ChiSquare`

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
	// Create LBP hasher with default settings (1x1 grid)
	hasher, err := imghash.NewLBP()
	if err != nil {
		panic(err)
	}

	// Or with custom options for spatial histograms
	customHasher, err := imghash.NewLBP(
		imghash.WithSize(256, 256),
		imghash.WithGridSize(4, 4), // 4x4 grid = 4096 byte hash
		imghash.WithInterpolation(imghash.Bicubic),
	)
	if err != nil {
		panic(err)
	}

	// Load images
	img1, _ := loadImage("texture1.jpg")
	img2, _ := loadImage("texture2.jpg")

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

	fmt.Printf("LBP distance: %.2f\n", distance)
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

| Parameter         | Default Value            |
| ----------------- | ------------------------ |
| Width             | 256 pixels               |
| Height            | 256 pixels               |
| Interpolation     | Bilinear                 |
| Grid Size         | 1×1                      |
| Distance Function | Chi-Square               |
| Hash Length       | 256 bytes (for 1×1 grid) |

## Technical Details

* **Hash Type**: `hashtype.UInt8`
* **Hash Length**: gridX × gridY × 256 bytes
* **Default Comparison**: Chi-Square distance
* **Color Processing**: Converted to grayscale
* **Neighborhood**: 8 neighbors (3×3 window)
* **Pattern Range**: 0-255 (256 possible patterns)
* **Border Handling**: Out-of-bounds neighbors treated as zero

## Grid Size Considerations

| Grid Size | Hash Length | Spatial Detail   | Use Case                      |
| --------- | ----------- | ---------------- | ----------------------------- |
| 1×1       | 256 bytes   | Global texture   | Overall texture matching      |
| 2×2       | 1024 bytes  | Basic spatial    | Quadrant-based matching       |
| 4×4       | 4096 bytes  | Detailed spatial | Regional texture analysis     |
| 8×8       | 16384 bytes | Very detailed    | Fine-grained spatial matching |

## References

* Ojala, T., Pietikäinen, M., & Mäenpää, T. (2002). "Multiresolution Gray-Scale and Rotation Invariant Texture Classification with Local Binary Patterns." IEEE Transactions on Pattern Analysis and Machine Intelligence.
* [IEEE Xplore](https://ieeexplore.ieee.org/document/1017623)
