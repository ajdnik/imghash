# Algorithms

imghash supports 17 perceptual hashing algorithms. Most are ported from [OpenCV Contrib](https://github.com/opencv/opencv_contrib) and tested against its implementations.

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

## Average Hash

Produces a binary hash. Compares using Hamming distance. Based on [Looks Like It](https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

## Difference Hash

Computes gradients instead of averages. Compares using Hamming distance. Based on [Kind of Like That](https://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

## Median Hash

Like average hash but uses the median threshold. Compares using Hamming distance.

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |

## PHash

Uses DCT and produces a binary hash. Compares using weighted Hamming distance. From [Rihamark: Perceptual image hash benchmarking](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 32, 32 |
| `WithInterpolation(i)` | `BilinearExact` |
| `WithWeights(w)` | `[1, 1, 1, 1, 1, 1, 1, 1]` |

## Wavelet Hash (WHash)

Uses a multi-level Haar DWT and median thresholding. Compares using Hamming distance.

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 8, 8 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithLevel(l)` | 3 |

## Color Moments Hash

Builds a 42-element `float64` vector from YCbCr and HSV moments. Compares using L2 distance. Based on [Perceptual hashing for color images using invariant moments](https://www.researchgate.net/publication/286870507_Perceptual_hashing_for_color_images_using_invariant_moments).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 512, 512 |
| `WithInterpolation(i)` | `Bicubic` |
| `WithKernelSize(k)` | 3 |
| `WithSigma(s)` | 0 |

## MPEG-7 Color Layout Descriptor (CLD)

Builds a compact 12-element `uint8` descriptor in YCbCr. Compares using L2 distance.

Whitepaper: [Color and Texture Descriptors (Manjunath et al., 2001)](https://doi.org/10.1109/76.927424).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 64, 64 |
| `WithInterpolation(i)` | `Bilinear` |

## MPEG-7 Edge Histogram Descriptor (EHD)

Builds an 80-element `uint8` descriptor from local edge categories. Compares using L1 distance.

Whitepaper: [MPEG-7 Texture Descriptors (Wu et al., 2001)](https://doi.org/10.1142/S0219467801000311).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |

## Marr-Hildreth Hash

Uses a 2D wavelet transform and produces a binary hash. Compares using Hamming distance.

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 512, 512 |
| `WithInterpolation(i)` | `Bicubic` |
| `WithScale(s)` | 1 |
| `WithAlpha(a)` | 2 |
| `WithKernelSize(k)` | 7 |
| `WithSigma(s)` | 0 |

## Block Mean Value Hash

Computes block means and produces a binary hash. Compares using Hamming distance. Based on [Block Mean Value Based Image Perceptual Hashing](https://ieeexplore.ieee.org/document/4041692).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `BilinearExact` |
| `WithBlockSize(w, h)` | 16, 16 |
| `WithBlockMeanMethod(m)` | `Direct` |

Block mean methods: `Direct`, `Overlap`, `Rotation`, `RotationOverlap`.

`Rotation` and `RotationOverlap` compute and concatenate hashes for 24 rotations (0 to 345 degrees in 15-degree steps), so the result is 24x larger than non-rotational mode.

## Local Binary Pattern (LBP) Hash

Computes LBP codes and builds normalized histograms into a `uint8` vector. Compares using chi-square distance.

Based on [Multiresolution Gray-Scale and Rotation Invariant Texture Classification with Local Binary Patterns](https://ieeexplore.ieee.org/document/1017623) by Ojala et al.

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithGridSize(x, y)` | 1, 1 |

With default `1x1` grid, output is 256 elements. `WithGridSize(4, 4)` yields a 4096-element spatially-aware hash.

## HOG Hash (Histogram of Oriented Gradients)

Computes gradient-orientation histograms per cell into a `uint8` vector. Compares using cosine distance.

Based on [Histograms of Oriented Gradients for Human Detection](https://ieeexplore.ieee.org/document/1467360) by Dalal and Triggs.

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithCellSize(s)` | 8 |
| `WithNumBins(n)` | 9 |

With defaults, output is a 9216-element vector (32x32 cells x 9 bins). `WithSize(32, 32)` produces a 144-element hash (4x4 cells x 9 bins).

## GIST Descriptor Hash

Computes an Oliva-Torralba style holistic descriptor by applying an oriented Gabor filter bank and pooling responses over a spatial grid. Compares using cosine distance.

Based on [Modeling the Shape of the Scene: A Holistic Representation of the Spatial Envelope](https://people.csail.mit.edu/torralba/code/spatialenvelope/).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 64, 64 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithGridSize(x, y)` | 4, 4 |

With defaults, output is a 320-element descriptor (`4x4` cells x `20` filter channels). `WithGridSize(2, 2)` produces an 80-element hash.

## Radial Variance Hash

Uses radial projections to produce a 40-element `uint8` vector. Compares using L1 distance.

From [Robust image hashing based on radial variance of pixels](https://www.researchgate.net/publication/4186555_Robust_image_hashing_based_on_radial_variance_of_pixels).

| Option | Default |
|--------|---------|
| `WithSigma(s)` | 1 |
| `WithAngles(n)` | 180 |

## Zernike Moments Hash

Computes magnitude-only Zernike moments into a `float64` vector. Compares using L2 distance.

Magnitude makes the descriptor rotation-invariant while preserving low-order shape structure.

Whitepaper: [Invariant Image Recognition by Zernike Moments (Khotanzad and Hong, 1990)](https://doi.org/10.1109/34.55109).

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 64, 64 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithDegree(d)` | 8 |

## RASH (Rotation Aware Spatial Hash)

Produces a rotation-robust binary hash and compares with Hamming distance.

This custom algorithm combines concentric ring sampling, 1-D DCT compaction, and median thresholding. It is inspired by [Robust Image Hashing with Ring Partition and Invariant Vector Distance](https://ieeexplore.ieee.org/document/7368930) and radial variance hashing literature.

| Option | Default |
|--------|---------|
| `WithSize(w, h)` | 256, 256 |
| `WithInterpolation(i)` | `Bilinear` |
| `WithSigma(s)` | 1 |
| `WithRings(n)` | 180 |

## PDQ Hash

Produces a 256-bit binary hash for large-scale image deduplication. Compares using Hamming distance.

From Meta [ThreatExchange PDQ](https://github.com/facebook/ThreatExchange/tree/main/pdq).

| Option | Default |
|--------|---------|
| `WithInterpolation(i)` | `Bilinear` |

The input size is fixed at 64x64 per the algorithm specification.

## Binary Hash Size with Custom Options

For binary hashers with configurable dimensions, bit count may not be a multiple of 8. In that case:

- `Binary` stores `ceil(bits/8)` bytes
- unused bit positions in the last byte remain zero

Examples:

- `WithSize(3, 3)` for Average/Median/Difference/WHash yields 9 bits stored in 2 bytes
- `RASH` uses `min(64, rings-1)` bits, so `WithRings(5)` yields 4 bits stored in 1 byte
