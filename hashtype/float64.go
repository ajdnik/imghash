package hashtype

import (
	"fmt"
	"math"
)

// Float64 represents a hash type where the smallest hash element is a float64.
type Float64 []float64

// String returns float64 hash formated as a slice of float64 values.
func (h Float64) String() string {
	return fmt.Sprintf("%v", []float64(h))
}

// Equal checks if two float64 hashes are the same by using epsilon based value comparrison.
func (h Float64) Equal(fh Float64) bool {
	if len(h) != len(fh) {
		return false
	}
	eps := math.Nextafter(1.0, 2.0) - 1.0
	for i := 0; i < len(fh); i++ {
		if math.Abs(h[i]-fh[i]) > eps {
			return false
		}
	}
	return true
}
