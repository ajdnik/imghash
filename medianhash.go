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
	// Resize interpoletion method.
	interp Interpolation
}

// NewMedian creates a new Median hash with the given options.
// Without options, sensible defaults are used.
func NewMedian(opts ...Option) Median {
	o := options{
		width:  8,
		height: 8,
		interp: Bilinear,
	}
	applyOptions(&o, opts)
	return Median{
		width:  o.width,
		height: o.height,
		interp: o.interp,
	}
}

// Calculate returns a perceptual image hash.
func (mh *Median) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(mh.width, mh.height, img, mh.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	med, err := imgproc.Median(g)
	if err != nil {
		return nil, err
	}
	return mh.computeHash(g, uint(math.Round(med))), nil
}

// Computes the binary hash based on the median value of the resized image.
func (mh *Median) computeHash(img *image.Gray, median uint) hashtype.Binary {
	size := mh.width * mh.height / 8
	hash := make(hashtype.Binary, size)
	bnds := img.Bounds()
	var c uint
	for i := bnds.Min.Y; i < bnds.Max.Y; i++ {
		for j := bnds.Min.X; j < bnds.Max.X; j++ {
			pix := img.GrayAt(j, i).Y
			if uint(pix) > median {
				_ = hash.Set(c)
			}
			c++
		}
	}
	return hash
}
