# BoVW (Bag of Visual Words)

> Flexible local-feature hash using ORB or AKAZE with histogram, MinHash, or SimHash representations

## Overview

BoVW (Bag of Visual Words) implements local feature-based hashing using ORB-like or AKAZE-like keypoint descriptors. It extracts local binary descriptors and represents them as either a normalized histogram, MinHash signature, or SimHash bit-signature.

This algorithm is ideal for:

* Object recognition and matching
* Finding images containing similar local features
* Partial image matching (finds similar regions)
* Handling occlusions and local transformations

> [!NOTE]
> BoVW can return **Float64** (Histogram/MinHash) or **Binary** (SimHash) hash types depending on storage configuration.

## How It Works

1. **Resize**: Image is resized to the specified dimensions (default: 256×256)
2. **Feature Detection**: Detects keypoints using ORB or AKAZE
   * **ORB**: FAST corner detector with orientation
   * **AKAZE**: Multi-scale Hessian detector
3. **Descriptor Extraction**: Computes binary descriptors at each keypoint
   * **ORB**: BRIEF-like 256-bit descriptors
   * **AKAZE**: 256-bit oriented descriptors
4. **Visual Words**: Maps descriptors to visual vocabulary (hash-based bucketing)
5. **Representation**: Creates final hash based on storage type:
   * **Histogram**: L2-normalized word frequency vector
   * **MinHash**: Locality-sensitive signature for set similarity
   * **SimHash**: Bit signature for approximate matching

## Constructor

```go
func NewBoVW(opts ...BoVWOption) (BoVW, error)
```

Creates a new BoVW hasher with optional configuration.

### Available Options

* `WithSize(width, height uint)` - Set resize dimensions (default: 256×256)
* `WithInterpolation(interp Interpolation)` - Set interpolation method (default: Bilinear)
* `WithBoVWFeature(feature BoVWFeatureType)` - Set feature detector (default: BoVWORB)
* `WithBoVWStorage(storage BoVWStorageType)` - Set storage type (default: BoVWHistogram)
* `WithVocabularySize(size uint)` - Set visual vocabulary size (default: 256)
* `WithMaxKeypoints(count uint)` - Set maximum keypoints (default: 500)
* `WithMinHashSize(size uint)` - Set MinHash signature size (default: 64)
* `WithSimHashBits(bits uint)` - Set SimHash bit length (default: 128)
* `WithDistance(fn DistanceFunc)` - Override default distance function

### Feature Types

* `imghash.BoVWORB` - ORB-like FAST + BRIEF pipeline (default)
* `imghash.BoVWAKAZE` - AKAZE-like Hessian detector with binary descriptors

### Storage Types

* `imghash.BoVWHistogram` - Normalized histogram (Float64, Cosine distance)
* `imghash.BoVWMinHash` - MinHash signature (Float64, Jaccard distance)
* `imghash.BoVWSimHash` - SimHash bit-signature (Binary, Jaccard distance)

## Usage Examples

### Histogram Storage (Default)

```go
package main

import (
    "fmt"
    "image"
    _ "image/jpeg"
    "os"

    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create BoVW hasher with histogram storage
    hasher, err := imghash.NewBoVW(
        imghash.WithBoVWFeature(imghash.BoVWORB),
        imghash.WithBoVWStorage(imghash.BoVWHistogram),
        imghash.WithVocabularySize(512),
        imghash.WithMaxKeypoints(1000),
    )
    if err != nil {
        panic(err)
    }

    img1, _ := loadImage("object1.jpg")
    img2, _ := loadImage("object2.jpg")

    hash1, _ := hasher.Calculate(img1)
    hash2, _ := hasher.Calculate(img2)

    // Uses Cosine distance for histogram
    distance, _ := hasher.Compare(hash1, hash2)
    fmt.Printf("Distance: %.4f\n", distance)
}
```

### MinHash Storage

```go
// Create BoVW with MinHash for fast approximate matching
hasher, err := imghash.NewBoVW(
    imghash.WithBoVWStorage(imghash.BoVWMinHash),
    imghash.WithMinHashSize(128),
    imghash.WithVocabularySize(256),
)

// Uses Jaccard distance for MinHash
distance, _ := hasher.Compare(hash1, hash2)
```

### SimHash Storage

```go
// Create BoVW with SimHash for compact binary signatures
hasher, err := imghash.NewBoVW(
    imghash.WithBoVWStorage(imghash.BoVWSimHash),
    imghash.WithSimHashBits(256),
    imghash.WithBoVWFeature(imghash.BoVWAKAZE),
)

// Returns Binary hash, uses Jaccard distance
distance, _ := hasher.Compare(hash1, hash2)
```

### Helper Function

```go
func loadImage(path string) (image.Image, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    img, _, err := image.Decode(f)
    return img, err
}
```

## Default Settings

- **`width`** `uint` (default: `256`) — Image resize width
- **`height`** `uint` (default: `256`) — Image resize height
- **`interpolation`** `Interpolation` (default: `Bilinear`) — Resize interpolation method
- **`featureType`** `BoVWFeatureType` (default: `BoVWORB`) — Local feature detector (BoVWORB or BoVWAKAZE)
- **`storageType`** `BoVWStorageType` (default: `BoVWHistogram`) — Storage representation (BoVWHistogram, BoVWMinHash, or BoVWSimHash)
- **`vocabularySize`** `uint` (default: `256`) — Visual vocabulary size (must be > 0)
- **`maxKeypoints`** `uint` (default: `500`) — Maximum number of keypoints to detect (must be > 0)
- **`minHashSize`** `uint` (default: `64`) — MinHash signature size (must be > 0, used when storageType = BoVWMinHash)
- **`simHashBits`** `uint` (default: `128`) — SimHash bit length (must be > 0, used when storageType = BoVWSimHash)

## Hash Type

Depends on storage type:

* **BoVWHistogram**: `hashtype.Float64` with size = `vocabularySize`
* **BoVWMinHash**: `hashtype.Float64` with size = `minHashSize`
* **BoVWSimHash**: `hashtype.Binary` with size = `simHashBits` bits

## Distance Metrics

Default comparison varies by storage type:

* **BoVWHistogram**: Cosine distance (lower = more similar)
* **BoVWMinHash**: Jaccard distance (estimates set similarity)
* **BoVWSimHash**: Jaccard distance (binary similarity)

All can be overridden with `WithDistance()`.

## Feature Detector Comparison

- **BoVWORB** — * FAST corner detector * Faster computation * Good for textured images * Orientation-compensated BRIEF
- **BoVWAKAZE** — * Multi-scale Hessian detector * More robust features * Better for blob-like structures * Scale-space analysis

## Storage Type Comparison

- **Histogram** — * Most accurate * Larger hash size * Cosine similarity * Best for quality
- **MinHash** — * Compact Float64 * Approximate matching * Good speed/quality * Jaccard estimation
- **SimHash** — * Binary hash * Smallest size * Fastest comparison * Good for large scale

## Performance Tuning

**Vocabulary Size**:

* Smaller (128-256): Faster, less discriminative
* Larger (512-1024): More discriminative, slower, larger hash

**Max Keypoints**:

* Fewer (100-500): Faster detection, may miss features
* More (1000+): Better coverage, slower

**Image Size**:

* Smaller (128×128): Faster, fewer features
* Larger (512×512): More features, better quality, slower
