package imghash

import (
	"image"
	"math"
	"sort"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// RASH is a Rotation Aware Spatial Hash — a perceptual hash designed to be
// robust against image rotation. It combines concentric ring sampling for
// rotation invariance, a 1-D DCT for frequency compaction, and median
// thresholding for binarisation.
//
// Because ring-mean features are inherently rotation-invariant (rotating the
// image only permutes pixels within a ring, leaving its mean unchanged), the
// resulting hash stays stable under arbitrary rotations.
//
// Inspired by ring-partition hashing literature:
//   - Tang et al., "Robust Image Hashing with Ring Partition and Invariant Vector Distance" (2016)
//   - De Roover et al., "Robust image hashing based on radial variance of pixels" (2005)
type RASH struct {
	baseConfig
	sigma    float64
	rings    int
	distFunc DistanceFunc
}

const rashHashBits = 64

// NewRASH creates a new RASH hash with the given options.
// Without options, sensible defaults are used.
func NewRASH(opts ...RASHOption) (RASH, error) {
	r := RASH{
		baseConfig: baseConfig{width: 256, height: 256, interp: Bilinear},
		sigma:      1,
		rings:      180,
	}
	for _, o := range opts {
		o.applyRASH(&r)
	}
	if r.width == 0 || r.height == 0 {
		return RASH{}, ErrInvalidSize
	}
	if r.sigma < 0 {
		return RASH{}, ErrInvalidSigma
	}
	if r.rings <= 0 {
		return RASH{}, ErrInvalidRings
	}
	return r, nil
}

// Calculate returns a perceptual image hash.
func (r RASH) Calculate(img image.Image) (hashtype.Hash, error) {
	resized := imgproc.Resize(r.width, r.height, img, r.interp.resizeType())
	g, err := imgproc.Grayscale(resized)
	if err != nil {
		return nil, err
	}
	blurred := imgproc.GaussianBlur(g, 0, r.sigma)
	means := r.ringMeans(blurred.(*image.Gray))
	return r.computeHash(means), nil
}

// ringMeans computes the mean pixel intensity for each concentric ring.
func (r RASH) ringMeans(img *image.Gray) []float64 {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	cx, cy := float64(w)/2.0, float64(h)/2.0
	maxRadius := cx
	if cy < maxRadius {
		maxRadius = cy
	}

	ringWidth := maxRadius / float64(r.rings)
	sums := make([]float64, r.rings)
	counts := make([]int, r.rings)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			ring := int(dist / ringWidth)
			if ring >= r.rings {
				ring = r.rings - 1
			}
			sums[ring] += float64(img.GrayAt(x+bounds.Min.X, y+bounds.Min.Y).Y)
			counts[ring]++
		}
	}

	means := make([]float64, r.rings)
	for i := 0; i < r.rings; i++ {
		if counts[i] > 0 {
			means[i] = sums[i] / float64(counts[i])
		}
	}
	return means
}

// computeHash applies a 1-D DCT to the ring means, keeps the first
// rashHashBits low-frequency coefficients (skipping DC), and binarises
// them against the median.
func (r RASH) computeHash(means []float64) hashtype.Binary {
	n := len(means)
	dct := make([]float64, n)
	c0 := math.Sqrt(1.0 / float64(n))
	c1 := math.Sqrt(2.0 / float64(n))
	for k := 0; k < n; k++ {
		var sum float64
		for i := 0; i < n; i++ {
			sum += means[i] * math.Cos(math.Pi*float64(2*i+1)*float64(k)/float64(2*n))
		}
		if k == 0 {
			dct[k] = sum * c0
		} else {
			dct[k] = sum * c1
		}
	}

	hashBits := rashHashBits
	if n-1 < hashBits {
		hashBits = n - 1
	}
	coeffs := dct[1 : hashBits+1]

	sorted := make([]float64, len(coeffs))
	copy(sorted, coeffs)
	sort.Float64s(sorted)
	median := sorted[len(sorted)/2]

	hash := make(hashtype.Binary, hashBits/8)
	for i, c := range coeffs {
		if c > median {
			_ = hash.Set(uint(i))
		}
	}
	return hash
}

// Compare computes the Hamming distance between two RASH hashes.
func (r RASH) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateBinaryCompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if r.distFunc != nil {
		return r.distFunc(h1, h2)
	}
	return similarity.Hamming(h1, h2)
}
