package imghash

import (
	"image"

	"gocv.io/x/gocv"
)

func resizeImageCV(width, height uint, img image.Image, resizeType ResizeType) image.Image {
	var interp gocv.InterpolationFlags
	switch resizeType {
	case NearestNeighbor:
		interp = gocv.InterpolationNearestNeighbor
	case Bilinear:
		interp = gocv.InterpolationLinear
	case Bicubic:
		interp = gocv.InterpolationCubic
	case MitchellNetravali:
		interp = gocv.InterpolationNearestNeighbor
	case Lanczos2:
		interp = gocv.InterpolationNearestNeighbor
	case Lanczos3:
		interp = gocv.InterpolationNearestNeighbor
	case BilinearExact:
		interp = 5
	default:
		interp = gocv.InterpolationNearestNeighbor
	}
	var imgMat gocv.Mat
	switch i := img.(type) {
	case *image.Gray:
		imgMat, _ = gocv.ImageGrayToMatGray(i)
	default:
		imgMat, _ = gocv.ImageToMatRGB(i)
	}
	dst := gocv.NewMat()
	gocv.Resize(imgMat, &dst, image.Point{int(width), int(height)}, 0, 0, interp)
	res, err := dst.ToImage()
	if err != nil {
		panic(err)
	}
	return res
}

func ReadImageCV(file string) (image.Image, error) {
	mat := gocv.IMRead(file, gocv.IMReadColor)
	return mat.ToImage()
}

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
