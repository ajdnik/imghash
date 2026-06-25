# HOGHash (Histogram of Oriented Gradients)

> HOG-based perceptual hash for shape and object recognition

## Overview

HOGHash is a perceptual hash algorithm based on Histogram of Oriented Gradients (HOG), a feature descriptor widely used in computer vision for object detection and recognition. It captures the distribution of gradient orientations in an image, making it effective for shape-based similarity.

## How It Works

The HOGHash algorithm processes images through the following steps:

1. **Preprocessing**: The image is resized and converted to grayscale
2. **Gradient computation**: For each pixel, compute:
   * Gradient magnitude: √(gx² + gy²)
   * Gradient orientation: atan2(gy, gx) mapped to 0-180° (unsigned)
3. **Cell division**: The image is divided into square cells of specified size
4. **Histogram computation**: For each cell:
   * Build an orientation histogram with specified number of bins (0-180°)
   * Weight each gradient by its magnitude
5. **Normalization**: Each cell histogram is normalized (0-255 range)
6. **Concatenation**: All cell histograms are concatenated into a single descriptor

> [!NOTE]
> HOGHash returns a **UInt8** hash type with size = (width/cellSize) × (height/cellSize) × numBins bytes.

## When to Use HOGHash

HOGHash is particularly effective for:

* Object detection and recognition
* Shape-based image similarity
* Pedestrian detection
* Logo and icon matching
* Structure and contour matching
* Images where edge orientation is discriminative

It is less suitable for:

* Color-based matching (uses grayscale only)
* Texture without directional structure
* Images with very uniform regions

## Constructor

```go
func NewHOGHash(opts ...HOGHashOption) (HOGHash, error)
```

### Options

- **`WithSize`** `func(width, height uint) SizeOption`
  Sets the resize dimensions before hash computation. The image is resized to this size before processing.
  **Default**: 256×256 pixels

- **`WithInterpolation`** `func(interp Interpolation) InterpolationOption`
  Sets the interpolation method used when resizing the image.
  **Default**: `Bilinear`
  Available options: `NearestNeighbor`, `Bilinear`, `Bicubic`, `Lanczos`

- **`WithCellSize`** `func(size uint) CellSizeOption`
  Sets the cell size in pixels (square cells) for HOG computation.
  **Default**: 8 pixels
  Smaller cells capture finer details but produce longer hashes.

- **`WithNumBins`** `func(bins uint) NumBinsOption`
  Sets the number of orientation histogram bins (0-180°).
  **Default**: 9 bins (20° per bin)
  More bins provide finer orientation resolution but increase hash size.

- **`WithDistance`** `func(fn DistanceFunc) DistanceOption`
  Overrides the default distance function used by the `Compare` method.
  **Default**: `similarity.Cosine`
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
	// Create HOGHash with default settings
	hasher, err := imghash.NewHOGHash()
	if err != nil {
		panic(err)
	}

	// Or with custom options
	customHasher, err := imghash.NewHOGHash(
		imghash.WithSize(256, 256),
		imghash.WithCellSize(16),   // 16x16 pixel cells
		imghash.WithNumBins(12),     // 12 orientation bins
		imghash.WithInterpolation(imghash.Bicubic),
	)
	if err != nil {
		panic(err)
	}

	// Load images
	img1, _ := loadImage("object1.jpg")
	img2, _ := loadImage("object2.jpg")

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

	fmt.Printf("HOGHash distance: %.2f\n", distance)
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

| Parameter         | Default Value                      |
| ----------------- | ---------------------------------- |
| Width             | 256 pixels                         |
| Height            | 256 pixels                         |
| Interpolation     | Bilinear                           |
| Cell Size         | 8 pixels                           |
| Number of Bins    | 9 bins                             |
| Distance Function | Cosine                             |
| Hash Length       | 2,304 bytes (32×32 cells × 9 bins) |

## Technical Details

* **Hash Type**: `hashtype.UInt8`
* **Hash Length**: (width/cellSize) × (height/cellSize) × numBins bytes
* **Default Comparison**: Cosine distance
* **Color Processing**: Converted to grayscale
* **Gradient Method**: Central differences
* **Orientation Range**: 0-180° (unsigned gradients)
* **Magnitude Weighting**: Histograms weighted by gradient magnitude
* **Normalization**: Per-cell normalization to 0-255 range

## Configuration Examples

| Image Size | Cell Size | Bins | Hash Length | Use Case                     |
| ---------- | --------- | ---- | ----------- | ---------------------------- |
| 256×256    | 8         | 9    | 2,304 bytes | Default, general purpose     |
| 256×256    | 16        | 9    | 576 bytes   | Faster, coarser features     |
| 256×256    | 4         | 9    | 9,216 bytes | Detailed, fine features      |
| 256×256    | 8         | 18   | 4,608 bytes | Finer orientation resolution |

## References

* Dalal, N., & Triggs, B. (2005). "Histograms of Oriented Gradients for Human Detection." IEEE Computer Society Conference on Computer Vision and Pattern Recognition (CVPR).
* [IEEE Xplore](https://ieeexplore.ieee.org/document/1467360)
