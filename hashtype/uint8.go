// Copyright 2020 Rok Ajdnik. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package hashtype implements data types used to represent hashes.
// It is used by various hashing algorithm implementations to
// represent the algorithm's results.
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
