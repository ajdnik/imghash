package imghash_test

import (
	"errors"
	"testing"

	"github.com/ajdnik/imghash/v2"
)

func TestWithInterpolation_invalidValue(t *testing.T) {
	invalid := imghash.Interpolation(99)
	tests := []struct {
		name string
		new  func() error
	}{
		{
			name: "Average",
			new: func() error {
				_, err := imghash.NewAverage(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "Difference",
			new: func() error {
				_, err := imghash.NewDifference(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "Median",
			new: func() error {
				_, err := imghash.NewMedian(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "PHash",
			new: func() error {
				_, err := imghash.NewPHash(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "BlockMean",
			new: func() error {
				_, err := imghash.NewBlockMean(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "MarrHildreth",
			new: func() error {
				_, err := imghash.NewMarrHildreth(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "ColorMoment",
			new: func() error {
				_, err := imghash.NewColorMoment(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "WHash",
			new: func() error {
				_, err := imghash.NewWHash(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "LBP",
			new: func() error {
				_, err := imghash.NewLBP(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "HOGHash",
			new: func() error {
				_, err := imghash.NewHOGHash(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "PDQ",
			new: func() error {
				_, err := imghash.NewPDQ(imghash.WithInterpolation(invalid))
				return err
			},
		},
		{
			name: "RASH",
			new: func() error {
				_, err := imghash.NewRASH(imghash.WithInterpolation(invalid))
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.new()
			if !errors.Is(err, imghash.ErrInvalidInterpolation) {
				t.Fatalf("expected ErrInvalidInterpolation, got %v", err)
			}
		})
	}
}
