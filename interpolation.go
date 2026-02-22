package imghash

import "github.com/ajdnik/imghash/v2/internal/imgproc"

// Interpolation specifies the resize interpolation method used during hash computation.
type Interpolation int

// Supported interpolation methods.
const (
	NearestNeighbor   Interpolation = Interpolation(imgproc.NearestNeighbor)
	Bilinear          Interpolation = Interpolation(imgproc.Bilinear)
	Bicubic           Interpolation = Interpolation(imgproc.Bicubic)
	MitchellNetravali Interpolation = Interpolation(imgproc.MitchellNetravali)
	Lanczos2          Interpolation = Interpolation(imgproc.Lanczos2)
	Lanczos3          Interpolation = Interpolation(imgproc.Lanczos3)
	BilinearExact     Interpolation = Interpolation(imgproc.BilinearExact)
)

var interpolationNames = [...]string{
	NearestNeighbor:   "NearestNeighbor",
	Bilinear:          "Bilinear",
	Bicubic:           "Bicubic",
	MitchellNetravali: "MitchellNetravali",
	Lanczos2:          "Lanczos2",
	Lanczos3:          "Lanczos3",
	BilinearExact:     "BilinearExact",
}

func (i Interpolation) valid() bool {
	return int(i) >= 0 && int(i) < len(interpolationNames)
}

// String returns the name of the interpolation method.
func (i Interpolation) String() string {
	if i.valid() {
		return interpolationNames[i]
	}
	return "Unknown"
}

func (i Interpolation) validate() error {
	if !i.valid() {
		return ErrInvalidInterpolation
	}
	return nil
}

func (i Interpolation) resizeType() imgproc.ResizeType {
	if !i.valid() {
		panic("imghash: invalid interpolation value")
	}
	return imgproc.ResizeType(i)
}
