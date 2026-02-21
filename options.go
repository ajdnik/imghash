package imghash

// Algorithm-specific option interfaces. Each uses an unexported method
// so that only option types in this package can implement them.

// AverageOption configures the Average hash algorithm.
type AverageOption interface{ applyAverage(*Average) }

// DifferenceOption configures the Difference hash algorithm.
type DifferenceOption interface{ applyDifference(*Difference) }

// MedianOption configures the Median hash algorithm.
type MedianOption interface{ applyMedian(*Median) }

// PHashOption configures the PHash algorithm.
type PHashOption interface{ applyPHash(*PHash) }

// BlockMeanOption configures the BlockMean hash algorithm.
type BlockMeanOption interface{ applyBlockMean(*BlockMean) }

// MarrHildrethOption configures the MarrHildreth hash algorithm.
type MarrHildrethOption interface{ applyMarrHildreth(*MarrHildreth) }

// RadialVarianceOption configures the RadialVariance hash algorithm.
type RadialVarianceOption interface{ applyRadialVariance(*RadialVariance) }

// ColorMomentOption configures the ColorMoment hash algorithm.
type ColorMomentOption interface{ applyColorMoment(*ColorMoment) }

// --- concrete option types ---

type sizeOption struct{ width, height uint }

func (o sizeOption) applyAverage(a *Average)          { a.width, a.height = o.width, o.height }
func (o sizeOption) applyDifference(d *Difference)    { d.width, d.height = o.width, o.height }
func (o sizeOption) applyMedian(m *Median)            { m.width, m.height = o.width, o.height }
func (o sizeOption) applyPHash(p *PHash)              { p.width, p.height = o.width, o.height }
func (o sizeOption) applyBlockMean(b *BlockMean)      { b.rWidth, b.rHeight = o.width, o.height }
func (o sizeOption) applyMarrHildreth(m *MarrHildreth) { m.width, m.height = o.width, o.height }
func (o sizeOption) applyColorMoment(c *ColorMoment)  { c.width, c.height = o.width, o.height }

type interpolationOption struct{ interp Interpolation }

func (o interpolationOption) applyAverage(a *Average)          { a.interp = o.interp }
func (o interpolationOption) applyDifference(d *Difference)    { d.interp = o.interp }
func (o interpolationOption) applyMedian(m *Median)            { m.interp = o.interp }
func (o interpolationOption) applyPHash(p *PHash)              { p.interp = o.interp }
func (o interpolationOption) applyBlockMean(b *BlockMean)      { b.interp = o.interp }
func (o interpolationOption) applyMarrHildreth(m *MarrHildreth) { m.interp = o.interp }
func (o interpolationOption) applyColorMoment(c *ColorMoment)  { c.interp = o.interp }

type kernelSizeOption struct{ size int }

func (o kernelSizeOption) applyMarrHildreth(m *MarrHildreth) { m.kernel = o.size }
func (o kernelSizeOption) applyColorMoment(c *ColorMoment)   { c.kernel = o.size }

type sigmaOption struct{ sigma float64 }

func (o sigmaOption) applyMarrHildreth(m *MarrHildreth)    { m.sigma = o.sigma }
func (o sigmaOption) applyColorMoment(c *ColorMoment)      { c.sigma = o.sigma }
func (o sigmaOption) applyRadialVariance(r *RadialVariance) { r.sigma = o.sigma }

type blockSizeOption struct{ width, height uint }

func (o blockSizeOption) applyBlockMean(b *BlockMean) { b.bWidth, b.bHeight = o.width, o.height }

type blockMeanMethodOption struct{ method BlockMeanMethod }

func (o blockMeanMethodOption) applyBlockMean(b *BlockMean) { b.method = o.method }

type scaleOption struct{ scale float64 }

func (o scaleOption) applyMarrHildreth(m *MarrHildreth) { m.scale = o.scale }

type alphaOption struct{ alpha float64 }

func (o alphaOption) applyMarrHildreth(m *MarrHildreth) { m.alpha = o.alpha }

type anglesOption struct{ angles int }

func (o anglesOption) applyRadialVariance(r *RadialVariance) { r.angles = o.angles }

// --- public constructors ---

// WithSize sets the resize dimensions used during hash computation.
// Applies to Average, Difference, Median, PHash, BlockMean, MarrHildreth, and ColorMoment.
func WithSize(width, height uint) sizeOption {
	return sizeOption{width, height}
}

// WithInterpolation sets the resize interpolation method.
// Applies to Average, Difference, Median, PHash, BlockMean, MarrHildreth, and ColorMoment.
func WithInterpolation(interp Interpolation) interpolationOption {
	return interpolationOption{interp}
}

// WithKernelSize sets the Gaussian kernel size.
// Applies to MarrHildreth and ColorMoment.
func WithKernelSize(size int) kernelSizeOption {
	return kernelSizeOption{size}
}

// WithSigma sets the Gaussian kernel standard deviation.
// Applies to MarrHildreth, ColorMoment, and RadialVariance.
func WithSigma(sigma float64) sigmaOption {
	return sigmaOption{sigma}
}

// WithBlockSize sets the block dimensions for block mean hashing.
// Applies to BlockMean.
func WithBlockSize(width, height uint) blockSizeOption {
	return blockSizeOption{width, height}
}

// WithBlockMeanMethod sets the block construction method.
// Applies to BlockMean.
func WithBlockMeanMethod(method BlockMeanMethod) blockMeanMethodOption {
	return blockMeanMethodOption{method}
}

// WithScale sets the scale parameter.
// Applies to MarrHildreth.
func WithScale(scale float64) scaleOption {
	return scaleOption{scale}
}

// WithAlpha sets the alpha parameter.
// Applies to MarrHildreth.
func WithAlpha(alpha float64) alphaOption {
	return alphaOption{alpha}
}

// WithAngles sets the number of projection angles.
// Applies to RadialVariance.
func WithAngles(angles int) anglesOption {
	return anglesOption{angles}
}
