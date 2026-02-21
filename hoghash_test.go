package imghash_test

import (
	"fmt"
	"testing"

	. "github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var hogHashCalculateTests = []struct {
	filename   string
	hash       hashtype.UInt8
	width      uint
	height     uint
	resizeType Interpolation
	cellSize   uint
	numBins    uint
}{
	{"assets/lena.jpg", hashtype.UInt8{255, 11, 15, 17, 32, 10, 6, 0, 159, 255, 145, 10, 79, 90, 62, 20, 45, 136, 96, 13, 5, 56, 50, 35, 255, 92, 106, 87, 255, 250, 33, 17, 9, 61, 103, 19, 160, 21, 4, 0, 4, 0, 0, 30, 255, 86, 57, 255, 45, 125, 29, 42, 41, 18, 135, 163, 255, 135, 143, 56, 38, 75, 101, 101, 255, 98, 4, 1, 2, 2, 3, 3, 249, 47, 43, 3, 3, 4, 2, 21, 255, 114, 255, 97, 77, 62, 0, 51, 130, 117, 186, 255, 10, 31, 18, 14, 6, 107, 92, 213, 255, 158, 59, 99, 0, 5, 25, 62, 255, 55, 11, 17, 3, 23, 7, 11, 121, 153, 255, 50, 12, 5, 3, 9, 102, 111, 115, 4, 13, 2, 5, 60, 231, 245, 254, 119, 255, 193, 17, 41, 30, 0, 33, 76}, 32, 32, Bilinear, 8, 9},
	{"assets/baboon.jpg", hashtype.UInt8{255, 116, 96, 53, 101, 61, 75, 156, 81, 80, 79, 255, 207, 249, 162, 204, 67, 51, 19, 45, 178, 130, 255, 26, 87, 150, 39, 255, 63, 94, 0, 57, 49, 49, 64, 174, 85, 49, 17, 101, 31, 94, 151, 255, 64, 255, 108, 66, 14, 6, 2, 0, 39, 167, 255, 47, 5, 0, 4, 0, 62, 72, 199, 99, 251, 255, 150, 30, 28, 16, 115, 76, 171, 43, 34, 69, 63, 55, 84, 76, 255, 76, 37, 9, 23, 22, 47, 32, 149, 254, 255, 209, 58, 54, 58, 16, 32, 90, 60, 68, 110, 46, 255, 76, 168, 142, 143, 74, 254, 103, 46, 34, 100, 49, 193, 163, 223, 52, 18, 5, 40, 81, 193, 255, 146, 71, 48, 64, 86, 255, 98, 32, 16, 21, 0, 67, 73, 255, 127, 75, 31, 38, 53, 75}, 32, 32, Bilinear, 8, 9},
	{"assets/cat.jpg", hashtype.UInt8{42, 41, 55, 40, 188, 216, 236, 255, 65, 255, 47, 80, 5, 111, 116, 123, 105, 66, 98, 47, 254, 131, 82, 86, 100, 45, 47, 214, 151, 131, 255, 205, 154, 196, 227, 67, 84, 216, 255, 221, 186, 170, 73, 89, 57, 51, 60, 161, 255, 71, 30, 63, 0, 71, 43, 46, 150, 193, 255, 29, 122, 53, 56, 58, 48, 125, 112, 255, 115, 77, 100, 133, 150, 255, 169, 107, 84, 60, 87, 50, 107, 89, 199, 255, 95, 77, 29, 36, 43, 57, 2, 9, 121, 116, 226, 255, 201, 68, 31, 33, 133, 172, 255, 87, 34, 131, 62, 12, 59, 60, 68, 92, 198, 254, 212, 30, 30, 25, 18, 60, 100, 254, 78, 57, 29, 0, 94, 35, 191, 143, 255, 128, 36, 78, 63, 44, 71, 99, 255, 118, 63, 3, 69, 35}, 32, 32, Bilinear, 8, 9},
	{"assets/monarch.jpg", hashtype.UInt8{255, 99, 111, 81, 41, 39, 61, 48, 245, 142, 88, 91, 109, 81, 255, 23, 116, 63, 244, 116, 89, 78, 30, 75, 143, 255, 123, 123, 28, 2, 7, 0, 7, 28, 195, 255, 73, 157, 180, 255, 32, 78, 226, 51, 85, 177, 194, 69, 123, 111, 0, 59, 195, 255, 119, 50, 96, 67, 254, 254, 114, 147, 46, 255, 193, 156, 41, 76, 93, 185, 119, 57, 84, 87, 90, 136, 72, 193, 123, 255, 198, 173, 131, 94, 38, 21, 46, 212, 255, 191, 120, 255, 208, 129, 102, 104, 75, 89, 162, 255, 236, 110, 166, 72, 43, 1, 29, 110, 182, 199, 162, 116, 141, 103, 176, 255, 118, 160, 255, 164, 154, 93, 92, 153, 210, 149, 132, 108, 100, 163, 213, 209, 139, 254, 147, 138, 198, 65, 7, 29, 3, 29, 255, 185}, 32, 32, Bilinear, 8, 9},
	{"assets/peppers.jpg", hashtype.UInt8{73, 49, 255, 79, 90, 80, 92, 80, 40, 157, 224, 94, 115, 255, 36, 67, 93, 80, 39, 35, 83, 65, 202, 136, 255, 242, 90, 255, 192, 116, 176, 92, 48, 47, 53, 111, 255, 144, 55, 81, 141, 72, 38, 88, 52, 255, 115, 81, 26, 75, 47, 65, 61, 175, 144, 191, 126, 114, 104, 115, 198, 106, 255, 193, 118, 82, 81, 106, 243, 139, 157, 254, 255, 35, 2, 14, 18, 12, 21, 15, 22, 254, 198, 23, 38, 17, 67, 28, 124, 139, 189, 151, 135, 174, 28, 48, 157, 255, 179, 255, 70, 8, 37, 53, 33, 13, 12, 123, 169, 77, 62, 16, 255, 154, 82, 95, 201, 188, 80, 5, 0, 91, 255, 145, 106, 117, 59, 65, 166, 126, 140, 143, 254, 67, 34, 246, 93, 56, 30, 188, 255, 145, 24, 82}, 32, 32, Bilinear, 8, 9},
	{"assets/tulips.jpg", hashtype.UInt8{101, 255, 244, 176, 52, 120, 46, 76, 148, 212, 134, 255, 100, 37, 53, 53, 94, 111, 99, 97, 70, 72, 89, 249, 105, 255, 142, 205, 81, 53, 100, 76, 238, 226, 254, 67, 111, 127, 201, 94, 36, 95, 255, 155, 105, 181, 105, 12, 156, 254, 77, 104, 111, 163, 120, 77, 78, 25, 43, 60, 104, 255, 211, 155, 126, 224, 79, 243, 146, 255, 99, 96, 183, 255, 147, 79, 93, 82, 146, 231, 118, 205, 255, 116, 76, 107, 105, 236, 165, 162, 255, 210, 148, 196, 59, 222, 131, 101, 144, 166, 97, 117, 255, 135, 105, 66, 45, 87, 117, 103, 43, 200, 255, 108, 215, 96, 18, 255, 167, 117, 123, 38, 40, 164, 127, 128, 255, 54, 159, 148, 75, 58, 43, 220, 95, 175, 21, 13, 31, 14, 14, 11, 136, 255}, 32, 32, Bilinear, 8, 9},
}

func TestHOGHash_Calculate(t *testing.T) {
	for _, tt := range hogHashCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := NewHOGHash(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType), WithCellSize(tt.cellSize), WithNumBins(tt.numBins))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img, err := OpenImage(tt.filename)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.filename, err)
			}
			result, err := hash.Calculate(img)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			res := result.(hashtype.UInt8)
			if !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleHOGHash_Calculate() {
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	hog, err := NewHOGHash(WithSize(32, 32))
	if err != nil {
		panic(err)
	}
	hash, err := hog.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [42 41 55 40 188 216 236 255 65 255 47 80 5 111 116 123 105 66 98 47 254 131 82 86 100 45 47 214 151 131 255 205 154 196 227 67 84 216 255 221 186 170 73 89 57 51 60 161 255 71 30 63 0 71 43 46 150 193 255 29 122 53 56 58 48 125 112 255 115 77 100 133 150 255 169 107 84 60 87 50 107 89 199 255 95 77 29 36 43 57 2 9 121 116 226 255 201 68 31 33 133 172 255 87 34 131 62 12 59 60 68 92 198 254 212 30 30 25 18 60 100 254 78 57 29 0 94 35 191 143 255 128 36 78 63 44 71 99 255 118 63 3 69 35]
}

var hogHashDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  Interpolation
	cellSize    uint
	numBins     uint
}{
	{"assets/lena.jpg", "assets/cat.jpg", 1530.9330488300263, 32, 32, Bilinear, 8, 9},
	{"assets/lena.jpg", "assets/monarch.jpg", 1315.5991030705366, 32, 32, Bilinear, 8, 9},
	{"assets/baboon.jpg", "assets/cat.jpg", 1403.76386903211, 32, 32, Bilinear, 8, 9},
	{"assets/peppers.jpg", "assets/baboon.jpg", 1183.0942481476275, 32, 32, Bilinear, 8, 9},
	{"assets/tulips.jpg", "assets/monarch.jpg", 1056.4710123803682, 32, 32, Bilinear, 8, 9},
}

func TestHOGHash_Distance(t *testing.T) {
	for _, tt := range hogHashDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := NewHOGHash(WithSize(tt.width, tt.height), WithInterpolation(tt.resizeType), WithCellSize(tt.cellSize), WithNumBins(tt.numBins))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img1, err := OpenImage(tt.firstImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.firstImage, err)
			}
			img2, err := OpenImage(tt.secondImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.secondImage, err)
			}
			h1, err := hash.Calculate(img1)
			if err != nil {
				t.Fatalf("failed to calculate hash for %s: %v", tt.firstImage, err)
			}
			h2, err := hash.Calculate(img2)
			if err != nil {
				t.Fatalf("failed to calculate hash for %s: %v", tt.secondImage, err)
			}
			dist := similarity.L2(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}

func TestNewHOGHash_defaults(t *testing.T) {
	hog, err := NewHOGHash()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := hog.Calculate(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 256/8 = 32 cells per axis, 32*32*9 = 9216
	if h.Len() != 9216 {
		t.Errorf("expected hash length 9216, got %d", h.Len())
	}
}

func TestNewHOGHash_invalidSize(t *testing.T) {
	_, err := NewHOGHash(WithSize(0, 0))
	if err != ErrInvalidSize {
		t.Errorf("expected ErrInvalidSize, got %v", err)
	}
}

func TestNewHOGHash_invalidCellSize(t *testing.T) {
	_, err := NewHOGHash(WithCellSize(0))
	if err != ErrInvalidCellSize {
		t.Errorf("expected ErrInvalidCellSize, got %v", err)
	}
}

func TestNewHOGHash_invalidNumBins(t *testing.T) {
	_, err := NewHOGHash(WithNumBins(0))
	if err != ErrInvalidNumBins {
		t.Errorf("expected ErrInvalidNumBins, got %v", err)
	}
}

func TestNewHOGHash_customCellSize(t *testing.T) {
	hog, err := NewHOGHash(WithSize(32, 32), WithCellSize(16))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := hog.Calculate(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 32/16 = 2 cells per axis, 2*2*9 = 36
	if h.Len() != 36 {
		t.Errorf("expected hash length 36, got %d", h.Len())
	}
}

func TestNewHOGHash_customNumBins(t *testing.T) {
	hog, err := NewHOGHash(WithSize(32, 32), WithNumBins(18))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := hog.Calculate(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 32/8 = 4 cells per axis, 4*4*18 = 288
	if h.Len() != 288 {
		t.Errorf("expected hash length 288, got %d", h.Len())
	}
}
