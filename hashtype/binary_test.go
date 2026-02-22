package hashtype_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/hashtype"
)

var binaryStringTests = []struct {
	name   string
	hash   hashtype.Binary
	result string
}{
	{"zero byte", hashtype.Binary{0}, "[0]"},
	{"1 byte", hashtype.Binary{1}, "[1]"},
	{"max byte", hashtype.Binary{255}, "[255]"},
	{"two zero bytes", hashtype.Binary{0, 0}, "[0 0]"},
	{"1 two bytes", hashtype.Binary{1, 0}, "[1 0]"},
	{"1 two bytes second", hashtype.Binary{0, 1}, "[0 1]"},
}

func TestBinary_String(t *testing.T) {
	for _, tt := range binaryStringTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.String(); res != tt.result {
				t.Errorf("got %v, want %v", res, tt.result)
			}
		})
	}
}

func ExampleBinary_String() {
	hash := hashtype.Binary{115, 247, 1}
	fmt.Println(hash.String())
	// Output: [115 247 1]
}

var binarySetErrorTests = []struct {
	name     string
	start    hashtype.Binary
	position uint
	expect   error
}{
	{"first position one byte", hashtype.Binary{0}, 0, nil},
	{"eight position one byte", hashtype.Binary{0}, 7, nil},
	{"first position two bytes", hashtype.Binary{0, 0}, 0, nil},
	{"last position two bytes", hashtype.Binary{0, 0}, 15, nil},
	{"out of bounds", hashtype.Binary{0, 0}, 17, hashtype.ErrOutOfBounds},
}

func TestBinary_Set_error(t *testing.T) {
	for _, tt := range binarySetErrorTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.start.Set(tt.position); !errors.Is(res, tt.expect) {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}

}

var binarySetTests = []struct {
	name     string
	start    hashtype.Binary
	position uint
	expect   hashtype.Binary
}{
	{"first position one byte", hashtype.Binary{0}, 0, hashtype.Binary{1}},
	{"eight position one byte", hashtype.Binary{0}, 7, hashtype.Binary{128}},
	{"first position two bytes", hashtype.Binary{0, 0}, 0, hashtype.Binary{1, 0}},
	{"last position two bytes", hashtype.Binary{0, 0}, 15, hashtype.Binary{0, 128}},
}

func TestBinary_Set(t *testing.T) {
	for _, tt := range binarySetTests {
		t.Run(tt.name, func(t *testing.T) {
			if _ = tt.start.Set(tt.position); !tt.start.Equal(tt.expect) {
				t.Errorf("got %08b, want %08b", tt.start, tt.expect)
			}
		})
	}
}

func ExampleBinary_Set() {
	hash := hashtype.Binary{0, 0}
	_ = hash.Set(0)
	_ = hash.Set(15)
	fmt.Println(hash.String())
	// Output: [1 128]
}

var binarySetReverseErrorTests = []struct {
	name     string
	start    hashtype.Binary
	position uint
	expect   error
}{
	{"first position one byte", hashtype.Binary{0}, 0, nil},
	{"eight position one byte", hashtype.Binary{0}, 7, nil},
	{"first position two bytes", hashtype.Binary{0, 0}, 0, nil},
	{"last position two bytes", hashtype.Binary{0, 0}, 15, nil},
	{"out of bounds", hashtype.Binary{0, 0}, 17, hashtype.ErrOutOfBounds},
}

func TestBinary_SetReverse_error(t *testing.T) {
	for _, tt := range binarySetReverseErrorTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.start.SetReverse(tt.position); !errors.Is(res, tt.expect) {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}

}

var binarySetReverseTests = []struct {
	name     string
	start    hashtype.Binary
	position uint
	expect   hashtype.Binary
}{
	{"first position one byte", hashtype.Binary{0}, 0, hashtype.Binary{128}},
	{"eight position one byte", hashtype.Binary{0}, 7, hashtype.Binary{1}},
	{"first position two bytes", hashtype.Binary{0, 0}, 0, hashtype.Binary{128, 0}},
	{"last position two bytes", hashtype.Binary{0, 0}, 15, hashtype.Binary{0, 1}},
}

func TestBinary_SetReverse(t *testing.T) {
	for _, tt := range binarySetReverseTests {
		t.Run(tt.name, func(t *testing.T) {
			if _ = tt.start.SetReverse(tt.position); !tt.start.Equal(tt.expect) {
				t.Errorf("got %08b, want %08b", tt.start, tt.expect)
			}
		})
	}
}

func ExampleBinary_SetReverse() {
	hash := hashtype.Binary{0, 0}
	_ = hash.SetReverse(0)
	_ = hash.SetReverse(15)
	fmt.Println(hash.String())
	// Output: [128 1]
}

var binaryEqualTests = []struct {
	name   string
	h1     hashtype.Binary
	h2     hashtype.Binary
	expect bool
}{
	{"match 1", hashtype.Binary{0}, hashtype.Binary{0}, true},
	{"match 2", hashtype.Binary{1, 2, 3, 4, 5, 6, 7, 8, 9}, hashtype.Binary{1, 2, 3, 4, 5, 6, 7, 8, 9}, true},
	{"match 3", hashtype.Binary{231, 145, 13, 91, 22}, hashtype.Binary{231, 145, 13, 91, 22}, true},
	{"mismatch 1", hashtype.Binary{0}, hashtype.Binary{1}, false},
	{"mismatch 2", hashtype.Binary{1, 2, 3, 4}, hashtype.Binary{1, 2, 3, 4, 5, 6}, false},
}

func TestBinary_Equal(t *testing.T) {
	for _, tt := range binaryEqualTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.h1.Equal(tt.h2); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func ExampleBinary_Equal() {
	h1 := hashtype.Binary{1, 2, 128}
	h2 := hashtype.Binary{2, 128}
	h3 := hashtype.Binary{1, 2, 128}
	fmt.Println(h1.Equal(h2))
	fmt.Println(h1.Equal(h3))
	// Output:
	// false
	// true
}

var binaryLenTests = []struct {
	name   string
	hash   hashtype.Binary
	expect int
}{
	{"empty", hashtype.Binary{}, 0},
	{"one byte", hashtype.Binary{0}, 1},
	{"three bytes", hashtype.Binary{1, 2, 3}, 3},
}

func TestBinary_Len(t *testing.T) {
	for _, tt := range binaryLenTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.Len(); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

var binaryValueAtTests = []struct {
	name   string
	hash   hashtype.Binary
	idx    int
	expect float64
}{
	{"first byte zero", hashtype.Binary{0, 128}, 0, 0},
	{"second byte 128", hashtype.Binary{0, 128}, 1, 128},
	{"byte 255", hashtype.Binary{255}, 0, 255},
}

func TestBinary_ValueAt(t *testing.T) {
	for _, tt := range binaryValueAtTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.ValueAt(tt.idx); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestNewBinary(t *testing.T) {
	tests := []struct {
		name string
		bits uint
		want int
	}{
		{"zero bits", 0, 0},
		{"8 bits", 8, 1},
		{"9 bits", 9, 2},
		{"63 bits", 63, 8},
		{"64 bits", 64, 8},
		{"65 bits", 65, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashtype.NewBinary(tt.bits).Len(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
