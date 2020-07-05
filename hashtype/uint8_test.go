package hashtype_test

import (
	"fmt"
	"testing"

	. "github.com/ajdnik/imghash/hashtype"
)

var uint8StringTests = []struct {
	name   string
	hash   UInt8
	result string
}{
	{"empty hash", UInt8{}, "[]"},
	{"single value", UInt8{112}, "[112]"},
	{"multiple values", UInt8{1, 2, 89, 113}, "[1 2 89 113]"},
}

func TestUInt8_String(t *testing.T) {
	for _, tt := range uint8StringTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.String(); res != tt.result {
				t.Errorf("got %v, want %v", res, tt.result)
			}
		})
	}
}

func ExampleUInt8_String() {
	hash := UInt8{0, 1, 2, 3, 4, 5}
	fmt.Println(hash.String())
	// Output: [0 1 2 3 4 5]
}

var uint8EqualTests = []struct {
	name   string
	h1     UInt8
	h2     UInt8
	expect bool
}{
	{"match 1", UInt8{}, UInt8{}, true},
	{"match 2", UInt8{1, 2, 3, 4}, UInt8{1, 2, 3, 4}, true},
	{"match 3", UInt8{123, 255, 61, 72}, UInt8{123, 255, 61, 72}, true},
	{"mismatch 1", UInt8{}, UInt8{156}, false},
	{"mismatch 2", UInt8{1, 2, 3, 4}, UInt8{1, 3, 2, 4}, false},
	{"mismatch 3", UInt8{1, 2, 3, 4}, UInt8{1, 2, 3, 4, 5}, false},
}

func TestUInt8_Equal(t *testing.T) {
	for _, tt := range uint8EqualTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.h1.Equal(tt.h2); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func ExampleUInt8_Equal() {
	h1 := UInt8{0, 1, 2, 3}
	h2 := UInt8{0, 1, 2}
	h3 := UInt8{0, 1, 2, 3}
	fmt.Println(h1.Equal(h2))
	fmt.Println(h1.Equal(h3))
	// Output:
	// false
	// true
}
