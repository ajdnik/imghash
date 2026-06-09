# RadialVariance

> Radial variance-based perceptual hash for robust image matching

## Overview

RadialVariance is a perceptual hash algorithm that uses radial projection analysis to create rotation-aware image descriptors. It computes the variance of pixel values along radial lines at different angles, making it effective for images that may be rotated.

## How It Works

The RadialVariance algorithm processes images through the following steps:

1. **Preprocessing**: The image is converted to grayscale
2. **Gaussian blur**: A Gaussian blur is applied with specified sigma to reduce noise
3. **Radial projections**: Pixel values are projected along lines at different angles through the image center
4. **Variance computation**: For each angle, compute the variance of projected pixel values
5. **Normalization**: Feature vector is normalized (mean subtraction and variance scaling)
6. **DCT transformation**: Discrete Cosine Transform is applied to the feature vector
7. **Quantization**: DCT coefficients are quantized to a 40-element UInt8 descriptor

> [!NOTE]
> RadialVariance returns a **UInt8** hash type (40 bytes), not Binary or Float64.

## When to Use RadialVariance

RadialVariance is particularly effective for:

* Images that may be rotated
* Circular or radially symmetric objects
* Logo matching with orientation variations
* Images where rotation invariance is important
* Robust matching under geometric transformations

It is less suitable for:

* Images requiring color information (uses grayscale only)
* Very small images
* Images where exact orientation matters

## Constructor

```go
func NewRadialVariance(opts ...RadialVarianceOption) (RadialVariance, error)
```

### Options

- **`WithSigma`** `func(sigma float64) SigmaOption` — Sets the Gaussian kernel standard deviation for blur preprocessing.

  **Default**: 1.0

  Higher values create more blur, reducing noise but potentially losing detail.

- **`WithAngles`** `func(angles int) AnglesOption` — Sets the number of projection angles to consider.

  **Default**: 180

  More angles provide finer angular resolution but increase computation time.

- **`WithDistance`** `func(fn DistanceFunc) DistanceOption` — Overrides the default distance function used by the `Compare` method.

  **Default**: `similarity.L1` (Manhattan distance)

  Available functions: `Hamming`, `L1`, `L2`, `Cosine`, `ChiSquare`, `PCC`, `Jaccard`

> [!IMPORTANT]
> Unlike most algorithms, RadialVariance does not support `WithSize` or `WithInterpolation` options. It processes images at their original dimensions.

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
	// Create RadialVariance hasher with default settings
	hasher, err := imghash.NewRadialVariance()
	if err != nil {
		panic(err)
	}

	// Or with custom options
	customHasher, err := imghash.NewRadialVariance(
		imghash.WithSigma(1.5),    // More blur
		imghash.WithAngles(360),   // Finer angular resolution
	)
	if err != nil {
		panic(err)
	}

	// Load images (possibly rotated versions)
	img1, _ := loadImage("logo.jpg")
	img2, _ := loadImage("logo_rotated.jpg")

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

	fmt.Printf("RadialVariance distance: %.2f\n", distance)
	fmt.Println("(Lower distance indicates similar images, even if rotated)")
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
| Sigma             | 1.0            |
| Angles            | 180            |
| Distance Function | L1 (Manhattan) |
| Hash Length       | 40 bytes       |

## Technical Details

* **Hash Type**: `hashtype.UInt8`
* **Hash Length**: 40 bytes (fixed)
* **Default Comparison**: L1 (Manhattan) distance
* **Color Processing**: Converted to grayscale
* **Blur Method**: Gaussian blur with configurable sigma
* **Projection Method**: Radial lines from image center
* **Feature Extraction**: Variance along each radial projection
* **Transform**: DCT (Discrete Cosine Transform)
* **Normalization**: Mean-variance normalization of features

## Configuration Guidelines

### Sigma Values

| Sigma   | Blur Amount | Use Case                         |
| ------- | ----------- | -------------------------------- |
| 0.5     | Minimal     | Clean, high-quality images       |
| 1.0     | Default     | General purpose                  |
| 1.5-2.0 | Moderate    | Noisy images                     |
| > 2.0   | Heavy       | Very noisy or low-quality images |

### Angle Count

| Angles | Angular Resolution | Trade-off            |
| ------ | ------------------ | -------------------- |
| 90     | 4° per angle       | Faster, less precise |
| 180    | 2° per angle       | Default, balanced    |
| 360    | 1° per angle       | Slower, more precise |

## Implementation Notes

* The algorithm includes a small epsilon (0.00001) to avoid division by zero
* Handles variable image dimensions without resizing
* Projection calculations use rounding for pixel mapping
* DCT coefficients are scaled to 0-255 range for quantization

## References

* De Roover, C., De Vleeschouwer, C., Lefèbvre, F., & Macq, B. (2005). "Robust image hashing based on radial variance of pixels."
* [ResearchGate](https://www.researchgate.net/publication/4186555_Robust_image_hashing_based_on_radial_variance_of_pixels)
