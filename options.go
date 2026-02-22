package imghash

// baseConfig holds the resize dimensions and interpolation method shared by
// most hash algorithms. Algorithms embed this struct so that WithSize and
// WithInterpolation options can target a single applyBase method.
type baseConfig struct {
	width, height uint
	interp        Interpolation
}

func (b baseConfig) validate() error {
	if b.width == 0 || b.height == 0 {
		return ErrInvalidSize
	}
	return b.interp.validate()
}

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

// WHashOption configures the WHash algorithm.
type WHashOption interface{ applyWHash(*WHash) }

// LBPOption configures the LBP hash algorithm.
type LBPOption interface{ applyLBP(*LBP) }

// HOGHashOption configures the HOGHash algorithm.
type HOGHashOption interface{ applyHOGHash(*HOGHash) }

// PDQOption configures the PDQ hash algorithm.
type PDQOption interface{ applyPDQ(*PDQ) }

// RASHOption configures the RASH hash algorithm.
type RASHOption interface{ applyRASH(*RASH) }

// DistanceOption overrides the comparison function used by Compare.
type DistanceOption struct{ fn DistanceFunc }

func (o DistanceOption) applyAverage(a *Average)               { a.distFunc = o.fn }
func (o DistanceOption) applyDifference(d *Difference)         { d.distFunc = o.fn }
func (o DistanceOption) applyMedian(m *Median)                 { m.distFunc = o.fn }
func (o DistanceOption) applyPHash(p *PHash)                   { p.distFunc = o.fn }
func (o DistanceOption) applyBlockMean(b *BlockMean)           { b.distFunc = o.fn }
func (o DistanceOption) applyMarrHildreth(m *MarrHildreth)     { m.distFunc = o.fn }
func (o DistanceOption) applyRadialVariance(r *RadialVariance) { r.distFunc = o.fn }
func (o DistanceOption) applyColorMoment(c *ColorMoment)       { c.distFunc = o.fn }
func (o DistanceOption) applyWHash(w *WHash)                   { w.distFunc = o.fn }
func (o DistanceOption) applyLBP(l *LBP)                       { l.distFunc = o.fn }
func (o DistanceOption) applyHOGHash(h *HOGHash)               { h.distFunc = o.fn }
func (o DistanceOption) applyPDQ(p *PDQ)                       { p.distFunc = o.fn }
func (o DistanceOption) applyRASH(r *RASH)                     { r.distFunc = o.fn }

// --- concrete option types ---

// SizeOption sets width and height for hash computation.
type SizeOption struct{ width, height uint }

func (o SizeOption) applyBase(b *baseConfig)           { b.width, b.height = o.width, o.height }
func (o SizeOption) applyAverage(a *Average)           { o.applyBase(&a.baseConfig) }
func (o SizeOption) applyDifference(d *Difference)     { o.applyBase(&d.baseConfig) }
func (o SizeOption) applyMedian(m *Median)             { o.applyBase(&m.baseConfig) }
func (o SizeOption) applyPHash(p *PHash)               { o.applyBase(&p.baseConfig) }
func (o SizeOption) applyBlockMean(b *BlockMean)       { o.applyBase(&b.baseConfig) }
func (o SizeOption) applyMarrHildreth(m *MarrHildreth) { o.applyBase(&m.baseConfig) }
func (o SizeOption) applyColorMoment(c *ColorMoment)   { o.applyBase(&c.baseConfig) }
func (o SizeOption) applyWHash(w *WHash)               { o.applyBase(&w.baseConfig) }
func (o SizeOption) applyLBP(l *LBP)                   { o.applyBase(&l.baseConfig) }
func (o SizeOption) applyHOGHash(h *HOGHash)           { o.applyBase(&h.baseConfig) }
func (o SizeOption) applyRASH(r *RASH)                 { o.applyBase(&r.baseConfig) }

// InterpolationOption sets the resize interpolation method.
type InterpolationOption struct{ interp Interpolation }

func (o InterpolationOption) applyBase(b *baseConfig)           { b.interp = o.interp }
func (o InterpolationOption) applyAverage(a *Average)           { o.applyBase(&a.baseConfig) }
func (o InterpolationOption) applyDifference(d *Difference)     { o.applyBase(&d.baseConfig) }
func (o InterpolationOption) applyMedian(m *Median)             { o.applyBase(&m.baseConfig) }
func (o InterpolationOption) applyPHash(p *PHash)               { o.applyBase(&p.baseConfig) }
func (o InterpolationOption) applyBlockMean(b *BlockMean)       { o.applyBase(&b.baseConfig) }
func (o InterpolationOption) applyMarrHildreth(m *MarrHildreth) { o.applyBase(&m.baseConfig) }
func (o InterpolationOption) applyColorMoment(c *ColorMoment)   { o.applyBase(&c.baseConfig) }
func (o InterpolationOption) applyWHash(w *WHash)               { o.applyBase(&w.baseConfig) }
func (o InterpolationOption) applyLBP(l *LBP)                   { o.applyBase(&l.baseConfig) }
func (o InterpolationOption) applyHOGHash(h *HOGHash)           { o.applyBase(&h.baseConfig) }
func (o InterpolationOption) applyPDQ(p *PDQ)                   { p.interp = o.interp }
func (o InterpolationOption) applyRASH(r *RASH)                 { o.applyBase(&r.baseConfig) }

// KernelSizeOption sets the Gaussian kernel size.
type KernelSizeOption struct{ size int }

func (o KernelSizeOption) applyMarrHildreth(m *MarrHildreth) { m.kernel = o.size }
func (o KernelSizeOption) applyColorMoment(c *ColorMoment)   { c.kernel = o.size }

// SigmaOption sets the Gaussian standard deviation.
type SigmaOption struct{ sigma float64 }

func (o SigmaOption) applyMarrHildreth(m *MarrHildreth)     { m.sigma = o.sigma }
func (o SigmaOption) applyColorMoment(c *ColorMoment)       { c.sigma = o.sigma }
func (o SigmaOption) applyRadialVariance(r *RadialVariance) { r.sigma = o.sigma }
func (o SigmaOption) applyRASH(r *RASH)                     { r.sigma = o.sigma }

// BlockSizeOption sets block dimensions for BlockMean.
type BlockSizeOption struct{ width, height uint }

func (o BlockSizeOption) applyBlockMean(b *BlockMean) { b.bWidth, b.bHeight = o.width, o.height }

// BlockMeanMethodOption sets the block construction method.
type BlockMeanMethodOption struct{ method BlockMeanMethod }

func (o BlockMeanMethodOption) applyBlockMean(b *BlockMean) { b.method = o.method }

// ScaleOption sets the scale parameter.
type ScaleOption struct{ scale float64 }

func (o ScaleOption) applyMarrHildreth(m *MarrHildreth) { m.scale = o.scale }

// AlphaOption sets the alpha parameter.
type AlphaOption struct{ alpha float64 }

func (o AlphaOption) applyMarrHildreth(m *MarrHildreth) { m.alpha = o.alpha }

// AnglesOption sets the number of projection angles.
type AnglesOption struct{ angles int }

func (o AnglesOption) applyRadialVariance(r *RadialVariance) { r.angles = o.angles }

// LevelOption sets the wavelet decomposition level.
type LevelOption struct{ level int }

func (o LevelOption) applyWHash(w *WHash) { w.level = o.level }

// GridSizeOption sets the grid cell count for spatial histograms.
type GridSizeOption struct{ x, y uint }

func (o GridSizeOption) applyLBP(l *LBP) { l.gridX, l.gridY = o.x, o.y }

// CellSizeOption sets the cell size in pixels for HOG computation.
type CellSizeOption struct{ size uint }

func (o CellSizeOption) applyHOGHash(h *HOGHash) { h.cellSize = o.size }

// NumBinsOption sets the number of orientation histogram bins.
type NumBinsOption struct{ bins uint }

func (o NumBinsOption) applyHOGHash(h *HOGHash) { h.numBins = o.bins }

// RingsOption sets the number of concentric rings.
type RingsOption struct{ rings int }

func (o RingsOption) applyRASH(r *RASH) { r.rings = o.rings }

// WeightsOption sets the per-byte weights for weighted distance.
type WeightsOption struct{ weights []float64 }

func (o WeightsOption) applyPHash(p *PHash) { p.weights = append([]float64(nil), o.weights...) }

// --- public constructors ---

// WithSize sets the resize dimensions used during hash computation.
// Applies to Average, Difference, Median, PHash, BlockMean, MarrHildreth, ColorMoment, WHash, LBP, HOGHash, and RASH.
func WithSize(width, height uint) SizeOption {
	return SizeOption{width, height}
}

// WithInterpolation sets the resize interpolation method.
// Applies to Average, Difference, Median, PHash, BlockMean, MarrHildreth, ColorMoment, WHash, LBP, HOGHash, PDQ, and RASH.
func WithInterpolation(interp Interpolation) InterpolationOption {
	return InterpolationOption{interp}
}

// WithKernelSize sets the Gaussian kernel size.
// Applies to MarrHildreth and ColorMoment.
func WithKernelSize(size int) KernelSizeOption {
	return KernelSizeOption{size}
}

// WithSigma sets the Gaussian kernel standard deviation.
// Applies to MarrHildreth, ColorMoment, RadialVariance, and RASH.
func WithSigma(sigma float64) SigmaOption {
	return SigmaOption{sigma}
}

// WithBlockSize sets the block dimensions for block mean hashing.
// Applies to BlockMean.
func WithBlockSize(width, height uint) BlockSizeOption {
	return BlockSizeOption{width, height}
}

// WithBlockMeanMethod sets the block construction method.
// Applies to BlockMean.
func WithBlockMeanMethod(method BlockMeanMethod) BlockMeanMethodOption {
	return BlockMeanMethodOption{method}
}

// WithScale sets the scale parameter.
// Applies to MarrHildreth.
func WithScale(scale float64) ScaleOption {
	return ScaleOption{scale}
}

// WithAlpha sets the alpha parameter.
// Applies to MarrHildreth.
func WithAlpha(alpha float64) AlphaOption {
	return AlphaOption{alpha}
}

// WithAngles sets the number of projection angles.
// Applies to RadialVariance.
func WithAngles(angles int) AnglesOption {
	return AnglesOption{angles}
}

// WithLevel sets the number of Haar wavelet decomposition levels.
// Applies to WHash.
func WithLevel(level int) LevelOption {
	return LevelOption{level}
}

// WithGridSize sets the number of grid cells used to divide the image for
// spatial histogram computation.
// Applies to LBP.
func WithGridSize(x, y uint) GridSizeOption {
	return GridSizeOption{x, y}
}

// WithCellSize sets the cell size in pixels (square cells) for HOG computation.
// Applies to HOGHash.
func WithCellSize(size uint) CellSizeOption {
	return CellSizeOption{size}
}

// WithNumBins sets the number of orientation histogram bins.
// Applies to HOGHash.
func WithNumBins(bins uint) NumBinsOption {
	return NumBinsOption{bins}
}

// WithRings sets the number of concentric rings used for spatial sampling.
// Applies to RASH.
func WithRings(rings int) RingsOption {
	return RingsOption{rings}
}

// WithWeights sets the per-byte weights used for weighted Hamming distance.
// The slice length must match the number of hash bytes (8 for default PHash).
// Applies to PHash.
func WithWeights(weights []float64) WeightsOption {
	return WeightsOption{append([]float64(nil), weights...)}
}

// WithDistance overrides the default distance function used by Compare.
// All functions in the similarity package (Hamming, L1, L2, Cosine,
// ChiSquare, PCC) satisfy DistanceFunc and can be passed directly.
// Applies to all algorithms.
func WithDistance(fn DistanceFunc) DistanceOption {
	return DistanceOption{fn}
}
