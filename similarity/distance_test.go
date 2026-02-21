package similarity_test

import (
	"fmt"
	"testing"

	. "github.com/ajdnik/imghash/v2/similarity"
)

var distanceEqualTests = []struct {
	name string
	dst1 Distance
	dst2 Distance
	out  bool
}{
	{"same integer nums", Distance(12), Distance(12), true},
	{"different integer nums", Distance(12), Distance(33), false},
	{"different integer nums reversed", Distance(33), Distance(12), false},
	{"same decimal nums", Distance(1.23456789), Distance(1.23456789), true},
	{"different decimal nums", Distance(1.23456789), Distance(1.2345678), false},
	{"different decimal nums reversed", Distance(1.2345678), Distance(1.23456789), false},
}

func TestDistance_Equal(t *testing.T) {
	for _, tt := range distanceEqualTests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.dst1.Equal(tt.dst2)
			if res != tt.out {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

func ExampleDistance_Equal() {
	num1 := Distance(17.2299478)
	num2 := Distance(184.909055172)
	num3 := Distance(184.909055172)

	fmt.Println(num1.Equal(num2))
	fmt.Println(num2.Equal(num3))
	// Output:
	// false
	// true
}
