package imghash_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

func TestBoVW_CalculateHistogram(t *testing.T) {
	b, err := imghash.NewBoVW(
		imghash.WithSize(96, 96),
		imghash.WithBoVWFeature(imghash.BoVWORB),
		imghash.WithBoVWStorage(imghash.BoVWHistogram),
		imghash.WithVocabularySize(128),
		imghash.WithMaxKeypoints(200),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := b.Calculate(img)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}

	got, ok := h.(hashtype.Float64)
	if !ok {
		t.Fatalf("expected hashtype.Float64, got %T", h)
	}
	if got.Len() != 128 {
		t.Fatalf("expected histogram length 128, got %d", got.Len())
	}
}

func TestBoVW_CalculateMinHash(t *testing.T) {
	b, err := imghash.NewBoVW(
		imghash.WithSize(96, 96),
		imghash.WithBoVWFeature(imghash.BoVWAKAZE),
		imghash.WithBoVWStorage(imghash.BoVWMinHash),
		imghash.WithMinHashSize(48),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := b.Calculate(img)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}
	got, ok := h.(hashtype.Float64)
	if !ok {
		t.Fatalf("expected hashtype.Float64, got %T", h)
	}
	if got.Len() != 48 {
		t.Fatalf("expected signature length 48, got %d", got.Len())
	}
}

func TestBoVW_CalculateSimHash(t *testing.T) {
	b, err := imghash.NewBoVW(
		imghash.WithSize(96, 96),
		imghash.WithBoVWStorage(imghash.BoVWSimHash),
		imghash.WithSimHashBits(192),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := b.Calculate(img)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}
	got, ok := h.(hashtype.Binary)
	if !ok {
		t.Fatalf("expected hashtype.Binary, got %T", h)
	}
	if got.Len() != 24 {
		t.Fatalf("expected 24 bytes, got %d", got.Len())
	}
}

func TestBoVW_CompareDefaults(t *testing.T) {
	img1, err := imghash.OpenImage("assets/lena.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	img2, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}

	tests := []struct {
		name string
		opts []imghash.BoVWOption
		want imghash.DistanceFunc
	}{
		{
			name: "histogram uses cosine",
			opts: []imghash.BoVWOption{
				imghash.WithSize(96, 96),
				imghash.WithBoVWStorage(imghash.BoVWHistogram),
				imghash.WithVocabularySize(64),
			},
			want: similarity.Cosine,
		},
		{
			name: "minhash uses jaccard",
			opts: []imghash.BoVWOption{
				imghash.WithSize(96, 96),
				imghash.WithBoVWStorage(imghash.BoVWMinHash),
				imghash.WithMinHashSize(32),
			},
			want: similarity.Jaccard,
		},
		{
			name: "simhash uses jaccard",
			opts: []imghash.BoVWOption{
				imghash.WithSize(96, 96),
				imghash.WithBoVWStorage(imghash.BoVWSimHash),
				imghash.WithSimHashBits(128),
			},
			want: similarity.Jaccard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, berr := imghash.NewBoVW(tt.opts...)
			if berr != nil {
				t.Fatalf("unexpected error: %v", berr)
			}
			h1, berr := b.Calculate(img1)
			if berr != nil {
				t.Fatalf("failed to calculate hash: %v", berr)
			}
			h2, berr := b.Calculate(img2)
			if berr != nil {
				t.Fatalf("failed to calculate hash: %v", berr)
			}

			got, berr := b.Compare(h1, h2)
			if berr != nil {
				t.Fatalf("failed to compare hashes: %v", berr)
			}
			want, berr := tt.want(h1, h2)
			if berr != nil {
				t.Fatalf("failed to compute expected distance: %v", berr)
			}
			if !got.Equal(want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}

func TestBoVW_CompareTypeValidation(t *testing.T) {
	minHash, err := imghash.NewBoVW(imghash.WithBoVWStorage(imghash.BoVWMinHash))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = minHash.Compare(hashtype.Binary{0xFF}, hashtype.Binary{0xFF})
	if !errors.Is(err, imghash.ErrIncompatibleHash) {
		t.Fatalf("got %v, want %v", err, imghash.ErrIncompatibleHash)
	}

	simHash, err := imghash.NewBoVW(imghash.WithBoVWStorage(imghash.BoVWSimHash))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = simHash.Compare(hashtype.Float64{1}, hashtype.Float64{1})
	if !errors.Is(err, imghash.ErrIncompatibleHash) {
		t.Fatalf("got %v, want %v", err, imghash.ErrIncompatibleHash)
	}
}

func TestNewBoVW_Errors(t *testing.T) {
	tests := []struct {
		name string
		opts []imghash.BoVWOption
		err  error
	}{
		{"invalid size", []imghash.BoVWOption{imghash.WithSize(0, 64)}, imghash.ErrInvalidSize},
		{"invalid feature", []imghash.BoVWOption{imghash.WithBoVWFeature(imghash.BoVWFeatureType(99))}, imghash.ErrInvalidBoVWFeatureType},
		{"invalid storage", []imghash.BoVWOption{imghash.WithBoVWStorage(imghash.BoVWStorageType(99))}, imghash.ErrInvalidBoVWStorageType},
		{"invalid vocab size", []imghash.BoVWOption{imghash.WithVocabularySize(0)}, imghash.ErrInvalidVocabularySize},
		{"invalid keypoints", []imghash.BoVWOption{imghash.WithMaxKeypoints(0)}, imghash.ErrInvalidKeypoints},
		{"invalid minhash size", []imghash.BoVWOption{imghash.WithMinHashSize(0)}, imghash.ErrInvalidSignatureSize},
		{"invalid simhash bits", []imghash.BoVWOption{imghash.WithSimHashBits(0)}, imghash.ErrInvalidSignatureSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := imghash.NewBoVW(tt.opts...)
			if !errors.Is(err, tt.err) {
				t.Fatalf("got %v, want %v", err, tt.err)
			}
		})
	}
}

func ExampleBoVW_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	bovw, err := imghash.NewBoVW(
		imghash.WithBoVWFeature(imghash.BoVWAKAZE),
		imghash.WithBoVWStorage(imghash.BoVWMinHash),
		imghash.WithMinHashSize(32),
	)
	if err != nil {
		panic(err)
	}
	hash, err := bovw.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash.Len())
	// Output: 32
}
