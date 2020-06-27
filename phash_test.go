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
	resizeType ResizeType
}{
	{"assets/lena.jpg", hashtype.Binary{152, 99, 43, 180, 174, 196, 101, 105}, 32, 32, BilinearExact},
	{"assets/baboon.jpg", hashtype.Binary{251, 4, 6, 190, 240, 5, 27, 249}, 32, 32, BilinearExact},
	{"assets/cat.jpg", hashtype.Binary{170, 211, 65, 57, 10, 130, 34, 68}, 32, 32, BilinearExact},
	{"assets/monarch.jpg", hashtype.Binary{150, 222, 38, 63, 16, 104, 136, 78}, 32, 32, BilinearExact},
	{"assets/peppers.jpg", hashtype.Binary{196, 253, 62, 8, 227, 136, 3, 155}, 32, 32, BilinearExact},
	{"assets/tulips.jpg", hashtype.Binary{162, 245, 194, 93, 55, 122, 52, 37}, 32, 32, BilinearExact},
}

func TestPHash_Calculate(t *testing.T) {
	for _, tt := range pHashCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewPHashWithParams(tt.width, tt.height, tt.resizeType)
			img, _ := ReadImageCV(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

var pHashDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  ResizeType
}{
	{"assets/lena.jpg", "assets/cat.jpg", 28, 32, 32, BilinearExact},
	{"assets/lena.jpg", "assets/monarch.jpg", 36, 32, 32, BilinearExact},
	{"assets/baboon.jpg", "assets/cat.jpg", 37, 32, 32, BilinearExact},
	{"assets/peppers.jpg", "assets/baboon.jpg", 32, 32, 32, BilinearExact},
	{"assets/tulips.jpg", "assets/monarch.jpg", 30, 32, 32, BilinearExact},
}

func TestPHash_Distance(t *testing.T) {
	for _, tt := range pHashDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewPHashWithParams(tt.width, tt.height, tt.resizeType)
			img1, _ := ReadImageCV(tt.firstImage)
			img2, _ := ReadImageCV(tt.secondImage)
			h1 := hash.Calculate(img1)
			h2 := hash.Calculate(img2)
			dist := similarity.Hamming(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
