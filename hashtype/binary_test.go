package hashtype_test

import (
	"fmt"
	"testing"

	. "github.com/ajdnik/imghash/hashtype"
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
			if tt.start.Set(tt.position); !tt.start.Equal(tt.expect) {
				t.Errorf("got %08b, want %08b", tt.start, tt.expect)
			}
		})
	}
}

func ExampleBinary_Set() {
	hash := Binary{0, 0}
	hash.Set(0)
	hash.Set(15)
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
			if tt.start.SetReverse(tt.position); !tt.start.Equal(tt.expect) {
				t.Errorf("got %08b, want %08b", tt.start, tt.expect)
			}
		})
	}
}

func ExampleBinary_SetReverse() {
	hash := Binary{0, 0}
	hash.SetReverse(0)
	hash.SetReverse(15)
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
