package imghash

import (
	"image"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// ColorMoment is a perceptual hash that uses the method described in
// Perceptual Hashing for Color Images Using Invariant Moments; Tang et. al.
//
// See https://www.researchgate.net/publication/286870507_Perceptual_hashing_for_color_images_using_invariant_moments for more information.
type ColorMoment struct {
	baseConfig
	// Gaussian kernel size.
	kernel int
	// Gaussian kernel sigma.
	sigma    float64
	distFunc DistanceFunc
}

// NewColorMoment creates a new ColorMoment hash with the given options.
// Without options, sensible defaults are used.
func NewColorMoment(opts ...ColorMomentOption) (ColorMoment, error) {
	c := ColorMoment{
		baseConfig: baseConfig{width: 512, height: 512, interp: Bicubic},
		kernel:     3,
		sigma:      0,
	}
	for _, o := range opts {
		o.applyColorMoment(&c)
	}
	if c.width == 0 || c.height == 0 {
		return ColorMoment{}, ErrInvalidSize
	}
	if c.kernel <= 0 {
		return ColorMoment{}, ErrInvalidKernelSize
	}
	if c.sigma < 0 {
		return ColorMoment{}, ErrInvalidSigma
	}
	return c, nil
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

// Compare computes the L2 (Euclidean) distance between two ColorMoment hashes.
func (ch ColorMoment) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateFloat64CompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if ch.distFunc != nil {
		return ch.distFunc(h1, h2)
	}
	return similarity.L2(h1, h2)
}
