package hashtype_test

import (
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

func TestBinary_SetError(t *testing.T) {
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
