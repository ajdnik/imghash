// Package imgproc provides low-level image processing primitives used
// by the perceptual hashing algorithms in the parent package.
package imgproc

import (
	"image"
	"math"
)

// GaussianBlur computes and returns an image with an applied Gaussian filter.
// Both kernel and sigma are parameters used for generating a Gaussian filter kernel.
func GaussianBlur(img image.Image, kernel int, sigma float64) image.Image {
	// If kernel size is zero compute it from sigma.
	if kernel == 0 && sigma > 0 {
		kernel = cvRound(sigma*3*2+1) | 1
	}
	k := getGaussianKernel(kernel, sigma)
	kernelInt := make([]int, len(k))
	for i := 0; i < len(k); i++ {
		kernelInt[i] = int(k[i] * 256)
	}
	switch i := img.(type) {
	case *image.Gray:
		return sepFilter2DGray(i, kernelInt)
	default:
		return sepFilter2D(i, kernelInt)
	}
}

// The function computes and returns a slice representing Gaussian filter coefficients.
// The size represents the filter aperture size and should be an odd and positive number.
// The sigma represents the Gaussian standard deviation.
func getGaussianKernel(size int, sigma float64) []float32 {
	smallGaussianSize := 7
	smallGaussianTab := [][]float32{
		{1, 0, 0, 0, 0, 0, 0},
		{0.25, 0.5, 0.25, 0, 0, 0, 0},
		{0.0625, 0.25, 0.375, 0.25, 0.0625, 0, 0},
		{0.03125, 0.109375, 0.21875, 0.28125, 0.21875, 0.109375, 0.03125},
	}
	sigmaX := (float64(size-1)*0.5-1)*0.3 + 0.8
	if sigma > 0 {
		sigmaX = sigma
	}
	scale2x := -0.5 / (sigmaX * sigmaX)
	kernel := make([]float32, size)
	var sum float64
	for i := 0; i < size; i++ {
		x := float64(i) - float64(size-1)*0.5
		t := math.Exp(scale2x * x * x)
		if size%2 == 1 && size <= smallGaussianSize && sigma <= 0 {
			t = float64(smallGaussianTab[size>>1][i])
		}
		kernel[i] = float32(t)
		sum += float64(kernel[i])
	}
	sum = 1 / sum
	for i := 0; i < size; i++ {
		kernel[i] = float32(float64(kernel[i]) * sum)
	}
	return kernel
}
