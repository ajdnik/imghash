package imgproc

import (
	"image"
	"image/color"
	"testing"
)

func TestEqualizeHist(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			img.SetGray(x, y, color.Gray{uint8(x*16 + y*4)})
		}
	}
	result := EqualizeHist(img)
	if result.Bounds() != img.Bounds() {
		t.Errorf("bounds mismatch: got %v, want %v", result.Bounds(), img.Bounds())
	}
}

func TestEqualizeHist_uniform(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 3, 3))
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			img.SetGray(x, y, color.Gray{100})
		}
	}
	img.SetGray(0, 0, color.Gray{50})
	result := EqualizeHist(img)
	if result.Bounds() != img.Bounds() {
		t.Errorf("bounds mismatch")
	}
}
