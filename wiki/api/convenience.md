# Convenience Functions

> Helper functions for common image hashing operations

## Overview

Imghash provides several convenience functions that simplify common operations like opening images, decoding from readers, and computing hashes directly from files.

## Image Loading Functions

### OpenImage

Reads and decodes an image from a file path.

```go
func OpenImage(path string) (image.Image, error)
```

- **`path`** `string` — The file path to the image. Supports JPEG, PNG, and GIF formats.
- **`image`** `image.Image` — The decoded image.
- **`error`** `error` — Returns an error if the file cannot be opened or decoded.

#### Example

```go
import "github.com/ajdnik/imghash/v2"

// Open an image file
img, err := imghash.OpenImage("photos/vacation.jpg")
if err != nil {
    log.Fatal(err)
}

// Use the image for hashing
hasher := imghash.Average{}
hash, _ := hasher.Calculate(img)
```

### DecodeImage

Decodes an image from an `io.Reader`.

```go
func DecodeImage(r io.Reader) (image.Image, error)
```

- **`r`** `io.Reader` — A reader containing image data. Supports JPEG, PNG, and GIF formats.
- **`image`** `image.Image` — The decoded image.
- **`error`** `error` — Returns an error if the image cannot be decoded.

#### Example

```go
import (
    "bytes"
    "io"
    "net/http"
    "github.com/ajdnik/imghash/v2"
)

// Decode from HTTP response
resp, _ := http.Get("https://example.com/image.jpg")
defer resp.Body.Close()

img, err := imghash.DecodeImage(resp.Body)
if err != nil {
    log.Fatal(err)
}

// Decode from bytes
data := []byte{...} // image bytes
img2, err := imghash.DecodeImage(bytes.NewReader(data))
```

## Hash Computation Functions

### HashFile

Opens an image file and computes its hash in a single operation.

```go
func HashFile(hasher Hasher, path string) (hashtype.Hash, error)
```

- **`hasher`** `Hasher` — The hash algorithm to use. Any type implementing the `Hasher` interface.
- **`path`** `string` — The file path to the image.
- **`hash`** `hashtype.Hash` — The computed hash.
- **`error`** `error` — Returns an error if the file cannot be opened, decoded, or hashed.

#### Example

```go
import "github.com/ajdnik/imghash/v2"

// Hash a file directly
hash, err := imghash.HashFile(imghash.Average{}, "photo.jpg")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Hash: %v\n", hash)

// Compare multiple files
hash1, _ := imghash.HashFile(imghash.Average{}, "photo1.jpg")
hash2, _ := imghash.HashFile(imghash.Average{}, "photo2.jpg")

algo := imghash.Average{}
dist, _ := algo.Compare(hash1, hash2)
fmt.Printf("Distance: %v\n", dist)
```

### HashReader

Decodes an image from a reader and computes its hash in a single operation.

```go
func HashReader(hasher Hasher, r io.Reader) (hashtype.Hash, error)
```

- **`hasher`** `Hasher` — The hash algorithm to use. Any type implementing the `Hasher` interface.
- **`r`** `io.Reader` — A reader containing image data.
- **`hash`** `hashtype.Hash` — The computed hash.
- **`error`** `error` — Returns an error if the image cannot be decoded or hashed.

#### Example

```go
import (
    "net/http"
    "github.com/ajdnik/imghash/v2"
)

// Hash from HTTP response
resp, _ := http.Get("https://example.com/image.jpg")
defer resp.Body.Close()

hash, err := imghash.HashReader(imghash.Average{}, resp.Body)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Hash: %v\n", hash)
```

## Hash Comparison Function

### Compare

Computes the distance between two hashes with optional custom distance function.

```go
func Compare(h1, h2 hashtype.Hash, fn ...DistanceFunc) (similarity.Distance, error)
```

- **`h1`** `hashtype.Hash` — The first hash to compare.
- **`h2`** `hashtype.Hash` — The second hash to compare.
- **`fn`** `...DistanceFunc` — Optional distance function. If not provided, uses the natural metric:
  * Hamming distance for Binary hashes
  * L2 (Euclidean) distance for UInt8 and Float64 hashes
- **`distance`** `similarity.Distance` — The computed distance. Lower values indicate more similar hashes.
- **`error`** `error` — Returns `ErrIncompatibleHash` if binary and non-binary hashes are mixed without a compatible distance function.

#### Example

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

// Calculate hashes
avg := imghash.Average{}
hash1, _ := avg.Calculate(img1)
hash2, _ := avg.Calculate(img2)

// Use default metric (Hamming for binary hashes)
dist1, _ := imghash.Compare(hash1, hash2)
fmt.Printf("Hamming distance: %v\n", dist1)

// Use custom metric
dist2, _ := imghash.Compare(hash1, hash2, similarity.L1)
fmt.Printf("L1 distance: %v\n", dist2)

// Compare with different metrics
hamming, _ := imghash.Compare(hash1, hash2, similarity.Hamming)
l2, _ := imghash.Compare(hash1, hash2, similarity.L2)
cosine, _ := imghash.Compare(hash1, hash2, similarity.Cosine)
```

## Complete Workflow Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/similarity"
)

func main() {
    // Method 1: Manual workflow
    img, err := imghash.OpenImage("photo.jpg")
    if err != nil {
        log.Fatal(err)
    }
    
    hasher := imghash.Average{}
    hash, err := hasher.Calculate(img)
    if err != nil {
        log.Fatal(err)
    }
    
    // Method 2: Convenient one-liner
    hash2, err := imghash.HashFile(imghash.Average{}, "photo2.jpg")
    if err != nil {
        log.Fatal(err)
    }
    
    // Compare using default metric
    dist, err := imghash.Compare(hash, hash2)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Distance: %v\n", dist)
    
    // Compare using custom metric
    customDist, _ := imghash.Compare(hash, hash2, similarity.L2)
    fmt.Printf("Custom distance: %v\n", customDist)
}
```

## Processing Multiple Files

```go
import (
    "fmt"
    "path/filepath"
    "github.com/ajdnik/imghash/v2"
)

func hashDirectory(dir string) (map[string]imghash.Hash, error) {
    hashes := make(map[string]imghash.Hash)
    hasher := imghash.Average{}
    
    files, _ := filepath.Glob(filepath.Join(dir, "*.jpg"))
    
    for _, file := range files {
        hash, err := imghash.HashFile(hasher, file)
        if err != nil {
            return nil, err
        }
        hashes[file] = hash
    }
    
    return hashes, nil
}

func findSimilar(hashes map[string]imghash.Hash, threshold float64) {
    algo := imghash.Average{}
    
    for path1, hash1 := range hashes {
        for path2, hash2 := range hashes {
            if path1 >= path2 {
                continue
            }
            
            dist, _ := algo.Compare(hash1, hash2)
            if float64(dist) < threshold {
                fmt.Printf("%s and %s are similar (distance: %v)\n", 
                    path1, path2, dist)
            }
        }
    }
}
```

## See Also

* [Hasher Interface](api/hasher) - Hash algorithm interface
* [Comparer Interface](api/comparer) - Hash comparison interface
* [Similarity Functions](api/similarity) - Distance metrics
