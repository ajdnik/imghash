package imgproc

import (
	"github.com/r9y9/gossp/dct"
)

func DCT(mat [][]float32) [][]float32 {
	mat64 := matf32Tof64(mat)
	// DCT on rows
	for i := 0; i < len(mat64); i++ {
		mat64[i] = dct.DCTOrthogonal(mat64[i])
	}
	// Transpose
	mat64 = transpose(mat64)
	// DCT on columns
	for i := 0; i < len(mat64); i++ {
		mat64[i] = dct.DCTOrthogonal(mat64[i])
	}
	// Transpose
	mat64 = transpose(mat64)
	res32 := matf64Tof32(mat64)
	return res32
}

func matf32Tof64(mat [][]float32) [][]float64 {
	res := make([][]float64, len(mat))
	for i := 0; i < len(mat); i++ {
		res[i] = make([]float64, len(mat[i]))
		for j := 0; j < len(mat[i]); j++ {
			res[i][j] = float64(mat[i][j])
		}
	}
	return res
}

func matf64Tof32(mat [][]float64) [][]float32 {
	res := make([][]float32, len(mat))
	for i := 0; i < len(mat); i++ {
		res[i] = make([]float32, len(mat[i]))
		for j := 0; j < len(mat[i]); j++ {
			res[i][j] = float32(mat[i][j])
		}
	}
	return res
}

func transpose(mat [][]float64) [][]float64 {
	res := make([][]float64, len(mat))
	width, height := len(mat[0]), len(mat)
	for i := 0; i < height; i++ {
		res[i] = make([]float64, width)
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			res[i][j] = mat[j][i]
		}
	}
	return res
}
