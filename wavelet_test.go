package imghash_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var wHashCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	level      int
	resizeType imghash.Interpolation
}{
	{"assets/lena.jpg", hashtype.Binary{125, 25, 189, 145, 208, 208, 241, 49}, 8, 8, 3, imghash.Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{128, 195, 252, 60, 61, 25, 71, 61}, 8, 8, 3, imghash.Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 31, 7, 3, 1, 0, 47}, 8, 8, 3, imghash.Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{1, 9, 19, 252, 191, 223, 230, 64}, 8, 8, 3, imghash.Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{247, 225, 134, 180, 62, 50, 34, 135}, 8, 8, 3, imghash.Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{13, 102, 92, 90, 250, 122, 60, 6}, 8, 8, 3, imghash.Bilinear},
}

func TestWHash_Calculate(t *testing.T) {
	for _, tt := range wHashCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := imghash.NewWHash(imghash.WithSize(tt.width, tt.height), imghash.WithLevel(tt.level), imghash.WithInterpolation(tt.resizeType))
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
			res := result.(hashtype.Binary)
			if !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleWHash_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	wh, err := imghash.NewWHash()
	if err != nil {
		panic(err)
	}
	hash, err := wh.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [255 255 31 7 3 1 0 47]
}

var wHashDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	level       int
	resizeType  imghash.Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 32, 8, 8, 3, imghash.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 34, 8, 8, 3, imghash.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 34, 8, 8, 3, imghash.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 30, 8, 8, 3, imghash.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 32, 8, 8, 3, imghash.Bilinear},
}

func TestWHash_Distance(t *testing.T) {
	for _, tt := range wHashDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := imghash.NewWHash(imghash.WithSize(tt.width, tt.height), imghash.WithLevel(tt.level), imghash.WithInterpolation(tt.resizeType))
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

func TestNewWHash_defaults(t *testing.T) {
	wh, err := imghash.NewWHash()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	_, err = wh.Calculate(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewWHash_invalidSize(t *testing.T) {
	_, err := imghash.NewWHash(imghash.WithSize(0, 8))
	if !errors.Is(err, imghash.ErrInvalidSize) {
		t.Errorf("got %v, want imghash.ErrInvalidSize", err)
	}
}

func TestNewWHash_invalidLevel(t *testing.T) {
	_, err := imghash.NewWHash(imghash.WithLevel(0))
	if !errors.Is(err, imghash.ErrInvalidLevel) {
		t.Errorf("got %v, want imghash.ErrInvalidLevel", err)
	}
}
