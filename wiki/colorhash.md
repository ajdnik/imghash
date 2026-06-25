# ColorHash

> Binary perceptual hash encoding color distribution into 14 hue bins.

## Overview

ColorHash classifies each pixel as black, gray, faint-color, or bright-color based on HSV saturation and value, then accumulates a 14-bin histogram. Bin counts are encoded into a binary hash. Compatible with [Johannes Buchner's Python imagehash.colorhash](https://github.com/JohannesBuchner/imagehash). Compares using Hamming distance.

## Pixel Categorization

- **Black:** intensity < 32
- **Gray:** saturation < 85
- **Faint color:** saturation between 85 and 170, hue split into 6 bins
- **Bright color:** saturation > 170, hue split into 6 bins

Bin layout: 1 black + 1 gray + 6 faint hue + 6 bright hue = 14 bins. Each bin uses `binBits` bits, producing `14 * binBits` total bits.

## Constructor

```go
func NewColorHash(opts ...ColorHashOption) (ColorHash, error)
```

## Available Options

- `WithBinBits(n uint)` — bits per histogram bin
- `WithDistance(fn DistanceFunc)` — override comparison metric

## Default Settings

- **Hash type:** `Binary` (14 * 3 = 42 bits at default `binBits=3`)
- **Default metric:** Hamming distance
- **BinBits:** 3

## Usage Example

```go
ch, err := imghash.NewColorHash()
if err != nil {
    panic(err)
}
hash, err := imghash.HashFile(ch, "image.jpg")
```

## References

- [Johannes Buchner's Python imagehash.colorhash](https://github.com/JohannesBuchner/imagehash)
