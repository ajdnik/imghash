package imghash_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var cldCalculateTests = []struct {
	filename   string
	hash       hashtype.UInt8
	width      uint
	height     uint
	resizeType imghash.Interpolation
}{
	{"assets/lena.jpg", hashtype.UInt8{31, 11, 19, 16, 20, 15, 29, 17, 14, 42, 18, 17}, 64, 64, imghash.Bilinear},
	{"assets/baboon.jpg", hashtype.UInt8{32, 16, 13, 15, 13, 14, 29, 14, 17, 33, 18, 12}, 64, 64, imghash.Bilinear},
	{"assets/cat.jpg", hashtype.UInt8{36, 26, 31, 27, 19, 16, 28, 15, 17, 35, 16, 13}, 64, 64, imghash.Bilinear},
	{"assets/monarch.jpg", hashtype.UInt8{26, 18, 15, 12, 19, 17, 27, 18, 16, 38, 19, 15}, 64, 64, imghash.Bilinear},
	{"assets/peppers.jpg", hashtype.UInt8{30, 16, 20, 16, 11, 16, 24, 14, 15, 37, 17, 17}, 64, 64, imghash.Bilinear},
	{"assets/tulips.jpg", hashtype.UInt8{26, 20, 17, 11, 15, 8, 30, 16, 16, 33, 17, 16}, 64, 64, imghash.Bilinear},
}

func TestCLD_Calculate(t *testing.T) {
	for _, tt := range cldCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := imghash.NewCLD(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType))
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

func ExampleCLD_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	cld, err := imghash.NewCLD()
	if err != nil {
		panic(err)
	}
	hash, err := cld.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [36 26 31 27 19 16 28 15 17 35 16 13]
}

var cldDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imghash.Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 24.49489742783178, 64, 64, imghash.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 11.874342087037917, 64, 64, imghash.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 25.199206336708304, 64, 64, imghash.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 11.532562594670797, 64, 64, imghash.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 12.206555615733702, 64, 64, imghash.Bilinear},
}

func TestCLD_Distance(t *testing.T) {
	for _, tt := range cldDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := imghash.NewCLD(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType))
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

func TestNewCLD_Errors(t *testing.T) {
	tests := []struct {
		name string
		opts []imghash.CLDOption
		err  error
	}{
		{"zero width", []imghash.CLDOption{imghash.WithSize(0, 64)}, imghash.ErrInvalidSize},
		{"zero height", []imghash.CLDOption{imghash.WithSize(64, 0)}, imghash.ErrInvalidSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := imghash.NewCLD(tt.opts...)
			if !errors.Is(err, tt.err) {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}
