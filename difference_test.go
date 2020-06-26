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

var differenceCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType ResizeType
}{
	{"assets/lena.jpg", hashtype.Binary{110, 78, 190, 218, 157, 200, 90, 62}, 8, 8, Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{201, 101, 21, 23, 150, 20, 42, 148}, 8, 8, Bilinear},
	{"assets/cat.jpg", hashtype.Binary{172, 168, 226, 74, 124, 172, 43, 52}, 8, 8, Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{234, 229, 240, 118, 213, 157, 54, 99}, 8, 8, Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{60, 188, 187, 234, 162, 171, 171, 192}, 8, 8, Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{22, 43, 107, 101, 105, 87, 75, 67}, 8, 8, Bilinear},
}

func TestDifference_Calculate(t *testing.T) {
	for _, tt := range differenceCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			hash := NewDifferenceWithParams(tt.width, tt.height, tt.resizeType)
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

var differenceDistanceTests = []struct {
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
	{"assets/peppers.jpg", "assets/baboon.jpg", 38, 8, 8, Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 35, 8, 8, Bilinear},
}

func TestDifference_Distance(t *testing.T) {
	for _, tt := range differenceDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			t.Parallel()
			hash := NewDifferenceWithParams(tt.width, tt.height, tt.resizeType)
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
