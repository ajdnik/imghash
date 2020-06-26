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
	resizeType ResizeType
}{
	{"assets/lena.jpg", hashtype.Binary{125, 121, 185, 149, 213, 197, 112, 52}, 8, 8, Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{137, 197, 239, 54, 53, 3, 87, 165}, 8, 8, Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 143, 3, 33, 65, 32, 27}, 8, 8, Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{13, 17, 165, 252, 132, 209, 225, 66}, 8, 8, Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{241, 197, 207, 182, 126, 54, 34, 135}, 8, 8, Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{29, 50, 76, 91, 161, 54, 58, 70}, 8, 8, Bilinear},
}

func TestAverage_Calculate(t *testing.T) {
	for _, tt := range averageCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			hash := NewAverageWithParams(tt.width, tt.height, tt.resizeType)
			img, _ := ReadImageCV(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

var averageDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  ResizeType
}{
	{"assets/lena.jpg", "assets/cat.jpg", 27, 8, 8, Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 26, 8, 8, Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 31, 8, 8, Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 21, 8, 8, Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 30, 8, 8, Bilinear},
}

func TestAverage_Distance(t *testing.T) {
	for _, tt := range averageDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			t.Parallel()
			hash := NewAverageWithParams(tt.width, tt.height, tt.resizeType)
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
