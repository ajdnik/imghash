package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
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
	interp imgproc.ResizeType
}

// NewMedian creates a new Median struct using default values.
func NewMedian() Median {
	return Median{
		width:  8,
		height: 8,
		interp: imgproc.Bilinear,
	}
}

// NewMedianWithParams creates a new Median struct using the supplied parameters.
func NewMedianWithParams(resizeWidth, resizeHeight uint, resizeType imgproc.ResizeType) Median {
	return Median{
		width:  resizeWidth,
		height: resizeHeight,
		interp: resizeType,
	}
}

// Calculate returns a perceptual image hash.
func (mh *Median) Calculate(img image.Image) hashtype.Binary {
	r := imgproc.Resize(mh.width, mh.height, img, mh.interp)
	g, _ := imgproc.Grayscale(r)
	med, _ := imgproc.Median(g)
	return mh.computeHash(g, uint(math.Round(med)))
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
				hash.Set(c)
			}
			c++
		}
	}
	return hash
}
