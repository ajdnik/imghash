package imgproc

import (
	"image"
	"image/color"
	"testing"
)

func TestBorderReflect101(t *testing.T) {
	bounds := image.Rect(0, 0, 10, 10)
	tests := []struct {
		name    string
		x, y    int
		expectX int
		expectY int
	}{
		{"inside", 5, 5, 5, 5},
		{"left edge", -1, 5, 1, 5},
		{"right edge", 10, 5, 8, 5},
		{"top edge", 5, -1, 5, 1},
		{"bottom edge", 5, 10, 5, 8},
		{"corner", -1, -1, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rx, ry := borderReflect101(tt.x, tt.y, bounds)
			if rx != tt.expectX || ry != tt.expectY {
				t.Errorf("got (%d, %d), want (%d, %d)", rx, ry, tt.expectX, tt.expectY)
			}
		})
	}
}

func TestFilter2DGray(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 3, 3))
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			img.SetGray(x, y, color.Gray{uint8((x + y) * 20)})
		}
	}
	kernel := [][]float32{
		{0, 0, 0},
		{0, 1, 0},
		{0, 0, 0},
	}
	result := Filter2DGray(img, kernel)
	if len(result) != 3 || len(result[0]) != 3 {
		t.Fatalf("expected 3x3 result, got %dx%d", len(result), len(result[0]))
	}
	if result[1][1] != float32(img.GrayAt(1, 1).Y) {
		t.Errorf("identity kernel: got %v, want %v", result[1][1], img.GrayAt(1, 1).Y)
	}
}

func TestSepFilter2DGray(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 5, 5))
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			img.SetGray(x, y, color.Gray{128})
		}
	}
	kernel := []int{64, 128, 64}
	result := sepFilter2DGray(img, kernel)
	bounds := result.Bounds()
	if bounds.Dx() != 5 || bounds.Dy() != 5 {
		t.Errorf("expected 5x5 result, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestSepFilter2D(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 5, 5))
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			img.Set(x, y, color.RGBA{128, 128, 128, 255})
		}
	}
	kernel := []int{64, 128, 64}
	result := sepFilter2D(img, kernel)
	bounds := result.Bounds()
	if bounds.Dx() != 5 || bounds.Dy() != 5 {
		t.Errorf("expected 5x5 result, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
