package vit

import "math"

// layerNorm applies LayerNormalization in place to each row of x ([rows, cols]).
// Mean and variance are computed per row. eps is added before sqrt for
// numerical stability. gamma (scale) and beta (bias) have length cols.
// DINOv2 uses eps = 1e-6.
func layerNorm(x []float32, gamma, beta []float32, rows, cols int, eps float32) {
	cf := float64(cols)
	for r := 0; r < rows; r++ {
		row := x[r*cols : r*cols+cols]
		var sum float64
		for _, v := range row {
			sum += float64(v)
		}
		mean := sum / cf
		var ss float64
		for _, v := range row {
			d := float64(v) - mean
			ss += d * d
		}
		variance := ss / cf
		inv := float32(1 / math.Sqrt(variance+float64(eps)))
		meanF := float32(mean)
		for i, v := range row {
			row[i] = (v-meanF)*inv*gamma[i] + beta[i]
		}
	}
}
