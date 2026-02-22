package imghash_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var ehdCalculateTests = []struct {
	filename   string
	hash       hashtype.UInt8
	width      uint
	height     uint
	resizeType imghash.Interpolation
}{
	{"assets/lena.jpg", hashtype.UInt8{4, 0, 1, 0, 0, 1, 1, 3, 1, 2, 2, 0, 0, 3, 1, 0, 0, 5, 3, 1, 5, 0, 0, 0, 0, 3, 1, 7, 3, 5, 2, 1, 5, 3, 2, 1, 0, 4, 0, 1, 6, 0, 4, 2, 2, 5, 1, 5, 3, 4, 6, 1, 4, 3, 1, 2, 0, 2, 0, 0, 6, 1, 3, 4, 2, 4, 1, 5, 4, 3, 2, 0, 1, 1, 0, 3, 1, 7, 2, 1}, 256, 256, imghash.Bilinear},
	{"assets/baboon.jpg", hashtype.UInt8{2, 3, 5, 6, 6, 3, 4, 5, 6, 5, 3, 4, 6, 5, 5, 2, 4, 6, 5, 6, 2, 4, 5, 6, 5, 5, 1, 4, 5, 2, 4, 1, 5, 3, 2, 2, 5, 5, 6, 5, 2, 4, 6, 5, 5, 4, 2, 5, 6, 3, 4, 2, 6, 4, 4, 2, 4, 4, 6, 5, 4, 2, 6, 3, 5, 3, 2, 6, 3, 5, 4, 2, 3, 6, 5, 3, 3, 4, 6, 4}, 256, 256, imghash.Bilinear},
	{"assets/cat.jpg", hashtype.UInt8{0, 2, 1, 3, 0, 1, 1, 1, 4, 0, 2, 2, 4, 4, 0, 2, 5, 4, 6, 1, 1, 4, 4, 6, 2, 2, 4, 6, 6, 4, 2, 4, 6, 6, 4, 4, 4, 5, 6, 3, 2, 5, 4, 6, 5, 3, 4, 6, 6, 4, 2, 5, 6, 6, 4, 1, 6, 6, 6, 3, 1, 6, 6, 5, 4, 2, 5, 7, 5, 4, 2, 6, 7, 5, 2, 2, 6, 6, 6, 2}, 256, 256, imghash.Bilinear},
	{"assets/monarch.jpg", hashtype.UInt8{5, 1, 5, 5, 2, 3, 1, 3, 4, 0, 1, 0, 1, 2, 0, 0, 0, 0, 0, 0, 5, 1, 5, 6, 2, 4, 2, 5, 4, 3, 2, 3, 2, 6, 2, 1, 0, 1, 3, 2, 5, 2, 6, 7, 2, 5, 1, 5, 6, 3, 4, 1, 3, 5, 3, 0, 0, 0, 0, 0, 5, 2, 6, 6, 2, 3, 1, 3, 3, 2, 1, 0, 1, 2, 1, 0, 0, 0, 0, 0}, 256, 256, imghash.Bilinear},
	{"assets/peppers.jpg", hashtype.UInt8{3, 2, 6, 2, 1, 1, 1, 3, 1, 1, 2, 1, 3, 4, 1, 3, 1, 4, 3, 1, 3, 1, 6, 2, 1, 3, 1, 3, 3, 1, 4, 1, 4, 3, 1, 4, 2, 4, 5, 1, 3, 1, 3, 3, 1, 4, 0, 3, 3, 0, 2, 1, 4, 5, 1, 3, 0, 4, 2, 1, 3, 1, 2, 3, 0, 2, 2, 1, 4, 1, 1, 2, 4, 4, 1, 2, 2, 4, 5, 1}, 256, 256, imghash.Bilinear},
	{"assets/tulips.jpg", hashtype.UInt8{5, 1, 4, 5, 2, 6, 1, 5, 3, 2, 5, 1, 5, 5, 2, 6, 1, 5, 4, 2, 5, 1, 6, 5, 2, 6, 1, 6, 4, 3, 5, 1, 5, 5, 3, 6, 1, 6, 4, 3, 5, 2, 6, 5, 4, 5, 1, 4, 6, 4, 6, 1, 3, 4, 2, 5, 1, 5, 4, 3, 3, 2, 5, 6, 4, 6, 1, 5, 4, 2, 5, 1, 4, 4, 2, 4, 0, 2, 2, 1}, 256, 256, imghash.Bilinear},
}

func TestEHD_Calculate(t *testing.T) {
	for _, tt := range ehdCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := imghash.NewEHD(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img, err := imghash.OpenImage(tt.filename)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.filename, err)
			}
			result, err := hash.Calculate(img)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			res := result.(hashtype.UInt8)
			if !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleEHD_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	ehd, err := imghash.NewEHD()
	if err != nil {
		panic(err)
	}
	hash, err := ehd.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [0 2 1 3 0 1 1 1 4 0 2 2 4 4 0 2 5 4 6 1 1 4 4 6 2 2 4 6 6 4 2 4 6 6 4 4 4 5 6 3 2 5 4 6 5 3 4 6 6 4 2 5 6 6 4 1 6 6 6 3 1 6 6 5 4 2 5 7 5 4 2 6 7 5 2 2 6 6 6 2]
}

var ehdDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imghash.Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 206, 256, 256, imghash.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 123, 256, 256, imghash.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 147, 256, 256, imghash.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 168, 256, 256, imghash.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 129, 256, 256, imghash.Bilinear},
}

func TestEHD_Distance(t *testing.T) {
	for _, tt := range ehdDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := imghash.NewEHD(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img1, err := imghash.OpenImage(tt.firstImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.firstImage, err)
			}
			img2, err := imghash.OpenImage(tt.secondImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.secondImage, err)
			}
			h1, err := hash.Calculate(img1)
			if err != nil {
				t.Fatalf("failed to calculate hash for %s: %v", tt.firstImage, err)
			}
			h2, err := hash.Calculate(img2)
			if err != nil {
				t.Fatalf("failed to calculate hash for %s: %v", tt.secondImage, err)
			}
			dist, err := hash.Compare(h1, h2)
			if err != nil {
				t.Fatalf("failed to compute distance: %v", err)
			}
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}

func TestNewEHD_defaults(t *testing.T) {
	ehd, err := imghash.NewEHD()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := ehd.Calculate(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 80 {
		t.Errorf("expected hash length 80, got %d", h.Len())
	}
}

func TestNewEHD_Errors(t *testing.T) {
	tests := []struct {
		name string
		opts []imghash.EHDOption
		err  error
	}{
		{"zero width", []imghash.EHDOption{imghash.WithSize(0, 256)}, imghash.ErrInvalidSize},
		{"zero height", []imghash.EHDOption{imghash.WithSize(256, 0)}, imghash.ErrInvalidSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := imghash.NewEHD(tt.opts...)
			if !errors.Is(err, tt.err) {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}
