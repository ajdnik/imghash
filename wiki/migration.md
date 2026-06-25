# Migration Guide

> Guide for migrating from v1 to v2 of imghash

## Overview

This guide helps you migrate from imghash v1 to v2. Version 2 introduces significant improvements in API design, type safety, and extensibility while maintaining the core functionality.

## Quick Migration Checklist

### Step 1: Update Import Path

Change your import path to include the `/v2` suffix:

```go
// v1
import "github.com/ajdnik/imghash"

// v2
import "github.com/ajdnik/imghash/v2"
```

### Step 2: Update Constructor Calls

Algorithm constructors now return errors and use functional options:

```go
// v1
avg := imghash.NewAverage()

// v2
avg, err := imghash.NewAverage()
if err != nil {
    panic(err)
}
```

### Step 3: Update Hash Method Calls

The `Hash` method is now called `Calculate`:

```go
// v1
hash := avg.Hash(img)

// v2
hash, err := avg.Calculate(img)
if err != nil {
    panic(err)
}
```

### Step 4: Update Distance Calculations

Distance metrics are now in the `similarity` package:

```go
// v1
dist := imghash.HammingDistance(h1, h2)

// v2
import "github.com/ajdnik/imghash/v2/similarity"

dist, err := similarity.Hamming(h1, h2)
if err != nil {
    panic(err)
}
```

## Major Changes

### Package Structure

```go
import "github.com/ajdnik/imghash"

// Everything in one package
avg := imghash.NewAverage()
hash := avg.Hash(img)
dist := imghash.HammingDistance(h1, h2)
```

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

// Organized into subpackages
avg, _ := imghash.NewAverage()
hash, _ := avg.Calculate(img)
dist, _ := similarity.Hamming(h1, h2)
```

### Error Handling

V2 introduces comprehensive error handling throughout the API:

```go
// Constructors never failed
avg := imghash.NewAverage()

// Methods could panic or return invalid results
hash := avg.Hash(img)

// No validation
dist := imghash.HammingDistance(h1, h2)
```

```go
// Constructors validate configuration
avg, err := imghash.NewAverage()
if err != nil {
    // Handle invalid configuration
    return err
}

// Methods return errors
hash, err := avg.Calculate(img)
if err != nil {
    // Handle hashing errors
    return err
}

// Distance functions validate inputs
dist, err := similarity.Hamming(h1, h2)
if err != nil {
    // Handle incompatible hash types
    return err
}
```

### Configuration Options

V2 uses functional options for cleaner, more extensible configuration:

```go
// Limited customization with struct fields
avg := imghash.NewAverage()
avg.Size = 16
avg.Interpolation = imghash.Bicubic
```

```go
// Functional options pattern
avg, err := imghash.NewAverage(
    imghash.WithSize(16, 16),
    imghash.WithInterpolation(imghash.Bicubic),
)
```

### Hash Types

V2 introduces type-safe hash types:

```go
// Generic []byte or []float64
type Hash interface{}

hash := avg.Hash(img) // Returns interface{}
binaryHash := hash.([]byte)
```

```go
// Type-safe hash types
type Hash interface {
    Len() int
    ValueAt(int) float64
}

type Binary []byte
type UInt8 []uint8
type Float64 []float64

hash, _ := avg.Calculate(img)
binaryHash := hash.(imghash.Binary)
```

### Distance Metrics

Distance functions moved to dedicated package with validation:

```go
import "github.com/ajdnik/imghash"

dist := imghash.HammingDistance(h1, h2)
dist2 := imghash.EuclideanDistance(h1, h2)
```

```go
import "github.com/ajdnik/imghash/v2/similarity"

// Hamming for binary hashes
dist, err := similarity.Hamming(h1, h2)

// L2 (Euclidean) for float hashes
dist2, err := similarity.L2(h1, h2)

// Cosine for directional similarity
dist3, err := similarity.Cosine(h1, h2)
```

## Algorithm-Specific Changes

### PDQ

```go
pdq := imghash.NewPDQ()
hash := pdq.Hash(img)
dist := imghash.HammingDistance(h1, h2)
```

```go
pdq, err := imghash.NewPDQ()
if err != nil {
    panic(err)
}

hash, err := pdq.Calculate(img)
if err != nil {
    panic(err)
}

dist, err := pdq.Compare(h1, h2)
if err != nil {
    panic(err)
}
```

### PHash

```go
phash := imghash.NewPHash()
phash.Size = 32
hash := phash.Hash(img)
```

```go
phash, err := imghash.NewPHash(
    imghash.WithSize(32, 32),
)
if err != nil {
    panic(err)
}

hash, err := phash.Calculate(img)
if err != nil {
    panic(err)
}
```

### ColorMoment

```go
cm := imghash.NewColorMoment()
hash := cm.Hash(img)
dist := imghash.EuclideanDistance(h1, h2)
```

```go
cm, err := imghash.NewColorMoment(
    imghash.WithKernelSize(5),
    imghash.WithSigma(1.0),
)
if err != nil {
    panic(err)
}

hash, err := cm.Calculate(img)
if err != nil {
    panic(err)
}

// Uses L2 (Euclidean) by default
dist, err := cm.Compare(h1, h2)
if err != nil {
    panic(err)
}
```

## Common Migration Patterns

### Pattern 1: Basic Duplicate Detection

```go
import "github.com/ajdnik/imghash"

func isDuplicate(img1, img2 string) bool {
    avg := imghash.NewAverage()

    h1 := avg.HashFile(img1)
    h2 := avg.HashFile(img2)

    dist := imghash.HammingDistance(h1, h2)
    return dist < 5
}
```

```go
import "github.com/ajdnik/imghash/v2"

func isDuplicate(img1, img2 string) (bool, error) {
    avg, err := imghash.NewAverage()
    if err != nil {
        return false, err
    }

    h1, err := imghash.HashFile(avg, img1)
    if err != nil {
        return false, err
    }

    h2, err := imghash.HashFile(avg, img2)
    if err != nil {
        return false, err
    }

    dist, err := avg.Compare(h1, h2)
    if err != nil {
        return false, err
    }

    return dist < 5, nil
}
```

### Pattern 2: Custom Configuration

```go
import "github.com/ajdnik/imghash"

func createCustomHasher() *imghash.MarrHildreth {
    mh := imghash.NewMarrHildreth()
    mh.Scale = 1
    mh.Alpha = 2
    mh.Size = 512
    return mh
}
```

```go
import "github.com/ajdnik/imghash/v2"

func createCustomHasher() (imghash.MarrHildreth, error) {
    return imghash.NewMarrHildreth(
        imghash.WithScale(1),
        imghash.WithAlpha(2),
        imghash.WithSize(512, 512),
        imghash.WithInterpolation(imghash.Bicubic),
    )
}
```

### Pattern 3: Similarity Search

```go
import "github.com/ajdnik/imghash"

func findSimilar(query string, database []string) []string {
    gist := imghash.NewGIST()
    queryHash := gist.HashFile(query)

    var results []string
    for _, img := range database {
        hash := gist.HashFile(img)
        dist := imghash.CosineDistance(queryHash, hash)
        if dist < 0.3 {
            results = append(results, img)
        }
    }
    return results
}
```

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

func findSimilar(query string, database []string) ([]string, error) {
    gist, err := imghash.NewGIST()
    if err != nil {
        return nil, err
    }

    queryHash, err := imghash.HashFile(gist, query)
    if err != nil {
        return nil, err
    }

    var results []string
    for _, img := range database {
        hash, err := imghash.HashFile(gist, img)
        if err != nil {
            continue // Skip errors
        }

        dist, err := similarity.Cosine(queryHash, hash)
        if err != nil {
            continue
        }

        if dist < 0.3 {
            results = append(results, img)
        }
    }
    return results, nil
}
```

## Breaking Changes Summary

> [!WARNING]
> **These changes require code updates when migrating from v1 to v2:**

1. **Import path**: Add `/v2` suffix
2. **Constructors**: Now return `(Algorithm, error)` instead of just `Algorithm`
3. **Hash method**: Renamed from `Hash()` to `Calculate()` and returns error
4. **Distance functions**: Moved to `similarity` package and return errors
5. **Configuration**: Use functional options instead of struct fields
6. **Hash types**: Strong typing with `Binary`, `UInt8`, `Float64` types
7. **Compare method**: Each algorithm has a `Compare()` method with default metric

## Benefits of V2

- **Type Safety** — Strong typing prevents runtime errors and provides better IDE support
- **Error Handling** — Comprehensive error handling for production reliability
- **Extensibility** — Functional options make it easy to add new configuration
- **API Clarity** — Cleaner separation of concerns with dedicated packages
- **DINOHash** — New deep perceptual hash backed by a frozen DINOv2 ViT-S/14+reg backbone. 96-bit Hamming hash robust to crops, recolors, lossy re-encoding, and adversarial edits that defeat classical hashes. Ships its ~85 MB model weights in the sibling `dinoweights` module so importers do not pay the embed cost unless they opt in. See the [DINOHash algorithm page](dinohash) and the [installation guide](installation#dinohash-weights-module).

## Getting Help

- **[Examples](examples)** — See v2 examples in action

## Automated Migration

For large codebases, consider these steps:

1. **Update imports**:
   ```bash
   find . -name '*.go' -exec sed -i 's|github.com/ajdnik/imghash"|github.com/ajdnik/imghash/v2"|g' {} +
   ```

2. **Run tests** to identify breaking changes

3. **Fix errors systematically** by category (constructors, methods, etc.)

4. **Add error handling** where needed

5. **Validate** with comprehensive test coverage

> [!TIP]
> Start with a small subset of your codebase to validate the migration approach before applying it broadly.
