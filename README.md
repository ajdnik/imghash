# imghash

[![CI](https://github.com/ajdnik/imghash/workflows/ci/badge.svg "CI status")](https://github.com/ajdnik/imghash/actions?query=workflow%3Aci)
[![Coverage Status](https://badge.coveralls.io/repos/github/ajdnik/imghash/badge.svg?branch=main)](https://coveralls.io/github/ajdnik/imghash?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/ajdnik/imghash/v2.svg)](https://pkg.go.dev/github.com/ajdnik/imghash/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/ajdnik/imghash/v2)](https://goreportcard.com/report/github.com/ajdnik/imghash/v2)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/ajdnik/imghash/blob/main/LICENSE)

Go implementation of multiple perceptual hash algorithms for images. Perceptual hash functions are analogous if features are similar, whereas cryptographic hashing relies on the avalanche effect of a small change in input value creating a drastic change in output value.

## Community

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Contributing Guide](CONTRIBUTING.md)
- [Security Policy](SECURITY.md)
- [Migration Guide](MIGRATION.md)

## Installing

Using imghash is easy. First, use `go get` to install the latest version
of the library. This command will install the library and its dependencies:

    go get -u github.com/ajdnik/imghash/v2

Next, include imghash in your application:

```go
import "github.com/ajdnik/imghash/v2"
```

Most consumers only need the top-level `imghash` package. The core types (`Hash`, `Binary`, `UInt8`, `Float64`, `Distance`) are re-exported there. The `similarity` sub-package is available when you need a specific metric like PCC.

## Quick Start

If you're not sure which hash algorithm to use, go with PDQ. It's generally the best all-around choice — robust to JPEG compression, rescaling, and minor edits. Pick a different algorithm only if you have specific needs (e.g. rotation invariance, color-aware comparison, or a particular hash size).

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash/v2"
)

func main() {
  pdq, err := imghash.NewPDQ()
  if err != nil {
    panic(err)
  }

  h1, err := imghash.HashFile(pdq, "image1.png")
  if err != nil {
    panic(err)
  }

  h2, err := imghash.HashFile(pdq, "image2.png")
  if err != nil {
    panic(err)
  }

  // Use the algorithm's recommended distance metric
  dist, err := pdq.Compare(h1, h2)
  if err != nil {
    panic(err)
  }
  fmt.Printf("Distance: %v\n", dist)
}
```

## Convenience Functions

The library provides helpers so you don't need image decoding boilerplate:

| Function | Description |
|----------|-------------|
| `OpenImage(path)` | Reads and decodes an image file (JPEG, PNG, GIF) |
| `DecodeImage(r)` | Decodes an image from any `io.Reader` |
| `HashFile(hasher, path)` | Opens a file and computes its hash in one call |
| `HashReader(hasher, r)` | Decodes from a reader and computes the hash |
| `Compare(h1, h2)` | Computes distance using the natural metric for the hash type |

Use the algorithm's `Compare` method for its recommended metric, or call top-level `imghash.Compare(h1, h2)` for generic type-based comparison (see [Similarity Metrics](#similarity-metrics)).

## Perceptual Hash Algorithms

The library supports 13 perceptual hashing algorithms. Most are ported from [OpenCV Contrib](https://github.com/opencv/opencv_contrib) and tested against its implementations.

Every constructor accepts functional options. Call with no arguments for defaults, or pass `With*` options to customize:

```go
// Defaults
avg, err := imghash.NewAverage()

// Custom parameters
avg, err := imghash.NewAverage(
  imghash.WithSize(16, 16),
  imghash.WithInterpolation(imghash.Bicubic),
)
```

#### Average Hash

Produces a 64-bit binary hash. Compares using Hamming distance. Based on [Looks Like It](https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html) by Dr. Neal Krawetz.

```go
avg, err := imghash.NewAverage()
h1, err := avg.Calculate(img1)
h2, err := avg.Calculate(img2)
dist, err := avg.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

#### Difference Hash

Very similar to Average but computes gradients instead of averages. Compares using Hamming distance. Based on [Kind of Like That](https://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html).

```go
diff, err := imghash.NewDifference()
h1, err := diff.Calculate(img1)
h2, err := diff.Calculate(img2)
dist, err := diff.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

#### Median Hash

Almost identical to Average hash but uses the median instead of the mean. Compares using Hamming distance.

```go
med, err := imghash.NewMedian()
h1, err := med.Calculate(img1)
h2, err := med.Calculate(img2)
dist, err := med.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

#### pHash

Uses a discrete cosine transform to produce a 64-bit binary hash. Compares using weighted Hamming distance, where each byte position can be assigned a different weight (defaults to uniform). From [Rihamark: Perceptual image hash benchmarking](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking).

```go
ph, err := imghash.NewPHash()
h1, err := ph.Calculate(img1)
h2, err := ph.Calculate(img2)
dist, err := ph.Compare(h1, h2) // Weighted Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 32, 32 |
| `WithInterpolation(i)` | `BilinearExact` |
| `WithWeights(w)` | `[1, 1, 1, 1, 1, 1, 1, 1]` |

#### Wavelet Hash (wHash)

Uses a multi-level Haar discrete wavelet transform to produce a 64-bit binary hash. Compares using Hamming distance. The image is resized to `(width * 2^level) x (height * 2^level)`, converted to grayscale, and decomposed via the Haar DWT. The low-frequency (LL) subband coefficients are thresholded against their median. See [Wavelet image hash](https://fullstackml.com/wavelet-image-hash-in-python-3504571f3b08) for more information.

```go
wh, err := imghash.NewWHash()
h1, err := wh.Calculate(img1)
h2, err := wh.Calculate(img2)
dist, err := wh.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithLevel(l)` | 3 |

#### Color Moments Hash

Builds a 42-element float64 vector from YCbCr and HSV color moments. Compares using L2 (Euclidean) distance. Based on [Perceptual hashing for color images using invariant moments](https://www.researchgate.net/publication/286870507_Perceptual_hashing_for_color_images_using_invariant_moments).

```go
cm, err := imghash.NewColorMoment()
h1, err := cm.Calculate(img1)
h2, err := cm.Calculate(img2)
dist, err := cm.Compare(h1, h2) // L2 distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 512, 512 |
| `WithInterpolation(i)` | `Bicubic` |
| `WithKernelSize(k)` | 3 |
| `WithSigma(s)` | 0 |

#### Marr-Hildreth Hash

Uses a 2D wavelet transform to produce a 576-bit binary hash. Compares using Hamming distance. From [Rihamark: Perceptual image hash benchmarking](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking).

```go
mh, err := imghash.NewMarrHildreth()
h1, err := mh.Calculate(img1)
h2, err := mh.Calculate(img2)
dist, err := mh.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 512, 512 |
| `WithInterpolation(i)` | `Bicubic` |
| `WithScale(s)` | 1 |
| `WithAlpha(a)` | 2 |
| `WithKernelSize(k)` | 7 |
| `WithSigma(s)` | 0 |

#### Block Mean Value Hash

Computes means using a sliding window to produce a 256-bit binary hash. Compares using Hamming distance. Based on [Block Mean Value Based Image Perceptual Hashing](https://ieeexplore.ieee.org/document/4041692).

```go
bm, err := imghash.NewBlockMean()
h1, err := bm.Calculate(img1)
h2, err := bm.Calculate(img2)
dist, err := bm.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `BilinearExact` |
| `WithBlockSize(w, h)` | 16, 16 |
| `WithBlockMeanMethod(m)` | `Direct` |

Block mean methods: `Direct`, `Overlap`, `Rotation`, `RotationOverlap`.

#### Local Binary Pattern (LBP) Hash

Computes a 3x3 Local Binary Pattern code for each pixel and builds a normalized 256-bin histogram per grid cell, producing a uint8 vector. Compares using chi-square distance, which is well suited for histogram comparison. The grid can be increased for spatially-aware hashing. Based on [Multiresolution Gray-Scale and Rotation Invariant Texture Classification with Local Binary Patterns](https://ieeexplore.ieee.org/document/1017623) by Ojala et al.

```go
lbp, err := imghash.NewLBP()
h1, err := lbp.Calculate(img1)
h2, err := lbp.Calculate(img2)
dist, err := lbp.Compare(h1, h2) // Chi-square distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithGridSize(x, y)` | 1, 1 |

With the default 1x1 grid the hash is a 256-element uint8 vector. Set `WithGridSize(4, 4)` for a 4096-element spatially-aware hash.

#### HOG Hash (Histogram of Oriented Gradients)

Computes gradient magnitudes and orientations at each pixel, divides the image into square cells, and builds a magnitude-weighted orientation histogram per cell. The histograms are normalized and concatenated into a uint8 vector. Compares using cosine distance, which measures orientation similarity independent of magnitude scaling. Based on [Histograms of Oriented Gradients for Human Detection](https://ieeexplore.ieee.org/document/1467360) by Dalal and Triggs.

```go
hog, err := imghash.NewHOGHash()
h1, err := hog.Calculate(img1)
h2, err := hog.Calculate(img2)
dist, err := hog.Compare(h1, h2) // Cosine distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithCellSize(s)` | 8 |
| `WithNumBins(n)` | 9 |

With default settings the hash is a 9216-element uint8 vector (32×32 cells × 9 bins). Use `WithSize(32, 32)` for a compact 144-element hash (4×4 cells × 9 bins).

#### Radial Variance Hash

Uses radial projections and feature vector computations to generate a 40-element uint8 vector. Compares using L1 (Manhattan) distance. From [Robust image hashing based on radial variance of pixels](https://www.researchgate.net/publication/4186555_Robust_image_hashing_based_on_radial_variance_of_pixels).

```go
rv, err := imghash.NewRadialVariance()
h1, err := rv.Calculate(img1)
h2, err := rv.Calculate(img2)
dist, err := rv.Compare(h1, h2) // L1 distance
```

| Option | Default |
|--------|---------|
| `WithSigma(s)` | 1 |
| `WithAngles(n)` | 180 |

#### RASH (Rotation Aware Spatial Hash)

Produces a binary hash (64 bits with default settings) designed to be robust against image rotation. Compares using Hamming distance. This is a custom algorithm that combines concentric ring sampling for rotation invariance, a 1-D DCT for frequency compaction, and median thresholding for binarisation. The algorithm resizes the image, converts to grayscale, applies Gaussian blur, then samples pixel intensities on concentric rings around the image centre. Because ring-mean features are inherently rotation-invariant (rotating the image only permutes pixels within a ring, leaving its mean unchanged), the resulting hash stays stable under arbitrary rotations. Inspired by ring-partition hashing literature such as [Robust Image Hashing with Ring Partition and Invariant Vector Distance](https://ieeexplore.ieee.org/document/7368930) by Tang et al. and [Robust image hashing based on radial variance of pixels](https://www.researchgate.net/publication/4186555_Robust_image_hashing_based_on_radial_variance_of_pixels) by De Roover et al.

```go
rash, err := imghash.NewRASH()
h1, err := rash.Calculate(img1)
h2, err := rash.Calculate(img2)
dist, err := rash.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithSigma(s)` | 1 |
| `WithRings(n)` | 180 |

#### PDQ Hash

Produces a 256-bit binary hash designed for large-scale image deduplication. Compares using Hamming distance. The algorithm resizes to 64×64, applies a Jarosz box filter for noise smoothing, computes a 2D DCT, and thresholds the 16×16 low-frequency coefficients against their median. Robust to JPEG compression, rescaling, and minor edits. From Facebook (Meta) [ThreatExchange PDQ](https://github.com/facebook/ThreatExchange/tree/main/pdq).

```go
pdq, err := imghash.NewPDQ()
h1, err := pdq.Calculate(img1)
h2, err := pdq.Calculate(img2)
dist, err := pdq.Compare(h1, h2) // Hamming distance
```

| Option | Default |
|--------|---------|
| `WithInterpolation(i)` | `Bilinear` |

The input size is fixed at 64×64 per the algorithm specification and is not configurable.

### Binary hash size with custom options

For binary hashers with configurable dimensions, the bit count may not be a multiple of 8. In that case:
- the `Binary` value stores `ceil(bits/8)` bytes
- unused bit positions in the last byte remain zero

Examples:
- `WithSize(3, 3)` for Average/Median/Difference/WHash yields 9 bits stored in 2 bytes
- `RASH` uses `min(64, rings-1)` bits, so `WithRings(5)` yields 4 bits stored in 1 byte

## Similarity Metrics

### Algorithm-specific `Compare`

Every hash algorithm has a `Compare` method that uses the recommended distance metric for that algorithm. This is the preferred way to compare hashes:

```go
pdq, _ := imghash.NewPDQ()
h1, _ := imghash.HashFile(pdq, "image1.png")
h2, _ := imghash.HashFile(pdq, "image2.png")
dist, err := pdq.Compare(h1, h2)
```

Algorithm `Compare` methods validate inputs before computing distance:
- hash values must be the expected type for that algorithm
- hash lengths must match
- otherwise they return `ErrIncompatibleHash` or `ErrHashLengthMismatch`

The default metric per algorithm:

| Algorithm | Hash type | Default `Compare` metric |
|-----------|-----------|--------------------------|
| Average | `Binary` | Hamming |
| Difference | `Binary` | Hamming |
| Median | `Binary` | Hamming |
| PHash | `Binary` | Weighted Hamming |
| WHash | `Binary` | Hamming |
| MarrHildreth | `Binary` | Hamming |
| BlockMean | `Binary` | Hamming |
| PDQ | `Binary` | Hamming |
| RASH | `Binary` | Hamming |
| ColorMoment | `Float64` | L2 (Euclidean) |
| LBP | `UInt8` | Chi-Square |
| HOGHash | `UInt8` | Cosine |
| RadialVariance | `UInt8` | L1 (Manhattan) |

### Generic `Compare`

Use `imghash.Compare(h1, h2)` for generic comparison based on hash type:

| Hash type | Default metric | Description |
|-----------|---------------|-------------|
| `Binary` | Hamming | Number of differing bits |
| `UInt8` | L2 (Euclidean) | Square root of sum of squared differences |
| `Float64` | L2 (Euclidean) | Square root of sum of squared differences |

```go
dist, err := imghash.Compare(h1, h2)
```

`imghash.Compare` returns `ErrIncompatibleHash` when one hash is `Binary` and the other is not.

### Available metrics

The `similarity` sub-package provides all metrics for direct use:

```go
import "github.com/ajdnik/imghash/v2/similarity"

dist, err := similarity.Hamming(h1, h2)          // bit-level Hamming distance (Binary only)
dist, err = similarity.WeightedHamming(h1, h2, weights) // weighted Hamming (Binary only)
dist, err = similarity.L1(h1, h2)                       // Manhattan distance
dist, err = similarity.L2(h1, h2)                       // Euclidean distance
dist, err = similarity.ChiSquare(h1, h2)                // Chi-square distance
dist, err = similarity.Cosine(h1, h2)                   // Cosine distance (1 - cos similarity)
dist, err = similarity.PCC(h1, h2)                      // Peak cross-correlation
```

## Interpolation Methods

All resize-based algorithms accept `WithInterpolation`. Available methods:

`NearestNeighbor`, `Bilinear`, `Bicubic`, `MitchellNetravali`, `Lanczos2`, `Lanczos3`, `BilinearExact`

## License

Imghash is released under the MIT license. See [LICENSE](https://github.com/ajdnik/imghash/blob/main/LICENSE)
