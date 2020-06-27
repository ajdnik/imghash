package imgproc

import (
	"fmt"
	"math"
	"testing"
)

var getGaussianKernelTests = []struct {
	kernel int
	sigma  float64
	out    []float32
}{
	{7, 0, []float32{0.03125, 0.109375, 0.21875, 0.28125, 0.21875, 0.109375, 0.03125}},
	{7, 1, []float32{0.004433048, 0.054005582, 0.24203624, 0.39905027, 0.24203624, 0.054005582, 0.004433048}},
	{6, 1, []float32{0.017559513, 0.12974823, 0.35269225, 0.35269225, 0.12974823, 0.017559513}},
	{5, 2, []float32{0.15246914, 0.22184129, 0.25137913, 0.22184129, 0.15246914}},
	{7, 2, []float32{0.07015932, 0.13107488, 0.19071282, 0.21610594, 0.19071282, 0.13107488, 0.07015932}},
	{3, 1, []float32{0.27406862, 0.45186275, 0.27406862}},
	{9, 1, []float32{1.3383062e-04, 4.4318615e-03, 5.3991128e-02, 2.4197145e-01, 3.9894345e-01, 2.4197145e-01, 5.3991128e-02, 4.4318615e-03, 1.3383062e-04}},
	{1, 0.5, []float32{1}},
	{3, 0.8, []float32{0.23899426, 0.52201146, 0.23899426}},
}

func equal(s1, s2 []float32) bool {
	if len(s1) != len(s2) {
		return false
	}
	eps := math.Nextafter(1.0, 2.0) - 1.0
	for i := 0; i < len(s1); i++ {
		if math.Abs(float64(s1[i])-float64(s2[i])) > eps {
			return false
		}
	}
	return true
}

func TestGetGaussianKernel(t *testing.T) {
	for _, tt := range getGaussianKernelTests {
		t.Run(fmt.Sprintf("kernel=%v;sigma=%v", tt.kernel, tt.sigma), func(t *testing.T) {
			res := getGaussianKernel(tt.kernel, tt.sigma)
			if !equal(res, tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}

}
