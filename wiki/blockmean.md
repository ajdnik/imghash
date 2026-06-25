# Block Mean Hash

> Block-based perceptual hash with rotation support

## Overview

Block Mean Hash is a perceptual hashing algorithm based on the method described in "Block Mean Value Based Image Perceptual Hashing" by Yang et al. It divides the image into blocks, computes the mean intensity of each block, and thresholds against the global mean.

The algorithm supports multiple methods including direct blocking, overlapping blocks, and rotation-aware variants.

## When to Use

Use Block Mean Hash when you need:

* **Block-based features** for structural matching
* **Rotation robustness** (with rotation methods)
* **Flexible block sizes** for different use cases
* **Overlapping blocks** for smoother features
* **High robustness** to various transformations

> [!NOTE]
> Block Mean with rotation variants produces very large hashes (3,840 bits for default settings) but offers exceptional rotation robustness.

## Constructor

```go
func NewBlockMean(opts ...BlockMeanOption) (BlockMean, error)
```

### Available Options

* `WithSize(width, height uint)` - Sets the resize dimensions
* `WithBlockSize(width, height uint)` - Sets the block dimensions
* `WithBlockMeanMethod(method BlockMeanMethod)` - Sets the block construction method
* `WithInterpolation(interp Interpolation)` - Sets the resize interpolation method
* `WithDistance(fn DistanceFunc)` - Overrides the default Hamming distance function

### Block Mean Methods

* `Direct` - Non-overlapping blocks (default)
* `Overlap` - Overlapping blocks with 50% overlap
* `Rotation` - Direct method with 24 rotated variants (15° steps)
* `RotationOverlap` - Overlap method with 24 rotated variants

### Supported Interpolation Methods

* `NearestNeighbor`
* `Bilinear`
* `BilinearExact` (default)
* `Bicubic`
* `MitchellNetravali`
* `Lanczos2`
* `Lanczos3`

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create Block Mean hasher with default settings
    bm, err := imghash.NewBlockMean()
    if err != nil {
        panic(err)
    }

    // Hash an image file
    hash, err := imghash.HashFile(bm, "photo.jpg")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Block Mean hash: %v\n", hash)
}
```

### With Overlapping Blocks

```go
// Create Block Mean hasher with overlapping blocks
bm, err := imghash.NewBlockMean(
    imghash.WithBlockMeanMethod(imghash.Overlap),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(bm, "photo.jpg")
```

### With Rotation Support

```go
// Create Block Mean hasher with rotation robustness
bm, err := imghash.NewBlockMean(
    imghash.WithBlockMeanMethod(imghash.Rotation),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(bm, "photo.jpg")
// Hash will be the same for rotated versions of the image
```

### With Custom Block Size

```go
// Create Block Mean hasher with custom blocks
bm, err := imghash.NewBlockMean(
    imghash.WithSize(256, 256),
    imghash.WithBlockSize(32, 32),  // Larger blocks
    imghash.WithBlockMeanMethod(imghash.Overlap),
)
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(bm, "photo.jpg")
```

## Default Settings

* **Hash size**: 256 bits (32 bytes) for Direct method
* **Resize dimensions**: 256×256 pixels
* **Block size**: 16×16 pixels
* **Method**: Direct (non-overlapping)
* **Interpolation**: BilinearExact
* **Distance metric**: Hamming distance

## Hash Sizes by Method

| Method              | Blocks   | Hash Bits | Hash Bytes |
| ------------------- | -------- | --------- | ---------- |
| **Direct**          | 16×16    | 256       | 32         |
| **Overlap**         | 31×31    | 961       | 121        |
| **Rotation**        | 16×16×24 | 6,144     | 768        |
| **RotationOverlap** | 31×31×24 | 23,064    | 2,884      |

## How It Works

### Direct Method

1. Resizes image to specified dimensions (default 256×256)
2. Converts to grayscale
3. Divides image into non-overlapping blocks (default 16×16 pixels)
4. Computes mean intensity of each block
5. Computes global mean of the image
6. For each block:
   * Sets bit to 1 if block mean ≥ global mean
   * Sets bit to 0 if block mean \< global mean
7. Produces a binary hash (number of blocks bits)

### Overlap Method

1. Same as Direct, but blocks overlap by 50%
2. Block stride = block\_size / 2
3. Creates 2×blocks-1 in each dimension
4. Results in 4× more blocks (e.g., 31×31 instead of 16×16)

### Rotation Methods

For `Rotation` or `RotationOverlap`:

1. Performs the same process as Direct/Overlap
2. Repeats for 24 rotated versions of the image (0°, 15°, 30°, ..., 345°)
3. Concatenates all rotation hashes
4. Produces a hash 24× larger

> [!TIP]
> Rotation methods make the hash rotation-invariant by including all rotations, but significantly increase hash size and computation time.

## Block Size Trade-offs

```go
// Smaller blocks: More detail, larger hash
bm_small, _ := imghash.NewBlockMean(
    imghash.WithSize(256, 256),
    imghash.WithBlockSize(8, 8),  // 32×32 = 1,024 blocks
)

// Larger blocks: Coarser features, smaller hash
bm_large, _ := imghash.NewBlockMean(
    imghash.WithSize(256, 256),
    imghash.WithBlockSize(32, 32),  // 8×8 = 64 blocks
)
```

## Comparison

Block Mean hashes are compared using Hamming distance:

```go
bm, _ := imghash.NewBlockMean()

h1, _ := imghash.HashFile(bm, "original.jpg")
h2, _ := imghash.HashFile(bm, "modified.jpg")

dist, err := bm.Compare(h1, h2)
if err != nil {
    panic(err)
}

if dist < 30 {
    fmt.Println("Images are similar")
}
```

> [!NOTE]
> Similarity thresholds depend on the hash size. For 256-bit Direct method, use thresholds around 20-40. For rotation methods with larger hashes, proportionally increase thresholds.

## Performance Characteristics

| Method              | Speed     | Memory    | Robustness |
| ------------------- | --------- | --------- | ---------- |
| **Direct**          | Fast      | Low       | Good       |
| **Overlap**         | Fast      | Low       | Better     |
| **Rotation**        | Very Slow | High      | Excellent  |
| **RotationOverlap** | Very Slow | Very High | Excellent  |

## Use Case Recommendations

* **Direct**: General-purpose hashing, good speed/accuracy balance
* **Overlap**: Smoother features, better robustness to minor shifts
* **Rotation**: When rotation-invariance is critical
* **RotationOverlap**: Maximum robustness (rotation + smooth features)

## References

* [Block Mean Value Based Image Perceptual Hashing - Yang et al.](https://ieeexplore.ieee.org/document/4041692)

> [!WARNING]
> Rotation methods (Rotation and RotationOverlap) are computationally expensive, requiring 24× the processing time and producing 24× larger hashes. Use them only when rotation-invariance is essential.
