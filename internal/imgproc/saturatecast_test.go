package imgproc

import "testing"

func TestSaturateCastF32ToUI8(t *testing.T) {
	tests := []struct {
		name   string
		input  float32
		expect uint8
	}{
		{"zero", 0, 0},
		{"mid", 127.5, 128},
		{"max", 255, 255},
		{"above max", 300, 255},
		{"below min", -10, 0},
		{"round up", 1.6, 2},
		{"round down", 1.4, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := saturateCastF32ToUI8(tt.input); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestSaturateCastIToUI8(t *testing.T) {
	tests := []struct {
		name   string
		input  int
		expect uint8
	}{
		{"zero", 0, 0},
		{"mid", 128, 128},
		{"max", 255, 255},
		{"above max", 1000, 255},
		{"below min", -5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := saturateCastIToUI8(tt.input); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}

func TestSaturateCastUI8ToF32(t *testing.T) {
	tests := []struct {
		name   string
		input  uint8
		expect float32
	}{
		{"zero", 0, 0},
		{"mid", 128, 128},
		{"max", 255, 255},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := saturateCastUI8ToF32(tt.input); res != tt.expect {
				t.Errorf("got %v, want %v", res, tt.expect)
			}
		})
	}
}
