package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

const (
	ehdGridSize      = 4
	ehdEdgeTypes     = 5
	ehdHashLen       = ehdGridSize * ehdGridSize * ehdEdgeTypes
	ehdEdgeThreshold = 11.0
)

// EHD quantisation tables from MPEG-7 Edge Histogram Descriptor.
// The descriptor uses 3-bit quantised local edge histogram bins.
var ehdQuantTable = [ehdEdgeTypes][8]float64{
	{0.010867, 0.057915, 0.099526, 0.144849, 0.195573, 0.260504, 0.358031, 0.530128}, // vertical
	{0.012266, 0.069934, 0.125879, 0.182307, 0.243396, 0.314563, 0.411728, 0.564319}, // horizontal
	{0.004193, 0.025852, 0.046860, 0.068519, 0.093286, 0.123490, 0.161505, 0.228960}, // 45 degree
	{0.004174, 0.025924, 0.046232, 0.067163, 0.089655, 0.115391, 0.151904, 0.217745}, // 135 degree
	{0.006778, 0.051667, 0.108650, 0.166257, 0.224226, 0.285691, 0.356375, 0.450972}, // non-directional
}

// EHD is an MPEG-7 Edge Histogram Descriptor style perceptual hash.
// The image is resized, converted to grayscale, split into a 4x4 grid, and
// each cell is described by a 5-bin histogram of local edge orientations.
// Each histogram bin is quantised to 3 bits, producing an 80-element hash.
type EHD struct {
	baseConfig
	distFunc DistanceFunc
}

// NewEHD creates a new EHD hash with the given options.
// Without options, sensible defaults are used.
func NewEHD(opts ...EHDOption) (EHD, error) {
	e := EHD{
		baseConfig: baseConfig{width: 256, height: 256, interp: Bilinear},
	}
	for _, o := range opts {
		o.applyEHD(&e)
	}
	if err := e.validate(); err != nil {
		return EHD{}, err
	}
	return e, nil
}

// Calculate returns an MPEG-7 EHD style perceptual hash.
func (e EHD) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(e.width, e.height, img, e.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	return e.computeHash(g), nil
}

func (e EHD) computeHash(img *image.Gray) hashtype.UInt8 {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	hash := make(hashtype.UInt8, ehdHashLen)
	ox, oy := b.Min.X, b.Min.Y

	for gy := 0; gy < ehdGridSize; gy++ {
		y0 := gy * h / ehdGridSize
		y1 := (gy + 1) * h / ehdGridSize
		for gx := 0; gx < ehdGridSize; gx++ {
			x0 := gx * w / ehdGridSize
			x1 := (gx + 1) * w / ehdGridSize

			var hist [ehdEdgeTypes]float64
			var blocks float64
			for y := y0; y+1 < y1; y += 2 {
				for x := x0; x+1 < x1; x += 2 {
					a := float64(img.GrayAt(ox+x, oy+y).Y)
					b := float64(img.GrayAt(ox+x+1, oy+y).Y)
					c := float64(img.GrayAt(ox+x, oy+y+1).Y)
					d := float64(img.GrayAt(ox+x+1, oy+y+1).Y)

					responses := [ehdEdgeTypes]float64{
						math.Abs(a - b + c - d),
						math.Abs(a + b - c - d),
						math.Abs(math.Sqrt2 * (a - d)),
						math.Abs(math.Sqrt2 * (b - c)),
						math.Abs(2 * (a - b - c + d)),
					}

					maxIdx := 0
					maxVal := responses[0]
					for i := 1; i < ehdEdgeTypes; i++ {
						if responses[i] > maxVal {
							maxVal = responses[i]
							maxIdx = i
						}
					}
					if maxVal >= ehdEdgeThreshold {
						hist[maxIdx]++
					}
					blocks++
				}
			}

			base := (gy*ehdGridSize + gx) * ehdEdgeTypes
			if blocks == 0 {
				continue
			}
			for i := 0; i < ehdEdgeTypes; i++ {
				norm := hist[i] / blocks
				hash[base+i] = quantizeEHD(norm, ehdQuantTable[i])
			}
		}
	}
	return hash
}

func quantizeEHD(v float64, table [8]float64) uint8 {
	best := 0
	minDiff := math.Abs(v - table[0])
	for i := 1; i < len(table); i++ {
		diff := math.Abs(v - table[i])
		if diff < minDiff {
			minDiff = diff
			best = i
		}
	}
	return uint8(best)
}

// Compare computes the L1 (Manhattan) distance between two EHD hashes.
func (e EHD) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateUInt8CompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if e.distFunc != nil {
		return e.distFunc(h1, h2)
	}
	return similarity.L1(h1, h2)
}
