package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

var gistOrientationsPerScale = []int{8, 8, 4}

// GIST is a holistic scene descriptor based on a bank of oriented Gabor
// filters pooled over a coarse spatial grid.
//
// Based on: A. Oliva and A. Torralba, "Modeling the Shape of the Scene:
// A Holistic Representation of the Spatial Envelope" (2001).
type GIST struct {
	baseConfig
	gridX, gridY uint
	distFunc     DistanceFunc
}

// NewGIST creates a new GIST hash with the given options.
// Without options, sensible defaults are used.
func NewGIST(opts ...GISTOption) (GIST, error) {
	g := GIST{
		baseConfig: baseConfig{width: 64, height: 64, interp: Bilinear},
		gridX:      4,
		gridY:      4,
	}
	for _, o := range opts {
		o.applyGIST(&g)
	}
	if err := g.validate(); err != nil {
		return GIST{}, err
	}
	if g.gridX == 0 || g.gridY == 0 {
		return GIST{}, ErrInvalidGridSize
	}
	return g, nil
}

// Calculate returns a perceptual image hash.
func (g GIST) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(g.width, g.height, img, g.interp.resizeType())
	gray, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	mat := g.normalizeGray(gray)
	desc := g.computeDescriptor(mat)
	return hashtype.Float64(desc), nil
}

func (g GIST) normalizeGray(img *image.Gray) [][]float64 {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	ox, oy := bounds.Min.X, bounds.Min.Y
	mat := make([][]float64, h)
	var mean float64
	for y := 0; y < h; y++ {
		mat[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			v := float64(img.GrayAt(x+ox, y+oy).Y) / 255
			mat[y][x] = v
			mean += v
		}
	}

	n := float64(w * h)
	if n == 0 {
		return mat
	}
	mean /= n

	var variance float64
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d := mat[y][x] - mean
			variance += d * d
		}
	}
	std := math.Sqrt(variance / n)
	if std < 1e-12 {
		std = 1
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			mat[y][x] = (mat[y][x] - mean) / std
		}
	}
	return mat
}

func (g GIST) computeDescriptor(img [][]float64) []float64 {
	h := len(img)
	if h == 0 {
		return []float64{}
	}
	w := len(img[0])
	if w == 0 {
		return []float64{}
	}

	gridX, gridY := int(g.gridX), int(g.gridY)
	cellCounts := make([]int, gridX*gridY)
	for y := 0; y < h; y++ {
		by := y * gridY / h
		for x := 0; x < w; x++ {
			bx := x * gridX / w
			cellCounts[by*gridX+bx]++
		}
	}

	totalOrientations := 0
	for _, o := range gistOrientationsPerScale {
		totalOrientations += o
	}
	descriptor := make([]float64, 0, gridX*gridY*totalOrientations)

	wavelengths := []float64{4, 8, 12}
	for s, nOrient := range gistOrientationsPerScale {
		wavelength := wavelengths[s]
		for o := 0; o < nOrient; o++ {
			theta := float64(o) * math.Pi / float64(nOrient)
			kReal, kImag := gistGaborKernel(theta, wavelength)
			cellSums := gistFilterPool(img, kReal, kImag, gridX, gridY)
			for i := range cellSums {
				if cellCounts[i] > 0 {
					descriptor = append(descriptor, cellSums[i]/float64(cellCounts[i]))
				} else {
					descriptor = append(descriptor, 0)
				}
			}
		}
	}

	var norm float64
	for _, v := range descriptor {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range descriptor {
			descriptor[i] /= norm
		}
	}

	return descriptor
}

func gistGaborKernel(theta, wavelength float64) ([][]float64, [][]float64) {
	sigma := 0.56 * wavelength
	gamma := 0.5
	radius := int(math.Ceil(2 * sigma))
	if radius < 3 {
		radius = 3
	}
	if radius > 9 {
		radius = 9
	}
	size := 2*radius + 1

	kernelReal := make([][]float64, size)
	kernelImag := make([][]float64, size)
	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)

	var norm float64
	for y := -radius; y <= radius; y++ {
		yy := y + radius
		kernelReal[yy] = make([]float64, size)
		kernelImag[yy] = make([]float64, size)
		for x := -radius; x <= radius; x++ {
			xx := x + radius
			xTheta := float64(x)*cosTheta + float64(y)*sinTheta
			yTheta := -float64(x)*sinTheta + float64(y)*cosTheta
			gauss := math.Exp(-(xTheta*xTheta + gamma*gamma*yTheta*yTheta) / (2 * sigma * sigma))
			phase := 2 * math.Pi * xTheta / wavelength
			r := gauss * math.Cos(phase)
			i := gauss * math.Sin(phase)
			kernelReal[yy][xx] = r
			kernelImag[yy][xx] = i
			norm += r*r + i*i
		}
	}

	norm = math.Sqrt(norm)
	if norm > 0 {
		for y := range kernelReal {
			for x := range kernelReal[y] {
				kernelReal[y][x] /= norm
				kernelImag[y][x] /= norm
			}
		}
	}

	return kernelReal, kernelImag
}

func gistFilterPool(img, kReal, kImag [][]float64, gridX, gridY int) []float64 {
	h := len(img)
	w := len(img[0])
	kh := len(kReal)
	kw := len(kReal[0])
	ry := kh / 2
	rx := kw / 2

	cellSums := make([]float64, gridX*gridY)
	for y := 0; y < h; y++ {
		by := y * gridY / h
		for x := 0; x < w; x++ {
			var sumR, sumI float64
			for ky := 0; ky < kh; ky++ {
				iy := gistReflect101(y+ky-ry, h)
				for kx := 0; kx < kw; kx++ {
					ix := gistReflect101(x+kx-rx, w)
					p := img[iy][ix]
					sumR += p * kReal[ky][kx]
					sumI += p * kImag[ky][kx]
				}
			}
			bx := x * gridX / w
			cellSums[by*gridX+bx] += math.Hypot(sumR, sumI)
		}
	}

	return cellSums
}

func gistReflect101(idx, size int) int {
	if size <= 1 {
		return 0
	}
	for idx < 0 || idx >= size {
		if idx < 0 {
			idx = -idx
			continue
		}
		idx = 2*size - idx - 2
	}
	return idx
}

// Compare computes the cosine distance between two GIST hashes.
func (g GIST) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateFloat64CompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if g.distFunc != nil {
		return g.distFunc(h1, h2)
	}
	return similarity.Cosine(h1, h2)
}
