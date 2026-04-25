package imghash

import (
	"image"
	"image/color"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

const colorHashBins = 14

// ColorHash is compatible with Johannes Buchner's Python imagehash.colorhash.
// It encodes black, gray, faint-color hue, and bright-color hue fractions into
// 14 bins, using binBits bits per bin.
type ColorHash struct {
	binBits  uint
	distFunc DistanceFunc
}

// NewColorHash creates a new ColorHash with the given options.
// Without options, it uses binbits=3 to match Python imagehash.colorhash.
func NewColorHash(opts ...ColorHashOption) (ColorHash, error) {
	ch := ColorHash{binBits: 3}
	for _, o := range opts {
		o.applyColorHash(&ch)
	}
	if ch.binBits == 0 {
		return ColorHash{}, ErrInvalidNumBins
	}
	return ch, nil
}

// Calculate returns a ColorHash compatible with Python imagehash.colorhash.
func (ch ColorHash) Calculate(img image.Image) (hashtype.Hash, error) {
	if img == nil {
		return nil, imgproc.ErrImageIsNil
	}

	bounds := img.Bounds()
	total := bounds.Dx() * bounds.Dy()
	if total == 0 {
		return nil, ErrInvalidSize
	}

	var black, gray int
	var colorCount int
	faintCounts := make([]int, 6)
	brightCounts := make([]int, 6)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			r8, g8, b8 := c.R, c.G, c.B
			intensity := pilLuma(r8, g8, b8)
			h, s := pilHSV(r8, g8, b8)

			if intensity < 256/8 {
				black++
				continue
			}
			if s < 256/3 {
				gray++
				continue
			}

			colorCount++
			bin := hueBin(h)
			if s < 256*2/3 {
				faintCounts[bin]++
			} else if s > 256*2/3 {
				brightCounts[bin]++
			}
		}
	}

	denomColors := colorCount
	if denomColors < 1 {
		denomColors = 1
	}
	maxValue := 1 << ch.binBits
	values := make([]int, 0, colorHashBins)
	values = append(values, colorHashFractionValue(black, total, maxValue))
	values = append(values, colorHashFractionValue(gray, total, maxValue))
	for _, count := range faintCounts {
		values = append(values, colorHashFractionValue(count, denomColors, maxValue))
	}
	for _, count := range brightCounts {
		values = append(values, colorHashFractionValue(count, denomColors, maxValue))
	}

	hash := hashtype.NewBinary(uint(colorHashBins) * ch.binBits)
	var pos uint
	for _, value := range values {
		for i := uint(0); i < ch.binBits; i++ {
			divisor := 1 << (ch.binBits - i - 1)
			modulus := 1 << (ch.binBits - i)
			if (value/divisor)%modulus > 0 {
				if err := hash.SetReverse(pos); err != nil {
					return nil, err
				}
			}
			pos++
		}
	}
	return hash, nil
}

// Compare computes the Hamming distance between two ColorHash hashes.
func (ch ColorHash) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateBinaryCompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if ch.distFunc != nil {
		return ch.distFunc(h1, h2)
	}
	return similarity.Hamming(h1, h2)
}

func colorHashFractionValue(count, total, maxValue int) int {
	value := count * maxValue / total
	if value >= maxValue {
		return maxValue - 1
	}
	return value
}

func hueBin(h uint8) int {
	// numpy.histogram(..., bins=numpy.linspace(0, 255, 7)) puts values equal to
	// an interior edge into the bin on its right, while the final edge is closed.
	bin := int(math.Floor(float64(h) * 6 / 255))
	if bin > 5 {
		return 5
	}
	return bin
}

func pilLuma(r, g, b uint8) uint8 {
	return uint8((299*int(r) + 587*int(g) + 114*int(b) + 500) / 1000)
}

func pilHSV(r, g, b uint8) (uint8, uint8) {
	maxc := max(r, g, b)
	minc := min(r, g, b)
	if maxc == minc {
		return 0, 0
	}

	rc := float64(maxc-r) / float64(maxc-minc)
	gc := float64(maxc-g) / float64(maxc-minc)
	bc := float64(maxc-b) / float64(maxc-minc)

	var hue float64
	switch maxc {
	case r:
		hue = bc - gc
	case g:
		hue = 2 + rc - bc
	default:
		hue = 4 + gc - rc
	}
	hue = math.Mod(hue/6, 1)
	if hue < 0 {
		hue += 1
	}
	saturation := float64(maxc-minc) / float64(maxc)

	return uint8(hue * 255), uint8(saturation * 255)
}

