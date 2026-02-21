package imghash

import "github.com/ajdnik/imghash/imgproc"

// Interpolation specifies the resize interpolation method used during hash computation.
type Interpolation int

const (
	NearestNeighbor   Interpolation = Interpolation(imgproc.NearestNeighbor)
	Bilinear          Interpolation = Interpolation(imgproc.Bilinear)
	Bicubic           Interpolation = Interpolation(imgproc.Bicubic)
	MitchellNetravali Interpolation = Interpolation(imgproc.MitchellNetravali)
	Lanczos2          Interpolation = Interpolation(imgproc.Lanczos2)
	Lanczos3          Interpolation = Interpolation(imgproc.Lanczos3)
	BilinearExact     Interpolation = Interpolation(imgproc.BilinearExact)
)

func (i Interpolation) resizeType() imgproc.ResizeType {
	return imgproc.ResizeType(i)
}
