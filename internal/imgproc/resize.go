package imgproc

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

// ResizeType selects the interpolation method used during image resizing.
type ResizeType int

// Supported interpolation methods for image resizing.
const (
	NearestNeighbor ResizeType = iota
	Bilinear
	Bicubic
	MitchellNetravali
	Lanczos2
	Lanczos3
	BilinearExact
)

// Resize scales img to the given dimensions using the specified interpolation.
func Resize(width, height uint, img image.Image, typ ResizeType) image.Image {
	dr := image.Rect(0, 0, int(width), int(height))
	sr := img.Bounds()
	interp := interpolatorFor(typ)

	var dst draw.Image
	switch img.(type) {
	case *image.Gray:
		dst = image.NewGray(dr)
	default:
		dst = image.NewRGBA(dr)
	}

	interp.Scale(dst, dr, img, sr, draw.Src, nil)
	return dst
}

var (
	mitchellNetravali = &draw.Kernel{
		Support: 2,
		At:      mitchellNetravaliAt,
	}
	lanczos2 = &draw.Kernel{
		Support: 2,
		At:      lanczosAt(2),
	}
	lanczos3 = &draw.Kernel{
		Support: 3,
		At:      lanczosAt(3),
	}
)

func interpolatorFor(typ ResizeType) draw.Interpolator {
	switch typ {
	case NearestNeighbor:
		return draw.NearestNeighbor
	case Bilinear, BilinearExact:
		return draw.BiLinear
	case Bicubic:
		return draw.CatmullRom
	case MitchellNetravali:
		return mitchellNetravali
	case Lanczos2:
		return lanczos2
	case Lanczos3:
		return lanczos3
	default:
		return draw.NearestNeighbor
	}
}

// Mitchell-Netravali kernel (B=1/3, C=1/3).
func mitchellNetravaliAt(t float64) float64 {
	const b, c = 1.0 / 3.0, 1.0 / 3.0
	t = math.Abs(t)
	if t < 1 {
		return ((12-9*b-6*c)*t*t*t + (-18+12*b+6*c)*t*t + (6 - 2*b)) / 6
	}
	if t < 2 {
		return ((-b-6*c)*t*t*t + (6*b+30*c)*t*t + (-12*b-48*c)*t + (8*b + 24*c)) / 6
	}
	return 0
}

// lanczosAt returns a Lanczos kernel function with the given support.
func lanczosAt(a float64) func(float64) float64 {
	return func(t float64) float64 {
		t = math.Abs(t)
		if t == 0 {
			return 1
		}
		if t >= a {
			return 0
		}
		pt := math.Pi * t
		return (math.Sin(pt) / pt) * (math.Sin(pt/a) / (pt / a))
	}
}
