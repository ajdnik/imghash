package vit

import (
	"math"
	"testing"
)

func TestPatchEmbed_SingleChannelIdentityWeight(t *testing.T) {
	// 4x4 input, 2x2 patches -> 2x2 grid = 4 patches. 1 channel, dim=4.
	// Weight is identity-like: each output channel selects one pixel of the patch.
	// Patch layout per row [p0, p1, p2, p3] = [(0,0), (0,1), (1,0), (1,1)].
	imageSize, patchSize, channels, dim := 4, 2, 1, 4
	input := []float32{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
	// weight [4, 4]: row i picks element i of the 4-wide patch vector.
	weight := []float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	bias := []float32{0, 0, 0, 0}
	got := patchEmbed(input, weight, bias, imageSize, patchSize, channels, dim)

	// Expected per patch: the four pixels of that patch (kc=0, ky, kx) in order.
	// Patch (0,0) covers rows 0-1, cols 0-1: [1, 2, 5, 6].
	// Patch (0,1) covers rows 0-1, cols 2-3: [3, 4, 7, 8].
	// Patch (1,0) covers rows 2-3, cols 0-1: [9, 10, 13, 14].
	// Patch (1,1) covers rows 2-3, cols 2-3: [11, 12, 15, 16].
	want := []float32{
		1, 2, 5, 6,
		3, 4, 7, 8,
		9, 10, 13, 14,
		11, 12, 15, 16,
	}
	for i, v := range want {
		if math.Abs(float64(got[i]-v)) > 1e-5 {
			t.Errorf("got[%d] = %v, want %v", i, got[i], v)
		}
	}
}

func TestPatchEmbed_BiasApplied(t *testing.T) {
	imageSize, patchSize, channels, dim := 2, 2, 1, 1
	input := []float32{0, 0, 0, 0}
	weight := []float32{1, 0, 0, 0} // sums first patch element
	bias := []float32{7}
	got := patchEmbed(input, weight, bias, imageSize, patchSize, channels, dim)
	if len(got) != 1 {
		t.Fatalf("want length 1, got %d", len(got))
	}
	if math.Abs(float64(got[0]-7)) > 1e-5 {
		t.Errorf("got %v, want 7 (zero input + bias)", got[0])
	}
}
