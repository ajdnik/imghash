package similarity_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/hashtype"
	. "github.com/ajdnik/imghash/v2/similarity"
)

var chiSquareUint8Tests = []struct {
	name  string
	hash1 hashtype.UInt8
	hash2 hashtype.UInt8
	out   Distance
}{
	{"same hashes", hashtype.UInt8{10, 20, 30}, hashtype.UInt8{10, 20, 30}, Distance(0)},
	{"simple case", hashtype.UInt8{10, 20, 30}, hashtype.UInt8{15, 25, 35}, Distance(1.9401709401709402)},
	{"with zero pair", hashtype.UInt8{0, 20, 30}, hashtype.UInt8{0, 25, 35}, Distance(0.9401709401709402)},
	{"one zero", hashtype.UInt8{0, 10}, hashtype.UInt8{5, 10}, Distance(5)},
}

func TestChiSquareUint8(t *testing.T) {
	for _, tt := range chiSquareUint8Tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ChiSquare(tt.hash1, tt.hash2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

var chiSquareFloat64Tests = []struct {
	name  string
	hash1 hashtype.Float64
	hash2 hashtype.Float64
	out   Distance
}{
	{"same hashes", hashtype.Float64{1.0, 2.0, 3.0}, hashtype.Float64{1.0, 2.0, 3.0}, Distance(0)},
	{"simple case", hashtype.Float64{1.0, 2.0, 3.0}, hashtype.Float64{2.0, 3.0, 4.0}, Distance(0.6761904761904762)},
	{"both zeros skipped", hashtype.Float64{0.0, 2.0}, hashtype.Float64{0.0, 4.0}, Distance(0.6666666666666666)},
}

func TestChiSquareFloat64(t *testing.T) {
	for _, tt := range chiSquareFloat64Tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ChiSquare(tt.hash1, tt.hash2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

func ExampleChiSquare() {
	hash1 := hashtype.UInt8{10, 20, 30, 40}
	hash2 := hashtype.UInt8{15, 25, 35, 45}

	dist, _ := ChiSquare(hash1, hash2)
	fmt.Println(dist)
	// Output:
	// 2.2342885872297638
}
