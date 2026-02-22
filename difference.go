package imghash

import (
	"image"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// Difference is a perceptual hash that uses the method described in Kind of Like That by Dr. Neal Krawetz.
//
// See https://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html for more information.
type Difference struct {
	baseConfig
	distFunc DistanceFunc
}

// NewDifference creates a new Difference hash with the given options.
// Without options, sensible defaults are used.
func NewDifference(opts ...DifferenceOption) (Difference, error) {
	d := Difference{
		baseConfig: baseConfig{width: 8, height: 8, interp: Bilinear},
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
	return dh.computeHash(g)
}

// Computes the binary hash based on the gradients in the resized image.
func (dh Difference) computeHash(img *image.Gray) (hashtype.Binary, error) {
	hash := hashtype.NewBinary(dh.width * dh.height)
	bnds := img.Bounds()
	var c uint
	for i := bnds.Min.Y; i < bnds.Max.Y; i++ {
		for j := bnds.Min.X + 1; j < bnds.Max.X; j++ {
			lft := img.GrayAt(j-1, i).Y
			pix := img.GrayAt(j, i).Y
			if pix > lft {
				if err := hash.Set(c); err != nil {
					return nil, err
				}
			}
			c++
		}
	}
	return hash, nil
}

// Compare computes the Hamming distance between two Difference hashes.
func (dh Difference) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateBinaryCompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if dh.distFunc != nil {
		return dh.distFunc(h1, h2)
	}
	return similarity.Hamming(h1, h2)
}
