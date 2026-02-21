package similarity

import (
	"errors"
	"math/bits"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// ErrWeightLengthMismatch is reported when the weight slice length doesn't match the hash byte length.
var ErrWeightLengthMismatch = errors.New("weight slice length must match number of hash bytes")

// WeightedHamming calculates a weighted bit-level hamming distance between two binary hashes.
// Each byte position is assigned a weight from the weights slice; the number of differing bits
// at that position is multiplied by the corresponding weight.
// The weights slice must have the same length as the shorter hash.
func WeightedHamming(h1, h2 hashtype.Hash, weights []float64) (Distance, error) {
	b1, ok := h1.(hashtype.Binary)
	if !ok {
		return 0, ErrNotBinaryHash
	}
	b2, ok := h2.(hashtype.Binary)
	if !ok {
		return 0, ErrNotBinaryHash
	}
	l := len(b1)
	if len(b2) < l {
		l = len(b2)
	}
	if len(weights) != l {
		return 0, ErrWeightLengthMismatch
	}
	var dist float64
	for i := 0; i < l; i++ {
		dist += float64(bits.OnesCount8(b1[i]^b2[i])) * weights[i]
	}
	return Distance(dist), nil
}
