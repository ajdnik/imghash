package imghash_test

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	. "github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
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
			hash, err := NewDifference(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
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

func ExampleDifference_Calculate() {
	// Read image from file
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	// Create new Difference Hash using default parameters
	diff, err := NewDifference()
	if err != nil {
		panic(err)
	}
	// Calculate hash
	hash, err := diff.Calculate(img)
	if err != nil {
		panic(err)
	}

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
			hash, err := NewDifference(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
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
