package imghash_test

import (
	"fmt"
	"math"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/similarity"
)

var radialVarianceCalculateTests = []struct {
	filename string
	hash     hashtype.UInt8
	sigma    float64
	angles   int
}{
	{"assets/lena.jpg", hashtype.UInt8{132, 0, 255, 247, 53, 127, 136, 143, 64, 76, 159, 158, 112, 146, 100, 142, 137, 156, 140, 89, 128, 124, 124, 179, 107, 137, 133, 122, 145, 125, 134, 129, 118, 146, 133, 129, 124, 133, 131, 133}, 1, 180},
	{"assets/baboon.jpg", hashtype.UInt8{71, 68, 255, 38, 179, 87, 159, 70, 0, 107, 39, 66, 73, 101, 62, 65, 82, 82, 85, 54, 55, 71, 55, 60, 67, 86, 64, 58, 72, 68, 75, 71, 65, 76, 83, 69, 55, 68, 65, 73}, 1, 180},
	{"assets/cat.jpg", hashtype.UInt8{166, 246, 10, 0, 124, 193, 255, 203, 219, 156, 116, 175, 226, 154, 138, 185, 195, 174, 155, 143, 213, 154, 170, 210, 125, 152, 173, 167, 181, 170, 165, 183, 157, 179, 174, 161, 161, 171, 157, 194}, 1, 180},
	{"assets/monarch.jpg", hashtype.UInt8{126, 181, 255, 0, 155, 228, 75, 108, 156, 162, 119, 63, 136, 172, 106, 108, 104, 140, 153, 121, 106, 129, 93, 134, 157, 144, 131, 89, 122, 148, 111, 132, 137, 112, 139, 125, 135, 135, 98, 118}, 1, 180},
	{"assets/peppers.jpg", hashtype.UInt8{75, 57, 255, 132, 44, 1, 119, 93, 0, 153, 75, 61, 52, 49, 50, 98, 71, 44, 86, 100, 52, 64, 93, 61, 99, 80, 74, 95, 72, 79, 69, 75, 88, 75, 64, 57, 62, 76, 65, 74}, 1, 180},
	{"assets/tulips.jpg", hashtype.UInt8{94, 111, 255, 58, 0, 98, 37, 3, 111, 71, 71, 50, 129, 131, 75, 86, 70, 161, 168, 86, 138, 115, 113, 99, 92, 109, 60, 104, 131, 105, 95, 104, 106, 86, 90, 85, 118, 108, 79, 114}, 1, 180},
}

func uint8ApproxEqual(a, b hashtype.UInt8, maxDiff int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		diff := int(a[i]) - int(b[i])
		if diff < -maxDiff || diff > maxDiff {
			return false
		}
	}
	return true
}

func TestRadialVariance_Calculate(t *testing.T) {
	for _, tt := range radialVarianceCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewRadialVariance(WithSigma(tt.sigma), WithAngles(tt.angles))
			img, err := OpenImage(tt.filename)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.filename, err)
			}
			result, err := hash.Calculate(img)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			res := result.(hashtype.UInt8)
			if !uint8ApproxEqual(res, tt.hash, 1) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleRadialVariance_Calculate() {
	// Read image from file
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	// Create new Radial Variance Hash using default parameters
	rad := NewRadialVariance()
	// Calculate hash
	hash, err := rad.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}

var radialVarianceDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	sigma       float64
	angles      int
}{
	{"assets/lena.jpg", "assets/cat.jpg", 0.247818888750317, 1, 180},
	{"assets/lena.jpg", "assets/monarch.jpg", 0.6339729738012116, 1, 180},
	{"assets/baboon.jpg", "assets/cat.jpg", 0.31577803436311935, 1, 180},
	{"assets/peppers.jpg", "assets/baboon.jpg", 0.5791217549142161, 1, 180},
	{"assets/tulips.jpg", "assets/monarch.jpg", 0.5199454464617648, 1, 180},
}

func TestRadialVariance_Distance(t *testing.T) {
	for _, tt := range radialVarianceDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewRadialVariance(WithSigma(tt.sigma), WithAngles(tt.angles))
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
			dist, err := similarity.PCC(h1, h2)
			if err != nil {
				t.Fatalf("failed to compute distance: %v", err)
			}
			if math.Abs(float64(dist)-float64(tt.distance)) > 0.01 {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
