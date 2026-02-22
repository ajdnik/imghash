package imghash_test

import (
	"errors"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
)

type compareCase struct {
	name              string
	buildDefault      func() (imghash.Comparer, error)
	buildWithDistance func(imghash.DistanceFunc) (imghash.Comparer, error)
	valid1            hashtype.Hash
	valid2            hashtype.Hash
	invalidType1      hashtype.Hash
	invalidType2      hashtype.Hash
	invalidLen1       hashtype.Hash
	invalidLen2       hashtype.Hash
}

var compareCases = []compareCase{
	{
		name:         "Average",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewAverage() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewAverage(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "Difference",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewDifference() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewDifference(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "Median",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewMedian() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewMedian(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "PHash",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewPHash() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewPHash(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "WHash",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewWHash() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewWHash(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "BlockMean",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewBlockMean() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewBlockMean(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "MarrHildreth",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewMarrHildreth() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewMarrHildreth(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "PDQ",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewPDQ() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewPDQ(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "RASH",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewRASH() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewRASH(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Binary{0x00, 0xFF},
		valid2:       hashtype.Binary{0xFF, 0x00},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Binary{0x00},
		invalidLen2:  hashtype.Binary{0x00, 0xFF},
	},
	{
		name:         "RadialVariance",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewRadialVariance() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewRadialVariance(imghash.WithDistance(fn))
		},
		valid1:       hashtype.UInt8{1, 2, 3},
		valid2:       hashtype.UInt8{3, 2, 1},
		invalidType1: hashtype.Binary{0x00, 0xFF},
		invalidType2: hashtype.Binary{0xFF, 0x00},
		invalidLen1:  hashtype.UInt8{1, 2},
		invalidLen2:  hashtype.UInt8{1, 2, 3},
	},
	{
		name:         "LBP",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewLBP() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewLBP(imghash.WithDistance(fn))
		},
		valid1:       hashtype.UInt8{1, 2, 3},
		valid2:       hashtype.UInt8{3, 2, 1},
		invalidType1: hashtype.Binary{0x00, 0xFF},
		invalidType2: hashtype.Binary{0xFF, 0x00},
		invalidLen1:  hashtype.UInt8{1, 2},
		invalidLen2:  hashtype.UInt8{1, 2, 3},
	},
	{
		name:         "HOGHash",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewHOGHash() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewHOGHash(imghash.WithDistance(fn))
		},
		valid1:       hashtype.UInt8{1, 2, 3},
		valid2:       hashtype.UInt8{3, 2, 1},
		invalidType1: hashtype.Binary{0x00, 0xFF},
		invalidType2: hashtype.Binary{0xFF, 0x00},
		invalidLen1:  hashtype.UInt8{1, 2},
		invalidLen2:  hashtype.UInt8{1, 2, 3},
	},
	{
		name:         "ColorMoment",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewColorMoment() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewColorMoment(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Float64{1.0, 2.0},
		valid2:       hashtype.Float64{3.0, 4.0},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Float64{1.0, 2.0},
		invalidLen2:  hashtype.Float64{1.0, 2.0, 3.0},
	},
	{
		name:         "Zernike",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewZernike() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewZernike(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Float64{1.0, 2.0},
		valid2:       hashtype.Float64{3.0, 4.0},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Float64{1.0, 2.0},
		invalidLen2:  hashtype.Float64{1.0, 2.0, 3.0},
	},
	{
		name:         "GIST",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewGIST() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewGIST(imghash.WithDistance(fn))
		},
		valid1:       hashtype.Float64{1.0, 2.0},
		valid2:       hashtype.Float64{3.0, 4.0},
		invalidType1: hashtype.UInt8{1, 2},
		invalidType2: hashtype.UInt8{3, 4},
		invalidLen1:  hashtype.Float64{1.0, 2.0},
		invalidLen2:  hashtype.Float64{1.0, 2.0, 3.0},
	},
	{
		name:         "CLD",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewCLD() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewCLD(imghash.WithDistance(fn))
		},
		valid1:       hashtype.UInt8{1, 2, 3},
		valid2:       hashtype.UInt8{3, 2, 1},
		invalidType1: hashtype.Binary{0x00, 0xFF},
		invalidType2: hashtype.Binary{0xFF, 0x00},
		invalidLen1:  hashtype.UInt8{1, 2},
		invalidLen2:  hashtype.UInt8{1, 2, 3},
	},
	{
		name:         "EHD",
		buildDefault: func() (imghash.Comparer, error) { return imghash.NewEHD() },
		buildWithDistance: func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewEHD(imghash.WithDistance(fn))
		},
		valid1:       hashtype.UInt8{1, 2, 3},
		valid2:       hashtype.UInt8{3, 2, 1},
		invalidType1: hashtype.Binary{0x00, 0xFF},
		invalidType2: hashtype.Binary{0xFF, 0x00},
		invalidLen1:  hashtype.UInt8{1, 2},
		invalidLen2:  hashtype.UInt8{1, 2, 3},
	},
}

func TestWithDistance_overrideAcrossHashers(t *testing.T) {
	want := imghash.Distance(123.456)
	calls := 0
	custom := func(_, _ hashtype.Hash) (imghash.Distance, error) {
		calls++
		return want, nil
	}

	for _, tc := range compareCases {
		t.Run(tc.name, func(t *testing.T) {
			cmp, err := tc.buildWithDistance(custom)
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}

			got, err := cmp.Compare(tc.valid1, tc.valid2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}

	if calls != len(compareCases) {
		t.Fatalf("custom distance called %d times, want %d", calls, len(compareCases))
	}
}

func TestWithDistance_validatesInputsBeforeOverride(t *testing.T) {
	for _, tc := range compareCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			custom := func(_, _ hashtype.Hash) (imghash.Distance, error) {
				calls++
				return 1, nil
			}
			cmp, err := tc.buildWithDistance(custom)
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}

			_, err = cmp.Compare(tc.invalidType1, tc.invalidType2)
			if !errors.Is(err, imghash.ErrIncompatibleHash) {
				t.Fatalf("type mismatch: got %v, want %v", err, imghash.ErrIncompatibleHash)
			}

			_, err = cmp.Compare(tc.invalidLen1, tc.invalidLen2)
			if !errors.Is(err, imghash.ErrHashLengthMismatch) {
				t.Fatalf("length mismatch: got %v, want %v", err, imghash.ErrHashLengthMismatch)
			}

			if calls != 0 {
				t.Fatalf("custom distance called %d times, want 0", calls)
			}
		})
	}
}

func TestCompare_defaultValidatesInputsAcrossHashers(t *testing.T) {
	for _, tc := range compareCases {
		t.Run(tc.name, func(t *testing.T) {
			cmp, err := tc.buildDefault()
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}

			_, err = cmp.Compare(tc.invalidType1, tc.invalidType2)
			if !errors.Is(err, imghash.ErrIncompatibleHash) {
				t.Fatalf("type mismatch: got %v, want %v", err, imghash.ErrIncompatibleHash)
			}

			_, err = cmp.Compare(tc.invalidLen1, tc.invalidLen2)
			if !errors.Is(err, imghash.ErrHashLengthMismatch) {
				t.Fatalf("length mismatch: got %v, want %v", err, imghash.ErrHashLengthMismatch)
			}
		})
	}
}

func TestWithDistance_errorPropagation(t *testing.T) {
	wantErr := errors.New("distance failure")
	custom := func(_, _ hashtype.Hash) (imghash.Distance, error) {
		return 0, wantErr
	}

	avg, err := imghash.NewAverage(imghash.WithDistance(custom))
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	_, err = avg.Compare(hashtype.Binary{0x00}, hashtype.Binary{0xFF})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got %v, want %v", err, wantErr)
	}
}

func TestWithWeights_affectsPHashCompare(t *testing.T) {
	ph, err := imghash.NewPHash(imghash.WithWeights([]float64{2}))
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	got, err := ph.Compare(hashtype.Binary{1}, hashtype.Binary{2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(4) {
		t.Fatalf("got %v, want 4", got)
	}
}

func TestWithWeights_clonesCallerSlice(t *testing.T) {
	weights := []float64{2}
	opt := imghash.WithWeights(weights)
	weights[0] = 100

	ph, err := imghash.NewPHash(opt)
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	got, err := ph.Compare(hashtype.Binary{1}, hashtype.Binary{2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(4) {
		t.Fatalf("got %v, want 4", got)
	}
}
