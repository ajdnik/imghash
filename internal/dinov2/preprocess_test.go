package dinov2

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestImageToTensor_KnownRGBPixel(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, ImageSize, ImageSize))
	for y := 0; y < ImageSize; y++ {
		for x := 0; x < ImageSize; x++ {
			img.Set(x, y, color.RGBA{R: 128, G: 64, B: 200, A: 255})
		}
	}
	out := ImageToTensor(img)
	wantR := (float32(128)/255 - imagenetMean[0]) / imagenetStd[0]
	wantG := (float32(64)/255 - imagenetMean[1]) / imagenetStd[1]
	wantB := (float32(200)/255 - imagenetMean[2]) / imagenetStd[2]
	for i := 0; i < PixelCount; i++ {
		if math.Abs(float64(out[i]-wantR)) > 1e-5 {
			t.Fatalf("R channel mismatch at idx %d: got %v want %v", i, out[i], wantR)
		}
		if math.Abs(float64(out[PixelCount+i]-wantG)) > 1e-5 {
			t.Fatalf("G channel mismatch at idx %d: got %v want %v", i, out[PixelCount+i], wantG)
		}
		if math.Abs(float64(out[2*PixelCount+i]-wantB)) > 1e-5 {
			t.Fatalf("B channel mismatch at idx %d: got %v want %v", i, out[2*PixelCount+i], wantB)
		}
	}
}

func TestImageToTensor_Length(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	out := ImageToTensor(img)
	if got, want := len(out), Channels*PixelCount; got != want {
		t.Errorf("tensor length = %d, want %d", got, want)
	}
}

func TestImageToTensor_NCHWLayout(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, ImageSize, ImageSize))
	for y := 0; y < ImageSize; y++ {
		for x := 0; x < ImageSize; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 128, B: 0, A: 255})
		}
	}
	out := ImageToTensor(img)
	wantR := (float32(255)/255 - imagenetMean[0]) / imagenetStd[0]
	wantG := (float32(128)/255 - imagenetMean[1]) / imagenetStd[1]
	wantB := (float32(0)/255 - imagenetMean[2]) / imagenetStd[2]
	if math.Abs(float64(out[0]-wantR)) > 1e-5 {
		t.Errorf("channel 0 head = %v, want %v (red)", out[0], wantR)
	}
	if math.Abs(float64(out[PixelCount]-wantG)) > 1e-5 {
		t.Errorf("channel 1 head = %v, want %v (green)", out[PixelCount], wantG)
	}
	if math.Abs(float64(out[2*PixelCount]-wantB)) > 1e-5 {
		t.Errorf("channel 2 head = %v, want %v (blue)", out[2*PixelCount], wantB)
	}
}

func TestImageToTensor_GraySubimage(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 40, 30))
	// A nonzero origin and padded stride catch incorrect pixel indexing.
	sub := img.SubImage(image.Rect(5, 7, 25, 17)).(*image.Gray)
	for y := 7; y < 17; y++ {
		for x := 5; x < 25; x++ {
			sub.SetGray(x, y, color.Gray{Y: 128})
		}
	}
	got := ImageToTensor(sub)
	if len(got) != Channels*PixelCount {
		t.Fatalf("tensor length = %d, want %d", len(got), Channels*PixelCount)
	}
	for c := 0; c < Channels; c++ {
		want := (float32(128)/255 - imagenetMean[c]) / imagenetStd[c]
		for i := 0; i < PixelCount; i++ {
			if math.Abs(float64(got[c*PixelCount+i]-want)) > 1e-5 {
				t.Fatalf("channel %d pixel %d = %v, want %v", c, i, got[c*PixelCount+i], want)
			}
		}
	}
}
