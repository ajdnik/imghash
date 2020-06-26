package imgproc

import (
	"image"
	"image/color"
)

// EqualizeHist equalizes the histogram of the input image to
// normalize the brightness and increases the contrast of the image.
func EqualizeHist(img *image.Gray) *image.Gray {
	const histogramSize = 256
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	histogram := [histogramSize]int{}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixel := img.GrayAt(x, y).Y
			histogram[pixel]++
		}
	}
	i := 0
	for histogram[i] == 0 {
		i++
	}
	total := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	scale := (float32(histogramSize) - 1) / float32(total-histogram[i])
	lut := [histogramSize]uint8{}
	sum := 0
	i++
	for i < histogramSize {
		sum += histogram[i]
		lutVal := float32(sum) * scale
		lut[i] = saturateCastF32ToUI8(lutVal)
		i++
	}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixel := img.GrayAt(x, y).Y
			gray.SetGray(x, y, color.Gray{lut[pixel]})
		}
	}
	return gray
}
