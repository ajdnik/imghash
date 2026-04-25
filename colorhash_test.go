package imghash_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var colorHashCalculateTests = []struct {
	filename string
	binBits  uint
	hash     hashtype.Binary
}{
	{"assets/lena.jpg", 3, hashtype.Binary{1, 128, 3, 0, 0, 64}},
	{"assets/baboon.jpg", 3, hashtype.Binary{25, 16, 64, 32, 0, 0}},
	{"assets/cat.jpg", 3, hashtype.Binary{15, 128, 0, 0, 0, 0}},
	{"assets/monarch.jpg", 3, hashtype.Binary{3, 128, 1, 0, 0, 0}},
	{"assets/peppers.jpg", 3, hashtype.Binary{0, 48, 0, 64, 0, 0}},
	{"assets/tulips.jpg", 3, hashtype.Binary{36, 132, 64, 32, 0, 0}},
	{"assets/header.png", 3, hashtype.Binary{96, 1, 192, 0, 0, 0}},
	{"assets/logo.png", 3, hashtype.Binary{224, 0, 64, 0, 96, 0}},
	{"assets/cat.jpg", 4, hashtype.Binary{7, 240, 0, 0, 0, 0, 0}},
}

func TestColorHash_Calculate(t *testing.T) {
	for _, tt := range colorHashCalculateTests {
		t.Run(fmt.Sprintf("%s binbits=%d", tt.filename, tt.binBits), func(t *testing.T) {
			hash, err := imghash.NewColorHash(imghash.WithBinBits(tt.binBits))
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

func TestColorHash_Distance(t *testing.T) {
	hash, err := imghash.NewColorHash()
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}
	img1, err := imghash.OpenImage("assets/lena.jpg")
	if err != nil {
		t.Fatalf("failed to open first image: %v", err)
	}
	img2, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open second image: %v", err)
	}
	h1, err := hash.Calculate(img1)
	if err != nil {
		t.Fatalf("failed to calculate first hash: %v", err)
	}
	h2, err := hash.Calculate(img2)
	if err != nil {
		t.Fatalf("failed to calculate second hash: %v", err)
	}
	dist, err := hash.Compare(h1, h2)
	if err != nil {
		t.Fatalf("failed to compute distance: %v", err)
	}
	if !dist.Equal(similarity.Distance(6)) {
		t.Errorf("got %v, want %v", dist, similarity.Distance(6))
	}
}

func TestNewColorHash_InvalidBinBits(t *testing.T) {
	if _, err := imghash.NewColorHash(imghash.WithBinBits(0)); !errors.Is(err, imghash.ErrInvalidNumBins) {
		t.Fatalf("got %v, want %v", err, imghash.ErrInvalidNumBins)
	}
}
