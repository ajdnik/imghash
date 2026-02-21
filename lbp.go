package imghash

import (
	"image"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// LBP is a perceptual hash based on Local Binary Patterns.
// It computes a basic 3×3 LBP code for each pixel, divides the image into
// a grid of cells, builds a 256-bin histogram per cell, and concatenates
// them into a single hash vector.
//
// Based on Multiresolution Gray-Scale and Rotation Invariant Texture
// Classification with Local Binary Patterns; Ojala et. al.
//
// See https://ieeexplore.ieee.org/document/1017623 for more information.
type LBP struct {
	// Resized image width.
	width uint
	// Resized image height.
	height uint
	// Resize interpolation method.
	interp Interpolation
	// Number of horizontal grid cells.
	gridX uint
	// Number of vertical grid cells.
	gridY uint
}

// NewLBP creates a new LBP hash with the given options.
// Without options, sensible defaults are used.
func NewLBP(opts ...LBPOption) (LBP, error) {
	l := LBP{
		width:  256,
		height: 256,
		interp: Bilinear,
		gridX:  1,
		gridY:  1,
	}
	for _, o := range opts {
		o.applyLBP(&l)
	}
	if l.width == 0 || l.height == 0 {
		return LBP{}, ErrInvalidSize
	}
	if l.gridX == 0 || l.gridY == 0 {
		return LBP{}, ErrInvalidGridSize
	}
	return l, nil
}

// Calculate returns a perceptual image hash.
func (lh LBP) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(lh.width, lh.height, img, lh.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	bounds := g.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	lbpImg := lh.computeLBP(g, w, h)
	return lh.computeHash(lbpImg, w, h), nil
}

// 8-neighbor offsets starting from the right, going clockwise.
var lbpDX = [8]int{1, 1, 0, -1, -1, -1, 0, 1}
var lbpDY = [8]int{0, 1, 1, 1, 0, -1, -1, -1}

// computeLBP builds the LBP code image from a grayscale image.
// Border pixels compare out-of-bounds neighbors as zero.
func (lh LBP) computeLBP(img *image.Gray, w, h int) []uint8 {
	ox := img.Bounds().Min.X
	oy := img.Bounds().Min.Y
	result := make([]uint8, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			center := img.GrayAt(x+ox, y+oy).Y
			var code uint8
			for k := 0; k < 8; k++ {
				nx, ny := x+lbpDX[k], y+lbpDY[k]
				if nx >= 0 && nx < w && ny >= 0 && ny < h {
					if img.GrayAt(nx+ox, ny+oy).Y >= center {
						code |= 1 << uint(k)
					}
				}
			}
			result[y*w+x] = code
		}
	}
	return result
}

// computeHash builds a UInt8 hash from the LBP code image by computing
// a normalized 256-bin histogram for each grid cell.
func (lh LBP) computeHash(lbpImg []uint8, w, h int) hashtype.UInt8 {
	gx, gy := int(lh.gridX), int(lh.gridY)
	hash := make(hashtype.UInt8, gx*gy*256)
	cellW := w / gx
	cellH := h / gy

	for cy := 0; cy < gy; cy++ {
		for cx := 0; cx < gx; cx++ {
			startX := cx * cellW
			startY := cy * cellH
			endX := startX + cellW
			endY := startY + cellH
			if cx == gx-1 {
				endX = w
			}
			if cy == gy-1 {
				endY = h
			}

			var hist [256]float64
			var maxVal float64
			for y := startY; y < endY; y++ {
				for x := startX; x < endX; x++ {
					hist[lbpImg[y*w+x]]++
				}
			}
			for i := range hist {
				if hist[i] > maxVal {
					maxVal = hist[i]
				}
			}

			offset := (cy*gx + cx) * 256
			if maxVal > 0 {
				for i := 0; i < 256; i++ {
					hash[offset+i] = uint8(hist[i] * 255 / maxVal)
				}
			}
		}
	}
	return hash
}

// Compare computes the chi-square distance between two LBP hashes.
func (lh LBP) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	return similarity.ChiSquare(h1, h2), nil
}
