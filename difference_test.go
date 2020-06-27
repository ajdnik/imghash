package imghash_test

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
	"github.com/ajdnik/imghash/similarity"
)

var differenceCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType imgproc.ResizeType
}{
	{"assets/lena.jpg", hashtype.Binary{46, 14, 158, 218, 220, 200, 90, 28}, 8, 8, imgproc.Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{248, 213, 23, 22, 22, 28, 96, 22}, 8, 8, imgproc.Bilinear},
	{"assets/cat.jpg", hashtype.Binary{6, 2, 194, 64, 124, 60, 16, 16}, 8, 8, imgproc.Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{204, 204, 138, 138, 204, 77, 113, 101}, 8, 8, imgproc.Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{56, 246, 211, 211, 187, 187, 41, 225}, 8, 8, imgproc.Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{164, 51, 111, 109, 105, 31, 19, 35}, 8, 8, imgproc.Bilinear},
}

func TestDifference_Calculate(t *testing.T) {
	for _, tt := range differenceCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewDifferenceWithParams(tt.width, tt.height, tt.resizeType)
			img, _ := imgproc.Read(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

var differenceDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imgproc.ResizeType
}{
	{"assets/lena.jpg", "assets/cat.jpg", 24, 8, 8, imgproc.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 24, 8, 8, imgproc.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 32, 8, 8, imgproc.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 32, 8, 8, imgproc.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 35, 8, 8, imgproc.Bilinear},
}

func TestDifference_Distance(t *testing.T) {
	for _, tt := range differenceDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewDifferenceWithParams(tt.width, tt.height, tt.resizeType)
			img1, _ := imgproc.Read(tt.firstImage)
			img2, _ := imgproc.Read(tt.secondImage)
			h1 := hash.Calculate(img1)
			h2 := hash.Calculate(img2)
			dist := similarity.Hamming(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
