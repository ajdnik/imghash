package similarity_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var l1Uint8Tests = []struct {
	name  string
	hash1 hashtype.UInt8
	hash2 hashtype.UInt8
	out   similarity.Distance
}{
	{"same hashes", hashtype.UInt8{10, 20, 30}, hashtype.UInt8{10, 20, 30}, similarity.Distance(0)},
	{"single element", hashtype.UInt8{10}, hashtype.UInt8{15}, similarity.Distance(5)},
	{"multiple elements", hashtype.UInt8{10, 20, 30}, hashtype.UInt8{13, 25, 22}, similarity.Distance(16)},
	{"different lengths", hashtype.UInt8{10, 20}, hashtype.UInt8{15, 25, 30}, similarity.Distance(10)},
}

func TestL1Uint8(t *testing.T) {
	for _, tt := range l1Uint8Tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := similarity.L1(tt.hash1, tt.hash2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

var l1Float64Tests = []struct {
	name  string
	hash1 hashtype.Float64
	hash2 hashtype.Float64
	out   similarity.Distance
}{
	{"same hashes", hashtype.Float64{1.5, 2.5, 3.5}, hashtype.Float64{1.5, 2.5, 3.5}, similarity.Distance(0)},
	{"simple difference", hashtype.Float64{1.0, 2.0, 3.0}, hashtype.Float64{4.0, 6.0, 3.0}, similarity.Distance(7)},
	{"negative values", hashtype.Float64{-1.0, -2.0}, hashtype.Float64{1.0, 2.0}, similarity.Distance(6)},
}

func TestL1Float64(t *testing.T) {
	for _, tt := range l1Float64Tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := similarity.L1(tt.hash1, tt.hash2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

func ExampleL1() {
	hash1 := hashtype.UInt8{60, 67, 86, 64, 58, 72, 68, 75}
	hash2 := hashtype.UInt8{143, 213, 154, 170, 209, 125, 152, 173}

	dist, _ := similarity.L1(hash1, hash2)
	fmt.Println(dist)
	// Output:
	// 789
}
