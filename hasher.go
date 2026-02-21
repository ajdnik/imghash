package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
)

// Hasher computes a perceptual hash from an image.
// It is implemented by all hash algorithms in this package:
// Average, Difference, PHash, Median, BlockMean, MarrHildreth,
// RadialVariance, and ColorMoment.
type Hasher interface {
	Calculate(image.Image) hashtype.Hash
}
