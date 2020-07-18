package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
	"github.com/ajdnik/imghash/similarity"
)

var averageCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType imgproc.ResizeType
}{
	{"assets/lena.jpg", hashtype.Binary{125, 57, 188, 177, 208, 208, 240, 112}, 8, 8, imgproc.Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{128, 192, 252, 60, 61, 25, 255, 61}, 8, 8, imgproc.Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 15, 7, 1, 0, 0, 0}, 8, 8, imgproc.Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{1, 11, 23, 252, 191, 255, 230, 64}, 8, 8, imgproc.Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{225, 224, 206, 244, 62, 54, 2, 7}, 8, 8, imgproc.Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{15, 102, 92, 90, 254, 126, 62, 6}, 8, 8, imgproc.Bilinear},
}

func TestAverage_Calculate(t *testing.T) {
	for _, tt := range averageCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewAverageWithParams(tt.width, tt.height, tt.resizeType)
			img, _ := imgproc.Read(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleAverage_Calculate() {
	// Read image from file
	img, _ := imgproc.Read("assets/cat.jpg")
	// Create new Average Hash using default parameters
	avg := NewAverage()
	// Calculate hash
	hash := avg.Calculate(img)

	fmt.Println(hash)
	// Output: [255 255 15 7 1 0 0 0]
}

var averageDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imgproc.ResizeType
}{
	{"assets/lena.jpg", "assets/cat.jpg", 30, 8, 8, imgproc.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 33, 8, 8, imgproc.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 44, 8, 8, imgproc.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 28, 8, 8, imgproc.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 27, 8, 8, imgproc.Bilinear},
}

func TestAverage_Distance(t *testing.T) {
	for _, tt := range averageDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewAverageWithParams(tt.width, tt.height, tt.resizeType)
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
