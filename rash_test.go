package imghash_test

import (
	"fmt"
	"testing"

	. "github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var rashCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType Interpolation
	sigma      float64
	rings      int
}{
	{"assets/lena.jpg", hashtype.Binary{6, 0, 99, 12, 0, 253, 255, 255}, 256, 256, Bilinear, 1, 180},
	{"assets/baboon.jpg", hashtype.Binary{241, 27, 3, 0, 0, 240, 255, 255}, 256, 256, Bilinear, 1, 180},
	{"assets/cat.jpg", hashtype.Binary{46, 0, 64, 32, 128, 255, 255, 255}, 256, 256, Bilinear, 1, 180},
	{"assets/monarch.jpg", hashtype.Binary{255, 63, 12, 0, 0, 144, 248, 255}, 256, 256, Bilinear, 1, 180},
	{"assets/peppers.jpg", hashtype.Binary{9, 0, 40, 32, 168, 251, 255, 255}, 256, 256, Bilinear, 1, 180},
	{"assets/tulips.jpg", hashtype.Binary{119, 224, 1, 128, 0, 240, 255, 255}, 256, 256, Bilinear, 1, 180},
}

func TestRASH_Calculate(t *testing.T) {
	for _, tt := range rashCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := NewRASH(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType), WithSigma(tt.sigma), WithRings(tt.rings))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img, err := OpenImage(tt.filename)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.filename, err)
			}
			result, err := hash.Calculate(img)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			res := result.(hashtype.Binary)
			if !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleRASH_Calculate() {
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	rash, err := NewRASH()
	if err != nil {
		panic(err)
	}
	hash, err := rash.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [46 0 64 32 128 255 255 255]
}

var rashDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  Interpolation
	sigma       float64
	rings       int
}{
	{"assets/lena.jpg", "assets/cat.jpg", 10, 256, 256, Bilinear, 1, 180},
	{"assets/lena.jpg", "assets/monarch.jpg", 28, 256, 256, Bilinear, 1, 180},
	{"assets/baboon.jpg", "assets/cat.jpg", 20, 256, 256, Bilinear, 1, 180},
	{"assets/peppers.jpg", "assets/baboon.jpg", 20, 256, 256, Bilinear, 1, 180},
	{"assets/tulips.jpg", "assets/monarch.jpg", 18, 256, 256, Bilinear, 1, 180},
}

func TestRASH_Distance(t *testing.T) {
	for _, tt := range rashDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := NewRASH(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType), WithSigma(tt.sigma), WithRings(tt.rings))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img1, err := OpenImage(tt.firstImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.firstImage, err)
			}
			img2, err := OpenImage(tt.secondImage)
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

func TestNewRASH_Errors(t *testing.T) {
	tests := []struct {
		name string
		opts []RASHOption
		err  error
	}{
		{"zero width", []RASHOption{WithSize(0, 256)}, ErrInvalidSize},
		{"zero height", []RASHOption{WithSize(256, 0)}, ErrInvalidSize},
		{"negative sigma", []RASHOption{WithSigma(-1)}, ErrInvalidSigma},
		{"zero rings", []RASHOption{WithRings(0)}, ErrInvalidRings},
		{"negative rings", []RASHOption{WithRings(-1)}, ErrInvalidRings},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRASH(tt.opts...)
			if err != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}
