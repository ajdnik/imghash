package imgproc

import (
	"image"
	"image/color"
	"testing"
)

func TestGetMoments_gray(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			img.SetGray(x, y, color.Gray{uint8(x*10 + y*5 + 50)})
		}
	}
	moments := GetMoments(img)
	if len(moments) != 1 {
		t.Fatalf("expected 1 moment for gray, got %d", len(moments))
	}
	if moments[0].m00 == 0 {
		t.Error("m00 should not be zero for non-black image")
	}
}

func TestGetMoments_rgba(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			img.Set(x, y, color.RGBA{uint8(x * 50), uint8(y * 80), 100, 255})
		}
	}
	moments := GetMoments(img)
	if len(moments) != 3 {
		t.Fatalf("expected 3 moments for RGBA, got %d", len(moments))
	}
}

func TestHuMoments(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			img.SetGray(x, y, color.Gray{uint8(x*20 + y*10 + 30)})
		}
	}
	moments := GetMoments(img)
	hu := HuMoments(moments)
	if len(hu) != 7 {
		t.Fatalf("expected 7 Hu moments, got %d", len(hu))
	}
}

func TestHuMoments_rgba(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			img.Set(x, y, color.RGBA{uint8(x*50 + 10), uint8(y*80 + 20), uint8(x*y*10 + 30), 255})
		}
	}
	moments := GetMoments(img)
	hu := HuMoments(moments)
	if len(hu) != 21 {
		t.Fatalf("expected 21 Hu moments (7*3 channels), got %d", len(hu))
	}
}

func TestGetMoments_blackImage(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 3, 3))
	moments := GetMoments(img)
	if moments[0].m00 != 0 {
		t.Errorf("m00 should be 0 for black image, got %v", moments[0].m00)
	}
}
