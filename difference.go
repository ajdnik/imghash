package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
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
	interp imgproc.ResizeType
}

// NewDifference creates a new Difference struct using default values.
func NewDifference() Difference {
	return Difference{
		width:  8,
		height: 8,
		interp: imgproc.Bilinear,
	}
}

// NewDifferenceWithParams creates a new Difference struct based on supplied parameters.
func NewDifferenceWithParams(resizeWidth, resizeHeight uint, resizeType imgproc.ResizeType) Difference {
	return Difference{
		width:  resizeWidth,
		height: resizeHeight,
		interp: resizeType,
	}
}

// Calculate returns a perceptual image hash.
func (dh *Difference) Calculate(img image.Image) hashtype.Binary {
	r := imgproc.Resize(dh.width+1, dh.height, img, dh.interp)
	g, _ := imgproc.Grayscale(r)
	return dh.computeHash(g)
}

// Computes the binary hash based on the gradients in the resized image.
func (dh *Difference) computeHash(img *image.Gray) hashtype.Binary {
	size := dh.width * dh.height / 8
	hash := make(hashtype.Binary, size)
	bnds := img.Bounds()
	var c uint
	for i := bnds.Min.Y; i < bnds.Max.Y; i++ {
		for j := bnds.Min.X + 1; j < bnds.Max.X; j++ {
			lft := img.GrayAt(j-1, i).Y
			pix := img.GrayAt(j, i).Y
			if pix > lft {
				hash.Set(c)
			}
			c++
		}
	}
	return hash
}
