package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/internal/imgproc"
)

// Size of the low-frequency DCT coefficient block used by PHash.
// An 8x8 block produces a 64-bit (8-byte) binary hash.
const dctCoefSize = 8

// PHash is a perceptual hash that uses the method described in
// Implementation and Benchmarking of Perceptual Image Hash Functions; Zauner et. al.
//
// See https://www.researchgate.net/publication/252340846_Rihamark_Perceptual_image_hash_benchmarking for more information.
type PHash struct {
	// Resized image width.
	width uint
	// Resized image height.
	height uint
	// Resize interpolation method.
	interp Interpolation
}

// NewPHash creates a new PHash with the given options.
// Without options, sensible defaults are used.
func NewPHash(opts ...PHashOption) (PHash, error) {
	p := PHash{
		width:  32,
		height: 32,
		interp: BilinearExact,
	}
	for _, o := range opts {
		o.applyPHash(&p)
	}
	if p.width == 0 || p.height == 0 {
		return PHash{}, ErrInvalidSize
	}
	return p, nil
}

// Calculate returns a perceptual image hash.
func (ph PHash) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(ph.width, ph.height, img, ph.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	fImg := imgproc.GrayToF32(g)
	dctImg := imgproc.DCT(fImg)
	tLeft := ph.topLeft(dctImg)
	// Remove the strongest frequency
	tLeft[0][0] = 0
	mean := ph.mean(tLeft)
	bitImg := ph.compare(tLeft, mean)
	return ph.computeHash(bitImg), nil
}

// Computes the binary hash based on the binary image supplied.
func (ph PHash) computeHash(img [][]float32) hashtype.Binary {
	hash := make(hashtype.Binary, dctCoefSize)
	var c uint
	for i := range img {
		for j := range img[i] {
			if img[i][j] != 0 {
				_ = hash.Set(c)
			}
			c++
		}
	}
	return hash
}

// Extract top left block from supplied image.
func (ph PHash) topLeft(img [][]float32) [][]float32 {
	tL := make([][]float32, dctCoefSize)
	for i := range tL {
		tL[i] = img[i][0:dctCoefSize]
	}
	return tL
}

// Compute mean of the supplied image.
func (ph PHash) mean(img [][]float32) float32 {
	var c int
	var s float32
	for i := range img {
		c += len(img[i])
		for j := range img[i] {
			s += img[i][j]
		}
	}
	return s / float32(c)
}

// Build a binary image by comparing the value to the supplied image.
func (ph PHash) compare(img [][]float32, val float32) [][]float32 {
	bit := make([][]float32, len(img))
	for i := range img {
		bit[i] = make([]float32, len(img[i]))
		for j := range img[i] {
			if img[i][j] > val {
				bit[i][j] = 1
			}
		}
	}
	return bit
}
