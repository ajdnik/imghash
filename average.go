// Package imghash provides perceptual image hashing algorithms and
// similarity metrics for comparing images by visual content.
package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// Average is a perceptual hash that uses the method described in Looks Like It by Dr. Neal Krawetz.
//
// See https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html for more information.
type Average struct {
	baseConfig
	distFunc DistanceFunc
}

// NewAverage creates a new Average hash with the given options.
// Without options, sensible defaults are used.
func NewAverage(opts ...AverageOption) (Average, error) {
	a := Average{
		baseConfig: baseConfig{width: 8, height: 8, interp: Bilinear},
	}
	for _, o := range opts {
		o.applyAverage(&a)
	}
	if a.width == 0 || a.height == 0 {
		return Average{}, ErrInvalidSize
	}
	return a, nil
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
	return thresholdHash(g, uint(math.Round(m)))
}

// Compare computes the Hamming distance between two Average hashes.
func (ah Average) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateBinaryCompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if ah.distFunc != nil {
		return ah.distFunc(h1, h2)
	}
	return similarity.Hamming(h1, h2)
}
