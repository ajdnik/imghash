package similarity

import (
	"errors"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// ErrNotBinaryHash is reported when a non-binary hash is passed to Hamming.
var ErrNotBinaryHash = errors.New("hamming distance requires binary hashes")

// Hamming calculates the bit-level hamming distance between two binary hashes.
func Hamming(h1, h2 hashtype.Hash) (Distance, error) {
	b1, ok := h1.(hashtype.Binary)
	if !ok {
		return 0, ErrNotBinaryHash
	}
	if _, ok := h2.(hashtype.Binary); !ok {
		return 0, ErrNotBinaryHash
	}
	d, _ := b1.Distance(h2)
	return Distance(d), nil
}
