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
	{"assets/lena.jpg", hashtype.Binary{135, 65, 224, 252, 126, 62, 29, 14, 31, 31, 143, 199, 227, 241, 216, 251, 54, 63, 233, 122, 140, 184, 164, 250, 44, 225, 28, 219, 49, 229, 156, 162, 171, 63, 234, 113, 135, 136, 137, 14, 71, 229, 197, 135, 204, 31, 142, 234, 88, 247, 198, 29, 184, 6, 49, 32, 183, 14, 48, 238, 1, 185, 99, 79, 61, 94, 155, 61, 41, 247, 195, 87}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/baboon.jpg", hashtype.Binary{148, 140, 150, 203, 58, 216, 30, 186, 163, 169, 170, 181, 75, 64, 248, 46, 149, 48, 90, 117, 7, 156, 168, 218, 245, 144, 250, 89, 44, 248, 223, 12, 197, 2, 7, 36, 73, 28, 128, 240, 120, 224, 61, 103, 52, 73, 36, 128, 238, 19, 189, 169, 11, 29, 145, 57, 211, 130, 229, 101, 105, 52, 188, 183, 92, 87, 73, 85, 222, 21, 238, 169}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/cat.jpg", hashtype.Binary{92, 190, 42, 111, 87, 107, 101, 164, 184, 24, 75, 41, 185, 54, 178, 162, 26, 236, 155, 150, 108, 98, 233, 112, 56, 235, 124, 177, 139, 159, 148, 66, 89, 38, 229, 47, 195, 44, 158, 180, 85, 115, 79, 165, 92, 131, 225, 252, 54, 148, 218, 61, 99, 92, 82, 141, 141, 96, 112, 186, 185, 208, 174, 112, 252, 150, 153, 164, 173, 206, 43, 130}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/monarch.jpg", hashtype.Binary{54, 78, 167, 106, 105, 209, 208, 227, 66, 224, 174, 34, 54, 26, 45, 80, 234, 233, 29, 12, 163, 110, 170, 164, 86, 29, 88, 44, 232, 254, 67, 145, 114, 27, 97, 63, 215, 214, 85, 84, 156, 105, 208, 86, 177, 45, 15, 220, 42, 135, 109, 76, 126, 62, 27, 190, 63, 11, 194, 251, 163, 146, 63, 99, 234, 135, 25, 151, 203, 15, 30, 53}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/peppers.jpg", hashtype.Binary{13, 206, 150, 173, 65, 222, 89, 210, 59, 182, 87, 231, 20, 124, 56, 245, 130, 146, 150, 33, 251, 129, 208, 30, 226, 241, 77, 155, 10, 64, 115, 102, 37, 128, 79, 76, 137, 121, 167, 38, 53, 8, 231, 142, 53, 107, 169, 160, 106, 150, 89, 111, 103, 244, 29, 234, 199, 236, 119, 197, 227, 224, 78, 31, 174, 85, 64, 246, 216, 28, 53, 178}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/tulips.jpg", hashtype.Binary{31, 112, 92, 183, 132, 127, 115, 197, 38, 240, 6, 156, 145, 56, 126, 28, 46, 170, 23, 124, 79, 35, 222, 25, 119, 128, 180, 49, 209, 90, 52, 69, 127, 227, 209, 80, 152, 169, 222, 196, 239, 34, 224, 196, 116, 103, 30, 71, 251, 15, 130, 231, 49, 7, 232, 198, 73, 246, 49, 196, 84, 86, 133, 172, 136, 248, 159, 169, 82, 196, 131, 104}, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
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

func ExampleMarrHildreth_calculate() {
	// Read image from file
	img, _ := imgproc.Read("assets/cat.jpg")
	// Create new Marr-Hildreth Hash using default parameters
	marr := NewMarrHildreth()
	// Calculate hash
	hash := marr.Calculate(img)

	fmt.Println(hash)
	// Output: [92 190 42 111 87 107 101 164 184 24 75 41 185 54 178 162 26 236 155 150 108 98 233 112 56 235 124 177 139 159 148 66 89 38 229 47 195 44 158 180 85 115 79 165 92 131 225 252 54 148 218 61 99 92 82 141 141 96 112 186 185 208 174 112 252 150 153 164 173 206 43 130]
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
	{"assets/lena.jpg", "assets/cat.jpg", 274, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/lena.jpg", "assets/monarch.jpg", 311, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/baboon.jpg", "assets/cat.jpg", 288, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/peppers.jpg", "assets/baboon.jpg", 257, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
	{"assets/tulips.jpg", "assets/monarch.jpg", 311, 512, 512, imgproc.Bicubic, 1, 2, 7, 0},
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
