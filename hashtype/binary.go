// Copyright 2020 Rok Ajdnik. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package hashtype implements data types used to represent hashes.
// It can be used by hashing algorithm implementations to
// represent the algorithm's results.
package hashtype

import (
	"bytes"
	"errors"
	"fmt"
)

// Binary represents a hash type where the smallest hash element is a bit.
type Binary []byte

// String returns a string representation of the binary hash.
// Is is formated as an array of bytes.
func (h Binary) String() string {
	return fmt.Sprintf("%v", []byte(h))
}

// ErrOutOfBounds is reported when the bit position is larger than the number of bits in the hash.
var ErrOutOfBounds = errors.New("position out of bounds")

// ErrIncompatibleHash is reported when Distance is called with an incompatible hash type.
var ErrIncompatibleHash = errors.New("incompatible hash type for distance calculation")

// Distance returns the Hamming distance (number of differing bits) to another Binary hash.
// Returns ErrIncompatibleHash if other is not a Binary hash.
func (h Binary) Distance(other Hash) (float64, error) {
	b, ok := other.(Binary)
	if !ok {
		return 0, ErrIncompatibleHash
	}
	l := len(h)
	if len(b) < l {
		l = len(b)
	}
	var dist int
	for i := 0; i < l; i++ {
		xor := h[i] ^ b[i]
		for xor != 0 {
			xor &= xor - 1
			dist++
		}
	}
	return float64(dist), nil
}

// Set turns a bit on in the binary hash.
// The position argument determines which bit should be turned on.
// Returns error if position is out of bounds.
func (h Binary) Set(position uint) error {
	bit := position % 8
	return h.setBit(position, bit)
}

// SetReverse turns a bit on in the binary hash.
// The position argument determines which bit should be
// turned on, where position is counted in reverse order.
// Returns error if position is out of bounds.
func (h Binary) SetReverse(position uint) error {
	bit := 7 - position%8
	return h.setBit(position, bit)
}

// Len returns the number of bytes in the binary hash.
func (h Binary) Len() int {
	return len(h)
}

// ValueAt returns the byte at the given index as a float64.
func (h Binary) ValueAt(idx int) float64 {
	return float64(h[idx])
}

// Equal checks if two binary hashes are the same.
// Returns true if the hashes match.
func (h Binary) Equal(bh Binary) bool {
	return bytes.Equal(h, bh)
}

// setBit sets a bit to 1 in the binary hash.
// If the bit position is larger than the hash size
// it returns an ErrOutOfBounds.
func (h Binary) setBit(position, bit uint) error {
	byt := position / 8
	if byt >= uint(len(h)) {
		return ErrOutOfBounds
	}
	h[byt] |= 1 << bit
	return nil
}
