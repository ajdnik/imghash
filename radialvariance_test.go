package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
	"github.com/ajdnik/imghash/similarity"
)

var radialVarianceCalculateTests = []struct {
	filename string
	hash     hashtype.UInt8
	sigma    float64
	angles   int
}{
	{"assets/lena.jpg", hashtype.UInt8{132, 0, 255, 247, 54, 127, 136, 143, 64, 77, 159, 158, 113, 146, 101, 142, 138, 156, 140, 89, 128, 124, 124, 179, 108, 137, 134, 122, 145, 126, 134, 129, 118, 146, 133, 130, 124, 133, 132, 133}, 1, 180},
	{"assets/baboon.jpg", hashtype.UInt8{71, 68, 254, 38, 179, 87, 159, 70, 0, 106, 39, 66, 72, 101, 61, 65, 82, 82, 85, 53, 55, 71, 55, 60, 67, 86, 64, 58, 72, 68, 75, 70, 65, 76, 83, 69, 55, 68, 64, 73}, 1, 180},
	{"assets/cat.jpg", hashtype.UInt8{166, 246, 10, 0, 124, 194, 254, 203, 219, 156, 116, 176, 226, 154, 138, 184, 195, 174, 155, 143, 213, 154, 170, 209, 125, 152, 173, 167, 181, 170, 165, 183, 157, 179, 174, 162, 161, 171, 157, 194}, 1, 180},
	{"assets/monarch.jpg", hashtype.UInt8{126, 181, 255, 0, 155, 229, 74, 109, 156, 162, 120, 63, 136, 172, 106, 108, 105, 141, 153, 121, 106, 129, 93, 134, 157, 144, 131, 89, 122, 148, 111, 132, 137, 112, 140, 125, 135, 135, 98, 118}, 1, 180},
	{"assets/peppers.jpg", hashtype.UInt8{75, 56, 255, 131, 44, 1, 119, 93, 0, 152, 75, 61, 52, 49, 51, 98, 71, 44, 85, 100, 52, 64, 93, 61, 99, 80, 73, 95, 72, 79, 69, 75, 87, 74, 64, 57, 62, 76, 65, 74}, 1, 180},
	{"assets/tulips.jpg", hashtype.UInt8{94, 111, 255, 58, 0, 98, 37, 4, 111, 71, 71, 50, 129, 131, 75, 86, 70, 161, 167, 86, 138, 115, 113, 99, 92, 109, 61, 104, 131, 105, 94, 104, 106, 86, 90, 85, 118, 108, 80, 114}, 1, 180},
}

func TestRadialVariance_Calculate(t *testing.T) {
	for _, tt := range radialVarianceCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewRadialVarianceWithParams(tt.sigma, tt.angles)
			img, _ := imgproc.Read(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

var radialVarianceDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	sigma       float64
	angles      int
}{
	{"assets/lena.jpg", "assets/cat.jpg", 0.2480884829959218, 1, 180},
	{"assets/lena.jpg", "assets/monarch.jpg", 0.6338612038925339, 1, 180},
	{"assets/baboon.jpg", "assets/cat.jpg", 0.3171472992223605, 1, 180},
	{"assets/peppers.jpg", "assets/baboon.jpg", 0.5788259802149233, 1, 180},
	{"assets/tulips.jpg", "assets/monarch.jpg", 0.5192886736515012, 1, 180},
}

func TestRadialVariance_Distance(t *testing.T) {
	for _, tt := range radialVarianceDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewRadialVarianceWithParams(tt.sigma, tt.angles)
			img1, _ := imgproc.Read(tt.firstImage)
			img2, _ := imgproc.Read(tt.secondImage)
			h1 := hash.Calculate(img1)
			h2 := hash.Calculate(img2)
			dist, _ := similarity.PCCUInt8(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
