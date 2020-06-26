package hashtype

import (
	"fmt"
)

// UInt8 represents a hash type where the smallest hash element is a uint8 value.
type UInt8 []uint8

// String returns uint8 hash formated as a uint8 slice.
func (h UInt8) String() string {
	return fmt.Sprintf("%v", []uint8(h))
}

// Equal checks if two uint8 hashes are the same.
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
