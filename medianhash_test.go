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

var medianCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType imgproc.ResizeType
}{
	{"assets/lena.jpg", hashtype.Binary{125, 57, 188, 144, 208, 208, 240, 112}, 8, 8, imgproc.Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{128, 192, 252, 60, 61, 25, 95, 29}, 8, 8, imgproc.Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 31, 7, 1, 1, 1, 39}, 8, 8, imgproc.Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{1, 3, 17, 252, 191, 255, 198, 0}, 8, 8, imgproc.Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{225, 224, 206, 244, 62, 54, 2, 7}, 8, 8, imgproc.Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{13, 38, 76, 90, 122, 62, 62, 6}, 8, 8, imgproc.Bilinear},
}

func TestMedian_Calculate(t *testing.T) {
	for _, tt := range medianCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewMedianWithParams(tt.width, tt.height, tt.resizeType)
			img, _ := imgproc.Read(tt.filename)
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
	img, _ := imgproc.Read("assets/cat.jpg")
	// Create new Median Hash using default parameters
	med := NewMedian()
	// Calculate hash
	hash, _ := med.Calculate(img)

	fmt.Println(hash)
	// Output: [255 255 31 7 1 1 1 39]
}

var medianDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imgproc.ResizeType
}{
	{"assets/lena.jpg", "assets/cat.jpg", 33, 8, 8, imgproc.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 36, 8, 8, imgproc.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 38, 8, 8, imgproc.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 25, 8, 8, imgproc.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 28, 8, 8, imgproc.Bilinear},
}

func TestMedian_Distance(t *testing.T) {
	for _, tt := range medianDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewMedianWithParams(tt.width, tt.height, tt.resizeType)
			img1, _ := imgproc.Read(tt.firstImage)
			img2, _ := imgproc.Read(tt.secondImage)
			h1, _ := hash.Calculate(img1)
			h2, _ := hash.Calculate(img2)
			dist, _ := similarity.Hamming(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
