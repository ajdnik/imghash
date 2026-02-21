package imghash_test

import (
	"errors"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
)

func TestWithDistance_overrideAcrossHashers(t *testing.T) {
	want := imghash.Distance(123.456)
	calls := 0
	custom := func(_, _ hashtype.Hash) (imghash.Distance, error) {
		calls++
		return want, nil
	}

	cases := []struct {
		name  string
		build func(imghash.DistanceFunc) (imghash.Comparer, error)
	}{
		{"Average", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewAverage(imghash.WithDistance(fn))
		}},
		{"Difference", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewDifference(imghash.WithDistance(fn))
		}},
		{"Median", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewMedian(imghash.WithDistance(fn))
		}},
		{"PHash", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewPHash(imghash.WithDistance(fn))
		}},
		{"WHash", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewWHash(imghash.WithDistance(fn))
		}},
		{"BlockMean", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewBlockMean(imghash.WithDistance(fn))
		}},
		{"MarrHildreth", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewMarrHildreth(imghash.WithDistance(fn))
		}},
		{"RadialVariance", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewRadialVariance(imghash.WithDistance(fn))
		}},
		{"ColorMoment", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewColorMoment(imghash.WithDistance(fn))
		}},
		{"LBP", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewLBP(imghash.WithDistance(fn))
		}},
		{"HOGHash", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewHOGHash(imghash.WithDistance(fn))
		}},
		{"PDQ", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewPDQ(imghash.WithDistance(fn))
		}},
		{"RASH", func(fn imghash.DistanceFunc) (imghash.Comparer, error) {
			return imghash.NewRASH(imghash.WithDistance(fn))
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmp, err := tc.build(custom)
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}

			got, err := cmp.Compare(hashtype.UInt8{1, 2, 3}, hashtype.UInt8{4, 5, 6})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}

	if calls != len(cases) {
		t.Fatalf("custom distance called %d times, want %d", calls, len(cases))
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
