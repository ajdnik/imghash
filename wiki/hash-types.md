# Hash Types

> Type definitions for representing perceptual hashes

## Overview

Imghash uses three concrete hash types to represent different kinds of perceptual hashes, all implementing the common `Hash` interface.

## Hash Interface

The base interface implemented by all hash types:

```go
type Hash interface {
    String() string
    Len() int
    ValueAt(idx int) float64
}
```

### Methods

- **`String`** `func() string` — Returns a string representation of the hash.
  - **`string`** `string` — A formatted string representation of the hash data.
- **`Len`** `func() int` — Returns the length of the hash.
  - **`length`** `int`
    * For `Binary`: number of bytes
    * For `UInt8`: number of uint8 elements
    * For `Float64`: number of float64 elements
- **`ValueAt`** `func(idx int) float64` — Returns the value at the given index as a float64.
  - **`idx`** `int` — The zero-based index of the element to retrieve.
  - **`value`** `float64` — The element value converted to float64.

## Binary

Represents a hash where the smallest element is a bit. Used by most traditional perceptual hash algorithms.

```go
type Binary []byte
```

### Constructor

- **`NewBinary`** `func(bits uint) Binary` — Allocates a binary hash with enough bytes to store the specified number of bits.
  - **`bits`** `uint` — The number of bits the hash should store.
  - **`hash`** `Binary` — A newly allocated binary hash with `(bits+7)/8` bytes.

### Methods

- **`Set`** `func(position uint) error` — Turns a bit on in the binary hash.
  - **`position`** `uint` — The bit position to set (0-indexed).
  - **`error`** `error` — Returns `ErrOutOfBounds` if position is beyond the hash size.
- **`SetReverse`** `func(position uint) error` — Turns a bit on in the binary hash, counting position in reverse order within each byte.
  - **`position`** `uint` — The bit position to set (0-indexed, reversed within bytes).
  - **`error`** `error` — Returns `ErrOutOfBounds` if position is beyond the hash size.
- **`Equal`** `func(bh Binary) bool` — Checks if two binary hashes are identical.
  - **`bh`** `Binary` — The binary hash to compare against.
  - **`equal`** `bool` — Returns `true` if the hashes are byte-for-byte identical.

### Usage Example

```go
import "github.com/ajdnik/imghash/v2/hashtype"

// Create a binary hash for 64 bits
hash := hashtype.NewBinary(64)

// Set some bits
hash.Set(0)  // Set bit 0
hash.Set(5)  // Set bit 5
hash.Set(63) // Set bit 63

// Get info
fmt.Printf("Hash: %v\n", hash)        // [65 0 0 0 0 0 0 128]
fmt.Printf("Length: %d bytes\n", hash.Len()) // 8

// Access values
val := hash.ValueAt(0) // Returns 65.0 (first byte as float64)

// Compare hashes
hash2 := hashtype.NewBinary(64)
if hash.Equal(hash2) {
    fmt.Println("Hashes are identical")
}
```

## UInt8

Represents a hash where the smallest element is a uint8 value. Used by histogram-based algorithms.

```go
type UInt8 []uint8
```

### Methods

- **`Equal`** `func(uh UInt8) bool` — Checks if two UInt8 hashes are identical.
  - **`uh`** `UInt8` — The UInt8 hash to compare against.
  - **`equal`** `bool` — Returns `true` if both hashes have the same length and all elements match.

### Usage Example

```go
import "github.com/ajdnik/imghash/v2/hashtype"

// Create a UInt8 hash
hash := hashtype.UInt8{23, 45, 67, 89, 12, 34, 56, 78}

// Get info
fmt.Printf("Hash: %v\n", hash)            // [23 45 67 89 12 34 56 78]
fmt.Printf("Length: %d elements\n", hash.Len()) // 8

// Access values
val := hash.ValueAt(0) // Returns 23.0

// Compare hashes
hash2 := hashtype.UInt8{23, 45, 67, 89, 12, 34, 56, 78}
if hash.Equal(hash2) {
    fmt.Println("Hashes are identical")
}
```

## Float64

Represents a hash where the smallest element is a float64. Used by feature-based algorithms.

```go
type Float64 []float64
```

### Methods

- **`Equal`** `func(fh Float64) bool` — Checks if two Float64 hashes are identical using epsilon-based comparison.
  - **`fh`** `Float64` — The Float64 hash to compare against.
  - **`equal`** `bool` — Returns `true` if both hashes have the same length and all elements are within machine epsilon of each other.

### Usage Example

```go
import "github.com/ajdnik/imghash/v2/hashtype"

// Create a Float64 hash
hash := hashtype.Float64{0.123, 0.456, 0.789, 0.234, 0.567}

// Get info
fmt.Printf("Hash: %v\n", hash)            // [0.123 0.456 0.789 0.234 0.567]
fmt.Printf("Length: %d elements\n", hash.Len()) // 5

// Access values
val := hash.ValueAt(0) // Returns 0.123

// Compare hashes (uses epsilon comparison)
hash2 := hashtype.Float64{0.123, 0.456, 0.789, 0.234, 0.567}
if hash.Equal(hash2) {
    fmt.Println("Hashes are identical")
}
```

## Distance Type

Represents a similarity measure between two hashes:

```go
type Distance float64
```

### Methods

- **`Equal`** `func(dst Distance) bool` — Checks if two distances are equal using epsilon comparison.
  - **`dst`** `Distance` — The distance to compare against.
  - **`equal`** `bool` — Returns `true` if distances are within 1e-12 of each other.

### Usage Example

```go
import "github.com/ajdnik/imghash/v2/similarity"

dist1 := similarity.Distance(5.5)
dist2 := similarity.Distance(5.5000000000001)

if dist1.Equal(dist2) {
    fmt.Println("Distances are equal (within epsilon)")
}
```

## Errors

```go
// Returned when bit position is out of bounds
var ErrOutOfBounds = errors.New("position out of bounds")

// Returned when comparing incompatible hash types
var ErrIncompatibleHash = errors.New("incompatible hash type for distance calculation")

// Returned when hash lengths don't match
var ErrHashLengthMismatch = errors.New("imghash: hash lengths must match")
```

## Type Assertions

You can use type assertions to work with specific hash types:

```go
func processHash(h imghash.Hash) {
    switch v := h.(type) {
    case hashtype.Binary:
        fmt.Printf("Binary hash with %d bytes\n", v.Len())
    case hashtype.UInt8:
        fmt.Printf("UInt8 hash with %d elements\n", v.Len())
    case hashtype.Float64:
        fmt.Printf("Float64 hash with %d elements\n", v.Len())
    }
}
```

## See Also

* [Hasher Interface](hasher) - For computing hashes
* [Similarity Functions](similarity) - For comparing hashes
