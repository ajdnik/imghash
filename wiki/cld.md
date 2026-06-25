# CLD (Color Layout Descriptor)

> MPEG-7 Color Layout Descriptor perceptual hash

## Overview

CLD (Color Layout Descriptor) is an MPEG-7 standard perceptual hash algorithm that captures the spatial distribution of colors in an image. It provides a compact 12-element descriptor that is effective for color-based image similarity.

## How It Works

The CLD algorithm processes images through the following steps:

1. **Color space conversion**: The image is converted to YCbCr color space
2. **Layout summarization**: The image is divided into an 8×8 grid, with average color values computed for each cell
3. **DCT transformation**: A 2-D Discrete Cosine Transform (DCT) is applied to each color channel
4. **Coefficient selection**: Low-frequency zig-zag DCT coefficients are extracted:
   * 6 Y (luminance) coefficients
   * 3 Cb (blue chrominance) coefficients
   * 3 Cr (red chrominance) coefficients
5. **Quantization**: Coefficients are quantized into a compact 12-element UInt8 descriptor

> [!NOTE]
> CLD returns a **UInt8** hash type (12 bytes), not Binary or Float64.

## When to Use CLD

CLD is particularly effective for:

* Color-based image retrieval
* Images with distinct color layouts
* Applications requiring compact color descriptors
* Thumbnail or icon matching
* Images where spatial color distribution is important

It is less suitable for:

* Grayscale images (though it will still work)
* Images where fine texture details are more important than color

## Constructor

```go
func NewCLD(opts ...CLDOption) (CLD, error)
```

### Options

- **`WithSize`** `func(width, height uint) SizeOption`
  Sets the resize dimensions before hash computation. The image is resized to this size before processing.
  **Default**: 64×64 pixels

- **`WithInterpolation`** `func(interp Interpolation) InterpolationOption`
  Sets the interpolation method used when resizing the image.
  **Default**: `Bilinear`
  Available options: `NearestNeighbor`, `Bilinear`, `Bicubic`, `Lanczos`

- **`WithDistance`** `func(fn DistanceFunc) DistanceOption`
  Overrides the default distance function used by the `Compare` method.
  **Default**: `similarity.L2` (Euclidean distance)
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
	// Create CLD hasher with default settings
	hasher, err := imghash.NewCLD()
	if err != nil {
		panic(err)
	}

	// Or with custom options
	customHasher, err := imghash.NewCLD(
		imghash.WithSize(128, 128),
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

	fmt.Printf("CLD distance: %.2f\n", distance)
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
| Width             | 64 pixels      |
| Height            | 64 pixels      |
| Interpolation     | Bilinear       |
| Distance Function | L2 (Euclidean) |
| Hash Length       | 12 bytes       |

## Technical Details

* **Hash Type**: `hashtype.UInt8`
* **Hash Length**: 12 bytes (6 Y + 3 Cb + 3 Cr coefficients)
* **Default Comparison**: L2 (Euclidean) distance
* **Color Space**: YCbCr
* **Grid Size**: 8×8
* **DC Quantization Range**: 0-63 (6 bits)
* **AC Quantization Range**: 0-31 (5 bits, centered at 16)

## References

* MPEG-7 Color Layout Descriptor standard
* Based on DCT transformation and zig-zag coefficient ordering
