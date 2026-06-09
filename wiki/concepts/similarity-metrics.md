# Similarity Metrics

> Complete guide to the 9 distance and similarity functions for comparing perceptual hashes

## Overview

imghash provides 9 distinct similarity metrics to measure how alike two image hashes are. Each metric is optimized for specific hash types and use cases. Lower distance values indicate more similar images.

```go
type Distance float64

type DistanceFunc func(hashtype.Hash, hashtype.Hash) (Distance, error)
```

All metrics are exported from the `similarity` package and can be used with the `WithDistance` option.

## Hamming Distance

Counts the number of differing bits between two binary hashes. The classic metric for bit-based perceptual hashes.

### Signature

```go
func Hamming(h1, h2 hashtype.Hash) (Distance, error)
```

### Requirements

> [!WARNING]
> **Binary hashes only.** Returns `ErrIncompatibleHash` for UInt8 or Float64 types.

### Algorithm

```
For each byte position:
  XOR the bytes
  Count the number of 1 bits
Return sum of all bit differences
```

### Example

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

algo, _ := imghash.NewAverage()
hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(img2)

dist, _ := similarity.Hamming(hash1, hash2)
fmt.Printf("Hamming distance: %.0f bits\n", dist)

// Output: Hamming distance: 5 bits
```

### Use Cases

* **Average, Difference, Median hashes** - Fast near-duplicate detection
* **PHash** - Robust image similarity
* **Large-scale search** - Optimized for speed with bit operations
* **DINOHash** - 96-bit semantic perceptual hash with adversarial robustness

> [!TIP]
> Hamming distance of 0-10 typically indicates very similar images for 64-bit hashes.

## Weighted Hamming

Like Hamming distance but applies per-byte weights to emphasize important hash regions.

### Signature

```go
func WeightedHamming(h1, h2 hashtype.Hash, weights []float64) (Distance, error)
```

### Requirements

* **Binary hashes only**
* Weights slice length must match hash byte length

### Algorithm

```
For each byte position i:
  XOR the bytes
  Count bits
  Multiply by weights[i]
Return weighted sum
```

### Example

```go
// Emphasize center of 64-bit (8-byte) hash
weights := []float64{0.5, 0.8, 1.2, 1.5, 1.5, 1.2, 0.8, 0.5}

algo, _ := imghash.NewPHash()
hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(img2)

// Can't use as algorithm default, must call directly
b1 := hash1.(imghash.Binary)
b2 := hash2.(imghash.Binary)
dist, _ := similarity.WeightedHamming(b1, b2, weights)

fmt.Printf("Weighted distance: %.2f\n", dist)
```

### Use Cases

* **PDQ hashes** - Weight perceptually important DCT coefficients
* **Custom applications** - When certain hash regions are more discriminative
* **Object detection** - Weight center pixels higher than edges

> [!NOTE]
> `WeightedHamming` cannot be set as the default distance for an algorithm via `WithDistance`. You must call it directly.

## L1 Distance (Manhattan)

Sum of absolute differences between corresponding elements. Works with all hash types.

### Signature

```go
func L1(h1, h2 hashtype.Hash) (Distance, error)
```

### Algorithm

```
sum = 0
for each element i:
  sum += |h1[i] - h2[i]|
return sum
```

### Example

```go
algo, _ := imghash.NewColorMoment()
hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(img2)

// ColorMoment uses L1 by default
dist, _ := algo.Compare(hash1, hash2)

// Or call directly
dist, _ = similarity.L1(hash1, hash2)
fmt.Printf("L1 distance: %.2f\n", dist)
```

### Characteristics

- **Hash Types** — Binary, UInt8, Float64 - works with all types
- **Complexity** — O(n) where n is hash length - very fast
- **Sensitivity** — Linear - proportional to element differences
- **Range** — \[0, ∞) - unbounded maximum

### Use Cases

* **ColorMoment** - Color distribution comparison
* **EHD** - Edge histogram matching
* **When outliers matter** - Doesn't square differences like L2

## L2 Distance (Euclidean)

Square root of sum of squared differences. The straight-line distance in n-dimensional space.

### Signature

```go
func L2(h1, h2 hashtype.Hash) (Distance, error)
```

### Algorithm

```
sum = 0
for each element i:
  diff = h1[i] - h2[i]
  sum += diff * diff
return sqrt(sum)
```

### Example

```go
algo, _ := imghash.NewGIST()
hash1, _ := algo.Calculate(img1)
hash2, _ := ago.Calculate(img2)

dist, _ := similarity.L2(hash1, hash2)
fmt.Printf("L2 distance: %.4f\n", dist)
```

### Characteristics

* **Penalizes large differences** - Squaring amplifies outliers
* **Geometric interpretation** - True Euclidean distance in feature space
* **Smooth gradients** - Better for optimization tasks

### Use Cases

* **GIST** - Global scene descriptor comparison
* **Zernike** - Shape matching with moment invariants
* **Feature vectors** - When semantic distance matters

> [!TIP]
> L2 is more sensitive to large differences than L1, making it better for detecting significant variations.

## Cosine Distance

Measures the angle between two vectors, independent of magnitude. Returns 1 - cosine similarity.

### Signature

```go
func Cosine(h1, h2 hashtype.Hash) (Distance, error)
```

### Algorithm

```
dot = h1 · h2  (dot product)
mag1 = ||h1||  (magnitude)
mag2 = ||h2||

if mag1 * mag2 == 0:
  return 0

cosine_similarity = dot / (mag1 * mag2)
return 1 - cosine_similarity
```

### Example

```go
algo, _ := imghash.NewHOGHash()
hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(img2)

// Use custom distance metric
algo2, _ := imghash.NewHOGHash(
    imghash.WithDistance(similarity.Cosine),
)
dist, _ := algo2.Compare(hash1, hash2)

fmt.Printf("Cosine distance: %.4f\n", dist)
```

### Characteristics

* **Range:** \[0, 2] where 0 = identical direction, 2 = opposite direction
* **Magnitude independent** - Only measures angle, not scale
* **Normalized** - Good for vectors with different scales

### Use Cases

* **High-dimensional features** - When magnitude is less important than direction
* **Normalized histograms** - Comparing distributions
* **Text-like features** - Similar to TF-IDF comparison

> [!NOTE]
> Cosine returns 0 when both hashes are zero vectors, avoiding division by zero.

## Chi-Square Distance

Statistical measure comparing probability distributions. Common for histogram comparison.

### Signature

```go
func ChiSquare(h1, h2 hashtype.Hash) (Distance, error)
```

### Algorithm

```
sum = 0
for each element i:
  a = h1[i]
  b = h2[i]
  if a + b == 0:
    continue  // skip to avoid division by zero
  sum += (a - b)² / (a + b)
return sum
```

### Example

```go
algo, _ := imghash.NewCLD()  // Color Layout Descriptor
hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(img2)

// CLD uses ChiSquare by default
dist, _ := algo.Compare(hash1, hash2)

fmt.Printf("Chi-Square distance: %.4f\n", dist)
```

### Characteristics

* **Histogram optimized** - Designed for comparing distributions
* **Non-negative** - Values must be >= 0
* **Asymmetric penalty** - Differences in small bins matter more

### Use Cases

* **CLD** - Color layout histograms
* **LBP** - Texture pattern distributions
* **Any histogram-based feature** - Standard choice for distribution comparison

> [!WARNING]
> Chi-Square automatically skips positions where both values are zero to prevent division by zero.

## PCC (Peak Cross-Correlation)

Finds maximum correlation across circular rotations. Useful for rotation-invariant matching.

### Signature

```go
func PCC(h1, h2 hashtype.Hash) (Distance, error)
```

### Requirements

> [!WARNING]
> **Equal length hashes only.** Returns `ErrNotSameLength` if hash lengths differ.

### Algorithm

```
1. Normalize both hashes (subtract mean)
2. For each rotation of h2:
   - Compute correlation = covariance / (std1 * std2)
   - Track maximum correlation
3. Return peak correlation value
```

### Example

```go
algo, _ := imghash.NewRadialVariance()
hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(rotatedImg)

// RadialVariance uses PCC by default for rotation invariance
dist, _ := algo.Compare(hash1, hash2)

fmt.Printf("PCC: %.4f\n", dist)
```

### Characteristics

* **Rotation invariant** - Tests all circular shifts
* **Computationally expensive** - O(n²) where n is hash length
* **Returns correlation** - Higher values mean more similar
* **Range:** (-∞, 1] where 1 = perfect match

### Use Cases

* **RadialVariance** - Radial projection features
* **Rotated images** - When rotation is expected
* **Circular features** - Angular histograms, polar transforms

> [!TIP]
> PCC is unique in that higher values indicate more similar images, unlike other distance metrics.

## Jaccard Distance

Set-based similarity with three modes depending on hash type.

### Signature

```go
func Jaccard(h1, h2 hashtype.Hash) (Distance, error)
```

### Binary Mode (Bitset)

Compares set bits using intersection over union.

```
intersection = popcount(h1 & h2)
union = popcount(h1 | h2)

if union == 0:
  return 0
return 1 - (intersection / union)
```

### UInt8/Float64 Mode (MinHash)

Treats hashes as MinHash signatures, counting matching positions.

```
matches = count of positions where h1[i] == h2[i]
return 1 - (matches / hash_length)
```

### Example

#### Binary (Bitset)

```go
// Use with any binary hash algorithm
algo, _ := imghash.NewAverage()
hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(img2)

dist, _ := similarity.Jaccard(hash1, hash2)
fmt.Printf("Jaccard distance: %.4f\n", dist)
// Range: [0, 1] where 0 = identical sets
```

#### UInt8 (MinHash)

```go
// Use with BoVW configured for MinHash
algo, _ := imghash.NewBoVW(
    imghash.WithBoVWStorage(imghash.SimHashUInt8),
    imghash.WithBoVWVocabularySize(256),
)

hash1, _ := algo.Calculate(img1)
hash2, _ := algo.Calculate(img2)

// Compare using Jaccard
dist, _ := similarity.Jaccard(hash1, hash2)
fmt.Printf("MinHash distance: %.4f\n", dist)
```

### Characteristics

- **Binary Mode** — Bitset intersection/union - standard Jaccard index
- **MinHash Mode** — Signature matching - approximates set similarity
- **Range** — \[0, 1] where 0 = identical, 1 = completely different
- **Symmetric** — Jaccard(A,B) = Jaccard(B,A)

### Use Cases

* **Set-based features** - When images are represented as feature sets
* **BoVW with MinHash** - Bag of Visual Words similarity
* **Sparse binary features** - When most bits are 0

## Using Custom Metrics

Override an algorithm's default distance function using `WithDistance`:

```go
import "github.com/ajdnik/imghash/v2/similarity"

// Average normally uses Hamming
algo1, _ := imghash.NewAverage()

// Override to use Jaccard instead
algo2, _ := imghash.NewAverage(
    imghash.WithDistance(similarity.Jaccard),
)

hash1, _ := algo1.Calculate(img1)
hash2, _ := algo1.Calculate(img2)

// Default Hamming distance
dist1, _ := algo1.Compare(hash1, hash2)

// Custom Jaccard distance
dist2, _ := algo2.Compare(hash1, hash2)

fmt.Printf("Hamming: %.0f, Jaccard: %.4f\n", dist1, dist2)
```

## Metric Selection Guide

<details>
<summary>Binary Hashes (Average, Difference, PHash)</summary>

**Primary:** `Hamming` - Fast, standard choice

**Alternative:** `Jaccard` - Better for sparse binary vectors

**Advanced:** `WeightedHamming` - When spatial importance varies

</details>

<details>
<summary>Histogram Features (ColorMoment, CLD, EHD)</summary>

**Primary:** `ChiSquare` - Statistical standard for histograms

**Alternative:** `L1` - Simpler, faster

**Alternative:** `L2` - Penalizes large differences more

</details>

<details>
<summary>High-Dimensional Vectors (GIST, Zernike, HOGHash)</summary>

**Primary:** `L2` - True geometric distance

**Alternative:** `Cosine` - Magnitude-independent

**Alternative:** `L1` - Less sensitive to outliers

</details>

<details>
<summary>Rotation-Variant Features (RadialVariance)</summary>

**Primary:** `PCC` - Handles circular rotations

**Note:** Much slower but rotation-invariant

</details>

## Distance Comparison

| Metric              | Hash Types     | Complexity | Range   | Best For                   |
| ------------------- | -------------- | ---------- | ------- | -------------------------- |
| **Hamming**         | Binary         | O(n)       | \[0, ∞) | Bit-based hashes, speed    |
| **WeightedHamming** | Binary         | O(n)       | \[0, ∞) | Spatial importance         |
| **L1**              | All            | O(n)       | \[0, ∞) | Histograms, robustness     |
| **L2**              | All            | O(n)       | \[0, ∞) | Feature vectors, geometry  |
| **Cosine**          | All            | O(n)       | \[0, 2] | Direction, normalized data |
| **ChiSquare**       | All            | O(n)       | \[0, ∞) | Probability distributions  |
| **PCC**             | All (same len) | O(n²)      | (-∞, 1] | Rotation invariance        |
| **Jaccard**         | All            | O(n)       | \[0, 1] | Sets, MinHash              |

## Error Handling

```go
import "errors"

// Type incompatibility
binary := hashtype.Binary{1, 2, 3}
float := hashtype.Float64{1.0, 2.0, 3.0}

_, err := similarity.Hamming(binary, float)
if errors.Is(err, similarity.ErrNotBinaryHash) {
    fmt.Println("Hamming requires binary hashes")
}

// Length mismatch for PCC
h1 := hashtype.Float64{1, 2, 3}
h2 := hashtype.Float64{1, 2, 3, 4}

_, err = similarity.PCC(h1, h2)
if errors.Is(err, similarity.ErrNotSameLength) {
    fmt.Println("PCC requires equal-length hashes")
}

// WeightedHamming weight mismatch
b1 := hashtype.Binary{1, 2, 3}
b2 := hashtype.Binary{1, 2, 3}
weights := []float64{1.0, 1.0} // Wrong length

_, err = similarity.WeightedHamming(b1, b2, weights)
if errors.Is(err, similarity.ErrWeightLengthMismatch) {
    fmt.Println("Weights must match hash byte length")
}
```

## Performance Tips

### Step 1: Use Hamming for Binary

It's highly optimized with bit operations - fastest option for Binary hashes.

### Step 2: Avoid PCC Unless Necessary

PCC is O(n²) - only use when rotation invariance is required.

### Step 3: Prefer L1 Over L2

L1 avoids the sqrt operation, making it slightly faster with similar results.

### Step 4: Cache Hash Magnitudes

If comparing one hash against many, pre-compute its magnitude for Cosine distance.

## Related

* [Hash Types](concepts/hash-types) - Understand Binary, UInt8, and Float64 representations
* [Choosing an Algorithm](guides/choosing-algorithm) - Choose the right algorithm and metric combination
* [Similarity API](api/similarity) - Complete similarity package documentation
