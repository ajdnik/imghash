package similarity

import (
	"math/bits"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// ErrNotBinaryHash is reported when a non-binary hash is passed to Hamming.
var ErrNotBinaryHash = hashtype.ErrIncompatibleHash

// Hamming calculates the bit-level hamming distance between two binary hashes.
func Hamming(h1, h2 hashtype.Hash) (Distance, error) {
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
	var dist int
	for i := 0; i < l; i++ {
		dist += bits.OnesCount8(b1[i] ^ b2[i])
	}
	return Distance(dist), nil
}
