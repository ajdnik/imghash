package hashtype_test

import (
	"fmt"
	"testing"

	. "github.com/ajdnik/imghash/v2/hashtype"
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

var uint8LenTests = []struct {
	name   string
	hash   UInt8
	expect int
}{
	{"empty", UInt8{}, 0},
	{"one element", UInt8{42}, 1},
	{"five elements", UInt8{1, 2, 3, 4, 5}, 5},
}

func TestUInt8_Len(t *testing.T) {
	for _, tt := range uint8LenTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.Len(); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

var uint8ValueAtTests = []struct {
	name   string
	hash   UInt8
	idx    int
	expect float64
}{
	{"first element", UInt8{100, 200}, 0, 100},
	{"second element", UInt8{100, 200}, 1, 200},
	{"zero", UInt8{0}, 0, 0},
}

func TestUInt8_ValueAt(t *testing.T) {
	for _, tt := range uint8ValueAtTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.ValueAt(tt.idx); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestUInt8_Distance(t *testing.T) {
	h1 := UInt8{0, 0}
	h2 := UInt8{3, 4}
	d, err := h1.Distance(h2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 5 {
		t.Errorf("got %v, want 5", d)
	}

	d2, err := h1.Distance(h1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d2 != 0 {
		t.Errorf("got %v, want 0", d2)
	}
}
