package similarity

import (
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// L2 calculates the L2 (Euclidean) distance between two hashes.
func L2(h1, h2 hashtype.Hash) (Distance, error) {
	l := h1.Len()
	if h2.Len() < l {
		l = h2.Len()
	}
	var s float64
	for i := 0; i < l; i++ {
		d := h1.ValueAt(i) - h2.ValueAt(i)
		s += d * d
	}
	return Distance(math.Sqrt(s)), nil
}
