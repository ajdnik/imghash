# Installation

> Install imghash and set up your Go project

# Installation

Get started with imghash in your Go project with a single command.

## Requirements

> [!NOTE]
> imghash requires **Go 1.25 or later**. Check your version with `go version`.

The library has minimal dependencies:

* `golang.org/x/image` - Extended image format support
* `gonum.org/v1/gonum` - Pure-Go BLAS used by the DINOHash ViT forward pass

Both are pure Go; there is no CGO requirement.

## Install with go get

### Step 1: Install the package

Run the following command in your project directory:

```bash
go get -u github.com/ajdnik/imghash/v2
```

This downloads the latest version of imghash v2 and adds it to your `go.mod` file.

### Step 2: Import the package

Add the import to your Go source files:

```go
import "github.com/ajdnik/imghash/v2"
```

Most consumers only need the top-level `imghash` package. Core types (`Hash`, `Binary`, `UInt8`, `Float64`, `Distance`) are re-exported there.

### Step 3: Verify installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    pdq, err := imghash.NewPDQ()
    if err != nil {
        panic(err)
    }
    fmt.Printf("PDQ hasher created successfully: %T\n", pdq)
}
```

Run it:

```bash
go run test.go
```

You should see:

```
PDQ hasher created successfully: imghash.PDQ
```

## Package Structure

imghash is organized into several packages:

### Top-Level Package

The main `imghash` package exports all algorithms and core types:

```go
import "github.com/ajdnik/imghash/v2"

// Hash algorithms
pdq, _ := imghash.NewPDQ()
avg, _ := imghash.NewAverage()
phash, _ := imghash.NewPHash()
// ... and more

// Re-exported core types
var h imghash.Hash      // Common hash interface
var b imghash.Binary    // Bit-level hash
var u imghash.UInt8     // Byte-level hash
var f imghash.Float64   // Float-level hash
var d imghash.Distance  // Distance measure
```

### Subpackages (Advanced Usage)

For advanced use cases, you can import specific subpackages:

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/v2/hashtype"    // Hash type implementations
    "github.com/ajdnik/imghash/v2/similarity"  // Distance functions
)

// Use distance functions directly
dist, _ := similarity.Hamming(hash1, hash2)
dist, _ = similarity.L2(hash1, hash2)
dist, _ = similarity.Cosine(hash1, hash2)

// Create hash types manually
binary := hashtype.NewBinary(256) // 256-bit hash
```

> [!NOTE]
> Most users won't need to import subpackages. The main `imghash` package re-exports everything you need.

## DINOHash Weights Module

DINOHash ships its model weights in a sibling Go module so importers of `imghash` do not pay the ~85 MB embed cost unless they actually use DINOHash. The `dinoweights` module is a pure-data package: it exports a single `Blob []byte` variable embedded via `go:embed` and has zero non-stdlib dependencies.

### Install dinoweights

```bash
go get -u github.com/ajdnik/imghash/dinoweights
```

This adds the module to your `go.mod` and pulls down the embedded safetensors blob (~85 MB) once.

### Import and Wire to DINOHash

```go
import (
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/dinoweights"
)

d, err := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
if err != nil {
    panic(err)
}

hash, err := imghash.HashFile(d, "image.jpg")
```

> [!NOTE]
> If you do not import `dinoweights`, you can still use every other algorithm in `imghash`. `NewDINOHash` will only fail at the first `Calculate` call if no weights source was configured via `WithSafetensorsBlob` or `WithDINOWeights`.

### Custom Weights

For alternate checkpoints, networked weights, or non-safetensors formats, implement the `WeightsProvider` interface and pass it via `WithDINOWeights`:

```go
type WeightsProvider interface {
    Tensors() (map[string]Tensor, error)
}

d, err := imghash.NewDINOHash(imghash.WithDINOWeights(myProvider))
```

See the [DINOHash algorithm page](algorithms/dinohash) for the canonical tensor name layout.

## Supported Image Formats

imghash automatically registers decoders for common image formats:

* **JPEG** - `.jpg`, `.jpeg`
* **PNG** - `.png`
* **GIF** - `.gif`

These are registered via standard library imports:

```go
import (
    _ "image/gif"  // register GIF decoder
    _ "image/jpeg" // register JPEG decoder
    _ "image/png"  // register PNG decoder
)
```

### Adding More Formats

To support additional formats like WebP or TIFF, import their decoders:

```go
import (
    "github.com/ajdnik/imghash/v2"
    _ "golang.org/x/image/webp"  // WebP support
    _ "golang.org/x/image/tiff"  // TIFF support
)
```

## Module Configuration

Your `go.mod` file should include:

```go
module your-project

go 1.25

require github.com/ajdnik/imghash/v2 v2.x.x
```

Transitive dependencies are added automatically:

```go
require (
    github.com/ajdnik/imghash/v2 v2.x.x
    golang.org/x/image v0.41.0 // indirect
    gonum.org/v1/gonum v0.17.0 // indirect
)
```

If you also use DINOHash, your `go.mod` will include the `dinoweights` module:

```go
require (
    github.com/ajdnik/imghash/v2 v2.x.x
    github.com/ajdnik/imghash/dinoweights v1.x.x
    golang.org/x/image v0.41.0 // indirect
    gonum.org/v1/gonum v0.17.0 // indirect
)
```

The two modules are versioned independently; the `dinoweights` tag scheme is `dinoweights/vX.Y.Z`.

## Version Pinning

To use a specific version:

```bash
go get github.com/ajdnik/imghash/v2@v2.1.0
```

To update to the latest version:

```bash
go get -u github.com/ajdnik/imghash/v2
go mod tidy
```

## Troubleshooting

<details>
<summary>Error: package github.com/ajdnik/imghash/v2 is not in GOROOT</summary>

This means the package isn't downloaded. Run:

```bash
go mod download
go mod tidy
```

</details>

<details>
<summary>Error: unknown revision v2</summary>

Make sure you're using `/v2` in the import path for version 2:

```go
import "github.com/ajdnik/imghash/v2"  // Correct
import "github.com/ajdnik/imghash"     // Wrong - this is v1
```

</details>

<details>
<summary>Error: unsupported image format</summary>

The image format isn't registered. Add the appropriate decoder import:

```go
import _ "golang.org/x/image/webp"
```

</details>

<details>
<summary>Go version too old</summary>

imghash v2 requires Go 1.25+. Upgrade Go or use an older version of imghash:

```bash
go get github.com/ajdnik/imghash@v1.x.x
```

</details>

## Next Steps

Now that you have imghash installed, proceed to the [Quick Start Guide](quickstart) to compute your first perceptual hash.
