package imghash

import (
	"image"
	"sort"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// WHash is a perceptual hash based on the Haar wavelet transform.
// It applies a multi-level 2-D discrete wavelet transform to the image
// and thresholds the low-frequency (LL) coefficients against their median.
//
// See https://fullstackml.com/wavelet-image-hash-in-python-3504571f3b08 for more information.
type WHash struct {
	// Hash output width (columns in the LL subband).
	width uint
	// Hash output height (rows in the LL subband).
	height uint
	// Number of Haar DWT decomposition levels.
	level int
	// Resize interpolation method.
	interp Interpolation
}

// NewWHash creates a new WHash with the given options.
// Without options, sensible defaults are used (8×8 hash, 3 levels, Bilinear).
func NewWHash(opts ...WHashOption) (WHash, error) {
	w := WHash{
		width:  8,
		height: 8,
		level:  3,
		interp: Bilinear,
	}
	for _, o := range opts {
		o.applyWHash(&w)
	}
	if w.width == 0 || w.height == 0 {
		return WHash{}, ErrInvalidSize
	}
	if w.level <= 0 {
		return WHash{}, ErrInvalidLevel
	}
	return w, nil
}

// Calculate returns a perceptual image hash.
func (wh WHash) Calculate(img image.Image) (hashtype.Hash, error) {
	// Resize to (width * 2^level) x (height * 2^level) so that after
	// `level` DWT passes the LL subband is exactly width×height.
	scale := uint(1) << uint(wh.level)
	rw := wh.width * scale
	rh := wh.height * scale

	r := imgproc.Resize(rw, rh, img, wh.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	mat := imgproc.GrayToF32(g)
	imgproc.HaarDWT2D(mat, wh.level)

	ll := wh.extractLL(mat)
	med := wh.median(ll)
	return wh.computeHash(ll, med), nil
}

func (wh WHash) extractLL(mat [][]float32) [][]float32 {
	ll := make([][]float32, wh.height)
	for r := uint(0); r < wh.height; r++ {
		ll[r] = make([]float32, wh.width)
		copy(ll[r], mat[r][:wh.width])
	}
	return ll
}

func (wh WHash) median(mat [][]float32) float32 {
	vals := make([]float64, 0, wh.width*wh.height)
	for _, row := range mat {
		for _, v := range row {
			vals = append(vals, float64(v))
		}
	}
	sort.Float64s(vals)
	n := len(vals)
	if n%2 == 0 {
		return float32((vals[n/2-1] + vals[n/2]) / 2)
	}
	return float32(vals[n/2])
}

func (wh WHash) computeHash(ll [][]float32, median float32) hashtype.Binary {
	hash := make(hashtype.Binary, wh.width*wh.height/8)
	var c uint
	for _, row := range ll {
		for _, v := range row {
			if v > median {
				_ = hash.Set(c)
			}
			c++
		}
	}
	return hash
}

// Compare computes the Hamming distance between two WHash hashes.
func (wh WHash) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	return similarity.Hamming(h1, h2)
}

// Ensure WHash satisfies the Hasher interface.
var _ Hasher = WHash{}
