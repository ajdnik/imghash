package dinov2

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"slices"
	"strings"
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

func TestParseSafetensors_InvalidHeadersAndOffsets(t *testing.T) {
	for _, tc := range []struct {
		name, header, want string
		dataBytes          int
	}{
		{"invalid_json", `{`, "header parse", 0},
		{"invalid_entry", `{"weight":false}`, `tensor "weight" parse`, 0},
		{"inverted_offsets", `{"weight":{"dtype":"F32","data_offsets":[4,0]}}`, "inverted data offsets", 4},
		{"unaligned_length", `{"weight":{"dtype":"F32","data_offsets":[0,3]}}`, "not a multiple of 4", 3},
		{"data_overrun", `{"weight":{"dtype":"F32","data_offsets":[0,8]}}`, "overruns blob", 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			blob := make([]byte, 8+len(tc.header)+tc.dataBytes)
			binary.LittleEndian.PutUint64(blob[:8], uint64(len(tc.header)))
			copy(blob[8:], tc.header)
			got, err := ParseSafetensors(blob)
			if got != nil || err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("ParseSafetensors = %v, %v; want error containing %q", got, err, tc.want)
			}
		})
	}
	t.Run("header_overrun", func(t *testing.T) {
		blob := make([]byte, 8)
		binary.LittleEndian.PutUint64(blob, 100)
		if _, err := ParseSafetensors(blob); err == nil || !strings.Contains(err.Error(), "header (100 bytes) overruns blob") {
			t.Fatalf("error = %v, want header overrun", err)
		}
	})
}

func TestParseSafetensors_Metadata(t *testing.T) {
	header := `{"__metadata__":{"format":"pt"},"weight":{"dtype":"F32","shape":[1],"data_offsets":[0,4]}}`
	blob := make([]byte, 8+len(header)+4)
	binary.LittleEndian.PutUint64(blob[:8], uint64(len(header)))
	copy(blob[8:], header)
	binary.LittleEndian.PutUint32(blob[8+len(header):], math.Float32bits(-1.25))
	got, err := ParseSafetensors(blob)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !slices.Equal(got["weight"].Shape, []int{1}) || !slices.Equal(got["weight"].Data, []float32{-1.25}) {
		t.Fatalf("parsed tensors = %v, want only weight with shape [1] and value -1.25", got)
	}
}
