package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
)

// ColorMoment is a perceptual hash that uses the method described in
// Perceptual Hashing for Color Images Using Invariant Moments; Tang et. al.
//
// See https://www.researchgate.net/publication/286870507_Perceptual_hashing_for_color_images_using_invariant_moments for more information.
type ColorMoment struct {
	// Resized image width.
	width uint
	// Resized image height.
	height uint
	// Resize interpolation method.
	interp ResizeType
	// Gaussian kernel size.
	kernel int
	// Gaussian kernel sigma.
	sigma float64
}

// NewColorMoment creates a new ColorMoment struct using default values.
func NewColorMoment() ColorMoment {
	return ColorMoment{
		width:  512,
		height: 512,
		interp: Bicubic,
		kernel: 3,
		sigma:  0,
	}
}

// NewColorMomentWithParams creates a new ColorMoment struct based on supplied parameters.
func NewColorMomentWithParams(resizeWidth, resizeHeight uint, resizeType ResizeType, kernelSize int, sigma float64) ColorMoment {
	return ColorMoment{
		width:  resizeWidth,
		height: resizeHeight,
		interp: resizeType,
		kernel: kernelSize,
		sigma:  sigma,
	}
}

// Calculate returns a perceptual image hash.
func (ch *ColorMoment) Calculate(img image.Image) hashtype.Float64 {
	r := resizeImageCV(ch.width, ch.height, img, ch.interp)
	b := imgproc.GaussianBlur(r, ch.kernel, ch.sigma)
	yrb, _ := imgproc.YCrCb(b)
	hsv, _ := imgproc.HSV(b)
	yrbMom := imgproc.GetMoments(yrb)
	// Switch R and B channels
	yrbMom[0], yrbMom[2] = yrbMom[2], yrbMom[0]
	hsvMom := imgproc.GetMoments(hsv)
	// Switch R and B channels
	hsvMom[0], hsvMom[2] = hsvMom[2], hsvMom[0]
	yHuMom := imgproc.HuMoments(yrbMom)
	hHuMom := imgproc.HuMoments(hsvMom)
	hash := make(hashtype.Float64, len(hHuMom)+len(yHuMom))
	var i int
	for i = 0; i < len(hHuMom); i++ {
		hash[i] = hHuMom[i]
	}
	for ; i < len(hHuMom)+len(yHuMom); i++ {
		hash[i] = yHuMom[i-len(hHuMom)]
	}
	return hash
}
