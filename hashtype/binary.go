package hashtype

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

// Binary represents a hash type where the smallest hash element is a bit.
type Binary []byte

// String returns binary hash formated as a hexadecimal string with zero padding.
func (h Binary) String() string {
	return fmt.Sprintf("%02X", []byte(h))
}

// ErrOutOfBounds is reported when the bit position is out of bounds.
var ErrOutOfBounds = errors.New("position out of bounds")

// Set turns a bit on or off at position.
// Returns error if position is out of bounds.
func (h Binary) Set(position uint) error {
	bit := position % 8
	return h.setBit(position, bit)
}

// SetReverse turns a bit on or off at position, where position is counted in reverse order.
// Returns error if position is out of bounds.
func (h Binary) SetReverse(position uint) error {
	bit := 7 - position%8
	return h.setBit(position, bit)
}

// Equal checks if two binary hashes are the same.
func (h Binary) Equal(bh Binary) bool {
	return bytes.Equal(h, bh)
}

func (h Binary) setBit(position, bit uint) error {
	byt := position / 8
	if byt >= uint(len(h)) {
		return ErrOutOfBounds
	}
	h[byt] |= byte(math.Pow(2, float64(bit)))
	return nil
}
