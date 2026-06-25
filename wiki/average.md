# Average Hash

> Simple and fast perceptual hash based on pixel mean

## Overview

Average Hash (also known as aHash) is a simple perceptual hashing algorithm described by Dr. Neal Krawetz in "Looks Like It". It compares each pixel against the mean brightness to produce a binary hash.

This is one of the fastest and simplest perceptual hashing algorithms, making it ideal for quick similarity checks.

## When to Use

Use Average Hash when you need:

* **Fast computation** with minimal overhead
* **Simple implementation** that's easy to understand
* **Basic similarity detection** for images
* **Low memory footprint** (64-bit default hash)
* **Quick duplicate detection** in small to medium datasets

> [!WARNING]
> Average Hash is less robust to transformations than DCT-based methods like PDQ or PHash. Use it for basic similarity detection rather than precise matching.

## Constructor

```go
func NewAverage(opts ...AverageOption) (Average, error)
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
    // Create Average hasher with default settings
    avg, err := imghash.NewAverage()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(avg, "image.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Average hash: %v\n", hash)
}
```

### With Custom Options

```go
// Create Average hasher with custom size
avg, err := imghash.NewAverage(
    imghash.WithSize(16, 16),  // 256-bit hash
    imghash.WithInterpolation(imghash.Lanczos3),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(avg, "image.jpg")
```

## Default Settings

* **Hash size**: 64 bits (8 bytes)
* **Resize dimensions**: 8×8 pixels
* **Interpolation**: Bilinear
* **Distance metric**: Hamming distance

## How It Works

The Average Hash algorithm:

1. Resizes the image to the specified dimensions (default 8×8)
2. Converts to grayscale
3. Computes the mean pixel value
4. For each pixel:
   * Sets bit to 1 if pixel value > mean
   * Sets bit to 0 if pixel value ≤ mean
5. Produces a binary hash (width × height bits)

## Comparison

Average hashes are compared using Hamming distance. Lower distances indicate more similar images.

```go
avg, _ := imghash.NewAverage()

h1, _ := imghash.HashFile(avg, "image1.jpg")
h2, _ := imghash.HashFile(avg, "image2.jpg")

dist, err := avg.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 10 {
    fmt.Println("Images are similar")
}
```

> [!TIP]
> For 8×8 hashes (64 bits), a Hamming distance of 0-10 typically indicates very similar images, 10-20 indicates some similarity, and >20 indicates different images.

## Performance

Average Hash is extremely fast:

* Simple resize operation
* Single-pass mean calculation
* Direct pixel comparison
* No complex transforms (DCT, DWT, etc.)

## References

* [Looks Like It by Dr. Neal Krawetz](https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html)
