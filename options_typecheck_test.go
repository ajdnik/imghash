package imghash

// Compile-time assertions: these assignments verify that each With*
// function returns a type satisfying only the intended algorithm options.
// If any line fails to compile, the option was wired to the wrong algorithm.

var _ AverageOption = WithSize(0, 0)
var _ AverageOption = WithInterpolation(Bilinear)
var _ AverageOption = WithDistance(nil)

var _ DifferenceOption = WithSize(0, 0)
var _ DifferenceOption = WithInterpolation(Bilinear)
var _ DifferenceOption = WithDistance(nil)

var _ MedianOption = WithSize(0, 0)
var _ MedianOption = WithInterpolation(Bilinear)
var _ MedianOption = WithDistance(nil)

var _ PHashOption = WithSize(0, 0)
var _ PHashOption = WithInterpolation(Bilinear)
var _ PHashOption = WithWeights(nil)
var _ PHashOption = WithDistance(nil)

var _ BlockMeanOption = WithSize(0, 0)
var _ BlockMeanOption = WithInterpolation(Bilinear)
var _ BlockMeanOption = WithBlockSize(0, 0)
var _ BlockMeanOption = WithBlockMeanMethod(Direct)
var _ BlockMeanOption = WithDistance(nil)

var _ MarrHildrethOption = WithSize(0, 0)
var _ MarrHildrethOption = WithInterpolation(Bilinear)
var _ MarrHildrethOption = WithKernelSize(0)
var _ MarrHildrethOption = WithSigma(0)
var _ MarrHildrethOption = WithScale(0)
var _ MarrHildrethOption = WithAlpha(0)
var _ MarrHildrethOption = WithDistance(nil)

var _ RadialVarianceOption = WithSigma(0)
var _ RadialVarianceOption = WithAngles(0)
var _ RadialVarianceOption = WithDistance(nil)

var _ ColorMomentOption = WithSize(0, 0)
var _ ColorMomentOption = WithInterpolation(Bilinear)
var _ ColorMomentOption = WithKernelSize(0)
var _ ColorMomentOption = WithSigma(0)
var _ ColorMomentOption = WithDistance(nil)

var _ CLDOption = WithSize(0, 0)
var _ CLDOption = WithInterpolation(Bilinear)
var _ CLDOption = WithDistance(nil)

var _ EHDOption = WithSize(0, 0)
var _ EHDOption = WithInterpolation(Bilinear)
var _ EHDOption = WithDistance(nil)

var _ WHashOption = WithSize(0, 0)
var _ WHashOption = WithInterpolation(Bilinear)
var _ WHashOption = WithLevel(0)
var _ WHashOption = WithDistance(nil)

var _ LBPOption = WithSize(0, 0)
var _ LBPOption = WithInterpolation(Bilinear)
var _ LBPOption = WithGridSize(0, 0)
var _ LBPOption = WithDistance(nil)

var _ HOGHashOption = WithSize(0, 0)
var _ HOGHashOption = WithInterpolation(Bilinear)
var _ HOGHashOption = WithCellSize(0)
var _ HOGHashOption = WithNumBins(0)
var _ HOGHashOption = WithDistance(nil)

var _ PDQOption = WithInterpolation(Bilinear)
var _ PDQOption = WithDistance(nil)

var _ RASHOption = WithSize(0, 0)
var _ RASHOption = WithInterpolation(Bilinear)
var _ RASHOption = WithSigma(0)
var _ RASHOption = WithRings(0)
var _ RASHOption = WithDistance(nil)

var _ ZernikeOption = WithSize(0, 0)
var _ ZernikeOption = WithInterpolation(Bilinear)
var _ ZernikeOption = WithDegree(0)
var _ ZernikeOption = WithDistance(nil)
