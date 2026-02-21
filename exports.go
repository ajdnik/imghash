package imghash

import (
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/similarity"
)

// Re-export core types so most consumers only need to import "imghash".

// Hash is the common interface for all hash representations.
type Hash = hashtype.Hash

// Binary represents a hash where the smallest element is a bit.
type Binary = hashtype.Binary

// UInt8 represents a hash where the smallest element is a uint8 value.
type UInt8 = hashtype.UInt8

// Float64 represents a hash where the smallest element is a float64.
type Float64 = hashtype.Float64

// Distance represents a similarity measure between two hashes.
type Distance = similarity.Distance

// ErrIncompatibleHash is reported when Distance is called with an incompatible hash type
// (e.g. comparing a Binary hash with a non-Binary hash using Hamming distance).
var ErrIncompatibleHash = hashtype.ErrIncompatibleHash
