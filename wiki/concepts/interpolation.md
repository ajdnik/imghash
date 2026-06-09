# Interpolation Methods

> Understanding image resizing interpolation algorithms used during hash computation

## Overview

Most perceptual hashing algorithms resize input images to a fixed dimension before computing the hash. The interpolation method used during this resize operation can significantly impact hash quality and computation speed.

imghash provides 7 interpolation algorithms, each with different trade-offs between quality and performance.

```go
type Interpolation int
```

## Available Methods

### NearestNeighbor

Selects the pixel value from the nearest source pixel without any blending.

```go
const NearestNeighbor Interpolation = 0
```

- **Speed** — **Fastest** - No computation, direct pixel lookup
- **Quality** — **Lowest** - Can produce blocky, aliased results

#### Example

```go
algo, _ := imghash.NewAverage(
    imghash.WithInterpolation(imghash.NearestNeighbor),
)

hash, _ := algo.Calculate(img)
```

#### Characteristics

* No anti-aliasing
* Sharp edges preserved (can be jagged)
* Best for: pixel art, binary images, maximum speed
* Avoid for: natural images with gradients

> [!WARNING]
> NearestNeighbor can produce inconsistent hashes when image dimensions change slightly, as different source pixels may be selected.

### Bilinear

**Default for most algorithms.** Blends the 4 nearest pixels using linear interpolation.

```go
const Bilinear Interpolation = 1
```

- **Speed** — **Fast** - Simple 4-pixel weighted average
- **Quality** — **Good** - Smooth results for most use cases

#### Example

```go
// Bilinear is the default
algo1, _ := imghash.NewPHash()

// Explicit specification
algo2, _ := imghash.NewPHash(
    imghash.WithInterpolation(imghash.Bilinear),
)

// Both are equivalent
```

#### Algorithm

```
For each output pixel at (x, y):
  Find 4 surrounding source pixels
  Compute weights based on distance
  result = w1*p1 + w2*p2 + w3*p3 + w4*p4
```

#### Characteristics

* Smooth gradients
* Minimal blur
* Good balance of speed and quality
* **Recommended for:** Average, Difference, Median, PHash, most algorithms

> [!TIP]
> Bilinear is the sweet spot for most applications - good quality without significant overhead.

### BilinearExact

More precise bilinear interpolation using exact floating-point math.

```go
const BilinearExact Interpolation = 6
```

- **Speed** — **Slightly slower** than standard Bilinear
- **Quality** — **More accurate** floating-point calculations

#### Example

```go
algo, _ := imghash.NewColorMoment(
    imghash.WithInterpolation(imghash.BilinearExact),
)
```

#### Use When

* You need precise color values (ColorMoment, CLD)
* Reproducibility across platforms is critical
* Computational cost is not a concern

### Bicubic

Uses 16 nearest pixels with cubic polynomial weighting for smoother results.

```go
const Bicubic Interpolation = 2
```

- **Speed** — **Slower** - 16 pixel computations per output pixel
- **Quality** — **Higher** - Smoother than bilinear, less aliasing

#### Example

```go
algo, _ := imghash.NewZernike(
    imghash.WithInterpolation(imghash.Bicubic),
)

hash, _ := algo.Calculate(highResImage)
```

#### Algorithm

```
For each output pixel at (x, y):
  Find 16 surrounding source pixels (4x4 grid)
  Apply cubic convolution kernel
  Weighted sum with negative lobes for sharpness
```

#### Characteristics

* Smoother gradients than bilinear
* Can introduce slight ringing near edges
* Better preservation of high-frequency details
* **Recommended for:** GIST, Zernike, high-quality requirements

> [!NOTE]
> Bicubic can produce slightly sharper results than bilinear, which may improve discriminability for some algorithms.

### MitchellNetravali

Advanced bicubic filter with balanced sharpness and smoothness.

```go
const MitchellNetravali Interpolation = 3
```

- **Speed** — **Similar to Bicubic** - 16 pixel kernel
- **Quality** — **Excellent** - Optimized trade-off of blur and ringing

#### Example

```go
algo, _ := imghash.NewGIST(
    imghash.WithInterpolation(imghash.MitchellNetravali),
)
```

#### Characteristics

* Uses Mitchell-Netravali filter with B=1/3, C=1/3
* Perceptually balanced - minimal blur and ringing
* Preferred by many image processing experts
* **Recommended for:** Production systems requiring high quality

> [!TIP]
> MitchellNetravali is often considered the "best looking" interpolation for natural images.

### Lanczos2

Sinc-based filter with 2-lobe kernel for high-quality resampling.

```go
const Lanczos2 Interpolation = 4
```

- **Speed** — **Slow** - 4x4 kernel with sinc calculations
- **Quality** — **Very High** - Excellent detail preservation

#### Example

```go
algo, _ := imghash.NewPDQ(
    imghash.WithInterpolation(imghash.Lanczos2),
)

hash, _ := algo.Calculate(photographImage)
```

#### Algorithm

```
Lanczos kernel:
  L(x) = sinc(x) * sinc(x/a)  where a=2
  
For each output pixel:
  Apply 4x4 windowed sinc kernel
  Normalized weighted sum
```

#### Characteristics

* Based on sinc function (ideal low-pass filter)
* Sharp transitions, minimal blur
* Can introduce ringing artifacts near edges
* **Recommended for:** Downsampling photographs, detail-critical applications

### Lanczos3

Sinc-based filter with 3-lobe kernel - highest quality, slowest speed.

```go
const Lanczos3 Interpolation = 5
```

- **Speed** — **Slowest** - 6x6 kernel with sinc calculations
- **Quality** — **Highest** - Maximum detail retention

#### Example

```go
// For maximum quality when speed is not critical
algo, _ := imghash.NewHOGHash(
    imghash.WithInterpolation(imghash.Lanczos3),
    imghash.WithHOGCellSize(8),
)
```

#### Algorithm

```
Lanczos kernel:
  L(x) = sinc(x) * sinc(x/a)  where a=3
  
For each output pixel:
  Apply 6x6 windowed sinc kernel
  Normalized weighted sum
```

#### Characteristics

* Widest kernel - most source pixels considered
* Sharpest results, best frequency preservation
* Most prone to ringing artifacts
* **Recommended for:** Critical applications, large downsample ratios, research

> [!WARNING]
> Lanczos3 is computationally expensive. Only use when quality justifies the cost.

## Setting Interpolation

All algorithms that perform resizing accept the `WithInterpolation` option:

```go
import "github.com/ajdnik/imghash/v2"

// Average hash with nearest neighbor (fastest)
avg, _ := imghash.NewAverage(
    imghash.WithInterpolation(imghash.NearestNeighbor),
)

// PHash with Lanczos3 (highest quality)
phash, _ := imghash.NewPHash(
    imghash.WithInterpolation(imghash.Lanczos3),
)

// GIST with default bilinear
gist, _ := imghash.NewGIST()
// No WithInterpolation specified = Bilinear
```

### Algorithms Using Interpolation

These algorithms resize images and support `WithInterpolation`:

* Average, Difference, Median, PHash
* BlockMean, MarrHildreth
* ColorMoment, CLD, EHD
* LBP, HOGHash
* RadialVariance, Zernike
* GIST, WHash

> [!NOTE]
> Some algorithms (PDQ, BoVW, RASH) use fixed internal resizing and don't expose interpolation as an option.

## Performance Comparison

| Method                | Relative Speed  | Quality | Kernel Size | Use Case                  |
| --------------------- | --------------- | ------- | ----------- | ------------------------- |
| **NearestNeighbor**   | 1.0x (baseline) | ⭐       | 1x1         | Speed-critical, pixel art |
| **Bilinear**          | 0.8x            | ⭐⭐⭐     | 2x2         | **Default - balanced**    |
| **BilinearExact**     | 0.75x           | ⭐⭐⭐     | 2x2         | Precise color values      |
| **Bicubic**           | 0.4x            | ⭐⭐⭐⭐    | 4x4         | Natural images            |
| **MitchellNetravali** | 0.4x            | ⭐⭐⭐⭐    | 4x4         | Production quality        |
| **Lanczos2**          | 0.3x            | ⭐⭐⭐⭐⭐   | 4x4         | Photographs               |
| **Lanczos3**          | 0.2x            | ⭐⭐⭐⭐⭐   | 6x6         | Maximum quality           |

> [!NOTE]
> Speed measurements are approximate and depend on image size and CPU architecture.

## Quality vs Speed Trade-offs

#### Speed Priority

```go
// Real-time applications, large batch processing
algo, _ := imghash.NewAverage(
    imghash.WithInterpolation(imghash.NearestNeighbor),
)

// or
algo, _ := imghash.NewAverage(
    imghash.WithInterpolation(imghash.Bilinear),
)
```

**Use when:**

* Processing millions of images
* Real-time video hashing
* Mobile/embedded devices

#### Balanced

```go
// Production systems - good quality, reasonable speed
algo, _ := imghash.NewPHash(
    imghash.WithInterpolation(imghash.Bilinear),
)

// or for better quality
algo, _ := imghash.NewPHash(
    imghash.WithInterpolation(imghash.Bicubic),
)
```

**Use when:**

* General-purpose image similarity
* Content moderation pipelines
* Typical web applications

#### Quality Priority

```go
// Research, forensics, critical applications
algo, _ := imghash.NewZernike(
    imghash.WithInterpolation(imghash.MitchellNetravali),
)

// or maximum quality
algo, _ := imghash.NewGIST(
    imghash.WithInterpolation(imghash.Lanczos3),
)
```

**Use when:**

* Copyright detection
* Medical imaging
* Research experiments
* Quality is critical

## Interpolation Impact on Hashes

Different interpolation methods produce slightly different hashes from the same image:

```go
img := loadImage("photo.jpg")

algoNN, _ := imghash.NewAverage(
    imghash.WithInterpolation(imghash.NearestNeighbor),
)
algoLanczos, _ := imghash.NewAverage(
    imghash.WithInterpolation(imghash.Lanczos3),
)

hashNN, _ := algoNN.Calculate(img)
hashLanczos, _ := algoLanczos.Calculate(img)

// Hashes will differ slightly
dist, _ := similarity.Hamming(hashNN, hashLanczos)
fmt.Printf("Distance between interpolations: %.0f bits\n", dist)
// Typically 1-5 bits difference for Average hash
```

> [!WARNING]
> If you need to compare hashes across systems, ensure all systems use the same interpolation method.

## Best Practices

### Step 1: Use Bilinear as Default

Unless you have specific requirements, Bilinear provides excellent quality with minimal overhead.

### Step 2: Match Interpolation to Algorithm

Simple algorithms (Average, Difference) don't benefit much from high-quality interpolation. Save Lanczos for complex algorithms like GIST or Zernike.

### Step 3: Consider Downsample Ratio

Larger downsample ratios (e.g., 4K → 8x8) benefit more from high-quality interpolation than small ratios.

### Step 4: Profile Your Application

Measure actual performance impact. Interpolation cost is often small compared to algorithm computation.

### Step 5: Stay Consistent

Use the same interpolation method throughout your application to ensure comparable hashes.

## Recommendations by Algorithm

<details>
<summary>Simple Bit-based (Average, Difference, Median)</summary>

**Recommended:** Bilinear (default)

These algorithms threshold pixels to binary values, so subtle interpolation differences matter less.

```go
algo, _ := imghash.NewAverage(
    imghash.WithInterpolation(imghash.Bilinear),
)
```

</details>

<details>
<summary>DCT-based (PHash, PDQ)</summary>

**Recommended:** Bicubic or MitchellNetravali

DCT benefits from smooth frequency domain representation.

```go
algo, _ := imghash.NewPHash(
    imghash.WithInterpolation(imghash.Bicubic),
)
```

</details>

<details>
<summary>Color/Histogram (ColorMoment, CLD, EHD)</summary>

**Recommended:** BilinearExact or Bicubic

Color precision matters for these algorithms.

```go
algo, _ := imghash.NewColorMoment(
    imghash.WithInterpolation(imghash.BilinearExact),
)
```

</details>

<details>
<summary>Gradient-based (HOGHash, LBP)</summary>

**Recommended:** Bicubic or Lanczos2

Edge preservation is important for gradient calculation.

```go
algo, _ := imghash.NewHOGHash(
    imghash.WithInterpolation(imghash.Lanczos2),
)
```

</details>

<details>
<summary>Advanced Features (GIST, Zernike, RadialVariance)</summary>

**Recommended:** Lanczos2 or Lanczos3

These algorithms analyze fine details and benefit from high-quality resampling.

```go
algo, _ := imghash.NewGIST(
    imghash.WithInterpolation(imghash.Lanczos3),
)
```

</details>

## Error Handling

```go
import "errors"

// Invalid interpolation value
algo, err := imghash.NewAverage(
    imghash.WithInterpolation(Interpolation(999)),
)

if errors.Is(err, imghash.ErrInvalidInterpolation) {
    fmt.Println("Invalid interpolation method")
}
```

Valid interpolation values are 0-6 (NearestNeighbor through BilinearExact).

## String Representation

```go
interp := imghash.Bilinear
fmt.Println(interp.String()) // "Bilinear"

invalid := imghash.Interpolation(999)
fmt.Println(invalid.String()) // "Unknown"
```

## Advanced: Custom Interpolation

For custom resize logic, you can pre-process images before hashing:

```go
import "image"

// Custom resize function
func customResize(img image.Image, width, height int) image.Image {
    // Your custom interpolation logic
    return resizedImage
}

// Pre-resize, then use NearestNeighbor (no-op)
resized := customResize(img, 8, 8)

algo, _ := imghash.NewAverage(
    imghash.WithSize(8, 8),
    imghash.WithInterpolation(imghash.NearestNeighbor),
)

hash, _ := algo.Calculate(resized)
```

> [!NOTE]
> This approach gives full control but requires more code and understanding of the algorithm's expected input size.

## Related

* [Hash Types](concepts/hash-types) - Understanding hash representations
* [Choosing an Algorithm](guides/choosing-algorithm) - Algorithm selection and configuration
