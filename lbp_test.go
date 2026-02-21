package imghash_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var lbpCalculateTests = []struct {
	filename   string
	hash       hashtype.UInt8
	width      uint
	height     uint
	resizeType imghash.Interpolation
	gridX      uint
	gridY      uint
}{
	{"assets/lena.jpg", hashtype.UInt8{82, 18, 10, 16, 40, 5, 19, 51, 15, 2, 1, 1, 21, 6, 14, 66, 14, 2, 1, 1, 4, 0, 2, 3, 16, 1, 1, 1, 33, 2, 39, 30, 10, 2, 2, 1, 3, 0, 1, 4, 1, 0, 0, 0, 2, 0, 1, 6, 16, 1, 1, 2, 7, 0, 4, 3, 46, 2, 2, 2, 131, 3, 35, 28, 39, 4, 3, 8, 10, 1, 3, 6, 6, 0, 0, 0, 4, 0, 1, 2, 7, 0, 0, 0, 0, 0, 0, 0, 11, 0, 1, 1, 5, 0, 3, 3, 18, 4, 2, 4, 4, 0, 2, 3, 2, 0, 0, 0, 3, 0, 2, 3, 47, 3, 5, 3, 5, 0, 3, 6, 174, 4, 9, 6, 113, 5, 33, 33, 19, 27, 1, 66, 7, 13, 2, 226, 3, 3, 0, 2, 3, 5, 1, 59, 3, 2, 0, 3, 0, 0, 0, 4, 3, 1, 0, 3, 2, 3, 2, 29, 1, 2, 0, 3, 0, 0, 0, 12, 0, 0, 0, 0, 0, 0, 0, 7, 2, 2, 0, 2, 0, 0, 0, 7, 2, 3, 0, 7, 6, 5, 4, 52, 24, 48, 3, 203, 3, 7, 3, 166, 3, 4, 0, 8, 2, 3, 1, 36, 8, 5, 0, 5, 0, 1, 0, 7, 5, 4, 0, 5, 3, 4, 3, 29, 22, 54, 1, 59, 2, 4, 2, 51, 2, 4, 0, 5, 2, 2, 1, 30, 66, 42, 8, 39, 4, 4, 5, 45, 59, 37, 6, 56, 33, 35, 30, 255}, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/baboon.jpg", hashtype.UInt8{197, 69, 24, 31, 40, 8, 28, 46, 25, 8, 3, 4, 30, 12, 31, 59, 71, 12, 10, 6, 9, 1, 13, 5, 37, 6, 3, 4, 44, 7, 65, 41, 26, 9, 4, 4, 5, 0, 5, 6, 3, 1, 0, 0, 4, 2, 6, 12, 36, 6, 5, 3, 9, 0, 7, 5, 41, 4, 4, 4, 75, 4, 42, 34, 42, 8, 7, 9, 6, 1, 4, 6, 6, 1, 1, 1, 5, 1, 4, 4, 10, 1, 0, 0, 0, 0, 1, 0, 10, 0, 1, 0, 5, 0, 4, 3, 30, 12, 3, 6, 3, 0, 4, 5, 5, 2, 0, 1, 4, 1, 6, 9, 45, 5, 6, 4, 4, 0, 5, 5, 91, 6, 10, 7, 56, 4, 45, 33, 25, 31, 3, 38, 5, 9, 4, 79, 5, 5, 0, 3, 4, 6, 5, 46, 8, 5, 1, 4, 0, 0, 2, 5, 5, 4, 0, 4, 6, 5, 11, 35, 2, 3, 0, 4, 0, 1, 0, 8, 1, 0, 0, 0, 0, 1, 1, 9, 4, 4, 0, 3, 1, 0, 0, 6, 4, 4, 0, 6, 9, 7, 8, 45, 27, 41, 4, 78, 4, 4, 4, 50, 4, 7, 0, 9, 4, 6, 4, 40, 13, 5, 2, 4, 0, 0, 1, 3, 7, 5, 1, 6, 4, 4, 8, 33, 33, 60, 4, 40, 3, 4, 6, 39, 5, 10, 1, 8, 5, 10, 10, 74, 62, 42, 11, 32, 4, 3, 9, 32, 49, 33, 8, 46, 41, 31, 73, 255}, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/cat.jpg", hashtype.UInt8{140, 97, 12, 39, 15, 6, 22, 68, 13, 10, 1, 2, 24, 14, 62, 193, 107, 18, 7, 6, 6, 0, 9, 7, 42, 5, 2, 2, 73, 7, 186, 105, 14, 8, 1, 3, 1, 0, 1, 4, 1, 0, 0, 0, 1, 0, 3, 16, 53, 7, 3, 3, 3, 0, 4, 3, 41, 2, 2, 1, 85, 2, 85, 41, 16, 6, 1, 2, 2, 0, 1, 3, 1, 0, 0, 0, 1, 0, 2, 2, 6, 1, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 3, 0, 3, 2, 28, 11, 2, 5, 2, 0, 3, 4, 2, 1, 0, 0, 3, 0, 5, 8, 87, 8, 4, 3, 2, 0, 4, 3, 70, 3, 3, 2, 59, 3, 61, 30, 18, 44, 0, 35, 1, 3, 1, 66, 3, 4, 0, 2, 2, 6, 5, 73, 11, 6, 0, 3, 0, 0, 1, 2, 4, 3, 0, 1, 5, 3, 9, 31, 1, 2, 0, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 5, 3, 3, 0, 2, 0, 0, 0, 1, 2, 1, 0, 2, 4, 2, 5, 19, 30, 82, 1, 85, 1, 2, 3, 55, 3, 5, 0, 3, 3, 5, 4, 46, 17, 6, 1, 3, 0, 0, 1, 3, 8, 3, 0, 2, 4, 3, 7, 23, 74, 198, 4, 94, 2, 3, 5, 58, 5, 9, 0, 5, 6, 7, 12, 103, 255, 121, 15, 45, 3, 3, 11, 30, 99, 36, 7, 22, 60, 23, 114, 217}, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/monarch.jpg", hashtype.UInt8{46, 9, 7, 13, 26, 2, 21, 54, 6, 1, 0, 1, 17, 1, 27, 55, 11, 2, 1, 1, 2, 0, 2, 2, 12, 1, 1, 1, 52, 2, 64, 41, 7, 0, 1, 1, 3, 0, 1, 2, 0, 0, 0, 0, 1, 0, 1, 3, 14, 0, 1, 1, 4, 0, 1, 1, 42, 1, 1, 2, 146, 3, 47, 27, 26, 2, 4, 6, 4, 0, 2, 4, 2, 0, 0, 0, 2, 0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 3, 0, 2, 4, 17, 2, 1, 2, 2, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 2, 55, 2, 3, 2, 3, 0, 2, 2, 129, 3, 5, 3, 97, 8, 30, 32, 6, 13, 0, 38, 2, 3, 1, 129, 0, 0, 0, 1, 1, 1, 1, 46, 0, 0, 0, 1, 0, 0, 0, 2, 0, 1, 0, 1, 2, 2, 4, 33, 0, 1, 0, 2, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 3, 1, 1, 0, 1, 0, 0, 0, 4, 2, 1, 0, 2, 9, 7, 3, 43, 18, 56, 1, 194, 1, 3, 1, 112, 1, 2, 0, 8, 1, 2, 1, 36, 1, 2, 0, 3, 0, 0, 0, 5, 2, 2, 0, 6, 1, 2, 2, 35, 29, 64, 1, 64, 1, 2, 1, 39, 1, 3, 0, 2, 1, 3, 2, 23, 57, 35, 2, 30, 2, 3, 2, 31, 50, 31, 3, 41, 42, 38, 24, 255}, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/peppers.jpg", hashtype.UInt8{72, 28, 12, 29, 25, 4, 24, 73, 8, 1, 1, 2, 16, 4, 37, 111, 22, 5, 2, 3, 4, 0, 4, 4, 18, 2, 2, 3, 48, 3, 130, 72, 9, 3, 1, 1, 2, 0, 2, 2, 0, 0, 0, 0, 2, 0, 2, 4, 19, 3, 2, 2, 4, 0, 1, 3, 44, 2, 2, 2, 111, 3, 45, 32, 27, 5, 3, 6, 6, 0, 4, 5, 3, 0, 0, 0, 3, 0, 2, 4, 4, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 4, 1, 3, 4, 22, 4, 2, 2, 3, 0, 2, 2, 1, 0, 0, 0, 2, 0, 2, 3, 61, 4, 2, 3, 4, 0, 2, 3, 191, 3, 4, 4, 79, 3, 33, 28, 14, 32, 2, 80, 4, 6, 3, 224, 1, 1, 0, 3, 2, 4, 2, 71, 2, 3, 0, 4, 0, 0, 0, 6, 1, 1, 0, 4, 2, 2, 5, 43, 1, 2, 0, 3, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 6, 2, 3, 0, 3, 0, 0, 0, 6, 2, 2, 0, 4, 4, 3, 4, 39, 31, 87, 3, 211, 4, 6, 5, 147, 2, 3, 0, 7, 1, 3, 2, 58, 5, 5, 0, 6, 0, 0, 0, 7, 3, 3, 0, 5, 2, 2, 3, 39, 49, 133, 2, 90, 2, 5, 3, 60, 3, 4, 0, 5, 2, 4, 4, 49, 119, 81, 5, 48, 4, 4, 3, 42, 63, 36, 5, 45, 34, 29, 34, 255}, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/tulips.jpg", hashtype.UInt8{101, 16, 17, 21, 78, 5, 28, 49, 16, 2, 1, 1, 24, 5, 23, 67, 17, 2, 2, 2, 4, 0, 4, 2, 19, 1, 1, 1, 45, 2, 60, 36, 16, 1, 2, 1, 11, 0, 2, 3, 1, 0, 0, 0, 2, 0, 1, 5, 21, 2, 2, 1, 14, 0, 3, 2, 52, 1, 3, 3, 217, 3, 54, 35, 89, 4, 10, 15, 23, 0, 5, 6, 10, 0, 0, 1, 5, 0, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 6, 0, 3, 3, 32, 5, 2, 3, 5, 0, 2, 3, 2, 0, 0, 0, 2, 0, 1, 1, 57, 2, 4, 2, 6, 0, 2, 2, 255, 3, 11, 8, 109, 4, 31, 28, 18, 22, 1, 52, 10, 16, 2, 225, 2, 1, 0, 3, 2, 4, 2, 64, 2, 1, 0, 1, 0, 0, 0, 3, 2, 1, 0, 2, 3, 2, 4, 37, 2, 1, 0, 3, 0, 1, 0, 10, 0, 0, 0, 0, 0, 0, 0, 4, 2, 2, 0, 2, 1, 0, 0, 8, 3, 2, 0, 8, 10, 8, 5, 82, 36, 57, 2, 233, 5, 6, 2, 106, 2, 3, 0, 10, 1, 1, 1, 30, 4, 3, 0, 3, 0, 0, 0, 3, 4, 3, 0, 7, 1, 2, 2, 28, 28, 74, 2, 67, 3, 3, 1, 35, 1, 4, 0, 4, 1, 1, 1, 18, 81, 46, 4, 45, 4, 4, 2, 30, 75, 47, 5, 90, 34, 30, 20, 228}, 256, 256, imghash.Bilinear, 1, 1},
}

func TestLBP_Calculate(t *testing.T) {
	for _, tt := range lbpCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := imghash.NewLBP(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType), imghash.WithGridSize(tt.gridX, tt.gridY))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img, err := imghash.OpenImage(tt.filename)
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

func ExampleLBP_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	lbp, err := imghash.NewLBP()
	if err != nil {
		panic(err)
	}
	hash, err := lbp.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [140 97 12 39 15 6 22 68 13 10 1 2 24 14 62 193 107 18 7 6 6 0 9 7 42 5 2 2 73 7 186 105 14 8 1 3 1 0 1 4 1 0 0 0 1 0 3 16 53 7 3 3 3 0 4 3 41 2 2 1 85 2 85 41 16 6 1 2 2 0 1 3 1 0 0 0 1 0 2 2 6 1 0 0 0 0 0 0 2 0 0 0 3 0 3 2 28 11 2 5 2 0 3 4 2 1 0 0 3 0 5 8 87 8 4 3 2 0 4 3 70 3 3 2 59 3 61 30 18 44 0 35 1 3 1 66 3 4 0 2 2 6 5 73 11 6 0 3 0 0 1 2 4 3 0 1 5 3 9 31 1 2 0 2 0 0 0 3 0 0 0 0 0 0 0 5 3 3 0 2 0 0 0 1 2 1 0 2 4 2 5 19 30 82 1 85 1 2 3 55 3 5 0 3 3 5 4 46 17 6 1 3 0 0 1 3 8 3 0 2 4 3 7 23 74 198 4 94 2 3 5 58 5 9 0 5 6 7 12 103 255 121 15 45 3 3 11 30 99 36 7 22 60 23 114 217]
}

var lbpDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imghash.Interpolation
	gridX       uint
	gridY       uint
}{
	{"assets/lena.jpg", "assets/cat.jpg", 1408.0953514028656, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/lena.jpg", "assets/monarch.jpg", 252.35802418352222, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/baboon.jpg", "assets/cat.jpg", 775.0799875837864, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/peppers.jpg", "assets/baboon.jpg", 783.2270861965849, 256, 256, imghash.Bilinear, 1, 1},
	{"assets/tulips.jpg", "assets/monarch.jpg", 407.0482265310037, 256, 256, imghash.Bilinear, 1, 1},
}

func TestLBP_Distance(t *testing.T) {
	for _, tt := range lbpDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := imghash.NewLBP(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType), imghash.WithGridSize(tt.gridX, tt.gridY))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img1, err := imghash.OpenImage(tt.firstImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.firstImage, err)
			}
			img2, err := imghash.OpenImage(tt.secondImage)
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
			dist, err := hash.Compare(h1, h2)
			if err != nil {
				t.Fatalf("failed to compute distance: %v", err)
			}
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}

func TestNewLBP_defaults(t *testing.T) {
	lbp, err := imghash.NewLBP()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := lbp.Calculate(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 256 {
		t.Errorf("expected hash length 256, got %d", h.Len())
	}
}

func TestNewLBP_invalidSize(t *testing.T) {
	_, err := imghash.NewLBP(imghash.WithSize(0, 0))
	if !errors.Is(err, imghash.ErrInvalidSize) {
		t.Errorf("expected imghash.ErrInvalidSize, got %v", err)
	}
}

func TestNewLBP_invalidGridSize(t *testing.T) {
	_, err := imghash.NewLBP(imghash.WithGridSize(0, 0))
	if !errors.Is(err, imghash.ErrInvalidGridSize) {
		t.Errorf("expected imghash.ErrInvalidGridSize, got %v", err)
	}
}

func TestNewLBP_gridSize(t *testing.T) {
	lbp, err := imghash.NewLBP(imghash.WithGridSize(2, 2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := lbp.Calculate(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 2*2*256 {
		t.Errorf("expected hash length %d, got %d", 2*2*256, h.Len())
	}
}
