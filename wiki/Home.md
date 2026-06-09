# Introduction

> Go implementation of multiple perceptual hash algorithms for images

# Introduction to imghash

imghash is a comprehensive Go library providing perceptual hash algorithms for images. Unlike cryptographic hashes, perceptual hashes are designed to be similar for visually similar images, making them ideal for duplicate detection, content moderation, and image similarity search.

## What are Perceptual Hashes?

Perceptual hashes generate fingerprints based on visual content rather than raw bytes. This means:

* **Similar images produce similar hashes** - Even with compression, resizing, or minor edits
* **Fast comparison** - Binary hashes use efficient Hamming distance calculations
* **Scalable** - Suitable for large-scale image deduplication and search
* **Multiple algorithms** - Choose the right algorithm for your use case

## Key Features

- **Multiple Hash Algorithms** — From simple Average and Difference hashes to sophisticated PDQ, GIST, and Zernike moments
- **Multiple Distance Metrics** — Hamming, L1, L2, Cosine, Chi-Square, Jaccard, and more
- **Production Ready** — Battle-tested algorithms including Facebook's PDQ for content moderation
- **Type Safe** — Strongly typed hash representations (Binary, UInt8, Float64) with compile-time guarantees

## Available Algorithms

imghash provides perceptual hash algorithms, each optimized for different use cases:

### Binary Hashes (Hamming Distance)

Fast bit-level comparison using Hamming distance:

* **Average** - Simple threshold-based hash (8x8 = 64 bits)
* **Difference** - Gradient-based hash
* **Median** - Median threshold for better noise resistance
* **PHash** - DCT-based perceptual hash with weighted Hamming
* **WHash** - Wavelet-based hash
* **MarrHildreth** - Edge detection hash
* **BlockMean** - Block-based average hash
* **PDQ** - Facebook's robust 256-bit hash for content moderation
* **RASH** - Radial histogram hash
* **DINOHash** - DINOv2 ViT-S/14+reg deep perceptual hash (96 bits, adversarial-robust)

### Float64 Hashes (L2/Cosine Distance)

Floating-point descriptors for complex visual features:

* **ColorMoment** - Color distribution moments (L2 distance)
* **Zernike** - Zernike moments for shape description (L2 distance)
* **GIST** - Scene descriptor (Cosine distance)
* **BoVW** - Bag of Visual Words (supports Histogram/Float64, MinHash/Float64, or SimHash/Binary storage)

### UInt8 Hashes (L1/L2 Distance)

Byte-level descriptors for efficient storage:

* **CLD** - Color Layout Descriptor (L2 distance)
* **EHD** - Edge Histogram Descriptor (L1 distance)
* **LBP** - Local Binary Patterns (Chi-Square distance)
* **HOGHash** - Histogram of Oriented Gradients (Cosine distance)
* **RadialVariance** - Radial variance signature (L1 distance)

## Quick Start

- **[Installation](installation)** — Install imghash and set up your Go project
- **[Quick Start Guide](quickstart)** — Compute your first hash in under 5 minutes

## Core Concepts

### Hash Types

imghash uses three strongly-typed hash representations:

```go
type Hash interface {
    String() string
    Len() int
    ValueAt(idx int) float64
}

type Binary []byte      // Bit-level hashes (e.g., PDQ, Average)
type UInt8 []uint8      // Byte-level descriptors (e.g., CLD, EHD)
type Float64 []float64  // Floating-point features (e.g., GIST, Zernike)
```

### Interfaces

All algorithms implement consistent interfaces:

```go
// Hasher computes a perceptual hash from an image
type Hasher interface {
    Calculate(image.Image) (Hash, error)
}

// Comparer measures similarity between two hashes
type Comparer interface {
    Compare(Hash, Hash) (Distance, error)
}

// HasherComparer combines both capabilities
type HasherComparer interface {
    Hasher
    Comparer
}
```

### Distance Metrics

Each algorithm has a default distance metric optimized for its hash type:

* **Hamming** - Bit differences for Binary hashes
* **L1 (Manhattan)** - Sum of absolute differences
* **L2 (Euclidean)** - Geometric distance
* **Cosine** - Angular similarity
* **Chi-Square** - Histogram comparison
* **Jaccard** - Set similarity

## When to Use imghash

> [!NOTE]
> imghash is ideal for:
>
> * **Duplicate Detection** - Find exact and near-duplicate images
> * **Content Moderation** - Match against known problematic content using PDQ
> * **Image Search** - Find visually similar images at scale
> * **Copyright Protection** - Detect unauthorized image use
> * **Clustering** - Group similar images together

## Performance Characteristics

* **Hash Generation** - Typically 1-50ms for classical algorithms, around one second for DINOHash (pure-Go ViT inference)
* **Comparison** - Sub-microsecond for binary hashes (Hamming distance)
* **Memory** - Binary hashes are extremely compact (8-256 bits typical); DINOHash weights add ~85 MB when imported via the sibling `dinoweights` module
* **Accuracy** - Algorithm-dependent; PDQ offers excellent robustness to standard transformations, DINOHash adds robustness to crops, recolors, and adversarial edits at higher CPU cost

## Next Steps

### Step 1: Install the library

Follow the [Installation Guide](installation) to add imghash to your project

### Step 2: Try the quick start

Run your first perceptual hash using the [Quick Start Guide](quickstart)

### Step 3: Choose an algorithm

Learn about different algorithms and pick the right one for your use case

### Step 4: Integrate into your app

Use imghash's convenient helpers like `HashFile()` and `Compare()`
