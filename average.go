package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/internal/imgproc"
)

// Average is a perceptual hash that uses the method described in Looks Like It by Dr. Neal Krawetz.
//
// See https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html for more information.
type Average struct {
	// Resized image width.
	width uint
	// Resized image height.
	height uint
	// Resize interpolation method.
	interp Interpolation
}

// NewAverage creates a new Average hash with the given options.
// Without options, sensible defaults are used.
func NewAverage(opts ...AverageOption) Average {
	a := Average{
		width:  8,
		height: 8,
		interp: Bilinear,
	}
	for _, o := range opts {
		o.applyAverage(&a)
	}
	return a
}

// Calculate returns a perceptual image hash.
func (ah Average) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(ah.width, ah.height, img, ah.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	m, err := imgproc.Mean(g)
	if err != nil {
		return nil, err
	}
	return thresholdHash(g, uint(math.Round(m))), nil
}
