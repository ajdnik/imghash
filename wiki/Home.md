# imghash Wiki

imghash is a Go implementation of multiple perceptual image hashing algorithms.

## Getting Started

Install:

```sh
go get -u github.com/ajdnik/imghash/v2
```

Import:

```go
import "github.com/ajdnik/imghash/v2"
```

Quick start (PDQ):

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

## Documentation Index

- [Algorithms](Algorithms)
- [Similarity Metrics](Similarity-Metrics)
- [Convenience Functions](Convenience-Functions)
- [Interpolation Methods](Interpolation-Methods)
- [Migration Guide](Migration-Guide)

## Community

- [Code of Conduct](https://github.com/ajdnik/imghash/blob/main/CODE_OF_CONDUCT.md)
- [Contributing Guide](https://github.com/ajdnik/imghash/blob/main/CONTRIBUTING.md)
- [Security Policy](https://github.com/ajdnik/imghash/blob/main/SECURITY.md)
- [Migration Guide](Migration-Guide)
