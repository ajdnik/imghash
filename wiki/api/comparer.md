# Comparer Interface

> Interface for comparing perceptual hashes using appropriate distance metrics

## Comparer

The `Comparer` interface defines the contract for measuring similarity between two hashes using a distance metric appropriate for the algorithm that produced them.

```go
type Comparer interface {
    Compare(hashtype.Hash, hashtype.Hash) (similarity.Distance, error)
}
```

### Methods

- **`Compare`** `func(hashtype.Hash, hashtype.Hash) (similarity.Distance, error)` — Computes the distance between two hashes using an algorithm-appropriate metric.
- **`h1`** `hashtype.Hash` — The first hash to compare.
- **`h2`** `hashtype.Hash` — The second hash to compare.
- **`distance`** `similarity.Distance` — The computed distance between the two hashes. Lower values indicate more similar images.
  * For binary hashes: typically Hamming distance (number of differing bits)
  * For numeric hashes: typically Euclidean (L2) distance
  * Some algorithms use specialized metrics (e.g., chi-square, cosine distance)
- **`error`** `error` — Returns an error if the comparison fails. Common errors:
  * `ErrIncompatibleHash` - Hash types are incompatible
  * `ErrHashLengthMismatch` - Hash lengths don't match

## HasherComparer

Most algorithms in imghash implement both `Hasher` and `Comparer` interfaces through the `HasherComparer` interface:

```go
type HasherComparer interface {
    Hasher
    Comparer
}
```

This allows a single algorithm instance to both compute and compare hashes.

## Usage Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Open two images
    img1, _ := imghash.OpenImage("photo1.jpg")
    img2, _ := imghash.OpenImage("photo2.jpg")
    
    // Create a hasher/comparer
    algo := imghash.Average{}
    
    // Calculate hashes
    hash1, _ := algo.Calculate(img1)
    hash2, _ := algo.Calculate(img2)
    
    // Compare using the algorithm's built-in metric
    distance, err := algo.Compare(hash1, hash2)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Distance: %v\n", distance)
    
    // Lower distance means more similar images
    if distance < 10 {
        fmt.Println("Images are very similar")
    } else if distance < 20 {
        fmt.Println("Images are somewhat similar")
    } else {
        fmt.Println("Images are different")
    }
}
```

## Algorithm-Specific Metrics

Each algorithm uses an appropriate distance metric for its hash type:

```go
// Binary hash algorithms use Hamming distance
avg := imghash.Average{}
hash1, _ := avg.Calculate(img1)
hash2, _ := avg.Calculate(img2)
dist, _ := avg.Compare(hash1, hash2) // Uses Hamming distance

// ColorMoment uses chi-square distance
cm, _ := imghash.NewColorMoment()
hash1, _ = cm.Calculate(img1)
hash2, _ = cm.Calculate(img2)
dist, _ = cm.Compare(hash1, hash2) // Uses chi-square distance

// GIST uses Euclidean distance
gist, _ := imghash.NewGIST()
hash1, _ = gist.Calculate(img1)
hash2, _ = gist.Calculate(img2)
dist, _ = gist.Compare(hash1, hash2) // Uses L2 distance

// DINOHash uses Hamming distance on its 96-bit Binary hash
dn, _ := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
hash1, _ = dn.Calculate(img1)
hash2, _ = dn.Calculate(img2)
dist, _ = dn.Compare(hash1, hash2) // Uses Hamming distance
```

## Custom Distance Functions

You can also use the generic `Compare` function with custom distance metrics:

```go
import "github.com/ajdnik/imghash/v2/similarity"

// Calculate hashes
avg := imghash.Average{}
hash1, _ := avg.Calculate(img1)
hash2, _ := avg.Calculate(img2)

// Use different distance metrics
hamming, _ := imghash.Compare(hash1, hash2, similarity.Hamming)
l1, _ := imghash.Compare(hash1, hash2, similarity.L1)
l2, _ := imghash.Compare(hash1, hash2, similarity.L2)
cosine, _ := imghash.Compare(hash1, hash2, similarity.Cosine)
```

## DistanceFunc Type

The `DistanceFunc` type allows passing distance functions as parameters:

```go
type DistanceFunc func(hashtype.Hash, hashtype.Hash) (similarity.Distance, error)
```

All functions in the similarity package satisfy this signature:

```go
var (
    _ DistanceFunc = similarity.Hamming
    _ DistanceFunc = similarity.L1
    _ DistanceFunc = similarity.L2
    _ DistanceFunc = similarity.Cosine
    _ DistanceFunc = similarity.ChiSquare
    _ DistanceFunc = similarity.PCC
    _ DistanceFunc = similarity.Jaccard
)
```

## Error Handling

```go
hash1, _ := avg.Calculate(img1)
hash2, _ := phash.Calculate(img2) // Different algorithm!

// This will fail - hashes from different algorithms
dist, err := avg.Compare(hash1, hash2)
if err != nil {
    if errors.Is(err, imghash.ErrIncompatibleHash) {
        fmt.Println("Cannot compare hashes from different algorithms")
    }
}

// Use generic Compare for flexibility
dist, err = imghash.Compare(hash1, hash2, similarity.L2)
// This might work if both hashes support the metric
```

## See Also

* [Hasher Interface](api/hasher) - For computing hashes
* [Similarity Functions](api/similarity) - Distance metric implementations
* [Hash Types](api/hash-types) - Hash type definitions
