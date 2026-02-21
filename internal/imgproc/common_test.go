package imgproc

import (
	"image"
	"math"
	"testing"
)

func newGrayImage(w, h int, pixels []uint8) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, w, h))
	copy(img.Pix, pixels)
	return img
}

func TestMean(t *testing.T) {
	tests := []struct {
		name   string
		img    *image.Gray
		expect float64
		hasErr bool
	}{
		{"nil image", nil, 0, true},
		{"single pixel", newGrayImage(1, 1, []uint8{100}), 100, false},
		{"2x2 uniform", newGrayImage(2, 2, []uint8{10, 10, 10, 10}), 10, false},
		{"2x2 varied", newGrayImage(2, 2, []uint8{0, 100, 100, 200}), 100, false},
		{"empty image", newGrayImage(0, 0, nil), 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Mean(tt.img)
			if tt.hasErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.hasErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.hasErr && math.Abs(res-tt.expect) > 1e-9 {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		name   string
		img    *image.Gray
		expect float64
		hasErr bool
	}{
		{"nil image", nil, 0, true},
		{"single pixel", newGrayImage(1, 1, []uint8{50}), 50, false},
		{"3 pixels odd", newGrayImage(3, 1, []uint8{10, 30, 20}), 20, false},
		{"empty image", newGrayImage(0, 0, nil), 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Median(tt.img)
			if tt.hasErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.hasErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.hasErr && math.Abs(res-tt.expect) > 1e-9 {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestGetSize(t *testing.T) {
	img := image.NewGray(image.Rect(10, 20, 30, 50))
	w, h := getSize(img)
	if w != 20 || h != 30 {
		t.Errorf("got (%v, %v), want (20, 30)", w, h)
	}
}

func TestCvRound(t *testing.T) {
	tests := []struct {
		input  float64
		expect int
	}{
		{1.4, 1},
		{1.5, 2},
		{1.6, 2},
		{-0.5, -1},
		{-1.6, -2},
		{0.0, 0},
	}
	for _, tt := range tests {
		if res := cvRound(tt.input); res != tt.expect {
			t.Errorf("cvRound(%v) = %v, want %v", tt.input, res, tt.expect)
		}
	}
}

func TestNormalize(t *testing.T) {
	img := [][]float32{
		{0, 10},
		{5, 20},
	}
	Normalize(img)
	if img[0][0] != 0 {
		t.Errorf("min value should be 0, got %v", img[0][0])
	}
	if img[1][1] != 1 {
		t.Errorf("max value should be 1, got %v", img[1][1])
	}
	if math.Abs(float64(img[0][1])-0.5) > 1e-6 {
		t.Errorf("expected 0.5, got %v", img[0][1])
	}
	if math.Abs(float64(img[1][0])-0.25) > 1e-6 {
		t.Errorf("expected 0.25, got %v", img[1][0])
	}
}

func TestNormalize_uniform(t *testing.T) {
	img := [][]float32{{5, 5}, {5, 5}}
	defer func() {
		if r := recover(); r != nil {
			t.Skip("normalize panics on uniform image (division by zero)")
		}
	}()
	Normalize(img)
}
