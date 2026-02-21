package hashtype_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/hashtype"
)

var float64StringTests = []struct {
	name   string
	hash   hashtype.Float64
	result string
}{
	{"float zero", hashtype.Float64{0}, "[0]"},
	{"float zero", hashtype.Float64{0.00112233}, "[0.00112233]"},
	{"pi value", hashtype.Float64{3.14159265358979323846264338327950288419716939937510582097494459}, "[3.141592653589793]"},
	{"multiple values", hashtype.Float64{0, 1.1, 2.22, 3.333, 4.4444}, "[0 1.1 2.22 3.333 4.4444]"},
	{"multiple values 2", hashtype.Float64{0.000000000000000012345678, 3892.1234567890123456789}, "[1.2345678e-17 3892.1234567890124]"},
}

func TestFloat64_String(t *testing.T) {
	for _, tt := range float64StringTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.String(); res != tt.result {
				t.Errorf("got %v, want %v", res, tt.result)
			}
		})
	}
}

func ExampleFloat64_String() {
	hash := hashtype.Float64{0.000000000000000012345678, 3892.1234567890123456789}
	fmt.Println(hash.String())
	// Output: [1.2345678e-17 3892.1234567890124]
}

var float64EqualTests = []struct {
	name   string
	h1     hashtype.Float64
	h2     hashtype.Float64
	expect bool
}{
	{"match 1", hashtype.Float64{1.123456789123456789}, hashtype.Float64{1.123456789123456789}, true},
	{"match 2", hashtype.Float64{0.693147180559945309417232121458176568075500134360255254120680009, 1.77245385090551602729816748334114518279754945612238712821380779, 1.61803398874989484820458683436563811772030917980576286213544862}, hashtype.Float64{0.693147180559945309417232121458176568075500134360255254120680009, 1.77245385090551602729816748334114518279754945612238712821380779, 1.61803398874989484820458683436563811772030917980576286213544862}, true},
	{"match 3", hashtype.Float64{0, 1.1, 2.22, 3.333}, hashtype.Float64{0, 1.1, 2.22, 3.333}, true},
	{"match 4", hashtype.Float64{}, hashtype.Float64{}, true},
	{"mismatch 1", hashtype.Float64{0, 1.1, 2.22, 3.333}, hashtype.Float64{0, 1.1, 2.22}, false},
	{"mismatch 2", hashtype.Float64{0, 1.1, 2.22}, hashtype.Float64{0, 2.22, 1.1}, false},
}

func TestFloat64_Equal(t *testing.T) {
	for _, tt := range float64EqualTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.h1.Equal(tt.h2); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func ExampleFloat64_Equal() {
	h1 := hashtype.Float64{0, 1.1, 2.22, 3.333}
	h2 := hashtype.Float64{0, 1.1, 2.22}
	h3 := hashtype.Float64{0, 1.1, 2.22, 3.333}
	fmt.Println(h1.Equal(h2))
	fmt.Println(h1.Equal(h3))
	// Output:
	// false
	// true
}

var float64LenTests = []struct {
	name   string
	hash   hashtype.Float64
	expect int
}{
	{"empty", hashtype.Float64{}, 0},
	{"one element", hashtype.Float64{1.5}, 1},
	{"four elements", hashtype.Float64{1.1, 2.2, 3.3, 4.4}, 4},
}

func TestFloat64_Len(t *testing.T) {
	for _, tt := range float64LenTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.Len(); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

var float64ValueAtTests = []struct {
	name   string
	hash   hashtype.Float64
	idx    int
	expect float64
}{
	{"first element", hashtype.Float64{3.14, 2.71}, 0, 3.14},
	{"second element", hashtype.Float64{3.14, 2.71}, 1, 2.71},
}

func TestFloat64_ValueAt(t *testing.T) {
	for _, tt := range float64ValueAtTests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.hash.ValueAt(tt.idx); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestFloat64_Distance(t *testing.T) {
	h1 := hashtype.Float64{1, 2, 3}
	h2 := hashtype.Float64{1, 2, 3}
	d, err := h1.Distance(h2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 0 {
		t.Errorf("got %v, want 0", d)
	}

	h3 := hashtype.Float64{4, 6, 3}
	d2, err := h1.Distance(h3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d2 != 5 {
		t.Errorf("got %v, want 5", d2)
	}
}
