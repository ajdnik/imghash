package imgproc

import (
	"math"
	"testing"
)

func TestDCT(t *testing.T) {
	input := [][]float32{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	result := DCT(input)
	if len(result) != 4 || len(result[0]) != 4 {
		t.Fatalf("expected 4x4 result, got %dx%d", len(result), len(result[0]))
	}
	if math.Abs(float64(result[0][0])-34) > 0.01 {
		t.Errorf("DC coefficient should be ~34, got %v", result[0][0])
	}
}

func TestDCT_1x1(t *testing.T) {
	result := DCT([][]float32{{42}})
	if math.Abs(float64(result[0][0])-42) > 0.01 {
		t.Errorf("expected ~42, got %v", result[0][0])
	}
}

func TestDctOrthogonal(t *testing.T) {
	input := []float64{1, 2, 3, 4}
	result := dctOrthogonal(input)
	if len(result) != 4 {
		t.Fatalf("expected length 4, got %d", len(result))
	}
	expected0 := (1 + 2 + 3 + 4) * math.Sqrt(1.0/4.0)
	if math.Abs(result[0]-expected0) > 1e-10 {
		t.Errorf("DC component: got %v, want %v", result[0], expected0)
	}
}

func TestTranspose(t *testing.T) {
	input := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	result := transpose(input)
	if result[0][1] != 4 || result[1][0] != 2 {
		t.Errorf("transpose incorrect: [0][1]=%v (want 4), [1][0]=%v (want 2)", result[0][1], result[1][0])
	}
	if result[2][0] != 3 || result[0][2] != 7 {
		t.Errorf("transpose incorrect: [2][0]=%v (want 3), [0][2]=%v (want 7)", result[2][0], result[0][2])
	}
}

func TestMatf32Tof64(t *testing.T) {
	input := [][]float32{{1.5, 2.5}, {3.5, 4.5}}
	result := matf32Tof64(input)
	if result[0][0] != 1.5 || result[1][1] != 4.5 {
		t.Errorf("conversion incorrect")
	}
}

func TestMatf64Tof32(t *testing.T) {
	input := [][]float64{{1.5, 2.5}, {3.5, 4.5}}
	result := matf64Tof32(input)
	if result[0][0] != 1.5 || result[1][1] != 4.5 {
		t.Errorf("conversion incorrect")
	}
}

func TestHaarDWT2D_oneLevel(t *testing.T) {
	mat := [][]float32{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	HaarDWT2D(mat, 1)
	// LL quadrant (top-left 2x2): averages of 2x2 blocks
	wantLL := [][]float32{
		{3.5, 5.5},
		{11.5, 13.5},
	}
	for r := 0; r < 2; r++ {
		for c := 0; c < 2; c++ {
			if math.Abs(float64(mat[r][c]-wantLL[r][c])) > 1e-5 {
				t.Errorf("LL[%d][%d] = %v, want %v", r, c, mat[r][c], wantLL[r][c])
			}
		}
	}
}

func TestHaarDWT2D_twoLevels(t *testing.T) {
	mat := [][]float32{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	HaarDWT2D(mat, 2)
	// After 2 levels the top-left 1x1 is the global average
	globalAvg := float32(8.5)
	if math.Abs(float64(mat[0][0]-globalAvg)) > 1e-5 {
		t.Errorf("LL after 2 levels = %v, want %v", mat[0][0], globalAvg)
	}
}

func TestHaarDWT2D_empty(t *testing.T) {
	HaarDWT2D(nil, 1)
	HaarDWT2D([][]float32{}, 1)
}
