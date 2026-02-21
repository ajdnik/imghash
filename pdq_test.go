package imghash_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var pdqCalculateTests = []struct {
	filename   string
	hash       hashtype.Binary
	resizeType imghash.Interpolation
}{
	{"assets/lena.jpg", hashtype.Binary{191, 0, 123, 5, 107, 81, 181, 32, 255, 82, 205, 10, 101, 173, 121, 166, 146, 229, 208, 180, 210, 252, 0, 219, 64, 127, 24, 219, 0, 253, 88, 125}, imghash.Bilinear},
	{"assets/baboon.jpg", hashtype.Binary{255, 0, 5, 234, 7, 169, 191, 1, 225, 248, 133, 43, 17, 108, 251, 128, 16, 127, 131, 39, 0, 255, 234, 85, 64, 255, 226, 221, 64, 255, 106, 85}, imghash.Bilinear},
	{"assets/cat.jpg", hashtype.Binary{255, 0, 255, 0, 233, 42, 61, 168, 171, 5, 179, 17, 119, 2, 85, 174, 152, 255, 140, 217, 96, 247, 72, 93, 0, 253, 2, 119, 4, 255, 0, 255}, imghash.Bilinear},
	{"assets/monarch.jpg", hashtype.Binary{191, 0, 255, 1, 38, 136, 63, 0, 153, 247, 237, 41, 128, 240, 71, 2, 198, 93, 128, 253, 104, 245, 0, 255, 64, 255, 0, 255, 64, 255, 224, 255}, imghash.Bilinear},
	{"assets/peppers.jpg", hashtype.Binary{239, 1, 253, 144, 63, 0, 9, 229, 227, 8, 131, 218, 131, 218, 187, 87, 192, 246, 148, 10, 80, 255, 6, 255, 0, 253, 38, 109, 24, 255, 118, 36}, imghash.Bilinear},
	{"assets/tulips.jpg", hashtype.Binary{171, 0, 127, 0, 126, 4, 127, 8, 63, 8, 123, 16, 52, 120, 165, 11, 160, 251, 47, 95, 200, 55, 202, 53, 128, 247, 128, 253, 212, 255, 192, 247}, imghash.Bilinear},
}

func TestPDQ_Calculate(t *testing.T) {
	for _, tt := range pdqCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := imghash.NewPDQ(imghash.WithInterpolation(tt.resizeType))
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
			res := result.(hashtype.Binary)
			if !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExamplePDQ_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	pdq, err := imghash.NewPDQ()
	if err != nil {
		panic(err)
	}
	hash, err := pdq.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
	// Output: [255 0 255 0 233 42 61 168 171 5 179 17 119 2 85 174 152 255 140 217 96 247 72 93 0 253 2 119 4 255 0 255]
}

var pdqDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	resizeType  imghash.Interpolation
}{
	{"assets/lena.jpg", "assets/cat.jpg", 90, imghash.Bilinear},
	{"assets/lena.jpg", "assets/monarch.jpg", 80, imghash.Bilinear},
	{"assets/baboon.jpg", "assets/cat.jpg", 102, imghash.Bilinear},
	{"assets/peppers.jpg", "assets/baboon.jpg", 98, imghash.Bilinear},
	{"assets/tulips.jpg", "assets/monarch.jpg", 84, imghash.Bilinear},
}

func TestPDQ_Distance(t *testing.T) {
	for _, tt := range pdqDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := imghash.NewPDQ(imghash.WithInterpolation(tt.resizeType))
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
