package imghash_test

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/similarity"
)

var medianCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType ResizeType
}{
	{"assets/lena.jpg", hashtype.Binary{101, 121, 185, 145, 209, 197, 112, 52}, 8, 8, Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{9, 197, 239, 54, 21, 3, 87, 165}, 8, 8, Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 143, 3, 33, 65, 32, 27}, 8, 8, Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{141, 157, 165, 252, 132, 209, 225, 66}, 8, 8, Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{113, 197, 206, 182, 62, 22, 2, 135}, 8, 8, Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{29, 50, 76, 91, 229, 54, 58, 70}, 8, 8, Bilinear},
}

func TestMedian_Calculate(t *testing.T) {
	for _, tt := range medianCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			hash := NewMedianWithParams(tt.width, tt.height, tt.resizeType)
			file, err := os.Open(tt.filename)
			if err != nil {
				t.Errorf("failed to open image: %s", err)
			}
			defer file.Close()
			img, _, err := image.Decode(bufio.NewReader(file))
			if err != nil {
				t.Errorf("failed to decode image: %s", err)
			}
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

var medianDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  ResizeType
}{
	{"assets/lena.jpg", "assets/cat.jpg", 27, 8, 8, Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 30, 8, 8, Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 33, 8, 8, Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 20, 8, 8, Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 34, 8, 8, Bilinear},
}

func TestMedian_Distance(t *testing.T) {
	for _, tt := range medianDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			t.Parallel()
			hash := NewMedianWithParams(tt.width, tt.height, tt.resizeType)
			file1, err := os.Open(tt.firstImage)
			if err != nil {
				t.Errorf("failed to open image: %s", err)
			}
			defer file1.Close()
			img1, _, err := image.Decode(bufio.NewReader(file1))
			if err != nil {
				t.Errorf("failed to decode image: %s", err)
			}
			file2, err := os.Open(tt.secondImage)
			if err != nil {
				t.Errorf("failed to open image: %s", err)
			}
			defer file2.Close()
			img2, _, err := image.Decode(bufio.NewReader(file2))
			if err != nil {
				t.Errorf("failed to decode image: %s", err)
			}
			h1 := hash.Calculate(img1)
			h2 := hash.Calculate(img2)
			dist := similarity.Hamming(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
