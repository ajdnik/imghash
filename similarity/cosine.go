package similarity

import (
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// Cosine calculates the cosine distance between two hashes.
// Cosine distance is defined as 1 - cos(theta), where cos(theta) is the
// cosine similarity (dot product divided by the product of magnitudes).
// Returns 0 when both hashes are zero vectors.
func Cosine(h1, h2 hashtype.Hash) Distance {
	l := h1.Len()
	if h2.Len() < l {
		l = h2.Len()
	}
	var dot, mag1, mag2 float64
	for i := 0; i < l; i++ {
		a := h1.ValueAt(i)
		b := h2.ValueAt(i)
		dot += a * b
		mag1 += a * a
		mag2 += b * b
	}
	denom := math.Sqrt(mag1) * math.Sqrt(mag2)
	if denom == 0 {
		return 0
	}
	return Distance(1 - dot/denom)
}
