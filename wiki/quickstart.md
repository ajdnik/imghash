# Quick Start

> Compute your first perceptual hash in under 5 minutes

# Quick Start Guide

Learn how to compute perceptual hashes and compare images using imghash.

## Choose an Algorithm

> [!NOTE]
> If you're unsure which hash to pick, **start with PDQ**. It's Facebook's production-grade algorithm designed for robust content moderation and duplicate detection.

PDQ (PhotoDNA-Quality) is ideal for:

* Content moderation and abuse detection
* Near-duplicate image detection
* Robustness to JPEG compression, rescaling, and minor edits
* Large-scale image deduplication

## Basic Example: Compare Two Images

Here's a complete working example using the PDQ algorithm:

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create a PDQ hasher
    pdq, err := imghash.NewPDQ()
    if err != nil {
        panic(err)
    }

    // Hash two image files
    h1, err := imghash.HashFile(pdq, "image1.png")
    if err != nil {
        panic(err)
    }

    h2, err := imghash.HashFile(pdq, "image2.png")
    if err != nil {
        panic(err)
    }

    // Compare the hashes
    dist, err := pdq.Compare(h1, h2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Distance: %v\n", dist)
    
    // Interpret the result
    if dist < 32 {
        fmt.Println("Images are very similar or near-duplicates")
    } else if dist < 64 {
        fmt.Println("Images are somewhat similar")
    } else {
        fmt.Println("Images are different")
    }
}
```

## Step-by-Step Breakdown

### Step 1: Create a hasher

Instantiate the hash algorithm you want to use:

```go
pdq, err := imghash.NewPDQ()
if err != nil {
    panic(err)
}
```

PDQ uses sensible defaults (256-bit binary hash, bilinear interpolation). You can customize behavior with options:

```go
pdq, err := imghash.NewPDQ(
    imghash.WithPDQInterpolation(imghash.Bicubic),
)
```

### Step 2: Hash your images

Use the convenience function `HashFile()` to hash image files:

```go
h1, err := imghash.HashFile(pdq, "image1.png")
if err != nil {
    panic(err)
}
```

The `HashFile()` function:

* Opens the image file
* Decodes it (supports JPEG, PNG, GIF)
* Computes the perceptual hash
* Returns a `Hash` interface

For PDQ, this returns a `Binary` hash of 256 bits (32 bytes).

### Step 3: Compare the hashes

Use the algorithm's `Compare()` method to measure similarity:

```go
dist, err := pdq.Compare(h1, h2)
if err != nil {
    panic(err)
}
fmt.Printf("Distance: %v\n", dist)
```

For PDQ (and most binary hashes), this computes **Hamming distance** - the number of differing bits.

* `Distance = 0` - Identical images
* `Distance < 32` - Very similar or near-duplicates (for 256-bit hashes)
* `Distance > 128` - Completely different images

## Alternative Hash Sources

### Hash from io.Reader

For streaming or HTTP sources:

```go
import (
    "net/http"
    "github.com/ajdnik/imghash/v2"
)

resp, err := http.Get("https://example.com/image.jpg")
if err != nil {
    panic(err)
}
defer resp.Body.Close()

h, err := imghash.HashReader(pdq, resp.Body)
if err != nil {
    panic(err)
}
```

### Hash from image.Image

If you already have an `image.Image` object:

```go
import (
    "image"
    "os"
    "github.com/ajdnik/imghash/v2"
)

f, err := os.Open("photo.jpg")
if err != nil {
    panic(err)
}
defer f.Close()

img, _, err := image.Decode(f)
if err != nil {
    panic(err)
}

// Hash the image directly
pdq, _ := imghash.NewPDQ()
h, err := pdq.Calculate(img)
if err != nil {
    panic(err)
}
```

## Try Other Algorithms

<details>
<summary>Average Hash</summary>

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Simple 64-bit average hash
    avg, err := imghash.NewAverage()
    if err != nil {
        panic(err)
    }

    h1, err := imghash.HashFile(avg, "image1.png")
    if err != nil {
        panic(err)
    }

    h2, err := imghash.HashFile(avg, "image2.png")
    if err != nil {
        panic(err)
    }

    dist, err := avg.Compare(h1, h2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Hamming distance: %v\n", dist)
}
```

</details>

<details>
<summary>Difference Hash</summary>

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Gradient-based difference hash
    diff, err := imghash.NewDifference()
    if err != nil {
        panic(err)
    }

    h1, err := imghash.HashFile(diff, "image1.png")
    if err != nil {
        panic(err)
    }

    h2, err := imghash.HashFile(diff, "image2.png")
    if err != nil {
        panic(err)
    }

    dist, err := diff.Compare(h1, h2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Hamming distance: %v\n", dist)
}
```

</details>

<details>
<summary>PHash</summary>

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // DCT-based perceptual hash
    phash, err := imghash.NewPHash()
    if err != nil {
        panic(err)
    }

    h1, err := imghash.HashFile(phash, "image1.png")
    if err != nil {
        panic(err)
    }

    h2, err := imghash.HashFile(phash, "image2.png")
    if err != nil {
        panic(err)
    }

    // PHash uses weighted Hamming by default
    dist, err := phash.Compare(h1, h2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Weighted Hamming distance: %v\n", dist)
}
```

</details>

<details>
<summary>ColorMoment</summary>

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Color distribution hash (Float64)
    cm, err := imghash.NewColorMoment()
    if err != nil {
        panic(err)
    }

    h1, err := imghash.HashFile(cm, "image1.png")
    if err != nil {
        panic(err)
    }

    h2, err := imghash.HashFile(cm, "image2.png")
    if err != nil {
        panic(err)
    }

    // ColorMoment uses L2 (Euclidean) distance
    dist, err := cm.Compare(h1, h2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("L2 distance: %v\n", dist)
}
```

</details>

<details>
<summary>DINOHash</summary>

DINOHash requires the sibling `dinoweights` module for its embedded DINOv2 ViT-S/14+reg model weights (~85 MB):

```bash
go get -u github.com/ajdnik/imghash/v2/dinoweights
```

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/dinoweights"
)

func main() {
    // Deep perceptual hash backed by a frozen DINOv2 ViT-S/14+reg backbone.
    dn, err := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
    if err != nil {
        panic(err)
    }

    h1, err := imghash.HashFile(dn, "image1.png")
    if err != nil {
        panic(err)
    }

    h2, err := imghash.HashFile(dn, "image2.png")
    if err != nil {
        panic(err)
    }

    // DINOHash returns a 96-bit Binary hash; default metric is Hamming.
    dist, err := dn.Compare(h1, h2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Hamming distance: %v / 96 bits\n", dist)
}
```

> [!NOTE]
> First `Calculate` call parses the embedded ~85 MB safetensors blob and constructs the ViT model. Subsequent calls reuse it. Expect around one second per image on a modern CPU.

</details>

## Working with Hash Types

### Binary Hashes

Binary hashes store bits efficiently:

```go
pdq, _ := imghash.NewPDQ()
h, _ := imghash.HashFile(pdq, "image.png")

// Type assert to Binary
if binaryHash, ok := h.(imghash.Binary); ok {
    fmt.Printf("Hash bytes: %v\n", binaryHash)
    fmt.Printf("Length: %d bytes\n", binaryHash.Len())
    fmt.Printf("String: %s\n", binaryHash.String())
}
```

### Inspecting Hashes

All hash types implement the `Hash` interface:

```go
type Hash interface {
    String() string          // String representation
    Len() int                // Number of elements
    ValueAt(idx int) float64 // Value at index
}
```

Example:

```go
pdq, _ := imghash.NewPDQ()
h, _ := imghash.HashFile(pdq, "image.png")

fmt.Printf("Hash: %s\n", h.String())
fmt.Printf("Length: %d\n", h.Len())
fmt.Printf("First value: %v\n", h.ValueAt(0))
```

## Custom Distance Metrics

Override the default distance metric:

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

// Use L1 distance instead of Hamming for comparison
dist, err := imghash.Compare(h1, h2, similarity.L1)
if err != nil {
    panic(err)
}
```

Available distance functions:

* `similarity.Hamming` - Bit differences (Binary only)
* `similarity.L1` - Manhattan distance
* `similarity.L2` - Euclidean distance
* `similarity.Cosine` - Cosine similarity
* `similarity.ChiSquare` - Chi-square distance
* `similarity.Jaccard` - Jaccard index
* `similarity.PCC` - Pearson correlation coefficient
* `similarity.WeightedHamming` - Weighted bit differences (Binary only)

## Error Handling

Common errors to handle:

```go
pdq, err := imghash.NewPDQ()
if err != nil {
    // Constructor errors (invalid options)
    panic(err)
}

h, err := imghash.HashFile(pdq, "image.png")
if err != nil {
    // File errors: file not found, unsupported format, corrupt image
    panic(err)
}

dist, err := pdq.Compare(h1, h2)
if err != nil {
    // Comparison errors: incompatible hash types, length mismatch
    panic(err)
}
```

Key error types:

* `imghash.ErrIncompatibleHash` - Comparing incompatible hash types
* `imghash.ErrHashLengthMismatch` - Hash lengths don't match
* `imghash.ErrInvalidSize` - Invalid dimension in constructor
* `imghash.ErrInvalidInterpolation` - Invalid interpolation method

## Performance Tips

> [!NOTE]
> **Reuse hashers**: Create the hasher once and reuse it for multiple images. Hashers are safe to use concurrently.
>
> ```go
> // Good - create once, use many times
> pdq, _ := imghash.NewPDQ()
> for _, imagePath := range images {
>     h, _ := imghash.HashFile(pdq, imagePath)
>     // process hash
> }
>
> // Bad - creates new hasher every iteration
> for _, imagePath := range images {
>     pdq, _ := imghash.NewPDQ()  // Wasteful!
>     h, _ := imghash.HashFile(pdq, imagePath)
> }
> ```

## Complete Example: Find Similar Images

Here's a practical example that finds similar images in a directory:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create PDQ hasher
    pdq, err := imghash.NewPDQ()
    if err != nil {
        panic(err)
    }

    // Hash all images in directory
    hashes := make(map[string]imghash.Hash)
    err = filepath.Walk("./images", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        
        // Try to hash the file
        h, err := imghash.HashFile(pdq, path)
        if err != nil {
            // Skip non-image files
            return nil
        }
        
        hashes[path] = h
        return nil
    })
    if err != nil {
        panic(err)
    }

    // Find similar pairs
    const threshold = 32 // Distance threshold for "similar"
    
    paths := make([]string, 0, len(hashes))
    for path := range hashes {
        paths = append(paths, path)
    }

    for i := 0; i < len(paths); i++ {
        for j := i + 1; j < len(paths); j++ {
            dist, err := pdq.Compare(hashes[paths[i]], hashes[paths[j]])
            if err != nil {
                continue
            }
            
            if dist < threshold {
                fmt.Printf("Similar images (distance=%v):\n", dist)
                fmt.Printf("  %s\n", paths[i])
                fmt.Printf("  %s\n", paths[j])
            }
        }
    }
}
```

## Next Steps

Now that you can compute and compare hashes:

* Learn about algorithm-specific options and customization
* Explore different algorithms for your use case
* Understand distance metrics and thresholds
* Build a production duplicate detection system
* Integrate with databases for large-scale search
