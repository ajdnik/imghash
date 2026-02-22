package imghash

import "github.com/ajdnik/imghash/v2/hashtype"

func validateBinaryCompareInputs(h1, h2 hashtype.Hash) error {
	b1, ok := h1.(hashtype.Binary)
	if !ok {
		return ErrIncompatibleHash
	}
	b2, ok := h2.(hashtype.Binary)
	if !ok {
		return ErrIncompatibleHash
	}
	if len(b1) != len(b2) {
		return ErrHashLengthMismatch
	}
	return nil
}

func validateUInt8CompareInputs(h1, h2 hashtype.Hash) error {
	u1, ok := h1.(hashtype.UInt8)
	if !ok {
		return ErrIncompatibleHash
	}
	u2, ok := h2.(hashtype.UInt8)
	if !ok {
		return ErrIncompatibleHash
	}
	if len(u1) != len(u2) {
		return ErrHashLengthMismatch
	}
	return nil
}

func validateFloat64CompareInputs(h1, h2 hashtype.Hash) error {
	f1, ok := h1.(hashtype.Float64)
	if !ok {
		return ErrIncompatibleHash
	}
	f2, ok := h2.(hashtype.Float64)
	if !ok {
		return ErrIncompatibleHash
	}
	if len(f1) != len(f2) {
		return ErrHashLengthMismatch
	}
	return nil
}
