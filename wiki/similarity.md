# Similarity Functions

> Distance metrics for comparing perceptual hashes

## Overview

The similarity package provides various distance metrics for comparing perceptual hashes. Each function computes a similarity measure between two hashes, where lower values indicate more similar images.

All functions satisfy the `DistanceFunc` signature:

```go
type DistanceFunc func(hashtype.Hash, hashtype.Hash) (similarity.Distance, error)
```

## Distance Type

```go
type Distance float64
```

Represents a similarity measure as a float64 value. Lower values indicate more similar images.

- **`Equal`** `func(dst Distance) bool` — Checks if two distances are equal using epsilon comparison (1e-12 tolerance).

## Binary Hash Metrics

These metrics work exclusively with binary hashes.

### Hamming

Calculates the bit-level Hamming distance between two binary hashes.

```go
func Hamming(h1, h2 hashtype.Hash) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Must be of type `hashtype.Binary`.
- **`h2`** `hashtype.Hash` — Second hash. Must be of type `hashtype.Binary`.
- **`distance`** `Distance` — The number of differing bits between the two hashes.
- **`error`** `error` — Returns `ErrNotBinaryHash` if either hash is not binary.

#### Example

```go
import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

avg := imghash.Average{}
hash1, _ := avg.Calculate(img1)
hash2, _ := avg.Calculate(img2)

dist, err := similarity.Hamming(hash1, hash2)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Hamming distance: %v bits\n", dist)

// Typical interpretation:
// 0-5 bits: Nearly identical
// 5-10 bits: Very similar
// 10-20 bits: Similar
// 20+ bits: Different
```

### WeightedHamming

Calculates a weighted bit-level Hamming distance between two binary hashes.

```go
func WeightedHamming(h1, h2 hashtype.Hash, weights []float64) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Must be of type `hashtype.Binary`.
- **`h2`** `hashtype.Hash` — Second hash. Must be of type `hashtype.Binary`.
- **`weights`** `[]float64` — Weight for each byte position. Must have the same length as the shorter hash.
- **`distance`** `Distance` — The weighted sum of differing bits.
- **`error`** `error` — Returns `ErrNotBinaryHash` if either hash is not binary, or `ErrWeightLengthMismatch` if weights length doesn't match.

#### Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

// Create weights (e.g., higher weight for central bytes)
weights := []float64{0.5, 1.0, 2.0, 2.0, 2.0, 2.0, 1.0, 0.5}

hash1, _ := avg.Calculate(img1)
hash2, _ := avg.Calculate(img2)

dist, _ := similarity.WeightedHamming(hash1, hash2, weights)
fmt.Printf("Weighted Hamming distance: %v\n", dist)
```

## Vector Distance Metrics

These metrics work with all hash types by treating them as vectors.

### L1 (Manhattan Distance)

Calculates the L1 distance by summing absolute differences.

```go
func L1(h1, h2 hashtype.Hash) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Can be any hash type.
- **`h2`** `hashtype.Hash` — Second hash. Can be any hash type.
- **`distance`** `Distance` — Sum of absolute differences of corresponding elements.
- **`error`** `error` — Returns an error if computation fails.

#### Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

cm, _ := imghash.NewColorMoment()
hash1, _ := cm.Calculate(img1)
hash2, _ := cm.Calculate(img2)

dist, _ := similarity.L1(hash1, hash2)
fmt.Printf("L1 distance: %v\n", dist)
```

### L2 (Euclidean Distance)

Calculates the L2 (Euclidean) distance.

```go
func L2(h1, h2 hashtype.Hash) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Can be any hash type.
- **`h2`** `hashtype.Hash` — Second hash. Can be any hash type.
- **`distance`** `Distance` — Square root of sum of squared differences.
- **`error`** `error` — Returns an error if computation fails.

#### Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

gist, _ := imghash.NewGIST()
hash1, _ := gist.Calculate(img1)
hash2, _ := gist.Calculate(img2)

dist, _ := similarity.L2(hash1, hash2)
fmt.Printf("L2 distance: %v\n", dist)
```

### Cosine Distance

Calculates the cosine distance (1 - cosine similarity).

```go
func Cosine(h1, h2 hashtype.Hash) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Can be any hash type.
- **`h2`** `hashtype.Hash` — Second hash. Can be any hash type.
- **`distance`** `Distance` — 1 - (dot product / (magnitude1 \* magnitude2)). Returns 0 if both hashes are zero vectors.
- **`error`** `error` — Returns an error if computation fails.

#### Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

cld, _ := imghash.NewCLD()
hash1, _ := cld.Calculate(img1)
hash2, _ := cld.Calculate(img2)

dist, _ := similarity.Cosine(hash1, hash2)
fmt.Printf("Cosine distance: %v\n", dist)
// 0.0 = identical direction, 1.0 = orthogonal, 2.0 = opposite
```

## Statistical Metrics

### ChiSquare

Calculates the chi-square distance, useful for comparing histograms.

```go
func ChiSquare(h1, h2 hashtype.Hash) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Can be any hash type.
- **`h2`** `hashtype.Hash` — Second hash. Can be any hash type.
- **`distance`** `Distance` — Sum of (a - b)^2 / (a + b) for each element pair, skipping positions where both are zero.
- **`error`** `error` — Returns an error if computation fails.

#### Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

// Chi-square is particularly good for histogram-based hashes
ehd, _ := imghash.NewEHD()
hash1, _ := ehd.Calculate(img1)
hash2, _ := ehd.Calculate(img2)

dist, _ := similarity.ChiSquare(hash1, hash2)
fmt.Printf("Chi-square distance: %v\n", dist)
```

### PCC (Peak Cross-Correlation)

Calculates the peak cross-correlation between two hashes.

```go
func PCC(h1, h2 hashtype.Hash) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Can be any hash type.
- **`h2`** `hashtype.Hash` — Second hash. Must be same length as h1.
- **`distance`** `Distance` — Peak correlation value across all circular shifts.
- **`error`** `error` — Returns `ErrNotSameLength` if hash lengths don't match.

#### Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

// PCC is rotation-invariant
rv, _ := imghash.NewRadialVariance()
hash1, _ := rv.Calculate(img1)
hash2, _ := rv.Calculate(img2)

dist, _ := similarity.PCC(hash1, hash2)
fmt.Printf("Peak cross-correlation: %v\n", dist)
```

### Jaccard

Calculates the Jaccard distance.

```go
func Jaccard(h1, h2 hashtype.Hash) (Distance, error)
```

- **`h1`** `hashtype.Hash` — First hash. Can be Binary, UInt8, or Float64.
- **`h2`** `hashtype.Hash` — Second hash. Must be same type as h1.
- **`distance`** `Distance` — For Binary: 1 - (bitset intersection / bitset union)
  For UInt8/Float64: 1 - (matching positions / signature length)
- **`error`** `error` — Returns `ErrIncompatibleHash` if types don't match, or `ErrNotSameLength` for UInt8/Float64 if lengths differ.

#### Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

// Jaccard for binary hashes (set similarity)
avg := imghash.Average{}
hash1, _ := avg.Calculate(img1)
hash2, _ := avg.Calculate(img2)
dist1, _ := similarity.Jaccard(hash1, hash2)

// Jaccard for MinHash signatures
bovw, _ := imghash.NewBoVW(
    imghash.WithBoVWStorage(imghash.BoVWMinHash),
)
hash3, _ := bovw.Calculate(img1)
hash4, _ := bovw.Calculate(img2)
dist2, _ := similarity.Jaccard(hash3, hash4)

fmt.Printf("Jaccard distances: %v, %v\n", dist1, dist2)
```

## Choosing the Right Metric

### Binary Hashes (Average, Difference, PHash, etc.)

```go
// Best: Hamming distance (counts bit differences)
dist, _ := similarity.Hamming(hash1, hash2)

// Alternative: Jaccard (set similarity)
dist, _ := similarity.Jaccard(hash1, hash2)
```

### Histogram Hashes (CLD, EHD)

```go
// Best: Chi-square (designed for histograms)
dist, _ := similarity.ChiSquare(hash1, hash2)

// Alternative: L1 or L2
dist, _ := similarity.L1(hash1, hash2)
```

### Feature Vectors (ColorMoment, GIST)

```go
// Common: L2 (Euclidean)
dist, _ := similarity.L2(hash1, hash2)

// Alternative: Cosine (angle-based)
dist, _ := similarity.Cosine(hash1, hash2)
```

### Rotation-Invariant (RadialVariance)

```go
// Best: PCC (handles circular shifts)
dist, _ := similarity.PCC(hash1, hash2)
```

### MinHash/SimHash Signatures

```go
// Best: Jaccard (designed for signatures)
dist, _ := similarity.Jaccard(hash1, hash2)
```

## Comparison Example

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

func main() {
    img1, _ := imghash.OpenImage("photo1.jpg")
    img2, _ := imghash.OpenImage("photo2.jpg")
    
    avg := imghash.Average{}
    hash1, _ := avg.Calculate(img1)
    hash2, _ := avg.Calculate(img2)
    
    // Try multiple metrics
    metrics := map[string]imghash.DistanceFunc{
        "Hamming":   similarity.Hamming,
        "L1":        similarity.L1,
        "L2":        similarity.L2,
        "Cosine":    similarity.Cosine,
        "ChiSquare": similarity.ChiSquare,
        "Jaccard":   similarity.Jaccard,
    }
    
    for name, metric := range metrics {
        dist, err := metric(hash1, hash2)
        if err != nil {
            fmt.Printf("%s: error - %v\n", name, err)
        } else {
            fmt.Printf("%s: %v\n", name, dist)
        }
    }
}
```

## Error Types

```go
// Hash types incompatible for the chosen metric
var ErrIncompatibleHash = hashtype.ErrIncompatibleHash

// Hashes must be same length (for PCC, Jaccard with UInt8/Float64)
var ErrNotSameLength = errors.New("hashes aren't the same length")

// Weights slice length doesn't match hash length (WeightedHamming)
var ErrWeightLengthMismatch = errors.New("weight slice length must match number of hash bytes")

// Non-binary hash passed to Hamming or WeightedHamming
var ErrNotBinaryHash = hashtype.ErrIncompatibleHash
```

## See Also

* [Comparer Interface](comparer) - Algorithm-specific comparison
* [Convenience Functions](convenience) - Helper functions including Compare
* [Hash Types](hash-types) - Hash type definitions
