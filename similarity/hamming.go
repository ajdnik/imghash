package similarity

import (
	"github.com/ajdnik/imghash/hashtype"
	"github.com/steakknife/hamming"
)

// Hamming calculates distance between two binary hashes.
func Hamming(h1, h2 hashtype.Binary) Distance {
	return Distance(hamming.Bytes(h1, h2))
}
