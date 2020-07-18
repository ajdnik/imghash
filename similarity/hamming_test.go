package similarity_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/hashtype"
	. "github.com/ajdnik/imghash/similarity"
)

var hammingTests = []struct {
	name     string
	hash1    hashtype.Binary
	hash2    hashtype.Binary
	distance Distance
}{
	{"two bit difference", hashtype.Binary{1}, hashtype.Binary{2}, 2},
	{"reverse hashes", hashtype.Binary{2}, hashtype.Binary{1}, 2},
	{"two bytes", hashtype.Binary{1, 1}, hashtype.Binary{2, 2}, 4},
	{"sample1 vs sample2", hashtype.Binary{15, 131, 192, 224, 192, 252, 255, 255}, hashtype.Binary{24, 60, 126, 126, 126, 126, 60, 0}, 42},
	{"sample1 vs sample3", hashtype.Binary{15, 131, 192, 224, 192, 252, 255, 255}, hashtype.Binary{63, 131, 192, 224, 192, 252, 255, 63}, 4},
	{"sample1 vs sample4", hashtype.Binary{15, 131, 192, 224, 192, 252, 255, 255}, hashtype.Binary{16, 60, 124, 126, 124, 124, 60, 24}, 38},
	{"lena vs cat", hashtype.Binary{125, 121, 185, 149, 213, 197, 112, 52}, hashtype.Binary{255, 255, 143, 3, 33, 65, 32, 27}, 27},
}

func TestHamming(t *testing.T) {
	for _, tt := range hammingTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if res := Hamming(tt.hash1, tt.hash2); !res.Equal(tt.distance) {
				t.Errorf("got %v, want %v", res, tt.distance)
			}
		})
	}
}

func ExampleHamming() {
	hash1 := hashtype.Binary{15, 131, 192, 224, 192, 252, 255, 255}
	hash2 := hashtype.Binary{24, 60, 126, 126, 126, 126, 60, 0}
	hash3 := hashtype.Binary{63, 131, 192, 224, 192, 252, 255, 63}

	fmt.Println(Hamming(hash1, hash2))
	fmt.Println(Hamming(hash1, hash3))
	// Output:
	// 42
	// 4
}
