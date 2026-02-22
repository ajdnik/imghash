package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// HOGHash is a perceptual hash based on Histogram of Oriented Gradients.
// It computes gradient magnitudes and orientations at each pixel, divides
// the image into cells, builds an orientation histogram per cell weighted
// by gradient magnitude, and concatenates them into a single hash vector.
//
// Based on Histograms of Oriented Gradients for Human Detection;
// Dalal and Triggs.
//
// See https://ieeexplore.ieee.org/document/1467360 for more information.
type HOGHash struct {
	baseConfig
	// Cell size in pixels (square cells).
	cellSize uint
	// Number of orientation bins (unsigned gradients, 0–180°).
	numBins  uint
	distFunc DistanceFunc
}

// NewHOGHash creates a new HOGHash with the given options.
// Without options, sensible defaults are used.
func NewHOGHash(opts ...HOGHashOption) (HOGHash, error) {
	h := HOGHash{
		baseConfig: baseConfig{width: 256, height: 256, interp: Bilinear},
		cellSize:   8,
		numBins:    9,
	}
	for _, o := range opts {
		o.applyHOGHash(&h)
	}
	if h.width == 0 || h.height == 0 {
		return HOGHash{}, ErrInvalidSize
	}
	if h.cellSize == 0 {
		return HOGHash{}, ErrInvalidCellSize
	}
	if h.numBins == 0 {
		return HOGHash{}, ErrInvalidNumBins
	}
	return h, nil
}

// Calculate returns a perceptual image hash.
func (hh HOGHash) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(hh.width, hh.height, img, hh.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	bounds := g.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	mag, orient := hh.computeGradients(g, w, h)
	return hh.computeHash(mag, orient, w, h), nil
}

// computeGradients computes gradient magnitude and unsigned orientation
// (0–180°) for each pixel using central differences.
func (hh HOGHash) computeGradients(img *image.Gray, w, h int) ([]float64, []float64) {
	ox := img.Bounds().Min.X
	oy := img.Bounds().Min.Y
	mag := make([]float64, w*h)
	orient := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var gx, gy float64
			if x > 0 && x < w-1 {
				gx = float64(img.GrayAt(x+1+ox, y+oy).Y) - float64(img.GrayAt(x-1+ox, y+oy).Y)
			} else if x == 0 {
				gx = float64(img.GrayAt(x+1+ox, y+oy).Y) - float64(img.GrayAt(x+ox, y+oy).Y)
			} else {
				gx = float64(img.GrayAt(x+ox, y+oy).Y) - float64(img.GrayAt(x-1+ox, y+oy).Y)
			}
			if y > 0 && y < h-1 {
				gy = float64(img.GrayAt(x+ox, y+1+oy).Y) - float64(img.GrayAt(x+ox, y-1+oy).Y)
			} else if y == 0 {
				gy = float64(img.GrayAt(x+ox, y+1+oy).Y) - float64(img.GrayAt(x+ox, y+oy).Y)
			} else {
				gy = float64(img.GrayAt(x+ox, y+oy).Y) - float64(img.GrayAt(x+ox, y-1+oy).Y)
			}
			mag[y*w+x] = math.Sqrt(gx*gx + gy*gy)
			angle := math.Atan2(gy, gx) * 180 / math.Pi
			if angle < 0 {
				angle += 180
			}
			if angle >= 180 {
				angle = 0
			}
			orient[y*w+x] = angle
		}
	}
	return mag, orient
}

// computeHash builds a UInt8 hash from the gradient data by computing
// magnitude-weighted orientation histograms for each cell.
func (hh HOGHash) computeHash(mag, orient []float64, w, h int) hashtype.UInt8 {
	cs := int(hh.cellSize)
	nb := int(hh.numBins)
	cellsX := w / cs
	cellsY := h / cs
	if cellsX == 0 {
		cellsX = 1
	}
	if cellsY == 0 {
		cellsY = 1
	}

	binWidth := 180.0 / float64(nb)
	hash := make(hashtype.UInt8, cellsX*cellsY*nb)

	for cy := 0; cy < cellsY; cy++ {
		for cx := 0; cx < cellsX; cx++ {
			startX := cx * cs
			startY := cy * cs
			endX := startX + cs
			endY := startY + cs
			if endX > w {
				endX = w
			}
			if endY > h {
				endY = h
			}

			hist := make([]float64, nb)
			for y := startY; y < endY; y++ {
				for x := startX; x < endX; x++ {
					idx := y*w + x
					bin := int(orient[idx] / binWidth)
					if bin >= nb {
						bin = nb - 1
					}
					hist[bin] += mag[idx]
				}
			}

			var maxVal float64
			for i := range hist {
				if hist[i] > maxVal {
					maxVal = hist[i]
				}
			}

			offset := (cy*cellsX + cx) * nb
			if maxVal > 0 {
				for i := 0; i < nb; i++ {
					hash[offset+i] = uint8(hist[i] * 255 / maxVal)
				}
			}
		}
	}
	return hash
}

// Compare computes the cosine distance between two HOGHash hashes.
func (hh HOGHash) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateUInt8CompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if hh.distFunc != nil {
		return hh.distFunc(h1, h2)
	}
	return similarity.Cosine(h1, h2)
}
