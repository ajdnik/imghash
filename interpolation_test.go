package imghash_test

import (
	"testing"

	. "github.com/ajdnik/imghash"
)

func TestInterpolation_String(t *testing.T) {
	tests := []struct {
		name   string
		interp Interpolation
		expect string
	}{
		{"NearestNeighbor", NearestNeighbor, "NearestNeighbor"},
		{"Bilinear", Bilinear, "Bilinear"},
		{"Bicubic", Bicubic, "Bicubic"},
		{"MitchellNetravali", MitchellNetravali, "MitchellNetravali"},
		{"Lanczos2", Lanczos2, "Lanczos2"},
		{"Lanczos3", Lanczos3, "Lanczos3"},
		{"BilinearExact", BilinearExact, "BilinearExact"},
		{"Unknown positive", Interpolation(99), "Unknown"},
		{"Unknown negative", Interpolation(-1), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.interp.String(); res != tt.expect {
				t.Errorf("got %q, want %q", res, tt.expect)
			}
		})
	}
}
