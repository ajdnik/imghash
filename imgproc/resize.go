package imgproc

import (
	"image"

	"github.com/nfnt/resize"
)

type ResizeType int

const (
	NearestNeighbor ResizeType = iota
	Bilinear
	Bicubic
	MitchellNetravali
	Lanczos2
	Lanczos3
	BilinearExact
)

func Resize(width, height uint, img image.Image, typ ResizeType) image.Image {
	var interp resize.InterpolationFunction
	switch typ {
	case NearestNeighbor:
		interp = resize.NearestNeighbor
	case Bilinear:
		interp = resize.Bilinear
	case Bicubic:
		interp = resize.Bicubic
	case MitchellNetravali:
		interp = resize.MitchellNetravali
	case Lanczos2:
		interp = resize.Lanczos2
	case Lanczos3:
		interp = resize.Lanczos3
	case BilinearExact:
		interp = resize.Bilinear
	default:
		interp = resize.NearestNeighbor
	}
	return resize.Resize(width, height, img, interp)
}
