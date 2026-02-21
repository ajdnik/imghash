package imghash_test

import (
	"testing"

	"github.com/ajdnik/imghash/v2"
)

func TestInterpolation_String(t *testing.T) {
	tests := []struct {
		name   string
		interp imghash.Interpolation
		expect string
	}{
		{"NearestNeighbor", imghash.NearestNeighbor, "NearestNeighbor"},
		{"Bilinear", imghash.Bilinear, "Bilinear"},
		{"Bicubic", imghash.Bicubic, "Bicubic"},
		{"MitchellNetravali", imghash.MitchellNetravali, "MitchellNetravali"},
		{"Lanczos2", imghash.Lanczos2, "Lanczos2"},
		{"Lanczos3", imghash.Lanczos3, "Lanczos3"},
		{"BilinearExact", imghash.BilinearExact, "BilinearExact"},
		{"Unknown positive", imghash.Interpolation(99), "Unknown"},
		{"Unknown negative", imghash.Interpolation(-1), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.interp.String(); res != tt.expect {
				t.Errorf("got %q, want %q", res, tt.expect)
			}
		})
	}
}
