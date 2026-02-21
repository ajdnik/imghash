package imghash

import (
	"errors"
	"image"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

// Hasher computes a perceptual hash from an image.
// It is implemented by all hash algorithms in this package:
// Average, Difference, PHash, Median, BlockMean, MarrHildreth,
// RadialVariance, ColorMoment, WHash, LBP, and HOGHash.
type Hasher interface {
	Calculate(image.Image) (hashtype.Hash, error)
}

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

// Constructor validation errors.
var (
	// ErrInvalidSize is returned when width or height is zero.
	ErrInvalidSize = errors.New("imghash: size dimensions must be greater than zero")
	// ErrInvalidBlockSize is returned when block width or height is zero.
	ErrInvalidBlockSize = errors.New("imghash: block size dimensions must be greater than zero")
	// ErrInvalidAngles is returned when the number of projection angles is not positive.
	ErrInvalidAngles = errors.New("imghash: angles must be greater than zero")
	// ErrInvalidKernelSize is returned when the Gaussian kernel size is not positive.
	ErrInvalidKernelSize = errors.New("imghash: kernel size must be greater than zero")
	// ErrInvalidScale is returned when the scale parameter is not positive.
	ErrInvalidScale = errors.New("imghash: scale must be greater than zero")
	// ErrInvalidAlpha is returned when the alpha parameter is not positive.
	ErrInvalidAlpha = errors.New("imghash: alpha must be greater than zero")
	// ErrInvalidSigma is returned when sigma is negative.
	ErrInvalidSigma = errors.New("imghash: sigma must not be negative")
	// ErrInvalidLevel is returned when the wavelet decomposition level is not positive.
	ErrInvalidLevel = errors.New("imghash: level must be greater than zero")
	// ErrInvalidGridSize is returned when grid width or height is zero.
	ErrInvalidGridSize = errors.New("imghash: grid size dimensions must be greater than zero")
	// ErrInvalidCellSize is returned when the cell size is zero.
	ErrInvalidCellSize = errors.New("imghash: cell size must be greater than zero")
	// ErrInvalidNumBins is returned when the number of histogram bins is zero.
	ErrInvalidNumBins = errors.New("imghash: number of bins must be greater than zero")
)
