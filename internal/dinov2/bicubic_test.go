package dinov2

import (
	"math"
	"testing"
)

func TestCubicKernel(t *testing.T) {
	cases := []struct {
		t    float64
		want float64
	}{
		{0, 1},
		{1, 0},
		{2, 0},
		{0.5, (cubicA+2)*0.125 - (cubicA+3)*0.25 + 1},
		{1.5, cubicA*3.375 - 5*cubicA*2.25 + 8*cubicA*1.5 - 4*cubicA},
		{3, 0},
		{-0.5, (cubicA+2)*0.125 - (cubicA+3)*0.25 + 1},
	}
	for _, c := range cases {
		got := cubicKernel(c.t)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("cubicKernel(%v) = %v, want %v", c.t, got, c.want)
		}
	}
}

func TestInterpolatePosEmbed_LengthAndCLS(t *testing.T) {
	raw := make([]float32, PosEmbedNativeLen*384)
	for i := range raw {
		raw[i] = float32(i % 13)
	}
	out := InterpolatePosEmbed(raw)
	if len(out) != 257*384 {
		t.Fatalf("len(out) = %d, want %d", len(out), 257*384)
	}
	for i := 0; i < 384; i++ {
		if out[i] != raw[i] {
			t.Fatalf("CLS[%d] = %v, want %v", i, out[i], raw[i])
		}
	}
}

func TestInterpolatePosEmbed_ConstantInput(t *testing.T) {
	raw := make([]float32, PosEmbedNativeLen*384)
	for i := range raw {
		raw[i] = 2.5
	}
	out := InterpolatePosEmbed(raw)
	for i, v := range out {
		if math.Abs(float64(v)-2.5) > 1e-4 {
			t.Fatalf("out[%d] = %v, want 2.5", i, v)
		}
	}
}
