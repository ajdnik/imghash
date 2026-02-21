package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/internal/imgproc"
)

// Difference is a perceptual hash that uses the method described in Kinf of Like That by Dr. Neal Krawetz.
//
// See https://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html for more information.
type Difference struct {
	// Resized image width.
	width uint
	// Resized image height.
	height uint
	// Resize interpolation method.
	interp Interpolation
}

// NewDifference creates a new Difference hash with the given options.
// Without options, sensible defaults are used.
func NewDifference(opts ...DifferenceOption) (Difference, error) {
	d := Difference{
		width:  8,
		height: 8,
		interp: Bilinear,
	}
	for _, o := range opts {
		o.applyDifference(&d)
	}
	if d.width == 0 || d.height == 0 {
		return Difference{}, ErrInvalidSize
	}
	return d, nil
}

// Calculate returns a perceptual image hash.
func (dh Difference) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(dh.width+1, dh.height, img, dh.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	return dh.computeHash(g), nil
}

// Computes the binary hash based on the gradients in the resized image.
func (dh Difference) computeHash(img *image.Gray) hashtype.Binary {
	size := dh.width * dh.height / 8
	hash := make(hashtype.Binary, size)
	bnds := img.Bounds()
	var c uint
	for i := bnds.Min.Y; i < bnds.Max.Y; i++ {
		for j := bnds.Min.X + 1; j < bnds.Max.X; j++ {
			lft := img.GrayAt(j-1, i).Y
			pix := img.GrayAt(j, i).Y
			if pix > lft {
				_ = hash.Set(c)
			}
			c++
		}
	}
	return hash
}
