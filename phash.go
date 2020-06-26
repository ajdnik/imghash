package imghash

import (
	"image"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
)

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
	interp ResizeType
}

// NewPHash creates a new PHash struct using default values.
func NewPHash() PHash {
	return PHash{
		width:  32,
		height: 32,
		interp: BilinearExact,
	}
}

// NewPHashWithParams creates a new PHash struct using the supplied parameters.
func NewPHashWithParams(resizeWidth, resizeHeight uint, resizeType ResizeType) PHash {
	return PHash{
		width:  resizeWidth,
		height: resizeHeight,
		interp: resizeType,
	}
}

// Calculate returns a percaptual image hash.
func (ph *PHash) Calculate(img image.Image) hashtype.Binary {
	r := resizeImageCV(ph.width, ph.height, img, ph.interp)
	g, _ := imgproc.Grayscale(r)
	fImg := imgproc.GrayToF32(g)
	dctImg := dctCV(fImg)
	tLeft := ph.topLeft(dctImg)
	// Remove the strongest frequency
	tLeft[0][0] = 0
	mean := ph.mean(tLeft)
	bitImg := ph.compare(tLeft, mean)
	return ph.computeHash(bitImg)
}

// Computes the binary hash based on the binary image supplied.
func (ph *PHash) computeHash(img [][]float32) hashtype.Binary {
	// TODO: Remove magic numbers
	hash := make(hashtype.Binary, 8)
	var c uint
	for i := range img {
		for j := range img[i] {
			if img[i][j] != 0 {
				hash.Set(c)
			}
			c++
		}
	}
	return hash
}

// Extract top left block from supplied image.
func (ph *PHash) topLeft(img [][]float32) [][]float32 {
	// TODO: Remove magic numbers
	tL := make([][]float32, 8)
	for i := range tL {
		tL[i] = img[i][0:8]
	}
	return tL
}

// Compute mean of the supplied image.
func (ph *PHash) mean(img [][]float32) float32 {
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

// Build a binary image by comparring the value to the supplied image.
func (ph *PHash) compare(img [][]float32, val float32) [][]float32 {
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
