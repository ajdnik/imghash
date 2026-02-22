package imghash

import (
	"image"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// pdqDCTSize is the dimension of the DCT input buffer.
const pdqDCTSize = 64

// pdqCoefSize is the dimension of the DCT coefficient block used for hashing.
// A 16x16 block produces a 256-bit (32-byte) binary hash.
const pdqCoefSize = 16

// pdqHashBytes is the number of bytes in the resulting hash.
const pdqHashBytes = pdqCoefSize * pdqCoefSize / 8

// pdqJaroszWindow is the half-width of the Jarosz box filter kernel.
const pdqJaroszWindow = 2

// pdqJaroszReps is the number of Jarosz filter iterations per axis.
const pdqJaroszReps = 2

// PDQ is a perceptual hash that uses the method described in
// PDQ and TMK+PDQF by Facebook (now Meta).
// It produces a 256-bit hash robust to JPEG compression, rescaling,
// and minor edits while remaining fast enough for large-scale deduplication.
//
// See https://github.com/facebook/ThreatExchange/tree/main/pdq for more information.
type PDQ struct {
	// Resize interpolation method.
	interp   Interpolation
	distFunc DistanceFunc
}

// NewPDQ creates a new PDQ hasher with the given options.
// Without options, sensible defaults are used.
func NewPDQ(opts ...PDQOption) (PDQ, error) {
	p := PDQ{
		interp: Bilinear,
	}
	for _, o := range opts {
		o.applyPDQ(&p)
	}
	return p, nil
}

// Calculate returns a 256-bit perceptual hash of the image.
func (p PDQ) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(pdqDCTSize, pdqDCTSize, img, p.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	buf := imgproc.GrayToF32(g)
	imgproc.JaroszFilter(buf, pdqJaroszWindow, pdqJaroszReps)
	dct := imgproc.DCT(buf)
	block := p.extractBlock(dct)
	med := p.median(block)
	return p.computeHash(block, med), nil
}

// extractBlock returns the top-left pdqCoefSize x pdqCoefSize block from the DCT output.
func (p PDQ) extractBlock(dct [][]float32) [][]float32 {
	block := make([][]float32, pdqCoefSize)
	for i := 0; i < pdqCoefSize; i++ {
		block[i] = make([]float32, pdqCoefSize)
		copy(block[i], dct[i][:pdqCoefSize])
	}
	return block
}

// median computes the median of all values in the block.
func (p PDQ) median(block [][]float32) float32 {
	return imgproc.MedianF32(block)
}

// computeHash thresholds the DCT block against the median to produce a 256-bit hash.
func (p PDQ) computeHash(block [][]float32, median float32) hashtype.Binary {
	hash := make(hashtype.Binary, pdqHashBytes)
	var c uint
	for i := range block {
		for j := range block[i] {
			if block[i][j] > median {
				_ = hash.Set(c)
			}
			c++
		}
	}
	return hash
}

// Compare computes the Hamming distance between two PDQ hashes.
func (p PDQ) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateBinaryCompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if p.distFunc != nil {
		return p.distFunc(h1, h2)
	}
	return similarity.Hamming(h1, h2)
}
