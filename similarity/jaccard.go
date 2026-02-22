package similarity

import (
	"math/bits"

	"github.com/ajdnik/imghash/v2/hashtype"
)

// Jaccard calculates the Jaccard distance.
//
// For Binary hashes it compares set bits (bitset intersection/union).
// For UInt8 and Float64 hashes it treats values as MinHash signatures and
// computes 1 - (matching positions / signature length).
func Jaccard(h1, h2 hashtype.Hash) (Distance, error) {
	switch v1 := h1.(type) {
	case hashtype.Binary:
		v2, ok := h2.(hashtype.Binary)
		if !ok {
			return 0, hashtype.ErrIncompatibleHash
		}
		return jaccardBinary(v1, v2), nil
	case hashtype.UInt8:
		v2, ok := h2.(hashtype.UInt8)
		if !ok {
			return 0, hashtype.ErrIncompatibleHash
		}
		return jaccardSignatureUInt8(v1, v2)
	case hashtype.Float64:
		v2, ok := h2.(hashtype.Float64)
		if !ok {
			return 0, hashtype.ErrIncompatibleHash
		}
		return jaccardSignatureFloat64(v1, v2)
	default:
		return 0, hashtype.ErrIncompatibleHash
	}
}

func jaccardBinary(h1, h2 hashtype.Binary) Distance {
	l := len(h1)
	if len(h2) < l {
		l = len(h2)
	}

	var inter, union int
	for i := 0; i < l; i++ {
		inter += bits.OnesCount8(h1[i] & h2[i])
		union += bits.OnesCount8(h1[i] | h2[i])
	}

	if len(h1) > l {
		for i := l; i < len(h1); i++ {
			union += bits.OnesCount8(h1[i])
		}
	} else if len(h2) > l {
		for i := l; i < len(h2); i++ {
			union += bits.OnesCount8(h2[i])
		}
	}

	if union == 0 {
		return 0
	}
	return Distance(1 - float64(inter)/float64(union))
}

func jaccardSignatureUInt8(h1, h2 hashtype.UInt8) (Distance, error) {
	if len(h1) != len(h2) {
		return 0, ErrNotSameLength
	}
	if len(h1) == 0 {
		return 0, nil
	}

	matches := 0
	for i := range h1 {
		if h1[i] == h2[i] {
			matches++
		}
	}
	return Distance(1 - float64(matches)/float64(len(h1))), nil
}

func jaccardSignatureFloat64(h1, h2 hashtype.Float64) (Distance, error) {
	if len(h1) != len(h2) {
		return 0, ErrNotSameLength
	}
	if len(h1) == 0 {
		return 0, nil
	}

	matches := 0
	for i := range h1 {
		if h1[i] == h2[i] {
			matches++
		}
	}
	return Distance(1 - float64(matches)/float64(len(h1))), nil
}
