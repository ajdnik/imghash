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

// CLDOption configures the CLD hash algorithm.
type CLDOption interface{ applyCLD(*CLD) }

// EHDOption configures the EHD hash algorithm.
type EHDOption interface{ applyEHD(*EHD) }

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

// ZernikeOption configures the Zernike hash algorithm.
type ZernikeOption interface{ applyZernike(*Zernike) }

// GISTOption configures the GIST hash algorithm.
type GISTOption interface{ applyGIST(*GIST) }

// BoVWOption configures the BoVW hash algorithm.
type BoVWOption interface{ applyBoVW(*BoVW) }

// Option interfaces returned by With* constructors.
// Concrete implementations are intentionally unexported.

// DistanceOption overrides the comparison function used by Compare.
type DistanceOption interface {
	AverageOption
	DifferenceOption
	MedianOption
	PHashOption
	BlockMeanOption
	MarrHildrethOption
	RadialVarianceOption
	ColorMomentOption
	CLDOption
	EHDOption
	WHashOption
	LBPOption
	HOGHashOption
	PDQOption
	RASHOption
	ZernikeOption
	GISTOption
	BoVWOption
}

type distanceOption struct{ fn DistanceFunc }

func (o distanceOption) applyAverage(a *Average)               { a.distFunc = o.fn }
func (o distanceOption) applyDifference(d *Difference)         { d.distFunc = o.fn }
func (o distanceOption) applyMedian(m *Median)                 { m.distFunc = o.fn }
func (o distanceOption) applyPHash(p *PHash)                   { p.distFunc = o.fn }
func (o distanceOption) applyBlockMean(b *BlockMean)           { b.distFunc = o.fn }
func (o distanceOption) applyMarrHildreth(m *MarrHildreth)     { m.distFunc = o.fn }
func (o distanceOption) applyRadialVariance(r *RadialVariance) { r.distFunc = o.fn }
func (o distanceOption) applyColorMoment(c *ColorMoment)       { c.distFunc = o.fn }
func (o distanceOption) applyCLD(c *CLD)                       { c.distFunc = o.fn }
func (o distanceOption) applyEHD(e *EHD)                       { e.distFunc = o.fn }
func (o distanceOption) applyWHash(w *WHash)                   { w.distFunc = o.fn }
func (o distanceOption) applyLBP(l *LBP)                       { l.distFunc = o.fn }
func (o distanceOption) applyHOGHash(h *HOGHash)               { h.distFunc = o.fn }
func (o distanceOption) applyPDQ(p *PDQ)                       { p.distFunc = o.fn }
func (o distanceOption) applyRASH(r *RASH)                     { r.distFunc = o.fn }
func (o distanceOption) applyZernike(z *Zernike)               { z.distFunc = o.fn }
func (o distanceOption) applyGIST(g *GIST)                     { g.distFunc = o.fn }
func (o distanceOption) applyBoVW(b *BoVW)                     { b.distFunc = o.fn }

// --- concrete option implementations ---

// SizeOption sets width and height for hash computation.
type SizeOption interface {
	AverageOption
	DifferenceOption
	MedianOption
	PHashOption
	BlockMeanOption
	MarrHildrethOption
	ColorMomentOption
	CLDOption
	EHDOption
	WHashOption
	LBPOption
	HOGHashOption
	RASHOption
	ZernikeOption
	GISTOption
	BoVWOption
}

type sizeOption struct{ width, height uint }

func (o sizeOption) applyBase(b *baseConfig)           { b.width, b.height = o.width, o.height }
func (o sizeOption) applyAverage(a *Average)           { o.applyBase(&a.baseConfig) }
func (o sizeOption) applyDifference(d *Difference)     { o.applyBase(&d.baseConfig) }
func (o sizeOption) applyMedian(m *Median)             { o.applyBase(&m.baseConfig) }
func (o sizeOption) applyPHash(p *PHash)               { o.applyBase(&p.baseConfig) }
func (o sizeOption) applyBlockMean(b *BlockMean)       { o.applyBase(&b.baseConfig) }
func (o sizeOption) applyMarrHildreth(m *MarrHildreth) { o.applyBase(&m.baseConfig) }
func (o sizeOption) applyColorMoment(c *ColorMoment)   { o.applyBase(&c.baseConfig) }
func (o sizeOption) applyCLD(c *CLD)                   { o.applyBase(&c.baseConfig) }
func (o sizeOption) applyEHD(e *EHD)                   { o.applyBase(&e.baseConfig) }
func (o sizeOption) applyWHash(w *WHash)               { o.applyBase(&w.baseConfig) }
func (o sizeOption) applyLBP(l *LBP)                   { o.applyBase(&l.baseConfig) }
func (o sizeOption) applyHOGHash(h *HOGHash)           { o.applyBase(&h.baseConfig) }
func (o sizeOption) applyRASH(r *RASH)                 { o.applyBase(&r.baseConfig) }
func (o sizeOption) applyZernike(z *Zernike)           { o.applyBase(&z.baseConfig) }
func (o sizeOption) applyGIST(g *GIST)                 { o.applyBase(&g.baseConfig) }
func (o sizeOption) applyBoVW(b *BoVW)                 { o.applyBase(&b.baseConfig) }

// InterpolationOption sets the resize interpolation method.
type InterpolationOption interface {
	AverageOption
	DifferenceOption
	MedianOption
	PHashOption
	BlockMeanOption
	MarrHildrethOption
	ColorMomentOption
	CLDOption
	EHDOption
	WHashOption
	LBPOption
	HOGHashOption
	PDQOption
	RASHOption
	ZernikeOption
	GISTOption
	BoVWOption
}

type interpolationOption struct{ interp Interpolation }

func (o interpolationOption) applyBase(b *baseConfig)           { b.interp = o.interp }
func (o interpolationOption) applyAverage(a *Average)           { o.applyBase(&a.baseConfig) }
func (o interpolationOption) applyDifference(d *Difference)     { o.applyBase(&d.baseConfig) }
func (o interpolationOption) applyMedian(m *Median)             { o.applyBase(&m.baseConfig) }
func (o interpolationOption) applyPHash(p *PHash)               { o.applyBase(&p.baseConfig) }
func (o interpolationOption) applyBlockMean(b *BlockMean)       { o.applyBase(&b.baseConfig) }
func (o interpolationOption) applyMarrHildreth(m *MarrHildreth) { o.applyBase(&m.baseConfig) }
func (o interpolationOption) applyColorMoment(c *ColorMoment)   { o.applyBase(&c.baseConfig) }
func (o interpolationOption) applyCLD(c *CLD)                   { o.applyBase(&c.baseConfig) }
func (o interpolationOption) applyEHD(e *EHD)                   { o.applyBase(&e.baseConfig) }
func (o interpolationOption) applyWHash(w *WHash)               { o.applyBase(&w.baseConfig) }
func (o interpolationOption) applyLBP(l *LBP)                   { o.applyBase(&l.baseConfig) }
func (o interpolationOption) applyHOGHash(h *HOGHash)           { o.applyBase(&h.baseConfig) }
func (o interpolationOption) applyPDQ(p *PDQ)                   { p.interp = o.interp }
func (o interpolationOption) applyRASH(r *RASH)                 { o.applyBase(&r.baseConfig) }
func (o interpolationOption) applyZernike(z *Zernike)           { o.applyBase(&z.baseConfig) }
func (o interpolationOption) applyGIST(g *GIST)                 { o.applyBase(&g.baseConfig) }
func (o interpolationOption) applyBoVW(b *BoVW)                 { o.applyBase(&b.baseConfig) }

// KernelSizeOption sets the Gaussian kernel size.
type KernelSizeOption interface {
	MarrHildrethOption
	ColorMomentOption
}

type kernelSizeOption struct{ size int }

func (o kernelSizeOption) applyMarrHildreth(m *MarrHildreth) { m.kernel = o.size }
func (o kernelSizeOption) applyColorMoment(c *ColorMoment)   { c.kernel = o.size }

// SigmaOption sets the Gaussian standard deviation.
type SigmaOption interface {
	MarrHildrethOption
	ColorMomentOption
	RadialVarianceOption
	RASHOption
}

type sigmaOption struct{ sigma float64 }

func (o sigmaOption) applyMarrHildreth(m *MarrHildreth)     { m.sigma = o.sigma }
func (o sigmaOption) applyColorMoment(c *ColorMoment)       { c.sigma = o.sigma }
func (o sigmaOption) applyRadialVariance(r *RadialVariance) { r.sigma = o.sigma }
func (o sigmaOption) applyRASH(r *RASH)                     { r.sigma = o.sigma }

// BlockSizeOption sets block dimensions for BlockMean.
type BlockSizeOption interface {
	BlockMeanOption
}

type blockSizeOption struct{ width, height uint }

func (o blockSizeOption) applyBlockMean(b *BlockMean) { b.bWidth, b.bHeight = o.width, o.height }

// BlockMeanMethodOption sets the block construction method.
type BlockMeanMethodOption interface {
	BlockMeanOption
}

type blockMeanMethodOption struct{ method BlockMeanMethod }

func (o blockMeanMethodOption) applyBlockMean(b *BlockMean) { b.method = o.method }

// ScaleOption sets the scale parameter.
type ScaleOption interface {
	MarrHildrethOption
}

type scaleOption struct{ scale float64 }

func (o scaleOption) applyMarrHildreth(m *MarrHildreth) { m.scale = o.scale }

// AlphaOption sets the alpha parameter.
type AlphaOption interface {
	MarrHildrethOption
}

type alphaOption struct{ alpha float64 }

func (o alphaOption) applyMarrHildreth(m *MarrHildreth) { m.alpha = o.alpha }

// AnglesOption sets the number of projection angles.
type AnglesOption interface {
	RadialVarianceOption
}

type anglesOption struct{ angles int }

func (o anglesOption) applyRadialVariance(r *RadialVariance) { r.angles = o.angles }

// LevelOption sets the wavelet decomposition level.
type LevelOption interface {
	WHashOption
}

type levelOption struct{ level int }

func (o levelOption) applyWHash(w *WHash) { w.level = o.level }

// GridSizeOption sets the grid cell count for spatial histograms.
type GridSizeOption interface {
	LBPOption
	GISTOption
}

type gridSizeOption struct{ x, y uint }

func (o gridSizeOption) applyLBP(l *LBP) { l.gridX, l.gridY = o.x, o.y }
func (o gridSizeOption) applyGIST(g *GIST) {
	g.gridX, g.gridY = o.x, o.y
}

// CellSizeOption sets the cell size in pixels for HOG computation.
type CellSizeOption interface {
	HOGHashOption
}

type cellSizeOption struct{ size uint }

func (o cellSizeOption) applyHOGHash(h *HOGHash) { h.cellSize = o.size }

// NumBinsOption sets the number of orientation histogram bins.
type NumBinsOption interface {
	HOGHashOption
}

type numBinsOption struct{ bins uint }

func (o numBinsOption) applyHOGHash(h *HOGHash) { h.numBins = o.bins }

// RingsOption sets the number of concentric rings.
type RingsOption interface {
	RASHOption
}

type ringsOption struct{ rings int }

func (o ringsOption) applyRASH(r *RASH) { r.rings = o.rings }

// DegreeOption sets the maximum Zernike degree.
type DegreeOption interface {
	ZernikeOption
}

type degreeOption struct{ degree int }

func (o degreeOption) applyZernike(z *Zernike) { z.degree = o.degree }

// WeightsOption sets the per-byte weights for weighted distance.
type WeightsOption interface {
	PHashOption
}

type weightsOption struct{ weights []float64 }

func (o weightsOption) applyPHash(p *PHash) { p.weights = append([]float64(nil), o.weights...) }

// BoVWFeatureOption sets the local feature extractor for BoVW.
type BoVWFeatureOption interface {
	BoVWOption
}

type bovwFeatureOption struct{ feature BoVWFeatureType }

func (o bovwFeatureOption) applyBoVW(b *BoVW) { b.featureType = o.feature }

// BoVWStorageOption sets the output storage type for BoVW.
type BoVWStorageOption interface {
	BoVWOption
}

type bovwStorageOption struct{ storage BoVWStorageType }

func (o bovwStorageOption) applyBoVW(b *BoVW) { b.storageType = o.storage }

// VocabularySizeOption sets the visual vocabulary size used by BoVW.
type VocabularySizeOption interface {
	BoVWOption
}

type vocabularySizeOption struct{ size uint }

func (o vocabularySizeOption) applyBoVW(b *BoVW) { b.vocabularySize = o.size }

// MaxKeypointsOption sets the maximum number of BoVW keypoints.
type MaxKeypointsOption interface {
	BoVWOption
}

type maxKeypointsOption struct{ count uint }

func (o maxKeypointsOption) applyBoVW(b *BoVW) { b.maxKeypoints = o.count }

// MinHashSizeOption sets the MinHash signature size used by BoVW.
type MinHashSizeOption interface {
	BoVWOption
}

type minHashSizeOption struct{ size uint }

func (o minHashSizeOption) applyBoVW(b *BoVW) { b.minHashSize = o.size }

// SimHashBitsOption sets the SimHash bit length used by BoVW.
type SimHashBitsOption interface {
	BoVWOption
}

type simHashBitsOption struct{ bits uint }

func (o simHashBitsOption) applyBoVW(b *BoVW) { b.simHashBits = o.bits }

// --- public constructors ---

// WithSize sets the resize dimensions used during hash computation.
// Applies to Average, Difference, Median, PHash, BlockMean, MarrHildreth, ColorMoment, CLD, EHD, WHash, LBP, HOGHash, BoVW, RASH, Zernike, and GIST.
func WithSize(width, height uint) SizeOption {
	return sizeOption{width, height}
}

// WithInterpolation sets the resize interpolation method.
// Applies to Average, Difference, Median, PHash, BlockMean, MarrHildreth, ColorMoment, CLD, EHD, WHash, LBP, HOGHash, BoVW, PDQ, RASH, Zernike, and GIST.
func WithInterpolation(interp Interpolation) InterpolationOption {
	return interpolationOption{interp}
}

// WithKernelSize sets the Gaussian kernel size.
// Applies to MarrHildreth and ColorMoment.
func WithKernelSize(size int) KernelSizeOption {
	return kernelSizeOption{size}
}

// WithSigma sets the Gaussian kernel standard deviation.
// Applies to MarrHildreth, ColorMoment, RadialVariance, and RASH.
func WithSigma(sigma float64) SigmaOption {
	return sigmaOption{sigma}
}

// WithBlockSize sets the block dimensions for block mean hashing.
// Applies to BlockMean.
func WithBlockSize(width, height uint) BlockSizeOption {
	return blockSizeOption{width, height}
}

// WithBlockMeanMethod sets the block construction method.
// Applies to BlockMean.
func WithBlockMeanMethod(method BlockMeanMethod) BlockMeanMethodOption {
	return blockMeanMethodOption{method}
}

// WithScale sets the scale parameter.
// Applies to MarrHildreth.
func WithScale(scale float64) ScaleOption {
	return scaleOption{scale}
}

// WithAlpha sets the alpha parameter.
// Applies to MarrHildreth.
func WithAlpha(alpha float64) AlphaOption {
	return alphaOption{alpha}
}

// WithAngles sets the number of projection angles.
// Applies to RadialVariance.
func WithAngles(angles int) AnglesOption {
	return anglesOption{angles}
}

// WithLevel sets the number of Haar wavelet decomposition levels.
// Applies to WHash.
func WithLevel(level int) LevelOption {
	return levelOption{level}
}

// WithGridSize sets the number of grid cells used to divide the image for
// spatial histogram computation.
// Applies to LBP and GIST.
func WithGridSize(x, y uint) GridSizeOption {
	return gridSizeOption{x, y}
}

// WithCellSize sets the cell size in pixels (square cells) for HOG computation.
// Applies to HOGHash.
func WithCellSize(size uint) CellSizeOption {
	return cellSizeOption{size}
}

// WithNumBins sets the number of orientation histogram bins.
// Applies to HOGHash.
func WithNumBins(bins uint) NumBinsOption {
	return numBinsOption{bins}
}

// WithRings sets the number of concentric rings used for spatial sampling.
// Applies to RASH.
func WithRings(rings int) RingsOption {
	return ringsOption{rings}
}

// WithDegree sets the maximum Zernike degree.
// Applies to Zernike.
func WithDegree(degree int) DegreeOption {
	return degreeOption{degree}
}

// WithWeights sets the per-byte weights used for weighted Hamming distance.
// The slice length must match the number of hash bytes (8 for default PHash).
// Applies to PHash.
func WithWeights(weights []float64) WeightsOption {
	return weightsOption{append([]float64(nil), weights...)}
}

// WithBoVWFeature sets the local feature extractor used by BoVW.
// Applies to BoVW.
func WithBoVWFeature(feature BoVWFeatureType) BoVWFeatureOption {
	return bovwFeatureOption{feature: feature}
}

// WithBoVWStorage sets the BoVW output storage representation.
// Applies to BoVW.
func WithBoVWStorage(storage BoVWStorageType) BoVWStorageOption {
	return bovwStorageOption{storage: storage}
}

// WithVocabularySize sets the visual vocabulary size used by BoVW.
// Applies to BoVW.
func WithVocabularySize(size uint) VocabularySizeOption {
	return vocabularySizeOption{size: size}
}

// WithMaxKeypoints sets the maximum number of BoVW keypoints.
// Applies to BoVW.
func WithMaxKeypoints(count uint) MaxKeypointsOption {
	return maxKeypointsOption{count: count}
}

// WithMinHashSize sets the MinHash signature length used by BoVW.
// Applies to BoVW.
func WithMinHashSize(size uint) MinHashSizeOption {
	return minHashSizeOption{size: size}
}

// WithSimHashBits sets the SimHash bit length used by BoVW.
// Applies to BoVW.
func WithSimHashBits(bits uint) SimHashBitsOption {
	return simHashBitsOption{bits: bits}
}

// WithDistance overrides the default distance function used by Compare.
// All functions in the similarity package (Hamming, L1, L2, Cosine,
// ChiSquare, PCC, Jaccard) satisfy DistanceFunc and can be passed directly.
// Applies to all algorithms.
func WithDistance(fn DistanceFunc) DistanceOption {
	return distanceOption{fn}
}
