package hashtype_test

import (
	"fmt"
	"testing"

	. "github.com/ajdnik/imghash/v2/hashtype"
)

var binaryStringTests = []struct {
	name   string
	hash   Binary
	result string
}{
	{"zero byte", Binary{0}, "[0]"},
	{"1 byte", Binary{1}, "[1]"},
	{"max byte", Binary{255}, "[255]"},
	{"two zero bytes", Binary{0, 0}, "[0 0]"},
	{"1 two bytes", Binary{1, 0}, "[1 0]"},
	{"1 two bytes second", Binary{0, 1}, "[0 1]"},
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
	hash := Binary{115, 247, 1}
	fmt.Println(hash.String())
	// Output: [115 247 1]
}

var binarySetErrorTests = []struct {
	name     string
	start    Binary
	position uint
	expect   error
}{
	{"first position one byte", Binary{0}, 0, nil},
	{"eight position one byte", Binary{0}, 7, nil},
	{"first position two bytes", Binary{0, 0}, 0, nil},
	{"last position two bytes", Binary{0, 0}, 15, nil},
	{"out of bounds", Binary{0, 0}, 17, ErrOutOfBounds},
}

func TestBinary_Set_error(t *testing.T) {
	for _, tt := range binarySetErrorTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.start.Set(tt.position); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}

}

var binarySetTests = []struct {
	name     string
	start    Binary
	position uint
	expect   Binary
}{
	{"first position one byte", Binary{0}, 0, Binary{1}},
	{"eight position one byte", Binary{0}, 7, Binary{128}},
	{"first position two bytes", Binary{0, 0}, 0, Binary{1, 0}},
	{"last position two bytes", Binary{0, 0}, 15, Binary{0, 128}},
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
	hash := Binary{0, 0}
	_ = hash.Set(0)
	_ = hash.Set(15)
	fmt.Println(hash.String())
	// Output: [1 128]
}

var binarySetReverseErrorTests = []struct {
	name     string
	start    Binary
	position uint
	expect   error
}{
	{"first position one byte", Binary{0}, 0, nil},
	{"eight position one byte", Binary{0}, 7, nil},
	{"first position two bytes", Binary{0, 0}, 0, nil},
	{"last position two bytes", Binary{0, 0}, 15, nil},
	{"out of bounds", Binary{0, 0}, 17, ErrOutOfBounds},
}

func TestBinary_SetReverse_error(t *testing.T) {
	for _, tt := range binarySetReverseErrorTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.start.SetReverse(tt.position); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}

}

var binarySetReverseTests = []struct {
	name     string
	start    Binary
	position uint
	expect   Binary
}{
	{"first position one byte", Binary{0}, 0, Binary{128}},
	{"eight position one byte", Binary{0}, 7, Binary{1}},
	{"first position two bytes", Binary{0, 0}, 0, Binary{128, 0}},
	{"last position two bytes", Binary{0, 0}, 15, Binary{0, 1}},
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
	hash := Binary{0, 0}
	_ = hash.SetReverse(0)
	_ = hash.SetReverse(15)
	fmt.Println(hash.String())
	// Output: [128 1]
}

var binaryEqualTests = []struct {
	name   string
	h1     Binary
	h2     Binary
	expect bool
}{
	{"match 1", Binary{0}, Binary{0}, true},
	{"match 2", Binary{1, 2, 3, 4, 5, 6, 7, 8, 9}, Binary{1, 2, 3, 4, 5, 6, 7, 8, 9}, true},
	{"match 3", Binary{231, 145, 13, 91, 22}, Binary{231, 145, 13, 91, 22}, true},
	{"mismatch 1", Binary{0}, Binary{1}, false},
	{"mismatch 2", Binary{1, 2, 3, 4}, Binary{1, 2, 3, 4, 5, 6}, false},
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
	h1 := Binary{1, 2, 128}
	h2 := Binary{2, 128}
	h3 := Binary{1, 2, 128}
	fmt.Println(h1.Equal(h2))
	fmt.Println(h1.Equal(h3))
	// Output:
	// false
	// true
}

var binaryLenTests = []struct {
	name   string
	hash   Binary
	expect int
}{
	{"empty", Binary{}, 0},
	{"one byte", Binary{0}, 1},
	{"three bytes", Binary{1, 2, 3}, 3},
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
	hash   Binary
	idx    int
	expect float64
}{
	{"first byte zero", Binary{0, 128}, 0, 0},
	{"second byte 128", Binary{0, 128}, 1, 128},
	{"byte 255", Binary{255}, 0, 255},
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

var binaryDistanceTests = []struct {
	name   string
	h1     Binary
	h2     Binary
	expect float64
}{
	{"identical", Binary{0xFF}, Binary{0xFF}, 0},
	{"one bit diff", Binary{0}, Binary{1}, 1},
	{"all bits diff", Binary{0x00}, Binary{0xFF}, 8},
	{"two bytes", Binary{0xFF, 0x00}, Binary{0x00, 0xFF}, 16},
	{"different lengths", Binary{0xFF, 0x00}, Binary{0xFF}, 0},
}

func TestBinary_Distance(t *testing.T) {
	for _, tt := range binaryDistanceTests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.h1.Distance(tt.h2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestBinary_Distance_incompatible(t *testing.T) {
	h1 := Binary{1, 2, 3}
	h2 := UInt8{1, 2, 3}
	_, err := h1.Distance(h2)
	if err != ErrIncompatibleHash {
		t.Errorf("got %v, want %v", err, ErrIncompatibleHash)
	}
}
