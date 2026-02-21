// Copyright 2020 Rok Ajdnik. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package hashtype

import (
	"fmt"
	"math"
)

// Float64 represents a hash type where the smallest hash element is a float64.
type Float64 []float64

// String returns a string representation of the float64 hash.
// Is is formated as a slice of float64 values.
func (h Float64) String() string {
	return fmt.Sprintf("%v", []float64(h))
}

// Len returns the number of elements in the hash.
func (h Float64) Len() int {
	return len(h)
}

// ValueAt returns the element at the given index.
func (h Float64) ValueAt(idx int) float64 {
	return h[idx]
}

// Equal checks if two float64 hashes are the same.
// It uses an epsilon based value comparrison.
// Returns true if each same index value pair
// is within an epsilon distance of each other.
// If the hashes aren't equal size the function returns false.
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
