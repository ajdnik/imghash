# Difference Hash

> Gradient-based perceptual hash for tracking image changes

## Overview

Difference Hash (also known as dHash) is a perceptual hashing algorithm described by Dr. Neal Krawetz in "Kind of Like That". Instead of comparing pixels to the mean, it compares adjacent pixels to capture horizontal gradients.

This approach is more robust to brightness adjustments and better captures structural information than Average Hash.

## When to Use

Use Difference Hash when you need:

* **Better robustness** than Average Hash
* **Resistance to brightness changes** and gamma correction
* **Fast computation** with low overhead
* **Gradient-based features** that capture image structure
* **Good balance** between speed and accuracy

> [!NOTE]
> Difference Hash is particularly effective at detecting images that differ only in brightness or exposure settings.

## Constructor

```go
func NewDifference(opts ...DifferenceOption) (Difference, error)
```

### Available Options

* `WithSize(width, height uint)` - Sets the resize dimensions (hash is width × height bits)
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
    // Create Difference hasher with default settings
    dhash, err := imghash.NewDifference()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(dhash, "image.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Difference hash: %v\n", hash)
}
```

### With Custom Options

```go
// Create Difference hasher with custom size
dhash, err := imghash.NewDifference(
    imghash.WithSize(16, 16),  // 256-bit hash
    imghash.WithInterpolation(imghash.Bicubic),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(dhash, "image.jpg")
```

## Default Settings

* **Hash size**: 64 bits (8 bytes)
* **Resize dimensions**: 8×8 pixels
* **Interpolation**: Bilinear
* **Distance metric**: Hamming distance

## How It Works

The Difference Hash algorithm:

1. Resizes the image to (width+1) × height pixels (e.g., 9×8 for default)
2. Converts to grayscale
3. For each row, compares adjacent pixels from left to right:
   * Sets bit to 1 if right pixel > left pixel
   * Sets bit to 0 if right pixel ≤ left pixel
4. Produces a width × height bit hash (64 bits by default)

> [!TIP]
> The algorithm compares pixel\[i] with pixel\[i-1] to detect gradients, which makes it more robust to uniform brightness changes than Average Hash.

## Comparison

Difference hashes are compared using Hamming distance. Lower distances indicate more similar images.

```go
dhash, _ := imghash.NewDifference()

h1, _ := imghash.HashFile(dhash, "original.jpg")
h2, _ := imghash.HashFile(dhash, "brightened.jpg")

dist, err := dhash.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 10 {
    fmt.Println("Images are similar (possibly same image with different brightness)")
}
```

## Advantages Over Average Hash

* **More robust to brightness/exposure changes**
* **Better structural information** through gradient comparison
* **Similar computational cost**
* **Still very fast and simple**

## Performance

Difference Hash is slightly more expensive than Average Hash but still very fast:

* Simple resize operation
* Single-pass gradient computation
* No mean calculation needed
* No complex transforms

## References

* [Kind of Like That by Dr. Neal Krawetz](https://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html)
* Source: `difference.go:11-13`

> [!WARNING]
> While more robust than Average Hash, Difference Hash is still susceptible to significant transformations like rotation, cropping, or perspective changes. Consider PDQ or PHash for more robust matching.
