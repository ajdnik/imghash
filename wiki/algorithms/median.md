# Median Hash

> Robust variant of Average Hash using median instead of mean

## Overview

Median Hash is a variant of Average Hash that uses the median pixel value instead of the mean for thresholding. This makes it more robust to outliers and extreme pixel values.

By using the median, the algorithm is less sensitive to bright spots, dark shadows, or other extreme values that might skew the mean.

## When to Use

Use Median Hash when you need:

* **Robustness to outliers** and extreme pixel values
* **Better handling** of images with bright spots or dark shadows
* **Simple, fast computation** like Average Hash
* **Resistance to noise** in certain scenarios
* **Median-based thresholding** instead of mean-based

> [!NOTE]
> Median Hash is particularly useful for images with uneven lighting, lens flare, or other artifacts that create extreme pixel values.

## Constructor

```go
func NewMedian(opts ...MedianOption) (Median, error)
```

### Available Options

* `WithSize(width, height uint)` - Sets the resize dimensions
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
    // Create Median hasher with default settings
    median, err := imghash.NewMedian()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(median, "image.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Median hash: %v\n", hash)
}
```

### With Custom Options

```go
// Create Median hasher with custom size
median, err := imghash.NewMedian(
    imghash.WithSize(12, 12),  // 144-bit hash
    imghash.WithInterpolation(imghash.Lanczos3),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(median, "image.jpg")
```

## Default Settings

* **Hash size**: 64 bits (8 bytes)
* **Resize dimensions**: 8×8 pixels
* **Interpolation**: Bilinear
* **Distance metric**: Hamming distance

## How It Works

The Median Hash algorithm:

1. Resizes the image to the specified dimensions (default 8×8)
2. Converts to grayscale
3. Computes the **median** pixel value (not mean)
4. For each pixel:
   * Sets bit to 1 if pixel value > median
   * Sets bit to 0 if pixel value ≤ median
5. Produces a binary hash (width × height bits)

> [!TIP]
> The median is the middle value when all pixels are sorted. Unlike the mean, it's not affected by extreme values, making it more robust to outliers.

## Comparison

Median hashes are compared using Hamming distance. Lower distances indicate more similar images.

```go
median, _ := imghash.NewMedian()

h1, _ := imghash.HashFile(median, "original.jpg")
h2, _ := imghash.HashFile(median, "with_flare.jpg")

dist, err := median.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 10 {
    fmt.Println("Images are similar (even with lens flare)")
}
```

## Median vs. Mean

| Aspect                  | Median Hash                | Average Hash           |
| ----------------------- | -------------------------- | ---------------------- |
| **Threshold**           | Median pixel value         | Mean pixel value       |
| **Outlier sensitivity** | Low (robust)               | High (sensitive)       |
| **Computation**         | Requires sorting           | Simple sum/count       |
| **Best for**            | Images with extreme values | Clean, well-lit images |
| **Speed**               | Slightly slower            | Fastest                |

## Performance

Median Hash is slightly slower than Average Hash due to median calculation:

* Simple resize operation
* Median computation requires sorting
* Direct pixel comparison
* No complex transforms

The overhead is minimal for small hash sizes (8×8 requires sorting only 64 values).

## References

* [DupImageLib implementation](https://github.com/Quickshot/DupImageLib/blob/3e914588958c4c1871d750de86b30446b9c07a3e/DupImageLib/ImageHashes.cs#L99)
* Source: `medianhash.go:12-14`

> [!WARNING]
> Like Average Hash, Median Hash is not robust to significant transformations like rotation, scaling, or cropping. Use DCT-based methods for more robust matching.
