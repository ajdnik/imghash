package imgproc

import (
	"image"
	"image/color"
	"testing"
)

func makeTestImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{uint8(x * 10), uint8(y * 10), 100, 255})
		}
	}
	return img
}

func makeGrayTestImage(w, h int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.SetGray(x, y, color.Gray{uint8((x + y) * 10)})
		}
	}
	return img
}

func TestResize(t *testing.T) {
	types := []struct {
		name string
		typ  ResizeType
	}{
		{"NearestNeighbor", NearestNeighbor},
		{"Bilinear", Bilinear},
		{"Bicubic", Bicubic},
		{"MitchellNetravali", MitchellNetravali},
		{"Lanczos2", Lanczos2},
		{"Lanczos3", Lanczos3},
		{"BilinearExact", BilinearExact},
	}
	for _, tt := range types {
		t.Run(tt.name+"_rgba", func(t *testing.T) {
			img := makeTestImage(10, 10)
			result := Resize(5, 5, img, tt.typ)
			bounds := result.Bounds()
			if bounds.Dx() != 5 || bounds.Dy() != 5 {
				t.Errorf("got %dx%d, want 5x5", bounds.Dx(), bounds.Dy())
			}
		})
		t.Run(tt.name+"_gray", func(t *testing.T) {
			img := makeGrayTestImage(10, 10)
			result := Resize(5, 5, img, tt.typ)
			bounds := result.Bounds()
			if bounds.Dx() != 5 || bounds.Dy() != 5 {
				t.Errorf("got %dx%d, want 5x5", bounds.Dx(), bounds.Dy())
			}
		})
	}
}

func TestResize_defaultInterpolator(t *testing.T) {
	img := makeTestImage(8, 8)
	result := Resize(4, 4, img, ResizeType(99))
	bounds := result.Bounds()
	if bounds.Dx() != 4 || bounds.Dy() != 4 {
		t.Errorf("got %dx%d, want 4x4", bounds.Dx(), bounds.Dy())
	}
}

func TestMitchellNetravaliAt(t *testing.T) {
	if mitchellNetravaliAt(0) == 0 {
		t.Error("expected non-zero at t=0")
	}
	if mitchellNetravaliAt(0.5) == 0 {
		t.Error("expected non-zero at t=0.5")
	}
	if mitchellNetravaliAt(1.5) == 0 {
		t.Error("expected non-zero at t=1.5")
	}
	if mitchellNetravaliAt(3) != 0 {
		t.Error("expected zero at t=3")
	}
	if mitchellNetravaliAt(-0.5) != mitchellNetravaliAt(0.5) {
		t.Error("kernel should be symmetric")
	}
}

func TestLanczosAt(t *testing.T) {
	l2 := lanczosAt(2)
	l3 := lanczosAt(3)

	if l2(0) != 1 {
		t.Errorf("lanczos2(0) = %v, want 1", l2(0))
	}
	if l3(0) != 1 {
		t.Errorf("lanczos3(0) = %v, want 1", l3(0))
	}
	if l2(2) != 0 {
		t.Errorf("lanczos2(2) = %v, want 0", l2(2))
	}
	if l3(3) != 0 {
		t.Errorf("lanczos3(3) = %v, want 0", l3(3))
	}
	if l2(1) == 0 {
		t.Error("lanczos2(1) should be non-zero")
	}
}
