package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/internal/imgproc"
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
	interp Interpolation
	// Gaussian kernel size.
	kernel int
	// Gaussian kernel sigma.
	sigma float64
}

// NewColorMoment creates a new ColorMoment hash with the given options.
// Without options, sensible defaults are used.
func NewColorMoment(opts ...Option) ColorMoment {
	o := options{
		width:  512,
		height: 512,
		interp: Bicubic,
		kernel: 3,
		sigma:  0,
	}
	applyOptions(&o, opts)
	return ColorMoment{
		width:  o.width,
		height: o.height,
		interp: o.interp,
		kernel: o.kernel,
		sigma:  o.sigma,
	}
}

// Calculate returns a perceptual image hash.
func (ch ColorMoment) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(ch.width, ch.height, img, ch.interp.resizeType())
	b := imgproc.GaussianBlur(r, ch.kernel, ch.sigma)
	yrb, err := imgproc.YCrCb(b)
	if err != nil {
		return nil, err
	}
	hsv, err := imgproc.HSV(b)
	if err != nil {
		return nil, err
	}
	yrbMom := imgproc.GetMoments(yrb)
	yrbMom[0], yrbMom[2] = yrbMom[2], yrbMom[0]
	hsvMom := imgproc.GetMoments(hsv)
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
	return hash, nil
}
