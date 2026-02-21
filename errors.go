package imghash

import "errors"

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
)
