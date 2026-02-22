package imghash_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

func TestGIST_Calculate(t *testing.T) {
	g, err := imghash.NewGIST()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := g.Calculate(img)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}
	got := h.(hashtype.Float64)
	if got.Len() != 320 {
		t.Fatalf("expected hash length 320, got %d", got.Len())
	}

	var energy float64
	for _, v := range got {
		energy += v * v
	}
	if energy <= 0 {
		t.Fatal("expected non-zero descriptor energy")
	}
}

func TestGIST_DistanceUsesCosineByDefault(t *testing.T) {
	g, err := imghash.NewGIST()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	img1, err := imghash.OpenImage("assets/lena.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	img2, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h1, err := g.Calculate(img1)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}
	h2, err := g.Calculate(img2)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}

	got, err := g.Compare(h1, h2)
	if err != nil {
		t.Fatalf("failed to compare: %v", err)
	}
	want, err := similarity.Cosine(h1, h2)
	if err != nil {
		t.Fatalf("failed to compute cosine distance: %v", err)
	}
	if !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestGIST_DistanceCosineBehavior(t *testing.T) {
	g, err := imghash.NewGIST()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := g.Compare(hashtype.Float64{1, 0}, hashtype.Float64{0, 1})
	if err != nil {
		t.Fatalf("failed to compare: %v", err)
	}
	if math.Abs(float64(got)-1) > 1e-12 {
		t.Fatalf("got %v, want 1", got)
	}
}

func TestNewGIST_CustomGridSize(t *testing.T) {
	g, err := imghash.NewGIST(imghash.WithGridSize(2, 2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := g.Calculate(img)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}
	if h.Len() != 80 {
		t.Fatalf("expected hash length 80, got %d", h.Len())
	}
}

func TestNewGIST_Errors(t *testing.T) {
	tests := []struct {
		name string
		opts []imghash.GISTOption
		err  error
	}{
		{"zero width", []imghash.GISTOption{imghash.WithSize(0, 64)}, imghash.ErrInvalidSize},
		{"zero height", []imghash.GISTOption{imghash.WithSize(64, 0)}, imghash.ErrInvalidSize},
		{"invalid interpolation", []imghash.GISTOption{imghash.WithInterpolation(imghash.Interpolation(999))}, imghash.ErrInvalidInterpolation},
		{"zero grid x", []imghash.GISTOption{imghash.WithGridSize(0, 4)}, imghash.ErrInvalidGridSize},
		{"zero grid y", []imghash.GISTOption{imghash.WithGridSize(4, 0)}, imghash.ErrInvalidGridSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := imghash.NewGIST(tt.opts...)
			if !errors.Is(err, tt.err) {
				t.Fatalf("got %v, want %v", err, tt.err)
			}
		})
	}
}

func ExampleGIST_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	g, err := imghash.NewGIST()
	if err != nil {
		panic(err)
	}
	hash, err := g.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash.Len())
	// Output: 320
}
