package imghash_test

import (
	"strings"
	"testing"

	"github.com/ajdnik/imghash/v2"
)

func TestOpenImage_nonexistent(t *testing.T) {
	_, err := imghash.OpenImage("nonexistent.jpg")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestDecodeImage_invalid(t *testing.T) {
	_, err := imghash.DecodeImage(strings.NewReader("not an image"))
	if err == nil {
		t.Error("expected error for invalid image data")
	}
}

func TestHashFile_nonexistent(t *testing.T) {
	avg, err := imghash.NewAverage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = imghash.HashFile(avg, "nonexistent.jpg")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestHashReader_invalid(t *testing.T) {
	avg, err := imghash.NewAverage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = imghash.HashReader(avg, strings.NewReader("not an image"))
	if err == nil {
		t.Error("expected error for invalid reader data")
	}
}

func TestCompare_binary(t *testing.T) {
	h1 := imghash.Binary{0xFF, 0x00}
	h2 := imghash.Binary{0x00, 0xFF}
	dist, err := imghash.Compare(h1, h2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dist == 0 {
		t.Error("expected non-zero distance")
	}
}

func TestCompare_uint8(t *testing.T) {
	h1 := imghash.UInt8{0, 0}
	h2 := imghash.UInt8{3, 4}
	dist, err := imghash.Compare(h1, h2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dist.Equal(5) {
		t.Errorf("got %v, want 5", dist)
	}
}

func TestCompare_incompatible(t *testing.T) {
	h1 := imghash.Binary{1, 2}
	h2 := imghash.UInt8{3, 4}
	_, err := imghash.Compare(h1, h2)
	if err == nil {
		t.Error("expected error for incompatible hash types")
	}
}
