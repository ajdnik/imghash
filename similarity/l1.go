package similarity

import (
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// L1 calculates the L1 (Manhattan) distance between two hashes.
// It sums the absolute differences of corresponding elements.
func L1(h1, h2 hashtype.Hash) (Distance, error) {
	l := h1.Len()
	if h2.Len() < l {
		l = h2.Len()
	}
	var s float64
	for i := 0; i < l; i++ {
		s += math.Abs(h1.ValueAt(i) - h2.ValueAt(i))
	}
	return Distance(s), nil
}
