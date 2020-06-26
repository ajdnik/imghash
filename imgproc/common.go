package imgproc

import (
	"image"
	"math"
	"sort"
)

// Mean ...
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

// Median ...
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
