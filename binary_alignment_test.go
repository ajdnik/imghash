package imghash_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
)

func TestBinaryHashers_NonByteAlignedCapacity(t *testing.T) {
	img := testGradientGray(21, 21)
	tests := []struct {
		name      string
		wantBytes int
		newHasher func() (imghash.Hasher, error)
	}{
		{
			name:      "average 3x3",
			wantBytes: 2,
			newHasher: func() (imghash.Hasher, error) {
				return imghash.NewAverage(imghash.WithSize(3, 3))
			},
		},
		{
			name:      "median 3x3",
			wantBytes: 2,
			newHasher: func() (imghash.Hasher, error) {
				return imghash.NewMedian(imghash.WithSize(3, 3))
			},
		},
		{
			name:      "difference 3x3",
			wantBytes: 2,
			newHasher: func() (imghash.Hasher, error) {
				return imghash.NewDifference(imghash.WithSize(3, 3))
			},
		},
		{
			name:      "wavelet 3x3",
			wantBytes: 2,
			newHasher: func() (imghash.Hasher, error) {
				return imghash.NewWHash(imghash.WithSize(3, 3), imghash.WithLevel(1))
			},
		},
		{
			name:      "rash rings 5",
			wantBytes: 1,
			newHasher: func() (imghash.Hasher, error) {
				return imghash.NewRASH(imghash.WithSize(32, 32), imghash.WithRings(5))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher, err := tt.newHasher()
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}

			hash, err := hasher.Calculate(img)
			if err != nil {
				t.Fatalf("failed to calculate hash: %v", err)
			}

			bh, ok := hash.(hashtype.Binary)
			if !ok {
				t.Fatalf("expected binary hash, got %T", hash)
			}

			if got := bh.Len(); got != tt.wantBytes {
				t.Fatalf("unexpected hash byte length: got %d, want %d", got, tt.wantBytes)
			}
		})
	}
}

func testGradientGray(width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{Y: uint8((x*17 + y*13) % 256)})
		}
	}
	return img
}
