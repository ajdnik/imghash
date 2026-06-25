# Choosing the Right Algorithm

> Decision guide for selecting the best perceptual hash algorithm for your use case

## Overview

Imghash provides different perceptual hashing algorithms, each optimized for different scenarios. This guide helps you select the right algorithm based on your requirements.

## Quick Recommendation

> [!NOTE]
> If you're unsure which algorithm to choose, **start with PDQ**. It provides excellent robustness to common transformations while remaining fast enough for most applications.

```go
pdq, err := imghash.NewPDQ()
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(pdq, "image.jpg")
```

## Decision Tree

### Step 1: Consider Your Primary Use Case

**Duplicate Detection & Content Moderation**

* Use **PDQ** for general-purpose deduplication
* Use **Average** or **Difference** for simple, fast duplicate detection
* Use **RASH** for rotation-invariant duplicate detection
* Use **DINOHash** when adversarial robustness or semantic-level matching matters (~1s/image CPU cost)

**Visual Similarity Search**

* Use **GIST** for scene-level similarity
* Use **ColorMoment** when color is important
* Use **BoVW** for feature-based similarity

**Copyright & Forensics**

* Use **PDQ** (developed by Meta for content matching)
* Use **PHash** for academic/research applications

**Texture & Pattern Matching**

* Use **LBP** for texture analysis
* Use **BlockMean** for block-based patterns

### Step 2: Evaluate Performance Requirements

**Speed Priority (Simple & Fast)**

* **Average**: Fastest, good for basic duplicate detection
* **Difference**: Very fast, detects horizontal changes
* **Median**: Fast, more robust than Average

**Balance (Robust & Reasonably Fast)**

* **PDQ**: Best balance of speed and robustness
* **WHash**: Wavelet-based, good middle ground
* **BlockMean**: Block-based approach, efficient

**Quality Priority (More Computation)**

* **GIST**: Rich scene descriptor
* **BoVW**: Feature-based, highly configurable
* **Zernike**: Moment-based, rotation invariant
* **DINOHash**: Deep semantic perceptual hash; highest robustness but ~1s/image

### Step 3: Consider Image Characteristics

**Color is Important**

* **ColorMoment**: Specifically designed for color images
* **CLD** (Color Layout Descriptor): MPEG-7 standard

**Grayscale/Structure Only**

* **Average**, **Difference**, **Median**: Simple structural hashes
* **PHash**: DCT-based structural hash
* **PDQ**: Advanced structural hash

**Rotation Invariance Needed**

* **RASH**: Rotation and Scale Hash
* **RadialVariance**: Radial-based approach
* **Zernike**: Moment-based rotation invariance

**Edge/Texture Focus**

* **MarrHildreth**: Edge-based hashing
* **EHD** (Edge Histogram Descriptor): MPEG-7 standard
* **LBP**: Local Binary Patterns for texture
* **HOGHash**: Histogram of Oriented Gradients

## Algorithm Comparison Table

> [!NOTE]
> All binary hash algorithms use **Hamming distance** by default. Float64 algorithms typically use **L2 (Euclidean)** or **Cosine** distance.

### Binary Hash Algorithms

Binary hashes are compact, fast to compare, and work well for duplicate detection.

| Algorithm          | Hash Type        | Use Case                       | Speed            | Robustness         |
| ------------------ | ---------------- | ------------------------------ | ---------------- | ------------------ |
| **Average**        | Binary (64-bit)  | Basic duplicate detection      | ⚡️⚡️⚡️ Very Fast | ⭐️⭐️ Moderate      |
| **Difference**     | Binary (64-bit)  | Horizontal gradient changes    | ⚡️⚡️⚡️ Very Fast | ⭐️⭐️ Moderate      |
| **Median**         | Binary (64-bit)  | Improved duplicate detection   | ⚡️⚡️⚡️ Very Fast | ⭐️⭐️⭐️ Good        |
| **PHash**          | Binary (64-bit)  | Academic/research standard     | ⚡️⚡️ Fast        | ⭐️⭐️⭐️⭐️ Excellent |
| **WHash**          | Binary (64-bit)  | Wavelet-based matching         | ⚡️⚡️ Fast        | ⭐️⭐️⭐️ Good        |
| **MarrHildreth**   | Binary (576-bit) | Edge-based detection           | ⚡️ Moderate      | ⭐️⭐️⭐️ Good        |
| **BlockMean**      | Binary (256-bit) | Block-based patterns           | ⚡️⚡️ Fast        | ⭐️⭐️⭐️ Good        |
| **PDQ**            | Binary (256-bit) | Production duplicate detection | ⚡️⚡️ Fast        | ⭐️⭐️⭐️⭐️⭐️ Best    |
| **RASH**           | Binary           | Rotation + scale invariant     | ⚡️ Moderate      | ⭐️⭐️⭐️⭐️ Excellent |
| **BoVW (SimHash)** | Binary           | Feature-based vocabulary       | ⚡️ Slow          | ⭐️⭐️⭐️⭐️ Excellent |
| **DINOHash**       | Binary (96-bit)  | Adversarial-robust matching    | 🐢 ~1s/image     | ⭐️⭐️⭐️⭐️⭐️ Best    |

### Float64 Hash Algorithms

Float64 hashes provide richer representations, suitable for similarity search.

| Algorithm            | Hash Type | Default Metric | Use Case                    |
| -------------------- | --------- | -------------- | --------------------------- |
| **ColorMoment**      | Float64   | L2 (Euclidean) | Color-aware similarity      |
| **Zernike**          | Float64   | L2 (Euclidean) | Rotation-invariant matching |
| **GIST**             | Float64   | Cosine         | Scene-level similarity      |
| **BoVW (Histogram)** | Float64   | Cosine         | Feature vocabulary          |
| **BoVW (MinHash)**   | Float64   | Jaccard        | Feature signatures          |

### UInt8 Hash Algorithms

UInt8 hashes are compact histograms, balancing size and expressiveness.

| Algorithm          | Hash Type | Default Metric | Use Case                |
| ------------------ | --------- | -------------- | ----------------------- |
| **CLD**            | UInt8     | L2 (Euclidean) | Color layout (MPEG-7)   |
| **EHD**            | UInt8     | L1 (Manhattan) | Edge histogram (MPEG-7) |
| **LBP**            | UInt8     | Chi-Square     | Texture analysis        |
| **HOGHash**        | UInt8     | Cosine         | Gradient-based matching |
| **RadialVariance** | UInt8     | L1 (Manhattan) | Radial patterns         |

## Use Case Examples

<details>
<summary>Detecting Near-Duplicate Images</summary>

**Best Choices**: PDQ, Average, Median, Difference

For production systems handling user uploads, JPEG compression, and minor edits:

```go
pdq, _ := imghash.NewPDQ()
h1, _ := imghash.HashFile(pdq, "original.jpg")
h2, _ := imghash.HashFile(pdq, "compressed.jpg")

dist, _ := pdq.Compare(h1, h2)

// PDQ is robust to JPEG compression, cropping, minor edits
if dist < 10 {
    fmt.Println("Likely duplicate")
}
```

For simpler cases with less variation:

```go
avg, _ := imghash.NewAverage()
h1, _ := imghash.HashFile(avg, "img1.png")
h2, _ := imghash.HashFile(avg, "img2.png")

dist, _ := avg.Compare(h1, h2)
if dist < 5 {
    fmt.Println("Duplicate detected")
}
```

</details>

<details>
<summary>Finding Similar Images by Content</summary>

**Best Choices**: GIST, ColorMoment, BoVW

For finding visually similar images (not exact duplicates):

```go
gist, _ := imghash.NewGIST()
h1, _ := imghash.HashFile(gist, "beach1.jpg")
h2, _ := imghash.HashFile(gist, "beach2.jpg")

dist, _ := gist.Compare(h1, h2)

// GIST captures scene-level features
// Lower cosine distance = more similar scenes
if dist < 0.3 {
    fmt.Println("Similar scenes")
}
```

When color matters:

```go
cm, _ := imghash.NewColorMoment()
h1, _ := imghash.HashFile(cm, "sunset1.jpg")
h2, _ := imghash.HashFile(cm, "sunset2.jpg")

dist, _ := cm.Compare(h1, h2)
if dist < 15.0 {
    fmt.Println("Similar color distribution")
}
```

</details>

<details>
<summary>Detecting Rotated Images</summary>

**Best Choices**: RASH, Zernike, RadialVariance

For images that may be rotated:

```go
rash, _ := imghash.NewRASH()
h1, _ := imghash.HashFile(rash, "original.jpg")
h2, _ := imghash.HashFile(rash, "rotated.jpg")

dist, _ := rash.Compare(h1, h2)

// RASH is designed to be rotation and scale invariant
if dist < 20 {
    fmt.Println("Same image despite rotation")
}
```

</details>

<details>
<summary>Texture and Pattern Matching</summary>

**Best Choices**: LBP, HOGHash, BlockMean

For texture analysis:

```go
lbp, _ := imghash.NewLBP()
h1, _ := imghash.HashFile(lbp, "fabric1.jpg")
h2, _ := imghash.HashFile(lbp, "fabric2.jpg")

dist, _ := lbp.Compare(h1, h2)

// LBP excels at texture patterns
// Uses Chi-Square distance by default
```

</details>

<details>
<summary>Content Moderation at Scale</summary>

**Best Choice**: PDQ

Facebook/Meta developed PDQ specifically for large-scale content moderation:

```go
pdq, _ := imghash.NewPDQ()

// Hash known harmful content
blocklist := make(map[string]imghash.Binary)
for _, path := range knownBadImages {
    hash, _ := imghash.HashFile(pdq, path)
    blocklist[path] = hash.(imghash.Binary)
}

// Check new upload
newHash, _ := imghash.HashFile(pdq, "upload.jpg")

for _, knownHash := range blocklist {
    dist, _ := pdq.Compare(newHash, knownHash)
    if dist < 10 {
        // Flag for review
        break
    }
}
```

</details>

<details>
<summary>Adversarial / Semantic Near-Duplicate Detection</summary>

**Best Choice**: DINOHash (with PDQ pre-filter)

When inputs may be deliberately edited to defeat classical hashes (e.g. heavy crops, recolors, watermarks, adversarial noise), the deep semantic features of DINOv2 hold up far better than DCT-based hashes. Pair with a faster pre-filter to amortize the per-image cost:

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/dinoweights"
)

pdq, _ := imghash.NewPDQ()
dn, _  := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))

refPDQ, _  := imghash.HashFile(pdq, "reference.jpg")
refDINO, _ := imghash.HashFile(dn, "reference.jpg")

candPDQ, _ := imghash.HashFile(pdq, "candidate.jpg")
pdqDist, _ := pdq.Compare(refPDQ, candPDQ)

// Cheap PDQ pre-filter, then expensive DINOHash verify on survivors.
if pdqDist < 100 {
    candDINO, _ := imghash.HashFile(dn, "candidate.jpg")
    dnDist, _ := dn.Compare(refDINO, candDINO)
    if dnDist < 25 {
        // Confirmed adversarial / semantic near-duplicate.
    }
}
```

</details>

## Algorithm Characteristics

### Simple Threshold-Based Algorithms

- **Average** — Compares each pixel to the image mean. Fast and simple, good baseline.
- **Difference** — Compares adjacent pixels horizontally. Detects gradient changes.
- **Median** — Uses median instead of mean. More robust to outliers than Average.

### Transform-Based Algorithms

- **PHash** — Discrete Cosine Transform (DCT) based. Academic standard for perceptual hashing.
- **PDQ** — Advanced DCT with Jarosz filtering. Industry standard for content moderation.
- **WHash** — Haar wavelet transform. Good frequency domain representation.

### Feature-Based Algorithms

- **GIST** — Scene-level descriptor. Captures spatial structure and frequency information.
- **BoVW** — Bag of Visual Words. Uses SIFT features and visual vocabulary.
- **HOGHash** — Histogram of Oriented Gradients. Shape and appearance descriptor.

### Specialized Algorithms

- **RASH** — Rotation and Scale Hash. Invariant to rotation and scaling.
- **ColorMoment** — Color-aware using Hu moments in YCrCb and HSV spaces.
- **LBP** — Local Binary Patterns. Excellent for texture classification.

### Deep Learning Algorithms

- **DINOHash** — Frozen DINOv2 ViT-S/14+reg backbone fused with a 96-bit PCA-aligned head. Captures high-level semantic features; robust to crops, recolors, lossy re-encoding, and adversarial edits that defeat classical hashes. Pure-Go ViT inference, ~1 second per image. Ships its ~85 MB model weights in the sibling [`dinoweights`](installation#dinohash-weights-module) module.

## Performance Considerations

> [!TIP]
> **Hash Size vs. Accuracy Trade-off**
>
> * Smaller hashes (64-bit): Faster comparison, less storage, lower accuracy
> * Larger hashes (256-bit+): Slower comparison, more storage, higher accuracy
> * Float64 hashes: Most expressive but require more sophisticated distance metrics

### Computation Speed (Relative)

```
Fastest:     Average, Difference, Median
Fast:        PHash, PDQ, WHash, BlockMean
Moderate:    MarrHildreth, RASH, ColorMoment, CLD, EHD, LBP, HOGHash
Slow:        GIST, BoVW, Zernike, RadialVariance
Very slow:   DINOHash (~1 second per image, pure-Go ViT inference)
```

### Memory Usage

```
Binary (8-64 bytes):    Average, Difference, Median, PHash, WHash
Binary (12 bytes):      DINOHash (96 bits)
Binary (32+ bytes):     BlockMean, PDQ, MarrHildreth, RASH
UInt8 (80-256 bytes):   CLD, EHD, LBP, HOGHash, RadialVariance
Float64 (128+ bytes):   ColorMoment, Zernike, GIST, BoVW
```

> [!NOTE]
> DINOHash hashes themselves are tiny (12 bytes). The cost is the embedded model weights (~85 MB) loaded once at construction via the sibling `dinoweights` module.

## Next Steps

- **[Comparing Hashes](comparing-hashes)** — Learn how to compare hashes and interpret distance values
- **[Practical Examples](examples)** — See real-world examples and use cases
