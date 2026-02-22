package imghash

import (
	"errors"
	"image"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

// DistanceFunc computes a distance between two hashes.
// All functions in the similarity package (Hamming, L1, L2, Cosine,
// ChiSquare, PCC) satisfy this signature and can be passed directly
// to WithDistance.
type DistanceFunc func(hashtype.Hash, hashtype.Hash) (similarity.Distance, error)

// Hasher computes a perceptual hash from an image.
// It is implemented by all hash algorithms in this package:
// Average, Difference, PHash, Median, BlockMean, MarrHildreth,
// RadialVariance, ColorMoment, CLD, EHD, WHash, LBP, HOGHash, PDQ, and RASH.
type Hasher interface {
	Calculate(image.Image) (hashtype.Hash, error)
}

// Comparer measures the similarity between two hashes using a
// distance metric appropriate for the algorithm that produced them.
// It is implemented by all hash algorithms in this package.
type Comparer interface {
	Compare(hashtype.Hash, hashtype.Hash) (similarity.Distance, error)
}

// HasherComparer combines Hasher and Comparer into a single interface
// for algorithms that can both compute and compare hashes.
type HasherComparer interface {
	Hasher
	Comparer
}

// Compile-time assertions: every algorithm satisfies HasherComparer.
var (
	_ HasherComparer = Average{}
	_ HasherComparer = Difference{}
	_ HasherComparer = Median{}
	_ HasherComparer = PHash{}
	_ HasherComparer = BlockMean{}
	_ HasherComparer = MarrHildreth{}
	_ HasherComparer = RadialVariance{}
	_ HasherComparer = ColorMoment{}
	_ HasherComparer = CLD{}
	_ HasherComparer = EHD{}
	_ HasherComparer = WHash{}
	_ HasherComparer = LBP{}
	_ HasherComparer = HOGHash{}
	_ HasherComparer = PDQ{}
	_ HasherComparer = RASH{}
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

// ErrIncompatibleHash is reported when a binary-only metric (Hamming or weighted Hamming)
// is used with incompatible hash types.
var ErrIncompatibleHash = hashtype.ErrIncompatibleHash

// ErrHashLengthMismatch is reported when two hashes of the expected type
// have different lengths and cannot be compared safely.
var ErrHashLengthMismatch = errors.New("imghash: hash lengths must match")

// Constructor validation errors.
var (
	// ErrInvalidSize is returned when width or height is zero.
	ErrInvalidSize = errors.New("imghash: size dimensions must be greater than zero")
	// ErrInvalidInterpolation is returned when an unsupported interpolation enum is supplied.
	ErrInvalidInterpolation = errors.New("imghash: invalid interpolation method")
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
	// ErrInvalidRings is returned when the number of concentric rings is not positive.
	ErrInvalidRings = errors.New("imghash: rings must be greater than zero")
)
