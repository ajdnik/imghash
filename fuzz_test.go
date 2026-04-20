package imghash

import (
	"image"
	"image/color"
	"testing"
)

// imageFromBytes builds a small NRGBA image from fuzz input.
// Returns nil if data is too short to form a valid image.
func imageFromBytes(data []byte) image.Image {
	if len(data) < 4 {
		return nil
	}
	w := int(data[0])%32 + 1
	h := int(data[1])%32 + 1
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	pixels := data[2:]
	for y := range h {
		for x := range w {
			i := (y*w + x) * 4
			var r, g, b, a byte
			if i < len(pixels) {
				r = pixels[i]
			}
			if i+1 < len(pixels) {
				g = pixels[i+1]
			}
			if i+2 < len(pixels) {
				b = pixels[i+2]
			}
			if i+3 < len(pixels) {
				a = pixels[i+3]
			}
			img.SetNRGBA(x, y, color.NRGBA{R: r, G: g, B: b, A: a})
		}
	}
	return img
}

func FuzzAverageCalculate(f *testing.F) {
	f.Add([]byte{8, 8, 0xFF, 0xAA, 0x55, 0x00})
	f.Fuzz(func(_ *testing.T, data []byte) {
		img := imageFromBytes(data)
		if img == nil {
			return
		}
		h, _ := NewAverage()
		h.Calculate(img) //nolint:errcheck
	})
}

func FuzzDifferenceCalculate(f *testing.F) {
	f.Add([]byte{8, 8, 0xFF, 0xAA, 0x55, 0x00})
	f.Fuzz(func(_ *testing.T, data []byte) {
		img := imageFromBytes(data)
		if img == nil {
			return
		}
		h, _ := NewDifference()
		h.Calculate(img) //nolint:errcheck
	})
}

func FuzzPHashCalculate(f *testing.F) {
	f.Add([]byte{8, 8, 0xFF, 0xAA, 0x55, 0x00})
	f.Fuzz(func(_ *testing.T, data []byte) {
		img := imageFromBytes(data)
		if img == nil {
			return
		}
		h, _ := NewPHash()
		h.Calculate(img) //nolint:errcheck
	})
}

func FuzzMedianCalculate(f *testing.F) {
	f.Add([]byte{8, 8, 0xFF, 0xAA, 0x55, 0x00})
	f.Fuzz(func(_ *testing.T, data []byte) {
		img := imageFromBytes(data)
		if img == nil {
			return
		}
		h, _ := NewMedian()
		h.Calculate(img) //nolint:errcheck
	})
}
