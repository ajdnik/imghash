# Hasher Interface

> Core interface for computing perceptual hashes from images

## Hasher

The `Hasher` interface defines the contract for computing perceptual hashes from images. All hash algorithms in imghash implement this interface.

```go
type Hasher interface {
    Calculate(image.Image) (hashtype.Hash, error)
}
```

### Methods

- **`Calculate`** `func(image.Image) (hashtype.Hash, error)` — Computes a perceptual hash from the provided image.
  - **`img`** `image.Image` — The input image to hash. Accepts any type that implements the standard `image.Image` interface.
  - **`hash`** `hashtype.Hash` — The computed hash. The concrete type depends on the algorithm:
    * Binary hash algorithms (Average, Difference, Median, PHash, etc.) return `hashtype.Binary`
    * Histogram-based algorithms (CLD, EHD, etc.) return `hashtype.UInt8`
    * Feature-based algorithms (ColorMoment, GIST, etc.) return `hashtype.Float64`
  - **`error`** `error` — Returns an error if the hash computation fails. Possible errors include invalid image dimensions or algorithm-specific failures.

## Implementing Algorithms

All hash algorithms in imghash implement the `Hasher` interface:

* **Average** - Average hash algorithm
* **Difference** - Difference hash algorithm
* **Median** - Median hash algorithm
* **PHash** - Perceptual hash using DCT
* **BlockMean** - Block mean hash
* **MarrHildreth** - Marr-Hildreth edge detection hash
* **RadialVariance** - Radial variance hash
* **ColorMoment** - Color moment hash
* **CLD** - Color Layout Descriptor
* **EHD** - Edge Histogram Descriptor
* **WHash** - Wavelet hash
* **LBP** - Local Binary Pattern hash
* **HOGHash** - Histogram of Oriented Gradients hash
* **BoVW** - Bag of Visual Words hash
* **PDQ** - Facebook's PDQ hash
* **RASH** - Robust Affine-invariant Spatial Hash
* **Zernike** - Zernike moments hash
* **GIST** - GIST descriptor hash
* **DINOHash** - DINOv2 ViT-S/14+reg deep perceptual hash (96 bits, adversarial-robust)

## Usage Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Open an image
    img, err := imghash.OpenImage("photo.jpg")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a hasher (Average hash in this example)
    hasher := imghash.Average{}
    
    // Calculate the hash
    hash, err := hasher.Calculate(img)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Hash: %v\n", hash)
}
```

## Using Different Algorithms

```go
// Use Difference hash
diffHasher := imghash.Difference{}
hash1, _ := diffHasher.Calculate(img)

// Use PHash with custom options
phasher, _ := imghash.NewPHash(
    imghash.WithPHashSize(32, 32),
    imghash.WithPHashInterpolation(imghash.Lanczos),
)
hash2, _ := phasher.Calculate(img)

// Use ColorMoment for color images
colorHasher, _ := imghash.NewColorMoment(
    imghash.WithColorMomentSize(128, 128),
)
hash3, _ := colorHasher.Calculate(img)
```

## Polymorphic Usage

```go
// Function that works with any hasher
func processImage(img image.Image, hasher imghash.Hasher) (imghash.Hash, error) {
    return hasher.Calculate(img)
}

// Can be called with any hash algorithm
hash1, _ := processImage(img, imghash.Average{})
hash2, _ := processImage(img, imghash.Difference{})
phasher, _ := imghash.NewPHash()
hash3, _ := processImage(img, phasher)
```

## See Also

* [Comparer Interface](api/comparer) - For comparing hashes
* [Hash Types](api/hash-types) - Hash type definitions
* [Convenience Functions](api/convenience) - Helper functions for common tasks
