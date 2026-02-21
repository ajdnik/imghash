package similarity

import (
	"github.com/ajdnik/imghash/v2/hashtype"
)

// ChiSquare calculates the chi-square distance between two hashes.
// For each element pair it computes (a - b)^2 / (a + b), skipping
// positions where both values are zero to avoid division by zero.
func ChiSquare(h1, h2 hashtype.Hash) Distance {
	l := h1.Len()
	if h2.Len() < l {
		l = h2.Len()
	}
	var s float64
	for i := 0; i < l; i++ {
		a := h1.ValueAt(i)
		b := h2.ValueAt(i)
		sum := a + b
		if sum == 0 {
			continue
		}
		d := a - b
		s += (d * d) / sum
	}
	return Distance(s)
}
