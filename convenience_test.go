package imghash_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
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
	tests := []struct {
		name string
		h1   imghash.Hash
		h2   imghash.Hash
	}{
		{
			name: "binary then uint8",
			h1:   imghash.Binary{1, 2},
			h2:   imghash.UInt8{3, 4},
		},
		{
			name: "uint8 then binary",
			h1:   imghash.UInt8{3, 4},
			h2:   imghash.Binary{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := imghash.Compare(tt.h1, tt.h2)
			if !errors.Is(err, imghash.ErrIncompatibleHash) {
				t.Fatalf("got %v, want %v", err, imghash.ErrIncompatibleHash)
			}
		})
	}
}

func TestCompare_override(t *testing.T) {
	want := imghash.Distance(77.7)
	called := false
	custom := func(_, _ hashtype.Hash) (imghash.Distance, error) {
		called = true
		return want, nil
	}

	got, err := imghash.Compare(imghash.UInt8{1, 2}, imghash.UInt8{3, 4}, custom)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected custom distance function to be called")
	}
	if !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCompare_overrideError(t *testing.T) {
	wantErr := errors.New("custom compare failed")
	custom := func(_, _ hashtype.Hash) (imghash.Distance, error) {
		return 0, wantErr
	}

	_, err := imghash.Compare(imghash.UInt8{1, 2}, imghash.UInt8{3, 4}, custom)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}

func TestCompare_nilOverrideFallsBack(t *testing.T) {
	got, err := imghash.Compare(imghash.UInt8{0, 0}, imghash.UInt8{3, 4}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(5) {
		t.Fatalf("got %v, want 5", got)
	}
}
