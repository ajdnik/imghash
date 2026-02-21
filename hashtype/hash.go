package hashtype

import "math"

// Hash is the common interface for all hash representations.
// It is implemented by Binary, UInt8, and Float64.
//
// Distance returns the natural similarity measure between two hashes:
// Hamming distance for Binary, L2 (Euclidean) distance for UInt8 and Float64.
type Hash interface {
	String() string
	Len() int
	ValueAt(idx int) float64
	Distance(Hash) (float64, error)
}

func l2(h1, h2 Hash) float64 {
	l := h1.Len()
	if h2.Len() < l {
		l = h2.Len()
	}
	var s float64
	for i := 0; i < l; i++ {
		d := h1.ValueAt(i) - h2.ValueAt(i)
		s += d * d
	}
	return math.Sqrt(s)
}
