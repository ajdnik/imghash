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
func NewAverage(opts ...Option) Average {
	o := options{
		width:  8,
		height: 8,
		interp: Bilinear,
	}
	applyOptions(&o, opts)
	return Average{
		width:  o.width,
		height: o.height,
		interp: o.interp,
	}
}

// Calculate returns a perceptual image hash.
func (ah *Average) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(ah.width, ah.height, img, ah.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	m, err := imgproc.Mean(g)
	if err != nil {
		return nil, err
	}
	return ah.computeHash(g, uint(math.Round(m))), nil
}

// Computes the binary hash based on the average value of resized image.
func (ah *Average) computeHash(img *image.Gray, mean uint) hashtype.Binary {
	size := ah.width * ah.height / 8
	hash := make(hashtype.Binary, size)
	bnds := img.Bounds()
	var c uint
	for i := bnds.Min.Y; i < bnds.Max.Y; i++ {
		for j := bnds.Min.X; j < bnds.Max.X; j++ {
			pix := img.GrayAt(j, i).Y
			if uint(pix) > mean {
				_ = hash.Set(c)
			}
			c++
		}
	}
	return hash
}
