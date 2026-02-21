package imgproc

import (
	"errors"
	"image"
	"image/color"
	"testing"
)

func TestGrayscale_nil(t *testing.T) {
	_, err := Grayscale(nil)
	if !errors.Is(err, ErrImageIsNil) {
		t.Errorf("got %v, want %v", err, ErrImageIsNil)
	}
}

func TestGrayscale_valid(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 3, 3))
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			rgba.Set(x, y, color.RGBA{uint8(x * 50), uint8(y * 80), 100, 255})
		}
	}
	gray, err := Grayscale(rgba)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gray.Bounds() != rgba.Bounds() {
		t.Errorf("bounds mismatch: got %v, want %v", gray.Bounds(), rgba.Bounds())
	}
}

func TestYCrCb_nil(t *testing.T) {
	_, err := YCrCb(nil)
	if !errors.Is(err, ErrImageIsNil) {
		t.Errorf("got %v, want %v", err, ErrImageIsNil)
	}
}

func TestYCrCb_valid(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 2))
	rgba.Set(0, 0, color.RGBA{255, 0, 0, 255})
	rgba.Set(1, 0, color.RGBA{0, 255, 0, 255})
	rgba.Set(0, 1, color.RGBA{0, 0, 255, 255})
	rgba.Set(1, 1, color.RGBA{255, 255, 255, 255})
	result, err := YCrCb(rgba)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Bounds() != rgba.Bounds() {
		t.Errorf("bounds mismatch")
	}
}

func TestHSV_nil(t *testing.T) {
	_, err := HSV(nil)
	if !errors.Is(err, ErrImageIsNil) {
		t.Errorf("got %v, want %v", err, ErrImageIsNil)
	}
}

func TestHSV_valid(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 2))
	rgba.Set(0, 0, color.RGBA{255, 0, 0, 255})
	rgba.Set(1, 0, color.RGBA{0, 255, 0, 255})
	rgba.Set(0, 1, color.RGBA{0, 0, 255, 255})
	rgba.Set(1, 1, color.RGBA{128, 128, 128, 255})
	result, err := HSV(rgba)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Bounds() != rgba.Bounds() {
		t.Errorf("bounds mismatch")
	}
}

func TestRgbToYCbCr(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
	}{
		{"white", 255, 255, 255},
		{"black", 0, 0, 0},
		{"red", 255, 0, 0},
		{"green", 0, 255, 0},
		{"blue", 0, 0, 255},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			y, cb, cr := rgbToYCbCr(tt.r, tt.g, tt.b)
			_ = y
			_ = cb
			_ = cr
		})
	}
}

func TestRgbToHSV(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
	}{
		{"white", 255, 255, 255},
		{"black", 0, 0, 0},
		{"red", 255, 0, 0},
		{"green", 0, 255, 0},
		{"blue", 0, 0, 255},
		{"yellow", 255, 255, 0},
		{"mid gray", 128, 128, 128},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			h, s, v := rgbToHSV(tt.r, tt.g, tt.b)
			_ = h
			_ = s
			_ = v
		})
	}
}
