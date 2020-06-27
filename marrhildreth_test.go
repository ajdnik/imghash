package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
	"github.com/ajdnik/imghash/similarity"
)

var marrHildrethCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	width      uint
	height     uint
	resizeType imgproc.ResizeType
	scale      float64
	alpha      float64
	kernelSize int
	sigma      float64
}{
	{"assets/lena.jpg", hashtype.Binary{208, 182, 63, 31, 143, 199, 227, 225, 224, 124, 126, 39, 3, 201, 136, 124, 126, 63, 228, 182, 135, 192, 246, 119, 145, 137, 195, 123, 118, 57, 45, 239, 3, 227, 225, 240, 15, 137, 137, 58, 9, 35, 103, 177, 240, 228, 112, 228, 129, 241, 249, 254, 79, 196, 230, 16, 193, 239, 143, 6, 12, 108, 216, 5, 189, 52, 4, 215, 167, 0, 30, 127}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/baboon.jpg", hashtype.Binary{240, 120, 59, 19, 49, 56, 156, 14, 91, 5, 137, 192, 228, 177, 255, 211, 144, 27, 37, 143, 188, 28, 14, 71, 98, 109, 192, 118, 77, 196, 227, 241, 57, 89, 242, 27, 164, 237, 191, 31, 143, 199, 100, 150, 115, 100, 236, 196, 193, 240, 56, 127, 108, 182, 108, 252, 60, 29, 142, 6, 147, 192, 75, 19, 137, 193, 229, 240, 126, 114, 112, 201}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/cat.jpg", hashtype.Binary{51, 1, 228, 237, 182, 92, 200, 237, 102, 109, 182, 60, 12, 178, 88, 222, 73, 182, 108, 124, 27, 36, 183, 79, 226, 73, 182, 76, 134, 68, 237, 229, 163, 242, 73, 182, 36, 159, 11, 115, 118, 76, 152, 151, 180, 100, 248, 11, 97, 121, 38, 145, 169, 45, 45, 146, 91, 31, 13, 33, 60, 95, 38, 47, 131, 201, 32, 253, 160, 108, 188, 27}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/monarch.jpg", hashtype.Binary{192, 50, 25, 35, 242, 35, 96, 158, 63, 31, 143, 253, 57, 32, 216, 95, 176, 75, 19, 201, 136, 224, 246, 23, 0, 149, 182, 210, 50, 94, 221, 142, 5, 32, 150, 38, 7, 137, 198, 205, 162, 79, 193, 147, 176, 35, 246, 118, 194, 113, 249, 211, 208, 217, 236, 118, 54, 128, 246, 198, 144, 63, 200, 208, 112, 56, 30, 78, 0, 60, 14, 15}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/peppers.jpg", hashtype.Binary{147, 0, 201, 242, 242, 123, 98, 108, 70, 201, 181, 196, 159, 143, 135, 131, 129, 164, 45, 252, 52, 28, 8, 24, 28, 30, 54, 109, 249, 55, 131, 190, 11, 39, 145, 182, 100, 241, 59, 28, 77, 7, 39, 134, 219, 242, 118, 91, 220, 99, 183, 99, 145, 218, 208, 19, 185, 27, 0, 56, 60, 30, 217, 252, 104, 6, 228, 182, 135, 227, 128, 111}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/tulips.jpg", hashtype.Binary{228, 22, 254, 156, 143, 135, 99, 100, 201, 15, 248, 57, 58, 103, 155, 47, 240, 73, 224, 13, 201, 248, 32, 254, 216, 54, 201, 35, 141, 164, 152, 78, 200, 156, 110, 15, 115, 114, 38, 192, 114, 252, 188, 240, 201, 220, 13, 183, 14, 105, 244, 156, 14, 216, 31, 205, 154, 57, 102, 201, 36, 14, 6, 218, 108, 7, 224, 112, 57, 61, 249, 240}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
}

func TestMarrHildreth_Calculate(t *testing.T) {
	for _, tt := range marrHildrethCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewMarrHildrethWithParams(tt.scale, tt.alpha, tt.width, tt.height, tt.resizeType, tt.kernelSize, tt.sigma)
			img, _ := imgproc.Read(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

var marrHildrethDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imgproc.ResizeType
	scale       float64
	alpha       float64
	kernelSize  int
	sigma       float64
}{
	{"assets/lena.jpg", "assets/cat.jpg", 276, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/lena.jpg", "assets/monarch.jpg", 266, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/baboon.jpg", "assets/cat.jpg", 284, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/peppers.jpg", "assets/baboon.jpg", 285, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/tulips.jpg", "assets/monarch.jpg", 304, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
}

func TestMarrHildreth_Distance(t *testing.T) {
	for _, tt := range marrHildrethDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewMarrHildrethWithParams(tt.scale, tt.alpha, tt.width, tt.height, tt.resizeType, tt.kernelSize, tt.sigma)
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
