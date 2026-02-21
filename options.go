package imghash

// Option configures hash algorithm parameters.
type Option func(*options)

type options struct {
	width       uint
	height      uint
	interp      Interpolation
	kernel      int
	sigma       float64
	blockWidth  uint
	blockHeight uint
	blockMethod BlockMeanMethod
	scale       float64
	alpha       float64
	angles      int
}

func applyOptions(o *options, opts []Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithSize sets the resize dimensions used during hash computation.
func WithSize(width, height uint) Option {
	return func(o *options) {
		o.width = width
		o.height = height
	}
}

// WithInterpolation sets the resize interpolation method.
func WithInterpolation(interp Interpolation) Option {
	return func(o *options) {
		o.interp = interp
	}
}

// WithKernelSize sets the Gaussian kernel size.
func WithKernelSize(size int) Option {
	return func(o *options) {
		o.kernel = size
	}
}

// WithSigma sets the Gaussian kernel standard deviation.
func WithSigma(sigma float64) Option {
	return func(o *options) {
		o.sigma = sigma
	}
}

// WithBlockSize sets the block dimensions for block mean hashing.
func WithBlockSize(width, height uint) Option {
	return func(o *options) {
		o.blockWidth = width
		o.blockHeight = height
	}
}

// WithBlockMeanMethod sets the block construction method.
func WithBlockMeanMethod(method BlockMeanMethod) Option {
	return func(o *options) {
		o.blockMethod = method
	}
}

// WithScale sets the scale parameter for Marr-Hildreth hashing.
func WithScale(scale float64) Option {
	return func(o *options) {
		o.scale = scale
	}
}

// WithAlpha sets the alpha parameter for Marr-Hildreth hashing.
func WithAlpha(alpha float64) Option {
	return func(o *options) {
		o.alpha = alpha
	}
}

// WithAngles sets the number of projection angles for radial variance hashing.
func WithAngles(angles int) Option {
	return func(o *options) {
		o.angles = angles
	}
}
