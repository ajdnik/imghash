package dinov2

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
)

// stEntry mirrors a single tensor descriptor in the safetensors JSON header.
type stEntry struct {
	Dtype       string `json:"dtype"`
	Shape       []int  `json:"shape"`
	DataOffsets [2]int `json:"data_offsets"`
}

// ParseSafetensors decodes a safetensors fp32 blob into a name->Tensor map.
// Only the "F32" dtype is accepted; fp32 is preserved on purpose so the
// returned values are bit-identical to the source initializers.
func ParseSafetensors(blob []byte) (map[string]Tensor, error) {
	if len(blob) < 8 {
		return nil, fmt.Errorf("imghash: safetensors blob too small (%d bytes)", len(blob))
	}
	headerLen := binary.LittleEndian.Uint64(blob[:8])
	if uint64(len(blob)) < 8+headerLen {
		return nil, fmt.Errorf("imghash: safetensors header (%d bytes) overruns blob (%d bytes)", headerLen, len(blob))
	}
	headerBytes := blob[8 : 8+headerLen]
	var header map[string]json.RawMessage
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("imghash: safetensors header parse: %w", err)
	}
	dataStart := 8 + headerLen
	out := make(map[string]Tensor, len(header))
	for name, raw := range header {
		if name == "__metadata__" {
			continue
		}
		var e stEntry
		if err := json.Unmarshal(raw, &e); err != nil {
			return nil, fmt.Errorf("imghash: tensor %q parse: %w", name, err)
		}
		if e.Dtype != "F32" {
			return nil, fmt.Errorf("imghash: tensor %q dtype %q, want F32", name, e.Dtype)
		}
		if e.DataOffsets[1] < e.DataOffsets[0] {
			return nil, fmt.Errorf("imghash: tensor %q has inverted data offsets", name)
		}
		byteLen := e.DataOffsets[1] - e.DataOffsets[0]
		if byteLen%4 != 0 {
			return nil, fmt.Errorf("imghash: tensor %q byte length %d not a multiple of 4", name, byteLen)
		}
		start := dataStart + uint64(e.DataOffsets[0])
		end := dataStart + uint64(e.DataOffsets[1])
		if end > uint64(len(blob)) {
			return nil, fmt.Errorf("imghash: tensor %q (offset %d..%d) overruns blob (%d)", name, start, end, len(blob))
		}
		section := blob[start:end]
		nFloats := byteLen / 4
		data := make([]float32, nFloats)
		for i := 0; i < nFloats; i++ {
			data[i] = math.Float32frombits(binary.LittleEndian.Uint32(section[i*4 : i*4+4]))
		}
		out[name] = Tensor{
			Shape: append([]int(nil), e.Shape...),
			Data:  data,
		}
	}
	return out, nil
}
