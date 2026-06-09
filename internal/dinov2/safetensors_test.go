package dinov2

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"testing"
)

// buildSafetensors assembles a minimal safetensors blob for tests.
func buildSafetensors(t *testing.T, entries map[string]struct {
	Shape []int
	Data  []float32
}) []byte {
	t.Helper()
	header := map[string]map[string]any{}
	offset := 0
	type rec struct {
		name string
		data []float32
		sh   []int
	}
	ordered := make([]rec, 0, len(entries))
	for name, e := range entries {
		ordered = append(ordered, rec{name: name, data: e.Data, sh: e.Shape})
	}
	for _, r := range ordered {
		n := len(r.data) * 4
		header[r.name] = map[string]any{
			"dtype":        "F32",
			"shape":        r.sh,
			"data_offsets": []int{offset, offset + n},
		}
		offset += n
	}
	hdrBytes, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	out := make([]byte, 8+len(hdrBytes)+offset)
	binary.LittleEndian.PutUint64(out[:8], uint64(len(hdrBytes)))
	copy(out[8:], hdrBytes)
	dataStart := 8 + len(hdrBytes)
	for _, r := range ordered {
		entry := header[r.name]
		off := entry["data_offsets"].([]int)
		dst := out[dataStart+off[0] : dataStart+off[1]]
		for i, v := range r.data {
			binary.LittleEndian.PutUint32(dst[i*4:i*4+4], math.Float32bits(v))
		}
	}
	return out
}

func TestParseSafetensors_RoundTrip(t *testing.T) {
	blob := buildSafetensors(t, map[string]struct {
		Shape []int
		Data  []float32
	}{
		"foo": {Shape: []int{2, 3}, Data: []float32{1, 2, 3, 4, 5, 6}},
		"bar": {Shape: []int{4}, Data: []float32{0.5, -1.5, 2.25, -3.75}},
	})
	got, err := ParseSafetensors(blob)
	if err != nil {
		t.Fatalf("ParseSafetensors: %v", err)
	}
	foo, ok := got["foo"]
	if !ok {
		t.Fatal("missing tensor foo")
	}
	if foo.Shape[0] != 2 || foo.Shape[1] != 3 {
		t.Errorf("foo shape = %v, want [2 3]", foo.Shape)
	}
	for i, v := range []float32{1, 2, 3, 4, 5, 6} {
		if foo.Data[i] != v {
			t.Errorf("foo.Data[%d] = %v, want %v", i, foo.Data[i], v)
		}
	}
	bar := got["bar"]
	for i, v := range []float32{0.5, -1.5, 2.25, -3.75} {
		if bar.Data[i] != v {
			t.Errorf("bar.Data[%d] = %v, want %v", i, bar.Data[i], v)
		}
	}
}

func TestParseSafetensors_RejectsNonF32(t *testing.T) {
	header := map[string]map[string]any{
		"t": {"dtype": "F16", "shape": []int{1}, "data_offsets": []int{0, 2}},
	}
	hb, _ := json.Marshal(header)
	blob := make([]byte, 8+len(hb)+2)
	binary.LittleEndian.PutUint64(blob[:8], uint64(len(hb)))
	copy(blob[8:], hb)
	if _, err := ParseSafetensors(blob); err == nil {
		t.Fatal("expected error for non-F32 dtype")
	}
}

func TestParseSafetensors_TooSmall(t *testing.T) {
	if _, err := ParseSafetensors([]byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for undersized blob")
	}
}
