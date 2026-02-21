package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/internal/imgproc"
)

const (
	mhBlockSize = 16 // pixel block size for summation
	mhNumBlocks = 31 // blocks per dimension (default 512/16 - 1)
	mhSubBlock  = 3  // sub-block dimension for hash construction
	mhStride    = 4  // step between hash sub-blocks
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
	interp Interpolation
	// Gaussian kernel size.
	kernel int
	// Gaussian kernel sigma parameter.
	sigma float64

	kernels [][]float32
}

// NewMarrHildreth creates a new MarrHildreth hash with the given options.
// Without options, sensible defaults are used.
func NewMarrHildreth(opts ...MarrHildrethOption) (MarrHildreth, error) {
	mh := MarrHildreth{
		scale:  1,
		alpha:  2,
		width:  512,
		height: 512,
		interp: Bicubic,
		kernel: 7,
		sigma:  0,
	}
	for _, o := range opts {
		o.applyMarrHildreth(&mh)
	}
	if mh.width == 0 || mh.height == 0 {
		return MarrHildreth{}, ErrInvalidSize
	}
	if mh.scale <= 0 {
		return MarrHildreth{}, ErrInvalidScale
	}
	if mh.alpha <= 0 {
		return MarrHildreth{}, ErrInvalidAlpha
	}
	if mh.kernel <= 0 {
		return MarrHildreth{}, ErrInvalidKernelSize
	}
	if mh.sigma < 0 {
		return MarrHildreth{}, ErrInvalidSigma
	}
	mh.kernels = computeMarrHildrethKernel(mh.alpha, mh.scale)
	return mh, nil
}

// Calculate returns a perceptual image hash.
func (mhh MarrHildreth) Calculate(img image.Image) (hashtype.Hash, error) {
	g, err := imgproc.Grayscale(img)
	if err != nil {
		return nil, err
	}
	b := imgproc.GaussianBlur(g, mhh.kernel, mhh.sigma)
	r := imgproc.Resize(mhh.width, mhh.height, b, mhh.interp.resizeType())
	eq := imgproc.EqualizeHist(r.(*image.Gray))
	f := imgproc.Filter2DGray(eq, mhh.kernels)
	blks := mhh.blocksSum(f)
	return mhh.createHash(blks), nil
}

// Compute sums of blocks.
func (mhh MarrHildreth) blocksSum(img [][]float32) [][]float32 {
	blocks := make([][]float32, mhNumBlocks)
	for r := 0; r < mhNumBlocks; r++ {
		blocks[r] = make([]float32, mhNumBlocks)
		for c := 0; c < mhNumBlocks; c++ {
			var sum float32
			for roiR := r * mhBlockSize; roiR < r*mhBlockSize+mhBlockSize; roiR++ {
				for roiC := c * mhBlockSize; roiC < c*mhBlockSize+mhBlockSize; roiC++ {
					sum += img[roiR][roiC]
				}
			}
			blocks[r][c] = sum
		}
	}
	return blocks
}

// Compute binary hash from block sums.
func (mhh MarrHildreth) createHash(blocks [][]float32) hashtype.Binary {
	gridLimit := mhNumBlocks - mhSubBlock + 1
	gridSteps := (gridLimit-1)/mhStride + 1
	hashBits := gridSteps * gridSteps * mhSubBlock * mhSubBlock
	hash := make(hashtype.Binary, hashBits/8)
	var count uint
	for r := 0; r < gridLimit; r += mhStride {
		for c := 0; c < gridLimit; c += mhStride {
			var sum float32
			for i := r; i < r+mhSubBlock; i++ {
				for j := c; j < c+mhSubBlock; j++ {
					sum += blocks[j][i]
				}
			}
			avg := sum / float32(mhSubBlock*mhSubBlock)
			for i := r; i < r+mhSubBlock; i++ {
				for j := c; j < c+mhSubBlock; j++ {
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
