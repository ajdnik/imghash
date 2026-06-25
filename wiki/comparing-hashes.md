# Comparing Hashes

> Learn how to compare perceptual hashes and interpret distance metrics

## Overview

After computing perceptual hashes, you need to compare them to determine image similarity. This guide covers distance metrics, threshold selection, and best practices.

## Basic Comparison

### Using Algorithm's Compare Method

Every algorithm implements the `Compare` method with its default distance metric:

```go
pdq, _ := imghash.NewPDQ()

h1, _ := imghash.HashFile(pdq, "image1.jpg")
h2, _ := imghash.HashFile(pdq, "image2.jpg")

// Use the algorithm's default metric (Hamming for PDQ)
dist, err := pdq.Compare(h1, h2)
if err != nil {
    panic(err)
}

fmt.Printf("Distance: %v\n", dist)
```

### Using the Convenience Function

The `Compare` function automatically selects an appropriate metric:

```go
avg, _ := imghash.NewAverage()
h1, _ := imghash.HashFile(avg, "lena.jpg")
h2, _ := imghash.HashFile(avg, "cat.jpg")

// Automatically uses Hamming for Binary hashes
dist, err := imghash.Compare(h1, h2)
if err != nil {
    panic(err)
}

fmt.Printf("Distance: %v\n", dist)
// Output: Distance: 29
```

## Distance Metrics

> [!NOTE]
> Different hash types use different distance metrics. Choosing the right metric is crucial for accurate similarity measurement.

### Binary Hash Metrics

Binary hashes (Average, Difference, Median, PHash, PDQ, etc.) primarily use **Hamming distance**.

#### Hamming Distance

Counts the number of differing bits between two binary hashes:

```go
import "github.com/ajdnik/imghash/v2/similarity"

avg, _ := imghash.NewAverage()
h1, _ := imghash.HashFile(avg, "img1.jpg")
h2, _ := imghash.HashFile(avg, "img2.jpg")

// Hamming distance: number of bit differences
dist, _ := similarity.Hamming(h1, h2)

fmt.Printf("Hamming distance: %v\n", dist)
// Lower values = more similar
// 0 = identical
// 64 = completely different (for 64-bit hash)
```

**Interpretation:**

* `0`: Identical images (or perceptually identical)
* `1-5`: Very similar (minor differences)
* `6-10`: Similar (noticeable differences)
* `11-15`: Somewhat similar
* `16+`: Likely different images

> [!WARNING]
> Hamming distance thresholds depend on hash size. A distance of 10 is significant for a 64-bit hash but minor for a 256-bit hash like PDQ.

#### Weighted Hamming Distance

PHash uses weighted Hamming by default, where different bit positions have different importance:

```go
phash, _ := imghash.NewPHash()
h1, _ := imghash.HashFile(phash, "img1.jpg")
h2, _ := imghash.HashFile(phash, "img2.jpg")

// Uses weighted Hamming automatically
dist, _ := phash.Compare(h1, h2)
```

You can also use it directly:

```go
import "github.com/ajdnik/imghash/v2/similarity"

weights := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
dist, _ := similarity.WeightedHamming(h1, h2, weights)
```

### Float64 Hash Metrics

Float64 hashes (ColorMoment, Zernike, GIST, BoVW Histogram/MinHash) use continuous distance metrics.

#### L2 (Euclidean) Distance

Standard Euclidean distance in n-dimensional space:

```go
import "github.com/ajdnik/imghash/v2/similarity"

cm, _ := imghash.NewColorMoment()
h1, _ := imghash.HashFile(cm, "img1.jpg")
h2, _ := imghash.HashFile(cm, "img2.jpg")

// L2 distance (default for ColorMoment)
dist, _ := similarity.L2(h1, h2)

fmt.Printf("L2 distance: %.2f\n", dist)
// Lower values = more similar
```

**Formula:** `sqrt(sum((h1[i] - h2[i])^2))`

**Typical ranges:**

* `0.0`: Identical
* `0.0-10.0`: Very similar
* `10.0-50.0`: Moderately similar
* `50.0+`: Different

#### L1 (Manhattan) Distance

Sum of absolute differences:

```go
import "github.com/ajdnik/imghash/v2/similarity"

ehd, _ := imghash.NewEHD()
h1, _ := imghash.HashFile(ehd, "img1.jpg")
h2, _ := imghash.HashFile(ehd, "img2.jpg")

// L1 distance (default for EHD)
dist, _ := similarity.L1(h1, h2)
```

**Formula:** `sum(abs(h1[i] - h2[i]))`

#### Cosine Distance

Measures the angle between two vectors (1 - cosine similarity):

```go
import "github.com/ajdnik/imghash/v2/similarity"

gist, _ := imghash.NewGIST()
h1, _ := imghash.HashFile(gist, "img1.jpg")
h2, _ := imghash.HashFile(gist, "img2.jpg")

// Cosine distance (default for GIST)
dist, _ := similarity.Cosine(h1, h2)

fmt.Printf("Cosine distance: %.4f\n", dist)
```

**Formula:** `1 - (dot(h1, h2) / (||h1|| * ||h2||))`

**Interpretation:**

* `0.0`: Identical direction (very similar)
* `0.0-0.2`: Similar
* `0.2-0.5`: Moderately similar
* `0.5-1.0`: Different
* `1.0`: Opposite direction

#### Jaccard Distance

For set-based comparisons (BoVW MinHash/SimHash):

```go
import "github.com/ajdnik/imghash/v2/similarity"

bovw, _ := imghash.NewBoVW(
    imghash.WithBoVWStorage(imghash.MinHashStorage),
)
h1, _ := imghash.HashFile(bovw, "img1.jpg")
h2, _ := imghash.HashFile(bovw, "img2.jpg")

// Jaccard distance
dist, _ := similarity.Jaccard(h1, h2)
```

**Formula (binary):** `1 - (intersection / union)`

**Formula (MinHash):** `1 - (matching positions / signature length)`

### UInt8 Hash Metrics

UInt8 hashes (CLD, EHD, LBP, HOGHash, RadialVariance) use histogram-appropriate metrics.

#### Chi-Square Distance

Ideal for histogram comparison (default for LBP):

```go
import "github.com/ajdnik/imghash/v2/similarity"

lbp, _ := imghash.NewLBP()
h1, _ := imghash.HashFile(lbp, "texture1.jpg")
h2, _ := imghash.HashFile(lbp, "texture2.jpg")

// Chi-Square distance (default for LBP)
dist, _ := similarity.ChiSquare(h1, h2)
```

**Formula:** `sum((h1[i] - h2[i])^2 / (h1[i] + h2[i]))`

## Overriding Distance Metrics

You can override the default distance metric when creating an algorithm:

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

// Use L2 distance instead of default Hamming for PDQ
pdq, _ := imghash.NewPDQ(
    imghash.WithDistance(similarity.L2),
)

h1, _ := imghash.HashFile(pdq, "img1.jpg")
h2, _ := imghash.HashFile(pdq, "img2.jpg")

// Now uses L2 distance
dist, _ := pdq.Compare(h1, h2)
```

Or pass it to the convenience function:

```go
avg, _ := imghash.NewAverage()
h1, _ := imghash.HashFile(avg, "img1.jpg")
h2, _ := imghash.HashFile(avg, "img2.jpg")

// Override with Cosine distance
dist, _ := imghash.Compare(h1, h2, similarity.Cosine)
```

## Threshold Selection

> [!TIP]
> There's no universal threshold. The right threshold depends on:
>
> * The algorithm used
> * The distance metric
> * Your specific use case
> * Acceptable false positive/negative rates

### Recommended Thresholds by Algorithm

<details>
<summary>Binary Hashes (Hamming Distance)</summary>

**Average, Difference, Median (64-bit)**

```go
if dist <= 5 {
    // Very likely duplicate
} else if dist <= 10 {
    // Probably similar
} else {
    // Different images
}
```

**PDQ (256-bit)**

```go
if dist <= 10 {
    // Highly likely duplicate (Meta's recommendation)
} else if dist <= 31 {
    // Possibly similar
} else {
    // Different images
}
```

**PHash (64-bit, weighted Hamming)**

```go
if dist <= 5.0 {
    // Very similar
} else if dist <= 10.0 {
    // Moderately similar
}
```

**DINOHash (96-bit)**

```go
if dist <= 15 {
    // Very likely semantic duplicate (high confidence)
} else if dist <= 25 {
    // Probably the same scene or subject
} else if dist <= 40 {
    // Loosely related; check use case
} else {
    // Unrelated
}
```

> [!NOTE]
> DINOHash distances cluster near the binomial null (~48 bits) for unrelated images. Sub-25 distances usually indicate semantic relatedness even when no classical hash agrees, since the underlying ViT embedding captures scene content rather than pixel structure.

</details>

<details>
<summary>Float64 Hashes</summary>

**ColorMoment (L2 distance)**

```go
if dist <= 10.0 {
    // Very similar color distribution
} else if dist <= 30.0 {
    // Somewhat similar
}
```

**GIST (Cosine distance)**

```go
if dist <= 0.1 {
    // Very similar scenes
} else if dist <= 0.3 {
    // Moderately similar scenes
}
```

**BoVW Histogram (Cosine distance)**

```go
if dist <= 0.2 {
    // Similar feature distribution
} else if dist <= 0.5 {
    // Somewhat related
}
```

</details>

<details>
<summary>UInt8 Hashes</summary>

**LBP (Chi-Square distance)**

```go
if dist <= 50.0 {
    // Similar texture
} else if dist <= 100.0 {
    // Moderately similar texture
}
```

**CLD (L2 distance)**

```go
if dist <= 20.0 {
    // Similar color layout
} else if dist <= 50.0 {
    // Somewhat similar
}
```

</details>

### Empirical Threshold Tuning

The best approach is to calibrate thresholds on your specific dataset:

### Step 1: Collect Test Data

Create a labeled dataset with:

* Known duplicate pairs
* Known similar (but not duplicate) pairs
* Known different image pairs

### Step 2: Compute Distances

```go
type TestPair struct {
    img1, img2 string
    label      string // "duplicate", "similar", "different"
}

pdq, _ := imghash.NewPDQ()

for _, pair := range testPairs {
    h1, _ := imghash.HashFile(pdq, pair.img1)
    h2, _ := imghash.HashFile(pdq, pair.img2)
    dist, _ := pdq.Compare(h1, h2)

    fmt.Printf("%s\t%s\t%s\t%.2f\n", 
        pair.img1, pair.img2, pair.label, dist)
}
```

### Step 3: Analyze Distribution

Plot the distance distribution for each label category. Look for separation between categories.

### Step 4: Select Threshold

Choose a threshold that balances false positives and false negatives for your use case.

## Practical Examples

### Example 1: Finding Duplicates

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    pdq, _ := imghash.NewPDQ()
    
    images := []string{
        "photo1.jpg",
        "photo1_compressed.jpg",
        "photo1_resized.jpg",
        "different_photo.jpg",
    }
    
    // Compute all hashes
    hashes := make([]imghash.Hash, len(images))
    for i, img := range images {
        hashes[i], _ = imghash.HashFile(pdq, img)
    }
    
    // Compare all pairs
    duplicateThreshold := 10.0
    
    for i := 0; i < len(images); i++ {
        for j := i + 1; j < len(images); j++ {
            dist, _ := pdq.Compare(hashes[i], hashes[j])
            
            if dist <= duplicateThreshold {
                fmt.Printf("DUPLICATE: %s <-> %s (distance: %.0f)\n",
                    images[i], images[j], dist)
            }
        }
    }
}
```

### Example 2: Similarity Search

```go
package main

import (
    "fmt"
    "sort"
    "github.com/ajdnik/imghash/v2"
)

type SimilarityResult struct {
    Image    string
    Distance float64
}

func findSimilar(query string, database []string, topK int) []SimilarityResult {
    gist, _ := imghash.NewGIST()
    
    queryHash, _ := imghash.HashFile(gist, query)
    
    results := make([]SimilarityResult, 0, len(database))
    
    for _, img := range database {
        hash, _ := imghash.HashFile(gist, img)
        dist, _ := gist.Compare(queryHash, hash)
        
        results = append(results, SimilarityResult{
            Image:    img,
            Distance: float64(dist),
        })
    }
    
    // Sort by distance (ascending)
    sort.Slice(results, func(i, j int) bool {
        return results[i].Distance < results[j].Distance
    })
    
    // Return top K
    if len(results) > topK {
        results = results[:topK]
    }
    
    return results
}

func main() {
    database := []string{
        "beach1.jpg", "beach2.jpg", "mountain1.jpg",
        "city1.jpg", "beach3.jpg",
    }
    
    results := findSimilar("query_beach.jpg", database, 3)
    
    fmt.Println("Top 3 similar images:")
    for i, r := range results {
        fmt.Printf("%d. %s (distance: %.4f)\n", 
            i+1, r.Image, r.Distance)
    }
}
```

### Example 3: Multi-Algorithm Voting

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func isDuplicate(img1, img2 string) bool {
    // Use multiple algorithms for more robust detection
    
    avg, _ := imghash.NewAverage()
    h1a, _ := imghash.HashFile(avg, img1)
    h2a, _ := imghash.HashFile(avg, img2)
    distAvg, _ := avg.Compare(h1a, h2a)
    
    pdq, _ := imghash.NewPDQ()
    h1p, _ := imghash.HashFile(pdq, img1)
    h2p, _ := imghash.HashFile(pdq, img2)
    distPDQ, _ := pdq.Compare(h1p, h2p)
    
    cm, _ := imghash.NewColorMoment()
    h1c, _ := imghash.HashFile(cm, img1)
    h2c, _ := imghash.HashFile(cm, img2)
    distCM, _ := cm.Compare(h1c, h2c)
    
    // Vote: at least 2 out of 3 must agree
    votes := 0
    
    if distAvg <= 5 {
        votes++
    }
    if distPDQ <= 10 {
        votes++
    }
    if distCM <= 15 {
        votes++
    }
    
    return votes >= 2
}

func main() {
    if isDuplicate("img1.jpg", "img2.jpg") {
        fmt.Println("Images are duplicates")
    } else {
        fmt.Println("Images are different")
    }
}
```

## Distance Metric Reference

| Metric               | Best For                       | Range                 | Formula                          |                 |    |       |    |   |    |   |     |
| -------------------- | ------------------------------ | --------------------- | -------------------------------- | --------------- | -- | ----- | -- | - | -- | - | --- |
| **Hamming**          | Binary hashes                  | `0` to `hash_bits`    | Count of differing bits          |                 |    |       |    |   |    |   |     |
| **Weighted Hamming** | Binary with importance weights | `0` to `weighted_max` | Weighted count of differing bits |                 |    |       |    |   |    |   |     |
| **L1 (Manhattan)**   | Histograms, robust to outliers | `0` to `∞`            | \`Σ                              | h1\[i] - h2\[i] | \` |       |    |   |    |   |     |
| **L2 (Euclidean)**   | Continuous features            | `0` to `∞`            | `√Σ(h1[i] - h2[i])²`             |                 |    |       |    |   |    |   |     |
| **Cosine**           | Directional similarity         | `0` to `2`            | \`1 - (h1·h2)/(                  |                 | h1 |       |    |   | h2 |   | )\` |
| **Chi-Square**       | Probability distributions      | `0` to `∞`            | `Σ(h1[i]-h2[i])²/(h1[i]+h2[i])`  |                 |    |       |    |   |    |   |     |
| **Jaccard**          | Set similarity                 | `0` to `1`            | \`1 -                            | intersection    | /  | union | \` |   |    |   |     |
| **PCC**              | Correlation                    | `-1` to `1`           | Pearson correlation coefficient  |                 |    |       |    |   |    |   |     |

## Next Steps

- **[Practical Examples](examples)** — See complete examples for common use cases
- **[API Reference](similarity)** — Explore distance metric documentation
