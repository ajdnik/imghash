// Copyright 2020 Rok Ajdnik. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package hashtype

import (
	"fmt"
)

// UInt8 represents a hash type where the smallest hash element is a uint8 value.
type UInt8 []uint8

// String returns a string representation of uint8 hash.
// Is is formated as a uint8 slice.
func (h UInt8) String() string {
	return fmt.Sprintf("%v", []uint8(h))
}

// Len returns the number of elements in the hash.
func (h UInt8) Len() int {
	return len(h)
}

// ValueAt returns the element at the given index as a float64.
func (h UInt8) ValueAt(idx int) float64 {
	return float64(h[idx])
}

// Equal checks if two uint8 hashes are the same.
// It checks each same index value pair.
// Returns true if all elements match.
// If the length of hashes isn't the same the function returns false.
func (h UInt8) Equal(uh UInt8) bool {
	if len(h) != len(uh) {
		return false
	}
	for i, v := range h {
		if v != uh[i] {
			return false
		}
	}
	return true
}
