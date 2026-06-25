# Hash Representations

> Understanding Binary, UInt8, and Float64 hash representations in imghash

## Overview

imghash uses three distinct hash types to represent perceptual hashes. Each type is optimized for different algorithms and comparison methods. All types implement the common `Hash` interface.

```go
type Hash interface {
    String() string
    Len() int
    ValueAt(idx int) float64
}
```

## Binary

Binary hashes represent image fingerprints where each element is a single bit. This is the most compact representation and is used by bit-based algorithms like Average, Difference, Median, and PHash.

### Type Definition

```go
type Binary []byte
```

Binary is a byte slice where bits are packed into bytes. A 64-bit hash uses 8 bytes of storage.

### Creating Binary Hashes

```go
// Allocate a binary hash with 64 bits
hash := hashtype.NewBinary(64)

// Allocate a binary hash with 256 bits
large := hashtype.NewBinary(256) // Uses 32 bytes
```

> [!NOTE]
> The `NewBinary` function automatically calculates the required number of bytes using `(bits+7)/8` to ensure proper storage.

### Setting Bits

Binary hashes provide two methods for setting bits:

#### Set (LSB first)

```go
hash := hashtype.NewBinary(16)
hash.Set(0)  // Sets bit 0 → [1 0]
hash.Set(15) // Sets bit 15 → [1 128]
fmt.Println(hash.String())
// Output: [1 128]
```

Bits are numbered with the least significant bit (LSB) first within each byte.

#### SetReverse (MSB first)

```go
hash := hashtype.NewBinary(16)
hash.SetReverse(0)  // Sets bit 0 (MSB) → [128 0]
hash.SetReverse(15) // Sets bit 15 (MSB) → [128 1]
fmt.Println(hash.String())
// Output: [128 1]
```

Bits are numbered with the most significant bit (MSB) first within each byte.

### Methods

- **`String()`** `string` — Returns a string representation formatted as a byte array.

  ```go
  hash := hashtype.Binary{115, 247, 1}
  fmt.Println(hash.String())
  // Output: [115 247 1]
  ```

- **`Len()`** `int` — Returns the number of bytes in the hash.

  ```go
  hash := hashtype.NewBinary(64)
  fmt.Println(hash.Len()) // 8
  ```

- **`ValueAt(idx)`** `float64` — Returns the byte at the given index as a float64.

  ```go
  hash := hashtype.Binary{0, 128, 255}
  fmt.Println(hash.ValueAt(1)) // 128.0
  ```

- **`Equal(other)`** `bool` — Checks if two binary hashes are identical.

  ```go
  h1 := hashtype.Binary{1, 2, 128}
  h2 := hashtype.Binary{1, 2, 128}
  h3 := hashtype.Binary{2, 128}

  fmt.Println(h1.Equal(h2)) // true
  fmt.Println(h1.Equal(h3)) // false
  ```

### Algorithms Using Binary

* **Average** - Compares pixels to mean value
* **Difference** - Compares adjacent pixels
* **Median** - Compares pixels to median value
* **PHash** - DCT-based perceptual hash
* **MarrHildreth** - Edge detection based
* **BlockMean** - Block-based averaging
* **WHash** - Wavelet-based hash
* **PDQ** - Facebook's robust hashing algorithm
* **DINOHash** - DINOv2 ViT-S/14+reg deep perceptual hash (96 bits, adversarial-robust)

> [!WARNING]
> Binary hashes can only be compared using `Hamming`, `WeightedHamming`, or `Jaccard` distance metrics. Other metrics will return `ErrIncompatibleHash`.

## UInt8

UInt8 hashes represent image fingerprints where each element is an unsigned 8-bit integer (0-255). Used by histogram-based and quantized algorithms.

### Type Definition

```go
type UInt8 []uint8
```

### Example

```go
// Create from algorithm output
algo, _ := imghash.NewColorMoment()
hash, _ := algo.Calculate(img)

// Cast to UInt8 to access methods
if u8, ok := hash.(imghash.UInt8); ok {
    fmt.Println(u8.String())
    fmt.Println("Length:", u8.Len())
}
```

### Methods

- **`String()`** `string` — Returns a string representation as a uint8 slice.

  ```go
  hash := hashtype.UInt8{10, 127, 255, 0}
  fmt.Println(hash.String())
  // Output: [10 127 255 0]
  ```

- **`Len()`** `int` — Returns the number of uint8 elements in the hash.

- **`ValueAt(idx)`** `float64` — Returns the element at the given index as float64.

  ```go
  hash := hashtype.UInt8{10, 127, 255}
  fmt.Println(hash.ValueAt(1)) // 127.0
  ```

- **`Equal(other)`** `bool` — Checks element-wise equality. Returns false if lengths differ.

  ```go
  h1 := hashtype.UInt8{10, 20, 30}
  h2 := hashtype.UInt8{10, 20, 30}
  h3 := hashtype.UInt8{10, 20}

  fmt.Println(h1.Equal(h2)) // true
  fmt.Println(h1.Equal(h3)) // false (different lengths)
  ```

### Algorithms Using UInt8

* **ColorMoment** - Color distribution moments
* **CLD** (Color Layout Descriptor) - Spatial color distribution
* **EHD** (Edge Histogram Descriptor) - Edge distribution
* **LBP** (Local Binary Patterns) - Texture descriptor
* **HOGHash** - Histogram of Oriented Gradients
* **BoVW** - Bag of Visual Words with integer storage

### Compatible Metrics

UInt8 hashes work with: `L1`, `L2`, `Cosine`, `ChiSquare`, `Jaccard` (MinHash mode)

## Float64

Float64 hashes represent image fingerprints as floating-point vectors. These provide the highest precision and are used by advanced feature descriptors and moment-based algorithms.

### Type Definition

```go
type Float64 []float64
```

### Example

```go
// Create from algorithm output
algo, _ := imghash.NewZernike()
hash, _ := algo.Calculate(img)

// Cast to Float64
if f64, ok := hash.(imghash.Float64); ok {
    fmt.Println("Hash dimension:", f64.Len())
    fmt.Println("First element:", f64.ValueAt(0))
}
```

### Methods

- **`String()`** `string` — Returns a string representation as a float64 slice.

  ```go
  hash := hashtype.Float64{0.5, 1.234, -0.75}
  fmt.Println(hash.String())
  // Output: [0.5 1.234 -0.75]
  ```

- **`Len()`** `int` — Returns the number of float64 elements in the hash.

- **`ValueAt(idx)`** `float64` — Returns the element at the given index.

  ```go
  hash := hashtype.Float64{0.1, 0.5, 0.9}
  fmt.Println(hash.ValueAt(1)) // 0.5
  ```

- **`Equal(other)`** `bool` — Checks equality using epsilon-based comparison for floating-point tolerance.

  ```go
  h1 := hashtype.Float64{0.1, 0.2, 0.3}
  h2 := hashtype.Float64{0.1, 0.2, 0.3}
  h3 := hashtype.Float64{0.1, 0.2, 0.30001}

  fmt.Println(h1.Equal(h2)) // true
  fmt.Println(h1.Equal(h3)) // false (exceeds epsilon)
  ```

> [!NOTE]
> Float64 uses `epsilon = math.Nextafter(1.0, 2.0) - 1.0` (machine epsilon) for equality comparison.

### Algorithms Using Float64

* **RadialVariance** - Radial projection features
* **Zernike** - Zernike moments
* **GIST** - Global scene descriptor
* **RASH** - Radon transform based
* **BoVW** - Bag of Visual Words with float storage

### Compatible Metrics

Float64 hashes work with: `L1`, `L2`, `Cosine`, `ChiSquare`, `PCC`, `Jaccard` (MinHash mode)

## Hash Interface

All three types implement the common `Hash` interface, allowing generic comparison functions:

```go
func compareHashes(h1, h2 imghash.Hash) {
    // Works with Binary, UInt8, or Float64
    fmt.Printf("Hash 1 length: %d\n", h1.Len())
    fmt.Printf("Hash 2 length: %d\n", h2.Len())
    
    // Access elements generically
    for i := 0; i < h1.Len(); i++ {
        fmt.Printf("Element %d: %f\n", i, h1.ValueAt(i))
    }
}
```

## Type Compatibility

> [!WARNING]
> Most similarity metrics require both hashes to be the same type. Comparing Binary with UInt8 will return `ErrIncompatibleHash`.

```go
binary := hashtype.Binary{1, 2, 3}
uint8 := hashtype.UInt8{1, 2, 3}

// This will error
dist, err := similarity.Hamming(binary, uint8)
if err != nil {
    fmt.Println(err) // incompatible hash type for distance calculation
}

// This works
dist, err := similarity.Hamming(binary, binary)
fmt.Println(dist) // 0.0
```

## Best Practices

- **Use Binary for Speed** — Binary hashes are fastest for comparison and use minimal memory. Ideal for large-scale image matching.
- **Use UInt8 for Histograms** — UInt8 is perfect for color and edge histograms where you need 256 bins of granularity.
- **Use Float64 for Precision** — Float64 provides maximum accuracy for moment-based and feature descriptors requiring sub-integer precision.
- **Match Hash Types** — Always ensure both hashes are the same type before comparison to avoid runtime errors.

## Error Handling

```go
// ErrOutOfBounds - Binary.Set() position out of range
hash := hashtype.NewBinary(8)
err := hash.Set(100) // Error: position exceeds 8 bits

// ErrIncompatibleHash - Type mismatch in comparison
binary := hashtype.Binary{1, 2}
float := hashtype.Float64{1.0, 2.0}
_, err = similarity.Hamming(binary, float)
// Error: incompatible hash type for distance calculation
```

## Related

* [Similarity Metrics](similarity-metrics) - Learn how to compare different hash types
* [Choosing an Algorithm](choosing-algorithm) - Choose the right algorithm and hash type
* [Hash Types API](hash-types) - Complete hashtype package documentation
