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
