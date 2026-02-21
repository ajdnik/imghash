package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
)

// MarrHildreth is a perceptual hash that uses the method described in
// Implementation and Benchmarking of Perceptual Image Hash Functions; Zauner et. al.
//
// See https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking for more information.
type MarrHildreth struct {
	// Scale parameter, used to compute Marr-Hildreth kernel.
	scale float64
	// Alpha parameter, used to compute Marr-Hildreth kernel.
	alpha float64
	// Resized image width.
	width uint
	// Resized image height.
	height uint
	// Resize interpolation method.
	interp imgproc.ResizeType
	// Gaussian kernel size.
	kernel int
	// Gaussian kernel sigma parameter.
	sigma float64

	kernels [][]float32
}

// NewMarrHildreth creates a new MarrHildreth hash with the given options.
// Without options, sensible defaults are used.
func NewMarrHildreth(opts ...Option) MarrHildreth {
	o := options{
		scale:  1,
		alpha:  2,
		width:  512,
		height: 512,
		interp: imgproc.Bicubic,
		kernel: 7,
		sigma:  0,
	}
	applyOptions(&o, opts)
	mh := MarrHildreth{
		scale:  o.scale,
		alpha:  o.alpha,
		width:  o.width,
		height: o.height,
		interp: o.interp,
		kernel: o.kernel,
		sigma:  o.sigma,
	}
	mh.kernels = computeMarrHildrethKernel(mh.alpha, mh.scale)
	return mh
}

// Calculate returns a perceptual image hash.
func (mhh *MarrHildreth) Calculate(img image.Image) (hashtype.Hash, error) {
	g, err := imgproc.Grayscale(img)
	if err != nil {
		return nil, err
	}
	b := imgproc.GaussianBlur(g, mhh.kernel, mhh.sigma)
	r := imgproc.Resize(mhh.width, mhh.height, b, mhh.interp)
	eq := imgproc.EqualizeHist(r.(*image.Gray))
	f := imgproc.Filter2DGray(eq, mhh.kernels)
	blks := mhh.blocksSum(f)
	return mhh.createHash(blks), nil
}

// Compute sums of blocks.
// TODO: Remove all magic numbers.
func (mhh *MarrHildreth) blocksSum(img [][]float32) [][]float32 {
	blocks := make([][]float32, 31)
	for r := 0; r < 31; r++ {
		blocks[r] = make([]float32, 31)
		for c := 0; c < 31; c++ {
			var sum float32
			for roiR := r * 16; roiR < r*16+16; roiR++ {
				for roiC := c * 16; roiC < c*16+16; roiC++ {
					sum += img[roiR][roiC]
				}
			}
			blocks[r][c] = sum
		}
	}
	return blocks
}

// Compute binary hash from block sums.
// TODO: Remove all magic numbers.
func (mhh *MarrHildreth) createHash(blocks [][]float32) hashtype.Binary {
	hash := make(hashtype.Binary, 72)
	var count uint
	for r := 0; r < 29; r += 4 {
		for c := 0; c < 29; c += 4 {
			var sum float32
			for i := r; i < r+3; i++ {
				for j := c; j < c+3; j++ {
					sum += blocks[j][i]
				}
			}
			avg := sum / 9.0
			for i := r; i < r+3; i++ {
				for j := c; j < c+3; j++ {
					if blocks[j][i] > avg {
						_ = hash.SetReverse(count)
					}
					count++
				}
			}
		}
	}
	return hash
}

// Compute 2D Marr-Hildreth kernel.
func computeMarrHildrethKernel(alpha, level float64) [][]float32 {
	sigma := int(4 * math.Pow(alpha, level))
	ratio := float32(math.Pow(alpha, -level))
	dim := 2*sigma + 1
	kernel := make([][]float32, dim)
	for i := range kernel {
		kernel[i] = make([]float32, dim)
		ydiff := float32(i - sigma)
		ypos := ratio * ydiff
		yposPow2 := ypos * ypos
		for j := range kernel[i] {
			xpos := ratio * float32(j-sigma)
			a := float64(xpos*xpos + yposPow2)
			kernel[i][j] = float32((2 - a) * math.Exp(a/-2))
		}
	}
	return kernel
}
