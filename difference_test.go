package imghash_test

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/similarity"
)

var differenceCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType Interpolation
}{
	{"assets/lena.jpg", hashtype.Binary{46, 14, 158, 218, 220, 200, 88, 28}, 8, 8, Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{248, 213, 23, 22, 22, 28, 72, 22}, 8, 8, Bilinear},
	{"assets/cat.jpg", hashtype.Binary{6, 2, 194, 64, 92, 60, 16, 16}, 8, 8, Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{204, 204, 138, 138, 204, 77, 113, 101}, 8, 8, Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{56, 242, 211, 211, 187, 187, 41, 225}, 8, 8, Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{150, 51, 111, 109, 105, 31, 19, 3}, 8, 8, Bilinear},
}

func TestDifference_Calculate(t *testing.T) {
	for _, tt := range differenceCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewDifference(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
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

func ExampleDifference_Calculate() {
	// Read image from file
	img, _ := OpenImage("assets/cat.jpg")
	// Create new Difference Hash using default parameters
	diff := NewDifference()
	// Calculate hash
	hash, _ := diff.Calculate(img)

	fmt.Println(hash)
	// Output: [6 2 194 64 92 60 16 16]
}

var differenceDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 22, 8, 8, Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 23, 8, 8, Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 31, 8, 8, Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 33, 8, 8, Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 37, 8, 8, Bilinear},
}

func TestDifference_Distance(t *testing.T) {
	for _, tt := range differenceDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewDifference(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
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
