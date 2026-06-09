package imghash_test

import (
	"errors"
	"image"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
)

func TestNewDINOHash_NoOptions(t *testing.T) {
	d, err := imghash.NewDINOHash()
	if err != nil {
		t.Fatalf("NewDINOHash: %v", err)
	}
	if d == nil {
		t.Fatal("NewDINOHash returned nil")
	}
}

func TestDINOHashCalculate_NoWeightsConfigured(t *testing.T) {
	d, _ := imghash.NewDINOHash()
	_, err := d.Calculate(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	if !errors.Is(err, imghash.ErrNoWeights) {
		t.Fatalf("got %v, want ErrNoWeights", err)
	}
}

func TestDINOHashCalculate_WeightsProviderError(t *testing.T) {
	want := errors.New("boom")
	d, _ := imghash.NewDINOHash(imghash.WithDINOWeights(errWeightsProvider{err: want}))
	_, err := d.Calculate(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	if !errors.Is(err, want) {
		t.Fatalf("got %v, want wrapped %v", err, want)
	}
}

func TestDINOHashCompare_LengthMismatch(t *testing.T) {
	d, _ := imghash.NewDINOHash()
	if _, err := d.Compare(hashtype.Binary{0, 0}, hashtype.Binary{0, 0, 0}); !errors.Is(err, imghash.ErrHashLengthMismatch) {
		t.Errorf("got %v, want ErrHashLengthMismatch", err)
	}
}

func TestDINOHashCompare_WrongType(t *testing.T) {
	d, _ := imghash.NewDINOHash()
	if _, err := d.Compare(hashtype.UInt8{0}, hashtype.Binary{0}); !errors.Is(err, imghash.ErrIncompatibleHash) {
		t.Errorf("got %v, want ErrIncompatibleHash", err)
	}
}

func TestDINOHashCompare_HammingDefault(t *testing.T) {
	d, _ := imghash.NewDINOHash()
	a := hashtype.Binary{0xFF, 0x00, 0xAA, 0x55, 0xFF, 0x00, 0xAA, 0x55, 0xFF, 0x00, 0xAA, 0x55}
	b := hashtype.Binary{0x00, 0xFF, 0x55, 0xAA, 0x00, 0xFF, 0x55, 0xAA, 0x00, 0xFF, 0x55, 0xAA}
	dist, err := d.Compare(a, b)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if dist != 96 {
		t.Errorf("distance = %v, want 96", dist)
	}
}

type errWeightsProvider struct{ err error }

func (e errWeightsProvider) Tensors() (map[string]imghash.Tensor, error) {
	return nil, e.err
}
