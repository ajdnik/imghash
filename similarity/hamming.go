package similarity

import (
	"errors"

	"github.com/ajdnik/imghash/hashtype"
	"github.com/steakknife/hamming"
)

// ErrNotBinaryHash is reported when a non-binary hash is passed to Hamming.
var ErrNotBinaryHash = errors.New("hamming distance requires binary hashes")

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
	return Distance(hamming.Bytes(b1, b2)), nil
}
