package imgproc

import (
	"image"
	"math"
	"sort"
)

// Mean returns the arithmetic mean of all pixel values in the grayscale image.
func Mean(img *image.Gray) (float64, error) {
	if img == nil {
		return 0, ErrImageIsNil
	}
	bounds := img.Bounds()
	width, height := getSize(img)
	totalPixels := float64(width * height)
	if totalPixels == 0 {
		return 0, nil
	}
	var sum float64
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixel := img.GrayAt(x, y).Y
			sum += float64(pixel)
		}
	}
	return sum / totalPixels, nil
}

// Median returns the median pixel value of the grayscale image.
func Median(img *image.Gray) (float64, error) {
	if img == nil {
		return 0, ErrImageIsNil
	}
	bounds := img.Bounds()
	width, height := getSize(img)
	totalPixels := width * height
	if totalPixels == 0 {
		return 0, nil
	}
	pixels := make([]int, totalPixels)
	count := 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixels[count] = int(img.GrayAt(x, y).Y)
			count++
		}
	}
	sort.Ints(pixels)
	if totalPixels%2 == 0 {
		return float64(pixels[totalPixels/2-1]+pixels[totalPixels/2+1]) / 2, nil
	}
	return float64(pixels[totalPixels/2]), nil
}

// MedianF32 returns the median of all values in a 2-D float32 matrix.
func MedianF32(mat [][]float32) float32 {
	var n int
	for _, row := range mat {
		n += len(row)
	}
	vals := make([]float32, 0, n)
	for _, row := range mat {
		vals = append(vals, row...)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	if n%2 == 0 {
		return (vals[n/2-1] + vals[n/2]) / 2
	}
	return vals[n/2]
}

func getSize(img image.Image) (int, int) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	return width, height
}

// cvRound rounds floating-point number to the nearest integer
func cvRound(value float64) int {
	return int(math.Round(value))
}

// Normalize normalize images to [0,1] range
func Normalize(img [][]float32) {
	var hi, lo float32
	hi = -1 << 30
	lo = 1 << 30
	for _, list := range img {
		for _, v := range list {
			if v > hi {
				hi = v
			}
			if v < lo {
				lo = v
			}
		}
	}

	diff := hi - lo
	for x, list := range img {
		for y, v := range list {
			(img)[x][y] = (v - lo) / diff
		}
	}
}
