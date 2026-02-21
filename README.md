# imghash

[![CI](https://github.com/ajdnik/imghash/workflows/ci/badge.svg "CI status")](https://github.com/ajdnik/imghash/actions?query=workflow%3Aci)
[![Coverage Status](https://coveralls.io/repos/github/ajdnik/imghash/badge.svg?branch=main)](https://coveralls.io/github/ajdnik/imghash?branch=main)
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

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash/v2"
)

func main() {
  phash, err := imghash.NewPHash()
  if err != nil {
    panic(err)
  }

  h1, err := imghash.HashFile(phash, "image1.png")
  if err != nil {
    panic(err)
  }

  h2, err := imghash.HashFile(phash, "image2.png")
  if err != nil {
    panic(err)
  }

  dist, err := imghash.Compare(h1, h2)
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

You can also call `h1.Distance(h2)` directly on any hash value.

## Perceptual Hash Algorithms

The library supports 12 perceptual hashing algorithms. Most are ported from [OpenCV Contrib](https://github.com/opencv/opencv_contrib) and tested against its implementations.

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
hash, err := avg.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

#### Difference Hash

Very similar to Average but computes gradients instead of averages. Based on [Kind of Like That](https://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html).

```go
diff, err := imghash.NewDifference()
hash, err := diff.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

#### Median Hash

Almost identical to Average hash but uses the median instead of the mean.

```go
med, err := imghash.NewMedian()
hash, err := med.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

#### pHash

Uses a discrete cosine transform to produce a 64-bit binary hash. From [Rihamark: Perceptual image hash benchmarking](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking).

```go
ph, err := imghash.NewPHash()
hash, err := ph.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 32, 32 |
| `WithInterpolation(i)` | `BilinearExact` |

#### Wavelet Hash (wHash)

Uses a multi-level Haar discrete wavelet transform to produce a 64-bit binary hash. The image is resized to `(width * 2^level) x (height * 2^level)`, converted to grayscale, and decomposed via the Haar DWT. The low-frequency (LL) subband coefficients are thresholded against their median. See [Wavelet image hash](https://fullstackml.com/wavelet-image-hash-in-python-3504571f3b08) for more information.

```go
wh, err := imghash.NewWHash()
hash, err := wh.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithLevel(l)` | 3 |

#### Color Moments Hash

Builds a 42-element float64 vector from YCbCr and HSV color moments. Based on [Perceptual hashing for color images using invariant moments](https://www.researchgate.net/publication/286870507_Perceptual_hashing_for_color_images_using_invariant_moments).

```go
cm, err := imghash.NewColorMoment()
hash, err := cm.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 512, 512 |
| `WithInterpolation(i)` | `Bicubic` |
| `WithKernelSize(k)` | 3 |
| `WithSigma(s)` | 0 |

#### Marr-Hildreth Hash

Uses a 2D wavelet transform to produce a 576-bit binary hash. From [Rihamark: Perceptual image hash benchmarking](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking).

```go
mh, err := imghash.NewMarrHildreth()
hash, err := mh.Calculate(img)
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

Computes means using a sliding window to produce a 256-bit binary hash. Based on [Block Mean Value Based Image Perceptual Hashing](https://ieeexplore.ieee.org/document/4041692).

```go
bm, err := imghash.NewBlockMean()
hash, err := bm.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `BilinearExact` |
| `WithBlockSize(w, h)` | 16, 16 |
| `WithBlockMeanMethod(m)` | `Direct` |

Block mean methods: `Direct`, `Overlap`, `Rotation`, `RotationOverlap`.

#### Local Binary Pattern (LBP) Hash

Computes a 3x3 Local Binary Pattern code for each pixel and builds a normalized 256-bin histogram per grid cell, producing a uint8 vector. The grid can be increased for spatially-aware hashing. Based on [Multiresolution Gray-Scale and Rotation Invariant Texture Classification with Local Binary Patterns](https://ieeexplore.ieee.org/document/1017623) by Ojala et al.

```go
lbp, err := imghash.NewLBP()
hash, err := lbp.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithGridSize(x, y)` | 1, 1 |

With the default 1x1 grid the hash is a 256-element uint8 vector. Set `WithGridSize(4, 4)` for a 4096-element spatially-aware hash.

#### HOG Hash (Histogram of Oriented Gradients)

Computes gradient magnitudes and orientations at each pixel, divides the image into square cells, and builds a magnitude-weighted orientation histogram per cell. The histograms are normalized and concatenated into a uint8 vector. Based on [Histograms of Oriented Gradients for Human Detection](https://ieeexplore.ieee.org/document/1467360) by Dalal and Triggs.

```go
hog, err := imghash.NewHOGHash()
hash, err := hog.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithCellSize(s)` | 8 |
| `WithNumBins(n)` | 9 |

With default settings the hash is a 9216-element uint8 vector (32×32 cells × 9 bins). Use `WithSize(32, 32)` for a compact 144-element hash (4×4 cells × 9 bins).

#### Radial Variance Hash

Uses radial projections and feature vector computations to generate a 40-element uint8 vector. From [Robust image hashing based on radial variance of pixels](https://www.researchgate.net/publication/4186555_Robust_image_hashing_based_on_radial_variance_of_pixels).

```go
rv, err := imghash.NewRadialVariance()
hash, err := rv.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithSigma(s)` | 1 |
| `WithAngles(n)` | 180 |

#### PDQ Hash

Produces a 256-bit binary hash designed for large-scale image deduplication. The algorithm resizes to 64×64, applies a Jarosz box filter for noise smoothing, computes a 2D DCT, and thresholds the 16×16 low-frequency coefficients against their median. Robust to JPEG compression, rescaling, and minor edits. From Facebook (Meta) [ThreatExchange PDQ](https://github.com/facebook/ThreatExchange/tree/main/pdq).

```go
pdq, err := imghash.NewPDQ()
hash, err := pdq.Calculate(img)
```

| Option | Default |
|--------|---------|
| `WithInterpolation(i)` | `Bilinear` |

The input size is fixed at 64×64 per the algorithm specification and is not configurable.

## Similarity Metrics

Every hash has a `Distance(other)` method that uses the natural metric for its type:

| Hash type | Default metric | Description |
|-----------|---------------|-------------|
| `Binary` | Hamming | Number of differing bits |
| `UInt8` | L2 (Euclidean) | Square root of sum of squared differences |
| `Float64` | L2 (Euclidean) | Square root of sum of squared differences |

For direct use:

```go
dist, err := h1.Distance(h2)
```

Or via the top-level convenience which returns a `Distance` value:

```go
dist, err := imghash.Compare(h1, h2)
```

For advanced use cases (e.g. Pearson Correlation), use the `similarity` sub-package:

```go
import "github.com/ajdnik/imghash/v2/similarity"

dist := similarity.L2(h1, h2)
dist, err := similarity.PCC(h1, h2)
dist, err := similarity.Hamming(h1, h2)
```

## Interpolation Methods

All resize-based algorithms accept `WithInterpolation`. Available methods:

`NearestNeighbor`, `Bilinear`, `Bicubic`, `MitchellNetravali`, `Lanczos2`, `Lanczos3`, `BilinearExact`

## License

Imghash is released under the MIT license. See [LICENSE](https://github.com/ajdnik/imghash/blob/main/LICENSE)
