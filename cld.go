package imghash

import (
	"image"
	"image/color"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

const (
	cldGridSize     = 8
	cldYCoeffs      = 6
	cldCCoeffs      = 3
	cldHashLen      = cldYCoeffs + 2*cldCCoeffs
	cldDCQuantRange = 63
	cldACQuantRange = 31
	cldACCenter     = 16
	cldACScale      = 16.0
)

// Zig-zag order for an 8x8 block, represented as row-major indices.
var cldZigZag = [64]uint8{
	0, 1, 8, 16, 9, 2, 3, 10,
	17, 24, 32, 25, 18, 11, 4, 5,
	12, 19, 26, 33, 40, 48, 41, 34,
	27, 20, 13, 6, 7, 14, 21, 28,
	35, 42, 49, 56, 57, 50, 43, 36,
	29, 22, 15, 23, 30, 37, 44, 51,
	58, 59, 52, 45, 38, 31, 39, 46,
	53, 60, 61, 54, 47, 55, 62, 63,
}

// CLD is an MPEG-7 Color Layout Descriptor style perceptual hash.
// The image is converted to YCbCr, summarised into an 8x8 layout, transformed
// with a 2-D DCT per channel, then low-frequency zig-zag coefficients are
// quantised into a compact 12-element descriptor.
type CLD struct {
	baseConfig
	distFunc DistanceFunc
}

// NewCLD creates a new CLD hash with the given options.
// Without options, sensible defaults are used.
func NewCLD(opts ...CLDOption) (CLD, error) {
	c := CLD{
		baseConfig: baseConfig{width: 64, height: 64, interp: Bilinear},
	}
	for _, o := range opts {
		o.applyCLD(&c)
	}
	if err := c.validate(); err != nil {
		return CLD{}, err
	}
	return c, nil
}

// Calculate returns an MPEG-7 CLD style perceptual hash.
func (c CLD) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(c.width, c.height, img, c.interp.resizeType())
	y, cb, cr := c.layoutYCbCr(r)

	yDCT := imgproc.DCT(y)
	cbDCT := imgproc.DCT(cb)
	crDCT := imgproc.DCT(cr)

	hash := make(hashtype.UInt8, cldHashLen)
	hash[0] = quantizeCLDDC(yDCT[0][0])
	for i := 1; i < cldYCoeffs; i++ {
		hash[i] = quantizeCLDAC(zigZagCoeff(yDCT, i))
	}

	offset := cldYCoeffs
	hash[offset] = quantizeCLDDC(cbDCT[0][0])
	for i := 1; i < cldCCoeffs; i++ {
		hash[offset+i] = quantizeCLDAC(zigZagCoeff(cbDCT, i))
	}

	offset += cldCCoeffs
	hash[offset] = quantizeCLDDC(crDCT[0][0])
	for i := 1; i < cldCCoeffs; i++ {
		hash[offset+i] = quantizeCLDAC(zigZagCoeff(crDCT, i))
	}

	return hash, nil
}

func (c CLD) layoutYCbCr(img image.Image) ([][]float32, [][]float32, [][]float32) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	y := make([][]float32, cldGridSize)
	cb := make([][]float32, cldGridSize)
	cr := make([][]float32, cldGridSize)
	for i := 0; i < cldGridSize; i++ {
		y[i] = make([]float32, cldGridSize)
		cb[i] = make([]float32, cldGridSize)
		cr[i] = make([]float32, cldGridSize)
	}

	for gy := 0; gy < cldGridSize; gy++ {
		y0 := gy * h / cldGridSize
		y1 := (gy + 1) * h / cldGridSize
		if y1 <= y0 {
			if y0 >= h {
				y0 = h - 1
			}
			y1 = y0 + 1
		}
		for gx := 0; gx < cldGridSize; gx++ {
			x0 := gx * w / cldGridSize
			x1 := (gx + 1) * w / cldGridSize
			if x1 <= x0 {
				if x0 >= w {
					x0 = w - 1
				}
				x1 = x0 + 1
			}

			var sy, scb, scr float64
			var n float64
			for yy := y0; yy < y1; yy++ {
				for xx := x0; xx < x1; xx++ {
					r, g, b, _ := img.At(bounds.Min.X+xx, bounds.Min.Y+yy).RGBA()
					yyc, cbc, crc := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
					sy += float64(yyc)
					scb += float64(cbc)
					scr += float64(crc)
					n++
				}
			}
			y[gy][gx] = float32(sy / n)
			cb[gy][gx] = float32(scb / n)
			cr[gy][gx] = float32(scr / n)
		}
	}

	return y, cb, cr
}

func zigZagCoeff(dct [][]float32, idx int) float64 {
	pos := cldZigZag[idx]
	return float64(dct[pos/cldGridSize][pos%cldGridSize])
}

func quantizeCLDDC(dc float32) uint8 {
	// For orthonormal 8x8 DCT the DC term equals mean*8; map mean [0,255] to 6 bits.
	mean := float64(dc) / cldGridSize
	q := int(math.Round(mean * cldDCQuantRange / 255.0))
	return uint8(clampInt(q, 0, cldDCQuantRange))
}

func quantizeCLDAC(ac float64) uint8 {
	// Signed 5-bit quantisation with midpoint at 16.
	q := int(math.Round(ac/cldACScale)) + cldACCenter
	return uint8(clampInt(q, 0, cldACQuantRange))
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Compare computes the L2 (Euclidean) distance between two CLD hashes.
func (c CLD) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateUInt8CompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if c.distFunc != nil {
		return c.distFunc(h1, h2)
	}
	return similarity.L2(h1, h2)
}
