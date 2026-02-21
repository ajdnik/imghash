package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/similarity"
)

var pHashCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType Interpolation
}{
	{"assets/lena.jpg", hashtype.Binary{152, 99, 42, 180, 174, 204, 69, 105}, 32, 32, BilinearExact},
	{"assets/baboon.jpg", hashtype.Binary{251, 4, 6, 190, 248, 133, 91, 241}, 32, 32, BilinearExact},
	{"assets/cat.jpg", hashtype.Binary{170, 195, 65, 29, 10, 2, 38, 84}, 32, 32, BilinearExact},
	{"assets/monarch.jpg", hashtype.Binary{150, 222, 38, 63, 25, 105, 128, 70}, 32, 32, BilinearExact},
	{"assets/peppers.jpg", hashtype.Binary{192, 245, 62, 8, 227, 136, 19, 155}, 32, 32, BilinearExact},
	{"assets/tulips.jpg", hashtype.Binary{34, 117, 194, 95, 55, 122, 48, 37}, 32, 32, BilinearExact},
}

func TestPHash_Calculate(t *testing.T) {
	for _, tt := range pHashCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewPHash(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
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

func ExamplePHash_Calculate() {
	// Read image from file
	img, _ := OpenImage("assets/cat.jpg")
	// Create new PHash using default parameters
	ph := NewPHash()
	// Calculate hash
	hash, _ := ph.Calculate(img)

	fmt.Println(hash)
	// Output: [170 195 65 29 10 2 38 84]
}

var pHashDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 31, 32, 32, BilinearExact},
	{"assets/lena.jpg", "assets/monarch.jpg", 34, 32, 32, BilinearExact},
	{"assets/baboon.jpg", "assets/cat.jpg", 35, 32, 32, BilinearExact},
	{"assets/peppers.jpg", "assets/baboon.jpg", 31, 32, 32, BilinearExact},
	{"assets/tulips.jpg", "assets/monarch.jpg", 29, 32, 32, BilinearExact},
}

func TestPHash_Distance(t *testing.T) {
	for _, tt := range pHashDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewPHash(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType))
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
