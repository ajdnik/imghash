package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/internal/imgproc"
)

// RadialVariance is a perceptual hash that uses the method described in
// Robust image hashing based on radial variance of pixels; De Roover et. al.
//
// See https://www.researchgate.net/publication/4186555_Robust_image_hashing_based_on_radial_variance_of_pixels for more information.
type RadialVariance struct {
	// Gaussian kernel standard deviation.
	sigma float64
	// Number of angles to consider.
	angles int
}

const hashSize = 40
const sqTwo = math.Sqrt2

// NewRadialVariance creates a new RadialVariance hash with the given options.
// Without options, sensible defaults are used.
func NewRadialVariance(opts ...RadialVarianceOption) RadialVariance {
	rv := RadialVariance{
		sigma:  1,
		angles: 180,
	}
	for _, o := range opts {
		o.applyRadialVariance(&rv)
	}
	return rv
}

// Calculate returns a perceptual image hash.
func (rv RadialVariance) Calculate(img image.Image) (hashtype.Hash, error) {
	g, err := imgproc.Grayscale(img)
	if err != nil {
		return nil, err
	}
	b := imgproc.GaussianBlur(g, 0, rv.sigma)
	proj, ppl, dim := rv.radialProjections(b.(*image.Gray))
	feat := rv.findFeatureVector(proj, ppl, dim)
	return rv.computeHash(feat), nil
}

func roundingFactor(val float32) float32 {
	if val >= 0 {
		return 0.5
	}
	return -0.5
}

func createOffset(len int) int {
	cen := float32(len / 2)
	return int(math.Floor(float64(cen + roundingFactor(cen))))
}

func (rv RadialVariance) radialProjections(img *image.Gray) ([]uint8, []int32, int) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	dim := h
	if w > h {
		dim = w
	}
	proj := make([]uint8, dim*rv.angles)
	pixPerLine := make([]int32, rv.angles)
	xOff, yOff := createOffset(w), createOffset(h)

	// First half of projections
	for k := 0; k < rv.angles/4+1; k++ {
		theta := float32(k) * math.Pi / float32(rv.angles)
		alpha := float32(math.Tan(float64(theta)))
		for x := 0; x < dim; x++ {
			y := alpha * float32(x-xOff)
			yd := int(math.Floor(float64(y + roundingFactor(y))))
			if (yd+yOff >= 0) && (yd+yOff < h) && (x < w) {
				proj[k*dim+x] = img.GrayAt(x, yd+yOff).Y
				pixPerLine[k]++
			}
			if (yd+xOff >= 0) && (yd+xOff < w) && (k != rv.angles/4) && (x < h) {
				proj[(rv.angles/2-k)*dim+x] = img.GrayAt(yd+xOff, x).Y
				pixPerLine[rv.angles/2-k]++
			}
		}
	}

	// Second half of projections
	init := 3 * rv.angles / 4
	for k, j := init, 0; k < rv.angles; k, j = k+1, j+2 {
		theta := float32(k) * math.Pi / float32(rv.angles)
		alpha := float32(math.Tan(float64(theta)))
		for x := 0; x < dim; x++ {
			y := alpha * float32(x-xOff)
			yd := int(math.Floor(float64(y + roundingFactor(y))))
			if (yd+yOff >= 0) && (yd+yOff < h) && (x < w) {
				proj[k*dim+x] = img.GrayAt(x, yd+yOff).Y
				pixPerLine[k]++
			}
			if (yOff-yd >= 0) && (yOff-yd < w) && (2*yOff-x >= 0) && (2*yOff-x < h) && (k != init) {
				proj[(k-j)*dim+x] = img.GrayAt(-yd+yOff, -(x-yOff)+yOff).Y
				pixPerLine[k-j]++
			}
		}
	}
	return proj, pixPerLine, dim
}

func (rv RadialVariance) findFeatureVector(proj []uint8, ppl []int32, dim int) []float64 {
	feat := make([]float64, rv.angles)
	var sum, sqSum float64
	for k := 0; k < rv.angles; k++ {
		var lSum, lSqSum float64
		// Original implementation of pHash may generate zero pixNum, this
		// will cause NaN value and make the features become less discriminative
		// to avoid this problem, I add a small value.
		pNum := float64(ppl[k]) + 0.00001
		for i := 0; i < dim; i++ {
			val := float64(proj[k*dim+i])
			lSum += val
			lSqSum += val * val
		}
		feat[k] = (lSqSum / pNum) - (lSum*lSum)/(pNum*pNum)
		sum += feat[k]
		sqSum += feat[k] * feat[k]
	}
	mean := sum / float64(rv.angles)
	vr := math.Sqrt((sqSum / float64(rv.angles)) - (sum*sum)/float64(rv.angles*rv.angles))
	for i := 0; i < rv.angles; i++ {
		feat[i] = (feat[i] - mean) / vr
	}
	return feat
}

func (rv RadialVariance) computeHash(feat []float64) hashtype.UInt8 {
	hash := make(hashtype.UInt8, hashSize)
	temp := make([]float64, hashSize)
	var max, min float64
	for i := 0; i < hashSize; i++ {
		var sum float64
		for j := 0; j < len(feat); j++ {
			sum += feat[j] * math.Cos((math.Pi*float64(2*j+1)*float64(i))/float64(2*len(feat)))
		}
		if i == 0 {
			temp[i] = sum / math.Sqrt(float64(len(feat)))
		} else {
			temp[i] = sum * sqTwo / math.Sqrt(float64(len(feat)))
		}
		if temp[i] > max {
			max = temp[i]
		} else if temp[i] < min {
			min = temp[i]
		}
	}
	r := max - min
	if r == 0 {
		return hash
	}
	for i := 0; i < hashSize; i++ {
		hash[i] = uint8((255 * (temp[i] - min) / r))
	}
	return hash
}
