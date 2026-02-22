# Similarity Metrics

## Algorithm-specific `Compare`

Every hash algorithm has a `Compare` method that uses the recommended distance metric for that algorithm.

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

Default metric per algorithm:

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
| Zernike | `Float64` | L2 (Euclidean) |
| CLD | `UInt8` | L2 (Euclidean) |
| EHD | `UInt8` | L1 (Manhattan) |
| LBP | `UInt8` | Chi-Square |
| HOGHash | `UInt8` | Cosine |
| RadialVariance | `UInt8` | L1 (Manhattan) |

## Generic `Compare`

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

## Available Metrics

The `similarity` sub-package provides all metrics for direct use:

```go
import "github.com/ajdnik/imghash/v2/similarity"

dist, err := similarity.Hamming(h1, h2)                 // bit-level Hamming distance (Binary only)
dist, err = similarity.WeightedHamming(h1, h2, weights) // weighted Hamming (Binary only)
dist, err = similarity.L1(h1, h2)                       // Manhattan distance
dist, err = similarity.L2(h1, h2)                       // Euclidean distance
dist, err = similarity.ChiSquare(h1, h2)                // Chi-square distance
dist, err = similarity.Cosine(h1, h2)                   // Cosine distance (1 - cos similarity)
dist, err = similarity.PCC(h1, h2)                      // Peak cross-correlation
```
