package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
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
	interp ResizeType
}

// NewAverage creates a new Average struct using default values.
func NewAverage() Average {
	return Average{
		width:  8,
		height: 8,
		interp: Bilinear,
	}

}

// NewAverageWithParams creates a new Average struct based on supplied parameters.
func NewAverageWithParams(resizeWidth, resizeHeight uint, resizeType ResizeType) Average {
	return Average{
		width:  resizeWidth,
		height: resizeHeight,
		interp: resizeType,
	}
}

// Calculate returns a perceptual image hash.
func (ah *Average) Calculate(img image.Image) hashtype.Binary {
	r := resizeImageCV(ah.width, ah.height, img, ah.interp)
	g, _ := imgproc.Grayscale(r)
	m, _ := imgproc.Mean(g)
	return ah.computeHash(g, uint(math.Round(m)))
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
				hash.Set(c)
			}
			c++
		}
	}
	return hash
}
