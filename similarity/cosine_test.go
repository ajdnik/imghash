package similarity_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/hashtype"
	. "github.com/ajdnik/imghash/v2/similarity"
)

var cosineUint8Tests = []struct {
	name  string
	hash1 hashtype.UInt8
	hash2 hashtype.UInt8
	out   Distance
}{
	{"identical hashes", hashtype.UInt8{10, 20, 30}, hashtype.UInt8{10, 20, 30}, Distance(0)},
	{"proportional hashes", hashtype.UInt8{1, 2, 3}, hashtype.UInt8{2, 4, 6}, Distance(0)},
	{"orthogonal-like", hashtype.UInt8{1, 0, 0}, hashtype.UInt8{0, 1, 0}, Distance(1)},
	{"different hashes", hashtype.UInt8{60, 67, 86, 64}, hashtype.UInt8{143, 213, 154, 170}, Distance(0.024536201729695284)},
}

func TestCosineUint8(t *testing.T) {
	for _, tt := range cosineUint8Tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Cosine(tt.hash1, tt.hash2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

var cosineFloat64Tests = []struct {
	name  string
	hash1 hashtype.Float64
	hash2 hashtype.Float64
	out   Distance
}{
	{"identical", hashtype.Float64{1.0, 2.0, 3.0}, hashtype.Float64{1.0, 2.0, 3.0}, Distance(0)},
	{"opposite", hashtype.Float64{1.0, 0.0}, hashtype.Float64{-1.0, 0.0}, Distance(2)},
	{"zero vectors", hashtype.Float64{0.0, 0.0}, hashtype.Float64{0.0, 0.0}, Distance(0)},
	{"one zero vector", hashtype.Float64{1.0, 2.0}, hashtype.Float64{0.0, 0.0}, Distance(0)},
}

func TestCosineFloat64(t *testing.T) {
	for _, tt := range cosineFloat64Tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Cosine(tt.hash1, tt.hash2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

func ExampleCosine() {
	hash1 := hashtype.UInt8{60, 67, 86, 64, 58, 72, 68, 75}
	hash2 := hashtype.UInt8{143, 213, 154, 170, 209, 125, 152, 173}

	dist, _ := Cosine(hash1, hash2)
	fmt.Println(dist)
	// Output:
	// 0.028706372204543307
}
