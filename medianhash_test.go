package imghash_test

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var medianCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType imghash.Interpolation
}{
	{"assets/lena.jpg", hashtype.Binary{125, 56, 188, 144, 208, 208, 240, 48}, 8, 8, imghash.Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{128, 128, 252, 60, 60, 25, 79, 63}, 8, 8, imghash.Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 31, 7, 1, 1, 1, 7}, 8, 8, imghash.Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{1, 3, 17, 252, 191, 255, 194, 64}, 8, 8, imghash.Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{241, 225, 206, 244, 62, 54, 2, 7}, 8, 8, imghash.Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{13, 38, 76, 90, 250, 62, 62, 6}, 8, 8, imghash.Bilinear},
}

func TestMedian_Calculate(t *testing.T) {
	for _, tt := range medianCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := imghash.NewMedian(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType))
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

func ExampleMedian_Calculate() {
	// Read image from file
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	// Create new Median Hash using default parameters
	med, err := imghash.NewMedian()
	if err != nil {
		panic(err)
	}
	// Calculate hash
	hash, err := med.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [255 255 31 7 1 1 1 7]
}

var medianDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imghash.Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 34, 8, 8, imghash.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 36, 8, 8, imghash.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 38, 8, 8, imghash.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 26, 8, 8, imghash.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 29, 8, 8, imghash.Bilinear},
}

func TestMedian_Distance(t *testing.T) {
	for _, tt := range medianDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := imghash.NewMedian(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType))
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
