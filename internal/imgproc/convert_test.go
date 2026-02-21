package imgproc

import (
	"image"
	"testing"
)

func TestGrayToF32(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 3, 2))
	img.Pix = []uint8{10, 20, 30, 40, 50, 60}
	result := GrayToF32(img)
	if len(result) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result))
	}
	if len(result[0]) != 3 {
		t.Fatalf("expected 3 cols, got %d", len(result[0]))
	}
	if result[0][0] != 10 {
		t.Errorf("expected 10 at [0][0], got %v", result[0][0])
	}
	if result[1][2] != 60 {
		t.Errorf("expected 60 at [1][2], got %v", result[1][2])
	}
}

func TestGrayToF32_singlePixel(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 1, 1))
	img.Pix = []uint8{128}
	result := GrayToF32(img)
	if result[0][0] != 128 {
		t.Errorf("expected 128, got %v", result[0][0])
	}
}
