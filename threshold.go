package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
)

// thresholdHash builds a binary hash by setting a bit for every pixel
// whose intensity exceeds the given threshold. Used by Average and Median.
func thresholdHash(img *image.Gray, threshold uint) hashtype.Binary {
	bnds := img.Bounds()
	size := bnds.Dx() * bnds.Dy() / 8
	hash := make(hashtype.Binary, size)
	var c uint
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			if uint(img.GrayAt(x, y).Y) > threshold {
				_ = hash.Set(c)
			}
			c++
		}
	}
	return hash
}
