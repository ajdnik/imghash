package imgproc

import "image"

// GrayToF32 converts grayscale image to float32 matrix
func GrayToF32(img *image.Gray) [][]float32 {
	bounds := img.Bounds()
	width, height := getSize(img)
	f32Img := make([][]float32, height)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		f32Img[y-bounds.Min.Y] = make([]float32, width)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.GrayAt(x, y).Y
			f32Img[y-bounds.Min.Y][x-bounds.Min.X] = saturateCastUI8ToF32(pixel)
		}
	}
	return f32Img
}
