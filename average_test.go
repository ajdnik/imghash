package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/similarity"
)

var averageCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType Interpolation
}{
	{"assets/lena.jpg", hashtype.Binary{125, 57, 189, 145, 208, 208, 240, 113}, 8, 8, Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{128, 193, 252, 60, 60, 25, 95, 63}, 8, 8, Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 15, 7, 1, 0, 0, 0}, 8, 8, Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{1, 11, 19, 252, 191, 255, 230, 192}, 8, 8, Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{241, 225, 206, 244, 62, 54, 2, 7}, 8, 8, Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{13, 102, 92, 90, 254, 126, 62, 6}, 8, 8, Bilinear},
}

func TestAverage_Calculate(t *testing.T) {
	for _, tt := range averageCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewAverage(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
			img, _ := OpenImage(tt.filename)
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

func ExampleAverage_Calculate() {
	// Read image from file
	img, _ := OpenImage("assets/cat.jpg")
	// Create new Average Hash using default parameters
	avg := NewAverage()
	// Calculate hash
	hash, _ := avg.Calculate(img)

	fmt.Println(hash)
	// Output: [255 255 15 7 1 0 0 0]
}

var averageDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 29, 8, 8, Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 36, 8, 8, Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 43, 8, 8, Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 25, 8, 8, Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 28, 8, 8, Bilinear},
}

func TestAverage_Distance(t *testing.T) {
	for _, tt := range averageDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewAverage(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
			img1, _ := OpenImage(tt.firstImage)
			img2, _ := OpenImage(tt.secondImage)
			h1, _ := hash.Calculate(img1)
			h2, _ := hash.Calculate(img2)
			dist, _ := similarity.Hamming(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
