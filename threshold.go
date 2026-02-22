package imghash

import (
	"image"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// thresholdHash builds a binary hash by setting a bit for every pixel
// whose intensity exceeds the given threshold. Used by Average and Median.
func thresholdHash(img *image.Gray, threshold uint) (hashtype.Binary, error) {
	bnds := img.Bounds()
	bits := uint(bnds.Dx() * bnds.Dy())
	hash := hashtype.NewBinary(bits)
	var c uint
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			if uint(img.GrayAt(x, y).Y) > threshold {
				if err := hash.Set(c); err != nil {
					return nil, err
				}
			}
			c++
		}
	}
	return hash, nil
}
