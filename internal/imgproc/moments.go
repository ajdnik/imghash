package imgproc

import (
	"image"
	"math"
)

// Moments holds raw, central, and normalized central moments for a single image channel.
type Moments struct {
	m10  float64
	m00  float64
	m01  float64
	m11  float64
	m12  float64
	m21  float64
	m02  float64
	m20  float64
	m03  float64
	m30  float64
	mu30 float64
	mu03 float64
	mu12 float64
	mu21 float64
	mu02 float64
	mu20 float64
	mu11 float64
	nu30 float64
	nu03 float64
	nu12 float64
	nu21 float64
	nu02 float64
	nu20 float64
	nu11 float64
}

// GetMoments computes image moments for each channel of the given image.
func GetMoments(img image.Image) []Moments {
	switch i := img.(type) {
	case *image.Gray:
		return getMomentsGray(i)
	default:
		return getMomentsDefault(i)
	}
}

// HuMoments computes the seven Hu invariant moments for each channel's Moments.
func HuMoments(m []Moments) []float64 {
	hu := make([]float64, len(m)*7)
	var cnt int
	for _, mom := range m {
		t0 := mom.nu30 + mom.nu12
		t1 := mom.nu21 + mom.nu03
		q0 := t0 * t0
		q1 := t1 * t1
		n4 := 4 * mom.nu11
		t2 := t0 * (q0 - 3*q1)
		t3 := t1 * (3*q0 - q1)
		q2 := mom.nu30 - 3*mom.nu12
		q3 := 3*mom.nu21 - mom.nu03
		s := mom.nu20 + mom.nu02
		d := mom.nu20 - mom.nu02
		hu[cnt] = s
		hu[cnt+1] = d*d + n4*mom.nu11
		hu[cnt+2] = q2*q2 + q3*q3
		hu[cnt+3] = q0 + q1
		hu[cnt+4] = q2*t2 + q3*t3
		hu[cnt+5] = d*(q0-q1) + n4*t0*t1
		hu[cnt+6] = q3*t2 - q2*t3
		cnt += 7
	}
	return hu
}

func getMomentsDefault(img image.Image) []Moments {
	mom := make([]Moments, 3)
	mom[0] = Moments{}
	mom[1] = Moments{}
	mom[2] = Moments{}
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var x0R, x0G, x0B, x1R, x1G, x1B, x2R, x2G, x2B, x3R, x3G, x3B float64
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r /= 0x101
			g /= 0x101
			b /= 0x101
			xpR := float64(x) * float64(r)
			xpG := float64(x) * float64(g)
			xpB := float64(x) * float64(b)
			x0R += float64(r)
			x0G += float64(g)
			x0B += float64(b)
			x1R += xpR
			x1G += xpG
			x1B += xpB
			xxpR := xpR * float64(x)
			xxpG := xpG * float64(x)
			xxpB := xpB * float64(x)
			x2R += xxpR
			x2G += xxpG
			x2B += xxpB
			x3R += xxpR * float64(x)
			x3G += xxpG * float64(x)
			x3B += xxpB * float64(x)
		}
		pyR := float64(y) * x0R
		pyG := float64(y) * x0G
		pyB := float64(y) * x0B
		sy := float64(y) * float64(y)
		mom[0].m00 += x0R
		mom[0].m10 += x1R
		mom[0].m01 += pyR
		mom[0].m20 += x2R
		mom[0].m11 += x1R * float64(y)
		mom[0].m02 += x0R * sy
		mom[0].m30 += x3R
		mom[0].m21 += x2R * float64(y)
		mom[0].m12 += x1R * sy
		mom[0].m03 += pyR * sy
		mom[1].m00 += x0G
		mom[1].m10 += x1G
		mom[1].m01 += pyG
		mom[1].m20 += x2G
		mom[1].m11 += x1G * float64(y)
		mom[1].m02 += x0G * sy
		mom[1].m30 += x3G
		mom[1].m21 += x2G * float64(y)
		mom[1].m12 += x1G * sy
		mom[1].m03 += pyG * sy
		mom[2].m00 += x0B
		mom[2].m10 += x1B
		mom[2].m01 += pyB
		mom[2].m20 += x2B
		mom[2].m11 += x1B * float64(y)
		mom[2].m02 += x0B * sy
		mom[2].m30 += x3B
		mom[2].m21 += x2B * float64(y)
		mom[2].m12 += x1B * sy
		mom[2].m03 += pyB * sy
	}

	var cxR, cxG, cxB, cyR, cyG, cyB float64
	var mu20R, mu20G, mu20B, mu11R, mu11G, mu11B, mu02R, mu02G, mu02B float64
	var invM00R, invM00G, invM00B float64
	epsilon := math.Nextafter(1.0, 2.0) - 1.0
	if math.Abs(mom[0].m00) > epsilon {
		invM00R = 1 / mom[0].m00
		cxR = mom[0].m10 * invM00R
		cyR = mom[0].m01 * invM00R
	}
	if math.Abs(mom[1].m00) > epsilon {
		invM00G = 1 / mom[1].m00
		cxG = mom[1].m10 * invM00G
		cyG = mom[1].m01 * invM00G
	}
	if math.Abs(mom[2].m00) > epsilon {
		invM00B = 1 / mom[2].m00
		cxB = mom[2].m10 * invM00B
		cyB = mom[2].m01 * invM00B
	}
	mu20R = mom[0].m20 - mom[0].m10*cxR
	mu11R = mom[0].m11 - mom[0].m10*cyR
	mu02R = mom[0].m02 - mom[0].m01*cyR
	mu20G = mom[1].m20 - mom[1].m10*cxG
	mu11G = mom[1].m11 - mom[1].m10*cyG
	mu02G = mom[1].m02 - mom[1].m01*cyG
	mu20B = mom[2].m20 - mom[2].m10*cxB
	mu11B = mom[2].m11 - mom[2].m10*cyB
	mu02B = mom[2].m02 - mom[2].m01*cyB

	mom[0].mu20 = mu20R
	mom[0].mu11 = mu11R
	mom[0].mu02 = mu02R
	mom[0].mu30 = mom[0].m30 - cxR*(3*mu20R+cxR*mom[0].m10)
	mu11R += mu11R
	mom[0].mu21 = mom[0].m21 - cxR*(mu11R+cxR*mom[0].m01) - cyR*mu20R
	mom[0].mu12 = mom[0].m12 - cyR*(mu11R+cyR*mom[0].m10) - cxR*mu02R
	mom[0].mu03 = mom[0].m03 - cyR*(3*mu02R+cyR*mom[0].m01)

	mom[1].mu20 = mu20G
	mom[1].mu11 = mu11G
	mom[1].mu02 = mu02G
	mom[1].mu30 = mom[1].m30 - cxG*(3*mu20G+cxG*mom[1].m10)
	mu11G += mu11G
	mom[1].mu21 = mom[1].m21 - cxG*(mu11G+cxG*mom[1].m01) - cyG*mu20G
	mom[1].mu12 = mom[1].m12 - cyG*(mu11G+cyG*mom[1].m10) - cxG*mu02G
	mom[1].mu03 = mom[1].m03 - cyG*(3*mu02G+cyG*mom[1].m01)

	mom[2].mu20 = mu20B
	mom[2].mu11 = mu11B
	mom[2].mu02 = mu02B
	mom[2].mu30 = mom[2].m30 - cxB*(3*mu20B+cxB*mom[2].m10)
	mu11B += mu11B
	mom[2].mu21 = mom[2].m21 - cxB*(mu11B+cxB*mom[2].m01) - cyB*mu20B
	mom[2].mu12 = mom[2].m12 - cyB*(mu11B+cyB*mom[2].m10) - cxB*mu02B
	mom[2].mu03 = mom[2].m03 - cyB*(3*mu02B+cyB*mom[2].m01)

	invSqrtM00R := math.Sqrt(math.Abs(invM00R))
	s2R := invM00R * invM00R
	s3R := s2R * invSqrtM00R
	invSqrtM00G := math.Sqrt(math.Abs(invM00G))
	s2G := invM00G * invM00G
	s3G := s2G * invSqrtM00G
	invSqrtM00B := math.Sqrt(math.Abs(invM00B))
	s2B := invM00B * invM00B
	s3B := s2B * invSqrtM00B

	mom[0].nu20 = mom[0].mu20 * s2R
	mom[0].nu11 = mom[0].mu11 * s2R
	mom[0].nu02 = mom[0].mu02 * s2R
	mom[0].nu30 = mom[0].mu30 * s3R
	mom[0].nu21 = mom[0].mu21 * s3R
	mom[0].nu12 = mom[0].mu12 * s3R
	mom[0].nu03 = mom[0].mu03 * s3R
	mom[1].nu20 = mom[1].mu20 * s2G
	mom[1].nu11 = mom[1].mu11 * s2G
	mom[1].nu02 = mom[1].mu02 * s2G
	mom[1].nu30 = mom[1].mu30 * s3G
	mom[1].nu21 = mom[1].mu21 * s3G
	mom[1].nu12 = mom[1].mu12 * s3G
	mom[1].nu03 = mom[1].mu03 * s3G
	mom[2].nu20 = mom[2].mu20 * s2B
	mom[2].nu11 = mom[2].mu11 * s2B
	mom[2].nu02 = mom[2].mu02 * s2B
	mom[2].nu30 = mom[2].mu30 * s3B
	mom[2].nu21 = mom[2].mu21 * s3B
	mom[2].nu12 = mom[2].mu12 * s3B
	mom[2].nu03 = mom[2].mu03 * s3B
	return mom
}

func getMomentsGray(img *image.Gray) []Moments {
	mom := Moments{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var x0, x1, x2, x3 float64 = 0, 0, 0, 0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p := uint32(img.GrayAt(x, y).Y)
			xp := float64(x) * float64(p)
			x0 += float64(p)
			x1 += xp
			xxp := xp * float64(x)
			x2 += xxp
			x3 += xxp * float64(x)
		}
		py := float64(y) * x0
		sy := float64(y) * float64(y)
		mom.m00 += x0
		mom.m10 += x1
		mom.m01 += py
		mom.m20 += x2
		mom.m11 += x1 * float64(y)
		mom.m02 += x0 * sy
		mom.m30 += x3
		mom.m21 += x2 * float64(y)
		mom.m12 += x1 * sy
		mom.m03 += py * sy
	}

	var cx, cy float64 = 0, 0
	var mu20, mu11, mu02 float64
	var invM00 float64
	epsilon := math.Nextafter(1.0, 2.0) - 1.0
	if math.Abs(mom.m00) > epsilon {
		invM00 = 1 / mom.m00
		cx = mom.m10 * invM00
		cy = mom.m01 * invM00
	}
	mu20 = mom.m20 - mom.m10*cx
	mu11 = mom.m11 - mom.m10*cy
	mu02 = mom.m02 - mom.m01*cy

	mom.mu20 = mu20
	mom.mu11 = mu11
	mom.mu02 = mu02
	mom.mu30 = mom.m30 - cx*(3*mu20+cx*mom.m10)
	mu11 += mu11
	mom.mu21 = mom.m21 - cx*(mu11+cx*mom.m01) - cy*mu20
	mom.mu12 = mom.m12 - cy*(mu11+cy*mom.m10) - cx*mu02
	mom.mu03 = mom.m03 - cy*(3*mu02+cy*mom.m01)

	invSqrtM00 := math.Sqrt(math.Abs(invM00))
	s2 := invM00 * invM00
	s3 := s2 * invSqrtM00

	mom.nu20 = mom.mu20 * s2
	mom.nu11 = mom.mu11 * s2
	mom.nu02 = mom.mu02 * s2
	mom.nu30 = mom.mu30 * s3
	mom.nu21 = mom.mu21 * s3
	mom.nu12 = mom.mu12 * s3
	mom.nu03 = mom.mu03 * s3
	return []Moments{mom}
}
