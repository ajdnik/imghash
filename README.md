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

## Perceptual Hash Algorithms

The library supports 8 different perceptual hashing algorithms which were ported from [OpenCV Contrib](https://github.com/opencv/opencv_contrib) and consequentally tested against its OpenCV implementations. 

#### Average Hash

One of the more basic perceptual hashing algorithms it produces a 64 bit binary hash which can be compared using the Hamming distance. The algorithm crushes the image into a grayscale 8x8 image and sets the 64 bits in the hash based on whether the pixel's value is greater than the average color for the image. You can read more about the implementation in the [Looks Like It](https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html) article.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  resizeWidth := 8
  resizeHeight := 8
  resizeType := imgproc.Bilinear

  ahash := imghash.NewAverageWithParams(resizeWidth, resizeHeight, resizeType)
  res := ahash.Calculate(img)
  fmt.Printf("Average Hash: %v\n", res)
}
```

If you want to create an Average hash with default parameters you can call the `NewAverage` instead of `NewAverageWithParams`.

#### Difference Hash
Difference hash is very similar in implementation then the average hash approach. Instead of computing averages it computes gradients instead. You can find a better explanation in the [Kind Of Like That](https://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html) article.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  resizeWidth := 8
  resizeHeight := 8
  resizeType := imgproc.Bilinear

  dhash := imghash.NewDifferenceWithParams(resizeWidth, resizeHeight, resizeType)
  res := dhash.Calculate(img)
  fmt.Printf("Difference Hash: %v\n", res)
}
```

If you want to create a Difference hash with default parameters you can call the `NewDifference` instead of `NewDifferenceWithParams`.

#### Median Hash

The hash is almost identical to the average hash. Instead of using the average metric it uses the median metric when computing the hash bits.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  resizeWidth := 8
  resizeHeight := 8
  resizeType := imgproc.Bilinear

  mhash := imghash.NewMedianWithParams(resizeWidth, resizeHeight, resizeType)
  res := mhash.Calculate(img)
  fmt.Printf("Median Hash: %v\n", res)
}
```

If you want to create a Median hash with default parameters you can call the `NewMedian` instead of `NewMedianWithParams`.

#### Color Moments Hash

The algorithm is an implementation of the [Perceptual hashing for color images using invariant moments](https://www.researchgate.net/publication/286870507_Perceptual_hashing_for_color_images_using_invariant_moments) article, which describes using color moments from YCbCr and HSV colorspaces to build a vector based perceptual hash. The vector consists of 42 floating point values.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  resizeWidth := 512
  resizeHeight := 512
  resizeType := imgproc.Bicubic
  gaussianKernelSize := 3
  gaussianKernelSigma := 0

  cmhash := imghash.NewColorMomentWithParams(resizeWidth, resizeHeight, resizeType, gaussianKernelSize, gaussianKernelSigma)
  res := cmhash.Calculate(img)
  fmt.Printf("Color Moment Hash: %v\n", res)
}
```

If you want to create a Color Moments hash with default parameters you can call the `NewColorMoment` instead of `NewColorMomentWithParams`.

#### Marr-Hildreth Hash

The algorithm is an implementation of the Marr-Hildreth Operator Based Hash from [Rihamark: Perceptual image hash benchmarking](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking). It uses a 2D wavelet transform to compute a 568 bit binary hash.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  resizeWidth := 512
  resizeHeight := 512
  resizeType := imgproc.Bicubic
  gaussianKernelSize := 7
  gaussianKernelSigma := 0
  marrHildrethScale := 1
  marrHildrethAlpha := 2

  mhhash := imghash.NewMarrHildrethWithParams(marrHildrethScale, marrHildrethAlpha, resizeWidth, resizeHeight, resizeType, gaussianKernelSize, gaussianKernelSigma)
  res := cmhash.Calculate(img)
  fmt.Printf("Marr Hildreth Hash: %v\n", res)
}
```

If you want to create a Marr-Hildreth hash with default parameters you can call the `NewMarrHildreth` instead of `NewMarrHildrethWithParams`.

#### Block Mean Value Hash
The algorithm is based on the implementation described in [Block Mean Value Based Image Perceptual Hashing](https://ieeexplore.ieee.org/document/4041692). It uses a similar approach as the [Average Hash]() but instead of resizing the image it computes means using a sliding window and builds the hash based on those means an a global mean value. The result is a 256 bit binary hash.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  resizeWidth := 256
  resizeHeight := 256
  resizeType := imgproc.BilinearExact
  blockWidth := 16
  blockHeight := 16
  method := imghash.Direct

  bmhash := imghash.NewBlockMeanWithParams(resizeWidth, resizeHeight, resizeType, blockWidth, blockHeight, method)
  res := bmhash.Calculate(img)
  fmt.Printf("Block Mean Hash: %v\n", res)
}
```

If you want to create a Block Mean hash with default parameters you can call the `NewBlockMean` instead of `NewBlockMeanWithParams`.

#### pHash
The pHash algorithm is implemented from the [Rihamark: Perceptual image hash benchmarking](https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking). It uses a discrete cosine transform to extract the relevant channels from the image and computes the hash based on the 8x8 DCT result. The hash is a vector of floating point values.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  resizeWidth := 32
  resizeHeight := 32
  resizeType := imgproc.BilinearExact

  phash := imghash.NewPHashWithParams(resizeWidth, resizeHeight, resizeType)
  res := phash.Calculate(img)
  fmt.Printf("pHash: %v\n", res)
}
```

If you want to create a pHash with default parameters you can call the `NewPHash` instead of `NewPHashWithParams`.

#### Radial Variance Hash
The algorithm is implemented from the article [Robust image hashing based on radial variance of pixels](https://www.researchgate.net/publication/4186555_Robust_image_hashing_based_on_radial_variance_of_pixels). The algorithm uses a combination of radial projections and feature vector computations to generate a vector based hash containing floating point values.

```go
package main

import (
  "fmt"

  "github.com/ajdnik/imghash"
  "github.com/ajdnik/imghash/imgproc"
)

func main() {
  img := ReadAndDecodeImage()

  // Default hash parameters
  gaussianKernelSigma := 1
  numberOfAngles := 180

  rvhash := imghash.NewRadialVarianceWithParams(gaussianKernelSigma, numberOfAngles)
  res := rvhash.Calculate(img)
  fmt.Printf("Radial Variance Hash: %v\n", res)
}
```

If you want to create a Radial Variance hash with default parameters you can call the `NewRadialVariance` instead of `NewRadialVarianceWithParams`.

## License

Imghash is released under the MIT license. See [LICENSE](https://github.com/ajdnik/imghash/blob/master/LICENSE)
