package imghash_test

import (
	"strings"
	"testing"

	. "github.com/ajdnik/imghash"
)

func TestOpenImage_nonexistent(t *testing.T) {
	_, err := OpenImage("nonexistent.jpg")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestDecodeImage_invalid(t *testing.T) {
	_, err := DecodeImage(strings.NewReader("not an image"))
	if err == nil {
		t.Error("expected error for invalid image data")
	}
}

func TestHashFile_nonexistent(t *testing.T) {
	avg, err := NewAverage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = HashFile(avg, "nonexistent.jpg")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestHashReader_invalid(t *testing.T) {
	avg, err := NewAverage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = HashReader(avg, strings.NewReader("not an image"))
	if err == nil {
		t.Error("expected error for invalid reader data")
	}
}

func TestCompare_binary(t *testing.T) {
	h1 := Binary{0xFF, 0x00}
	h2 := Binary{0x00, 0xFF}
	dist, err := Compare(h1, h2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dist == 0 {
		t.Error("expected non-zero distance")
	}
}

func TestCompare_uint8(t *testing.T) {
	h1 := UInt8{0, 0}
	h2 := UInt8{3, 4}
	dist, err := Compare(h1, h2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dist.Equal(5) {
		t.Errorf("got %v, want 5", dist)
	}
}

func TestCompare_incompatible(t *testing.T) {
	h1 := Binary{1, 2}
	h2 := UInt8{3, 4}
	_, err := Compare(h1, h2)
	if err == nil {
		t.Error("expected error for incompatible hash types")
	}
}
