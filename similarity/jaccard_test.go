package similarity_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

func TestJaccardBinary(t *testing.T) {
	tests := []struct {
		name string
		h1   hashtype.Binary
		h2   hashtype.Binary
		out  similarity.Distance
	}{
		{"same", hashtype.Binary{0xFF}, hashtype.Binary{0xFF}, 0},
		{"disjoint", hashtype.Binary{0xF0}, hashtype.Binary{0x0F}, 1},
		{"partial overlap", hashtype.Binary{0b11110000}, hashtype.Binary{0b11001100}, similarity.Distance(0.6666666666666667)},
		{"both empty", hashtype.Binary{}, hashtype.Binary{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := similarity.Jaccard(tt.h1, tt.h2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tt.out) {
				t.Fatalf("got %v, want %v", got, tt.out)
			}
		})
	}
}

func TestJaccardMinHashSignatures(t *testing.T) {
	tests := []struct {
		name string
		h1   hashtype.Float64
		h2   hashtype.Float64
		out  similarity.Distance
	}{
		{"same", hashtype.Float64{1, 2, 3, 4}, hashtype.Float64{1, 2, 3, 4}, 0},
		{"half match", hashtype.Float64{1, 2, 3, 4}, hashtype.Float64{1, 9, 3, 8}, 0.5},
		{"no match", hashtype.Float64{1, 2, 3, 4}, hashtype.Float64{9, 8, 7, 6}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := similarity.Jaccard(tt.h1, tt.h2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tt.out) {
				t.Fatalf("got %v, want %v", got, tt.out)
			}
		})
	}
}

func TestJaccardUInt8Signatures(t *testing.T) {
	got, err := similarity.Jaccard(hashtype.UInt8{1, 2, 3}, hashtype.UInt8{1, 9, 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(similarity.Distance(0.33333333333333337)) {
		t.Fatalf("got %v, want 0.33333333333333337", got)
	}
}

func TestJaccardErrors(t *testing.T) {
	_, err := similarity.Jaccard(hashtype.Float64{1}, hashtype.Float64{1, 2})
	if !errors.Is(err, similarity.ErrNotSameLength) {
		t.Fatalf("got %v, want %v", err, similarity.ErrNotSameLength)
	}

	_, err = similarity.Jaccard(hashtype.Binary{0xFF}, hashtype.UInt8{1})
	if !errors.Is(err, hashtype.ErrIncompatibleHash) {
		t.Fatalf("got %v, want %v", err, hashtype.ErrIncompatibleHash)
	}
}

func ExampleJaccard() {
	dist, _ := similarity.Jaccard(hashtype.Binary{0b11110000}, hashtype.Binary{0b11001100})
	fmt.Println(dist)
	// Output:
	// 0.6666666666666667
}
