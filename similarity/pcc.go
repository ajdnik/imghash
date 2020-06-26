package similarity

import (
	"errors"
	"math"

	"github.com/ajdnik/imghash/hashtype"
)

// PCCFloat64 calculates PCC distance for two float64 hashes.
func PCCFloat64(h1, h2 hashtype.Float64) (Distance, error) {
	return pccCalculate(h1, h2)
}

// PCCUInt8 calculates PCC distance for two uint8 hashes.
func PCCUInt8(h1, h2 hashtype.UInt8) (Distance, error) {
	return pccCalculate(h1, h2)
}

// ErrNotSameLength is reported when the length of hashes doesn't match.
var ErrNotSameLength = errors.New("hashes aren't the same length")

// Calculates peak correlation coefficient for a generic slice.
func pccCalculate(h1, h2 interface{}) (Distance, error) {
	len1, len2 := genLen(h1), genLen(h2)
	if len1 != len2 {
		return 0, ErrNotSameLength
	}
	hf1, hf2 := convert(h1), convert(h2)
	mn1, std1 := meanStdDev(hf1)
	mn2, std2 := meanStdDev(hf2)
	for i := range hf1 {
		hf1[i] -= float32(mn1)
		hf2[i] -= float32(mn2)
	}
	max := math.SmallestNonzeroFloat64
	for range hf1 {
		prod, _ := dotProd(hf1, hf2)
		covar := prod / float64(len(hf1))
		corre := covar / (std1*std2 + 1e-20)
		max = math.Max(corre, max)
		hf2 = rotate(hf2, 1)
	}
	return Distance(max), nil
}

// Convert generic slice to float32 slice.
func convert(slice interface{}) []float32 {
	var s []float32
	switch ss := slice.(type) {
	case hashtype.Float64:
		s = make([]float32, len(ss))
		for i := range ss {
			s[i] = float32(ss[i])
		}
	case hashtype.UInt8:
		s = make([]float32, len(ss))
		for i := range ss {
			s[i] = float32(ss[i])
		}
	}
	return s
}

// Compute slice mean ans standard deviation.
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

// Compute dot product of two slices.
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

// Rotate slice by distance k.
func rotate(slice []float32, k int) []float32 {
	if k < 0 || len(slice) == 0 {
		return slice
	}
	r := len(slice) - k%len(slice)
	slice = append(slice[r:], slice[:r]...)
	return slice
}
