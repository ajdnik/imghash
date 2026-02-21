package similarity

import (
	"errors"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// ErrNotSameLength is reported when the length of hashes doesn't match.
var ErrNotSameLength = errors.New("hashes aren't the same length")

// PCC calculates the peak cross-correlation between two hashes.
func PCC(h1, h2 hashtype.Hash) (Distance, error) {
	if h1.Len() != h2.Len() {
		return 0, ErrNotSameLength
	}
	hf1 := hashToFloat32(h1)
	hf2 := hashToFloat32(h2)
	mn1, std1 := meanStdDev(hf1)
	mn2, std2 := meanStdDev(hf2)
	for i := range hf1 {
		hf1[i] -= float32(mn1)
		hf2[i] -= float32(mn2)
	}
	peak := math.SmallestNonzeroFloat64
	for range hf1 {
		prod, _ := dotProd(hf1, hf2)
		covar := prod / float64(len(hf1))
		corre := covar / (std1*std2 + 1e-20)
		peak = math.Max(corre, peak)
		hf2 = rotate(hf2, 1)
	}
	return Distance(peak), nil
}

func hashToFloat32(h hashtype.Hash) []float32 {
	s := make([]float32, h.Len())
	for i := range s {
		s[i] = float32(h.ValueAt(i))
	}
	return s
}

func meanStdDev(slice []float32) (float64, float64) {
	var sum, sqSum float64
	for i := range slice {
		sum += float64(slice[i])
		sqSum += math.Pow(float64(slice[i]), 2)
	}
	inv := 1 / float64(len(slice))
	mean := sum * inv
	stdDev := math.Sqrt(math.Max(sqSum*inv-mean*mean, 0.0))
	return mean, stdDev
}

func dotProd(s1 []float32, s2 []float32) (float64, error) {
	if len(s1) != len(s2) {
		return 0, ErrNotSameLength
	}
	var sum float64
	for i := range s1 {
		sum += float64(s1[i]) * float64(s2[i])
	}
	return sum, nil
}

func rotate(slice []float32, k int) []float32 {
	if k < 0 || len(slice) == 0 {
		return slice
	}
	r := len(slice) - k%len(slice)
	slice = append(slice[r:], slice[:r]...)
	return slice
}
