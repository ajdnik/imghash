package similarity_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/v2/similarity"
)

var distanceEqualTests = []struct {
	name string
	dst1 similarity.Distance
	dst2 similarity.Distance
	out  bool
}{
	{"same integer nums", similarity.Distance(12), similarity.Distance(12), true},
	{"different integer nums", similarity.Distance(12), similarity.Distance(33), false},
	{"different integer nums reversed", similarity.Distance(33), similarity.Distance(12), false},
	{"same decimal nums", similarity.Distance(1.23456789), similarity.Distance(1.23456789), true},
	{"different decimal nums", similarity.Distance(1.23456789), similarity.Distance(1.2345678), false},
	{"different decimal nums reversed", similarity.Distance(1.2345678), similarity.Distance(1.23456789), false},
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
	num1 := similarity.Distance(17.2299478)
	num2 := similarity.Distance(184.909055172)
	num3 := similarity.Distance(184.909055172)

	fmt.Println(num1.Equal(num2))
	fmt.Println(num2.Equal(num3))
	// Output:
	// false
	// true
}
