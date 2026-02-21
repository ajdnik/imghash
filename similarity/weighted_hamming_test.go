package similarity_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/hashtype"
	. "github.com/ajdnik/imghash/v2/similarity"
)

var weightedHammingTests = []struct {
	name     string
	hash1    hashtype.Binary
	hash2    hashtype.Binary
	weights  []float64
	distance Distance
}{
	{"uniform weights", hashtype.Binary{1}, hashtype.Binary{2}, []float64{1.0}, 2},
	{"double weight", hashtype.Binary{1}, hashtype.Binary{2}, []float64{2.0}, 4},
	{"zero weight", hashtype.Binary{1}, hashtype.Binary{2}, []float64{0.0}, 0},
	{"two bytes uniform", hashtype.Binary{1, 1}, hashtype.Binary{2, 2}, []float64{1.0, 1.0}, 4},
	{"two bytes varied", hashtype.Binary{1, 1}, hashtype.Binary{2, 2}, []float64{1.0, 3.0}, 8},
	{"sample hashes", hashtype.Binary{15, 131, 192, 224, 192, 252, 255, 255}, hashtype.Binary{24, 60, 126, 126, 126, 126, 60, 0}, []float64{1, 1, 1, 1, 1, 1, 1, 1}, 42},
	{"sample hashes weighted", hashtype.Binary{15, 131, 192, 224, 192, 252, 255, 255}, hashtype.Binary{24, 60, 126, 126, 126, 126, 60, 0}, []float64{2, 1, 1, 1, 1, 1, 1, 0.5}, 42},
}

func TestWeightedHamming(t *testing.T) {
	for _, tt := range weightedHammingTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := WeightedHamming(tt.hash1, tt.hash2, tt.weights)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.distance) {
				t.Errorf("got %v, want %v", res, tt.distance)
			}
		})
	}
}

func TestWeightedHamming_notBinary(t *testing.T) {
	h1 := hashtype.UInt8{1, 2, 3}
	h2 := hashtype.UInt8{4, 5, 6}
	_, err := WeightedHamming(h1, h2, []float64{1, 1, 1})
	if err != ErrNotBinaryHash {
		t.Errorf("got %v, want %v", err, ErrNotBinaryHash)
	}
}

func TestWeightedHamming_lengthMismatch(t *testing.T) {
	h1 := hashtype.Binary{1, 2}
	h2 := hashtype.Binary{3, 4}
	_, err := WeightedHamming(h1, h2, []float64{1.0})
	if err != ErrWeightLengthMismatch {
		t.Errorf("got %v, want %v", err, ErrWeightLengthMismatch)
	}
}

func ExampleWeightedHamming() {
	hash1 := hashtype.Binary{15, 131, 192, 224, 192, 252, 255, 255}
	hash2 := hashtype.Binary{24, 60, 126, 126, 126, 126, 60, 0}
	weights := []float64{1, 1, 1, 1, 1, 1, 1, 1}

	res, err := WeightedHamming(hash1, hash2, weights)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// Output:
	// 42
}
