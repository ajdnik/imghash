package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
	"github.com/ajdnik/imghash/similarity"
)

var blockMeanCalculateTests = []struct {
	filename    string
	hash        hashtype.Binary
	width       uint
	height      uint
	resizeType  imgproc.ResizeType
	blockWidth  uint
	blockHeight uint
	method      BlockMeanMethod
}{
	{"assets/lena.jpg", hashtype.Binary{243, 61, 243, 61, 194, 31, 226, 151, 226, 159, 122, 206, 18, 199, 130, 231, 130, 226, 130, 227, 130, 243, 2, 243, 2, 83, 2, 127, 2, 47, 130, 15}, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/baboon.jpg", hashtype.Binary{2, 96, 2, 64, 6, 96, 6, 230, 102, 238, 116, 126, 240, 111, 240, 15, 241, 7, 97, 7, 67, 3, 139, 3, 143, 177, 63, 190, 255, 63, 247, 15}, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 255, 255, 255, 255, 255, 255, 255, 225, 127, 0, 63, 0, 15, 0, 7, 2, 3, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 116, 4}, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/monarch.jpg", hashtype.Binary{7, 0, 199, 64, 199, 1, 147, 1, 27, 0, 144, 207, 144, 215, 248, 223, 215, 239, 247, 231, 254, 243, 190, 123, 121, 120, 40, 120, 4, 112, 0, 96}, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/peppers.jpg", hashtype.Binary{57, 191, 7, 188, 35, 216, 33, 230, 189, 230, 188, 232, 56, 225, 56, 255, 248, 31, 216, 159, 216, 159, 24, 27, 28, 28, 29, 8, 15, 96, 127, 128}, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/tulips.jpg", hashtype.Binary{231, 8, 99, 2, 56, 6, 56, 62, 248, 56, 242, 51, 194, 3, 206, 57, 158, 127, 126, 125, 108, 127, 244, 15, 248, 7, 120, 7, 126, 2, 56, 0}, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/lena.jpg", hashtype.Binary{135, 255, 243, 135, 195, 255, 249, 199, 193, 255, 252, 227, 224, 127, 254, 97, 128, 127, 63, 48, 192, 255, 25, 24, 240, 255, 12, 13, 249, 255, 199, 134, 252, 255, 99, 227, 254, 255, 184, 241, 167, 63, 222, 240, 192, 15, 111, 24, 240, 129, 55, 12, 252, 224, 27, 0, 127, 240, 13, 192, 28, 252, 14, 96, 14, 126, 7, 240, 7, 63, 3, 240, 131, 159, 1, 248, 225, 239, 0, 252, 240, 119, 16, 60, 252, 59, 0, 30, 254, 29, 0, 14, 255, 12, 0, 143, 59, 7, 192, 223, 143, 1, 224, 255, 199, 0, 240, 255, 99, 0, 248, 231, 48, 0, 254, 35, 24, 0, 255, 3, 0}, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/baboon.jpg", hashtype.Binary{12, 0, 0, 60, 7, 0, 0, 156, 3, 0, 0, 206, 3, 0, 128, 231, 3, 0, 224, 243, 1, 64, 240, 251, 32, 112, 248, 57, 56, 124, 252, 28, 62, 126, 126, 152, 159, 255, 31, 196, 175, 255, 15, 227, 255, 223, 7, 240, 255, 231, 1, 248, 254, 99, 0, 124, 255, 1, 2, 190, 255, 0, 3, 255, 63, 128, 129, 223, 31, 192, 128, 207, 15, 96, 192, 243, 3, 112, 128, 249, 0, 56, 195, 126, 0, 156, 129, 63, 0, 206, 192, 3, 12, 247, 192, 3, 207, 255, 1, 224, 231, 255, 3, 255, 243, 255, 255, 255, 249, 255, 255, 127, 248, 252, 255, 15, 124, 252, 255, 1, 0}, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/cat.jpg", hashtype.Binary{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 127, 252, 255, 255, 9, 252, 255, 127, 0, 192, 255, 15, 0, 224, 255, 7, 0, 240, 255, 0, 0, 248, 31, 0, 0, 252, 3, 2, 0, 254, 0, 0, 0, 63, 0, 4, 128, 15, 0, 24, 193, 3, 0, 0, 160, 1, 0, 0, 144, 0, 0, 0, 88, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 88, 0, 64, 0, 224, 246, 96, 0, 0}, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/monarch.jpg", hashtype.Binary{31, 0, 0, 128, 135, 25, 0, 208, 7, 60, 0, 236, 3, 30, 0, 246, 1, 31, 0, 250, 128, 31, 0, 60, 12, 15, 0, 158, 7, 0, 0, 239, 3, 0, 128, 183, 193, 0, 176, 192, 248, 63, 28, 96, 254, 127, 14, 56, 254, 59, 7, 28, 255, 223, 131, 223, 255, 239, 243, 239, 255, 251, 63, 243, 127, 252, 143, 249, 31, 254, 199, 253, 7, 255, 231, 254, 131, 207, 247, 255, 240, 231, 251, 127, 254, 241, 189, 31, 255, 242, 174, 193, 127, 227, 23, 224, 191, 179, 7, 240, 31, 24, 1, 248, 15, 14, 0, 248, 7, 3, 0, 240, 131, 0, 0, 240, 0, 0, 0, 120, 0}, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/peppers.jpg", hashtype.Binary{195, 135, 255, 239, 11, 0, 255, 255, 15, 0, 254, 249, 131, 0, 126, 124, 64, 0, 62, 31, 48, 0, 254, 15, 24, 112, 248, 39, 14, 120, 252, 243, 111, 28, 254, 248, 55, 48, 63, 252, 27, 24, 31, 248, 65, 132, 15, 252, 48, 192, 7, 126, 56, 225, 131, 63, 252, 247, 193, 223, 252, 3, 224, 255, 223, 1, 240, 251, 255, 0, 248, 252, 127, 16, 124, 254, 63, 8, 62, 255, 31, 4, 31, 188, 15, 128, 7, 156, 7, 192, 3, 196, 3, 240, 1, 240, 129, 248, 0, 112, 192, 124, 0, 56, 96, 30, 0, 152, 240, 15, 0, 224, 251, 63, 0, 192, 253, 255, 0, 0, 1}, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/tulips.jpg", hashtype.Binary{63, 126, 192, 128, 15, 31, 0, 192, 131, 143, 3, 224, 241, 225, 19, 0, 254, 224, 9, 0, 127, 240, 63, 128, 63, 240, 31, 192, 31, 128, 31, 224, 255, 192, 15, 224, 255, 224, 7, 225, 255, 225, 193, 1, 255, 112, 96, 128, 127, 0, 240, 195, 31, 4, 240, 227, 143, 31, 248, 243, 199, 127, 252, 225, 223, 63, 255, 240, 252, 159, 127, 239, 254, 143, 159, 243, 255, 131, 199, 249, 255, 129, 255, 255, 63, 192, 254, 255, 1, 0, 255, 127, 0, 192, 255, 63, 0, 224, 255, 31, 0, 248, 239, 15, 224, 255, 231, 1, 224, 255, 99, 0, 224, 255, 0, 0, 0, 63, 0, 0, 0}, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
}

func TestBlockMean_Calculate(t *testing.T) {
	for _, tt := range blockMeanCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewBlockMeanWithParams(tt.width, tt.height, tt.resizeType, tt.blockWidth, tt.blockHeight, tt.method)
			img, _ := imgproc.Read(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleBlockMean_Calculate() {
	// Read image from file
	img, _ := imgproc.Read("assets/cat.jpg")
	// Create new Block Mean Hash using default parameters
	block := NewBlockMean()
	// Calculate hash
	hash := block.Calculate(img)

	fmt.Println(hash)
	// Output: [255 255 255 255 255 255 255 255 255 225 127 0 63 0 15 0 7 2 3 0 1 0 1 0 0 0 0 0 0 0 116 4]
}

var blockMeanDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imgproc.ResizeType
	blockWidth  uint
	blockHeight uint
	method      BlockMeanMethod
}{
	{"assets/lena.jpg", "assets/cat.jpg", 119, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/lena.jpg", "assets/monarch.jpg", 119, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/baboon.jpg", "assets/cat.jpg", 153, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/peppers.jpg", "assets/baboon.jpg", 119, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/tulips.jpg", "assets/monarch.jpg", 121, 256, 256, imgproc.BilinearExact, 16, 16, Direct},
	{"assets/lena.jpg", "assets/cat.jpg", 455, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/lena.jpg", "assets/monarch.jpg", 458, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/baboon.jpg", "assets/cat.jpg", 575, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/peppers.jpg", "assets/baboon.jpg", 426, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
	{"assets/tulips.jpg", "assets/monarch.jpg", 491, 256, 256, imgproc.BilinearExact, 16, 16, Overlap},
}

func TestBlockMean_Distance(t *testing.T) {
	for _, tt := range blockMeanDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewBlockMeanWithParams(tt.width, tt.height, tt.resizeType, tt.blockWidth, tt.blockHeight, tt.method)
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
