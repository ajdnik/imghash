package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

const zernikeEpsilon = 1e-12

type zernikeOrder struct {
	n int
	m int
}

// Zernike is a magnitude-based perceptual hash built from Zernike moments.
//
// The descriptor stores |A(n,m)| values for all valid orders up to a maximum
// degree (excluding A(0,0)). Magnitude makes the hash rotation-invariant.
//
// Based on A. Khotanzad and Y. H. Hong, "Invariant Image Recognition by Zernike
// Moments" (1990).
type Zernike struct {
	baseConfig
	// Maximum order n used when building Zernike moments.
	degree   int
	distFunc DistanceFunc
}

// NewZernike creates a new Zernike hash with the given options.
// Without options, sensible defaults are used.
func NewZernike(opts ...ZernikeOption) (Zernike, error) {
	z := Zernike{
		baseConfig: baseConfig{width: 64, height: 64, interp: Bilinear},
		degree:     8,
	}
	for _, o := range opts {
		o.applyZernike(&z)
	}
	if err := z.validate(); err != nil {
		return Zernike{}, err
	}
	if z.degree <= 0 {
		return Zernike{}, ErrInvalidDegree
	}
	return z, nil
}

// Calculate returns a perceptual image hash.
func (z Zernike) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(z.width, z.height, img, z.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	return z.computeHash(g), nil
}

func zernikeOrders(maxDegree int) []zernikeOrder {
	orders := make([]zernikeOrder, 0, (maxDegree+1)*(maxDegree+2)/4)
	for n := 0; n <= maxDegree; n++ {
		for m := n % 2; m <= n; m += 2 {
			orders = append(orders, zernikeOrder{n: n, m: m})
		}
	}
	return orders
}

func zernikeLogFactorials(limit int) []float64 {
	logFact := make([]float64, limit+1)
	for i := 1; i <= limit; i++ {
		logFact[i] = logFact[i-1] + math.Log(float64(i))
	}
	return logFact
}

func zernikeRadial(n, m int, rho float64, logFact []float64) float64 {
	if n < m || (n-m)%2 != 0 {
		return 0
	}
	halfPlus := (n + m) / 2
	halfMinus := (n - m) / 2
	maxS := halfMinus

	var sum float64
	for s := 0; s <= maxS; s++ {
		coeff := math.Exp(logFact[n-s] - logFact[s] - logFact[halfPlus-s] - logFact[halfMinus-s])
		term := coeff * math.Pow(rho, float64(n-2*s))
		if s%2 == 0 {
			sum += term
		} else {
			sum -= term
		}
	}
	return sum
}

func (z Zernike) computeHash(img *image.Gray) hashtype.Float64 {
	orders := zernikeOrders(z.degree)
	if len(orders) <= 1 {
		return hashtype.Float64{}
	}
	logFact := zernikeLogFactorials(z.degree)
	sumRe := make([]float64, len(orders))
	sumIm := make([]float64, len(orders))

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w == 0 || h == 0 {
		return make(hashtype.Float64, len(orders)-1)
	}
	scale := w
	if h < scale {
		scale = h
	}
	if scale == 0 {
		return make(hashtype.Float64, len(orders)-1)
	}
	radius := float64(scale) / 2
	cx := float64(w-1) / 2
	cy := float64(h-1) / 2
	ox := bounds.Min.X
	oy := bounds.Min.Y

	for y := 0; y < h; y++ {
		ny := (float64(y) - cy) / radius
		for x := 0; x < w; x++ {
			nx := (float64(x) - cx) / radius
			rho2 := nx*nx + ny*ny
			if rho2 > 1 {
				continue
			}
			intensity := float64(img.GrayAt(x+ox, y+oy).Y) / 255
			if intensity <= 0 {
				continue
			}
			rho := math.Sqrt(rho2)
			theta := math.Atan2(ny, nx)
			for i, ord := range orders {
				radial := zernikeRadial(ord.n, ord.m, rho, logFact)
				if radial == 0 {
					continue
				}
				angle := float64(ord.m) * theta
				c := math.Cos(angle)
				s := math.Sin(angle)
				sumRe[i] += intensity * radial * c
				sumIm[i] -= intensity * radial * s
			}
		}
	}

	descriptor := make(hashtype.Float64, len(orders)-1)
	dc := math.Hypot(sumRe[0], sumIm[0]) / math.Pi
	if dc <= zernikeEpsilon {
		return descriptor
	}
	for i := 1; i < len(orders); i++ {
		mag := math.Hypot(sumRe[i], sumIm[i]) * (float64(orders[i].n) + 1) / math.Pi
		descriptor[i-1] = mag / dc
	}
	return descriptor
}

// Compare computes the L2 (Euclidean) distance between two Zernike hashes.
func (z Zernike) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateFloat64CompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if z.distFunc != nil {
		return z.distFunc(h1, h2)
	}
	return similarity.L2(h1, h2)
}
