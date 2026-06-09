<p align="center">
  <img src="assets/logo.png" alt="imghash logo" width="220">
</p>

<p align="center">Go implementation of multiple perceptual hash algorithms for images.</p>

<p align="center">
  <a href="https://github.com/ajdnik/imghash/actions?query=workflow%3Aci"><img src="https://github.com/ajdnik/imghash/workflows/ci/badge.svg" alt="CI status"></a>
  <a href="https://coveralls.io/github/ajdnik/imghash?branch=main"><img src="https://badge.coveralls.io/repos/github/ajdnik/imghash/badge.svg?branch=main" alt="Coverage status"></a>
  <a href="https://pkg.go.dev/github.com/ajdnik/imghash/v2"><img src="https://pkg.go.dev/badge/github.com/ajdnik/imghash/v2.svg" alt="Go reference"></a>
  <a href="https://goreportcard.com/report/github.com/ajdnik/imghash/v2"><img src="https://goreportcard.com/badge/github.com/ajdnik/imghash/v2" alt="Go report card"></a>
  <a href="https://github.com/ajdnik/imghash/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-lightgrey.svg" alt="MIT license"></a>
  <a href="https://www.bestpractices.dev/projects/12605"><img src="https://www.bestpractices.dev/projects/12605/badge?v=2"></a>
</p>

## Documentation

Full documentation lives on the [wiki](https://github.com/ajdnik/imghash/wiki). It covers installation, a quick start guide, a per-algorithm reference (options, defaults, references, examples) for every hasher shipped with the library, an API reference for the `Hasher` / `Comparer` interfaces and the `similarity` distance functions, conceptual guides on hash types and interpolation methods, and practical guides for choosing an algorithm, comparing hashes, real-world examples, and migrating from v1 to v2.

## Installing

```sh
go get -u github.com/ajdnik/imghash/v2
```

```go
import "github.com/ajdnik/imghash/v2"
```

Most consumers only need the top-level `imghash` package. Core types (`Hash`, `Binary`, `UInt8`, `Float64`, `Distance`) are re-exported there.

## Quick Start

If you're unsure which hash to pick, start with PDQ.

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

  h1, err := imghash.HashFile(pdq, "image1.png")
  if err != nil {
    panic(err)
  }

  h2, err := imghash.HashFile(pdq, "image2.png")
  if err != nil {
    panic(err)
  }

  dist, err := pdq.Compare(h1, h2)
  if err != nil {
    panic(err)
  }

  fmt.Printf("Distance: %v\n", dist)
}
```

## License

Imghash is released under the MIT license. See [LICENSE](LICENSE).
