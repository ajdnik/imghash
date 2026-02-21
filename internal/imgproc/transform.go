package imgproc

import "math"

// DCT computes a 2-D orthogonal DCT-II on a float32 matrix.
func DCT(mat [][]float32) [][]float32 {
	mat64 := matf32Tof64(mat)
	for i := range mat64 {
		mat64[i] = dctOrthogonal(mat64[i])
	}
	mat64 = transpose(mat64)
	for i := range mat64 {
		mat64[i] = dctOrthogonal(mat64[i])
	}
	mat64 = transpose(mat64)
	return matf64Tof32(mat64)
}

// dctOrthogonal computes a 1-D orthogonal DCT-II.
func dctOrthogonal(x []float64) []float64 {
	n := len(x)
	c0 := math.Sqrt(1.0 / float64(n))
	c1 := math.Sqrt(2.0 / float64(n))
	result := make([]float64, n)
	for k := 0; k < n; k++ {
		var sum float64
		for i := 0; i < n; i++ {
			sum += x[i] * math.Cos(math.Pi*float64(2*i+1)*float64(k)/float64(2*n))
		}
		if k == 0 {
			result[k] = sum * c0
		} else {
			result[k] = sum * c1
		}
	}
	return result
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

// HaarDWT2D applies a multi-level 2-D Haar discrete wavelet transform
// in-place on the top-left region of mat. After `levels` iterations the
// top-left (rows/2^levels)×(cols/2^levels) block holds the LL coefficients.
func HaarDWT2D(mat [][]float32, levels int) {
	rows := len(mat)
	if rows == 0 {
		return
	}
	cols := len(mat[0])
	for l := 0; l < levels; l++ {
		h := rows >> uint(l)
		w := cols >> uint(l)
		if h < 2 || w < 2 {
			break
		}
		haarRows(mat, h, w)
		haarCols(mat, h, w)
	}
}

func haarRows(mat [][]float32, rows, cols int) {
	half := cols / 2
	tmp := make([]float32, cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < half; c++ {
			tmp[c] = (mat[r][2*c] + mat[r][2*c+1]) / 2
			tmp[half+c] = (mat[r][2*c] - mat[r][2*c+1]) / 2
		}
		copy(mat[r][:cols], tmp)
	}
}

func haarCols(mat [][]float32, rows, cols int) {
	half := rows / 2
	tmp := make([]float32, rows)
	for c := 0; c < cols; c++ {
		for r := 0; r < half; r++ {
			tmp[r] = (mat[2*r][c] + mat[2*r+1][c]) / 2
			tmp[half+r] = (mat[2*r][c] - mat[2*r+1][c]) / 2
		}
		for r := 0; r < rows; r++ {
			mat[r][c] = tmp[r]
		}
	}
}

func transpose(mat [][]float64) [][]float64 {
	height := len(mat)
	width := len(mat[0])
	res := make([][]float64, height)
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
