package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/internal/imgproc"
)

// Median is a perceptual hash that uses a similar approach as Average hash.
// But instead of using mean it uses median to compute the average value.
// See https://github.com/Quickshot/DupImageLib/blob/3e914588958c4c1871d750de86b30446b9c07a3e/DupImageLib/ImageHashes.cs#L99 for more information.
type Median struct {
	// Resized image width.
	width uint
	// Resized image height.
	height uint
	// Resize interpolation method.
	interp Interpolation
}

// NewMedian creates a new Median hash with the given options.
// Without options, sensible defaults are used.
func NewMedian(opts ...MedianOption) Median {
	m := Median{
		width:  8,
		height: 8,
		interp: Bilinear,
	}
	for _, o := range opts {
		o.applyMedian(&m)
	}
	return m
}

// Calculate returns a perceptual image hash.
func (mh Median) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(mh.width, mh.height, img, mh.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	med, err := imgproc.Median(g)
	if err != nil {
		return nil, err
	}
	return thresholdHash(g, uint(math.Round(med))), nil
}
