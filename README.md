# imghash

[![Build Status](https://github.com/ajdnik/imghash/workflows/tests-and-coverage/badge.svg "GitHub Actions status")](https://github.com/ajdnik/imghash/actions?query=workflow%3Atests-and-coverage)
[![Coverage Status](https://coveralls.io/repos/github/ajdnik/imghash/badge.svg?branch=master)](https://coveralls.io/github/ajdnik/imghash?branch=master)
[![GoDoc](https://godoc.org/github.com/ajdnik/imghash?status.svg "GoDoc")](https://godoc.org/github.com/ajdnik/imghash)
[![Go Report Card](https://goreportcard.com/badge/github.com/ajdnik/imghash)](https://goreportcard.com/report/github.com/ajdnik/imghash)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/ajdnik/imghash/blob/master/LICENSE)

Go implementation of multiple perceptual hash algorithms for images. Perceptual hash functions are analogous if features are similar, whereas cryptographic hashing relies on the avalanche effect of a small change in input value creating a drastic change in output value. 

## Installing

Using imghash is easy. First, use `go get` to install the latest version
of the library. This command will install the library and its dependencies:

    go get -u github.com/ajdnik/imghash

Next, include imghash in your application:

```go
import "github.com/ajdnik/imghash"
```

## Example

In the example below, we read an image from the filesystem and decode it. Afterwards we use the pHash to compute its perceptual hash. Next we load another image and compute its pHash, and lastly, we compare both hashes to get a difference score.

```go
package main

import (
  "fmt"
  "image"
	_ "image/png"
	"os"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/similarity"
)

func main() {
  // Open the first image
  fimg1, err := os.Open("image1.png")
  if err != nil {
    panic(err)
  }
  defer fimg1.Close()

  // Decode the image
  img1, _, err := image.Decode(fimg1)
  if err != nil {
    panic(err)
  }

  // Create pHash object
  phash := imghash.NewPHash()
  
  // Compute the hash of the image
  h1 := phash.Calculate(img1)
  fmt.Printf("First hash: %v\n", h1)

  // Open the second image
  fimg2, err := os.Open("image2.png")
  if err != nil {
    panic(err)
  }
  defer fimg2.Close()

  // Decode the second image
  img2, _, err := image.Decode(fimg2)
  if err != nil {
    panic(err)
  }

  // Compute the hash of the second image
  h2 := phash.Calculate(img2)
  fmt.Printf("Second hash: %v\n", h2)

  // Compute hash similarity score
  d := similarity.Hamming(h1, h2)
  fmt.Printf("Hash similarity: %v\n", d)
}
```

## License

Imghash is released under the MIT license. See [LICENSE](https://github.com/ajdnik/imghash/blob/master/LICENSE)
