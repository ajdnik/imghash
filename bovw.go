package imghash

import (
	"image"
	"math"
	"sort"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// BoVWFeatureType selects the local feature extractor used by BoVW.
type BoVWFeatureType uint8

const (
	// BoVWORB uses an ORB-like FAST + BRIEF pipeline.
	BoVWORB BoVWFeatureType = iota + 1
	// BoVWAKAZE uses an AKAZE-like Hessian detector with binary descriptors.
	BoVWAKAZE
)

func (f BoVWFeatureType) valid() bool {
	return f == BoVWORB || f == BoVWAKAZE
}

// BoVWStorageType selects the global BoVW representation format.
type BoVWStorageType uint8

const (
	// BoVWHistogram stores a normalized visual-word histogram.
	BoVWHistogram BoVWStorageType = iota + 1
	// BoVWMinHash stores a MinHash signature over visual words.
	BoVWMinHash
	// BoVWSimHash stores a SimHash bit-signature over visual words.
	BoVWSimHash
)

func (s BoVWStorageType) valid() bool {
	return s == BoVWHistogram || s == BoVWMinHash || s == BoVWSimHash
}

type bovwKeypoint struct {
	x, y     int
	response float64
	angle    float64
}

type bovwPair struct {
	x1, y1 int
	x2, y2 int
}

var (
	bovwORBPairs   = bovwGeneratePairs(256, 15, 0x6A09E667F3BCC909)
	bovwAKAZEPairs = bovwGeneratePairs(256, 19, 0xBB67AE8584CAA73B)
)

// BoVW implements bag-of-visual-words hashing using ORB-like or AKAZE-like
// local descriptors. The resulting representation can be a normalized
// histogram, MinHash signature, or SimHash bit-signature.
type BoVW struct {
	baseConfig
	featureType    BoVWFeatureType
	storageType    BoVWStorageType
	vocabularySize uint
	maxKeypoints   uint
	minHashSize    uint
	simHashBits    uint
	distFunc       DistanceFunc
}

// NewBoVW creates a new BoVW hasher with the given options.
// Without options, sensible defaults are used.
func NewBoVW(opts ...BoVWOption) (BoVW, error) {
	b := BoVW{
		baseConfig:     baseConfig{width: 256, height: 256, interp: Bilinear},
		featureType:    BoVWORB,
		storageType:    BoVWHistogram,
		vocabularySize: 256,
		maxKeypoints:   500,
		minHashSize:    64,
		simHashBits:    128,
	}
	for _, o := range opts {
		o.applyBoVW(&b)
	}
	if err := b.validate(); err != nil {
		return BoVW{}, err
	}
	if !b.featureType.valid() {
		return BoVW{}, ErrInvalidBoVWFeatureType
	}
	if !b.storageType.valid() {
		return BoVW{}, ErrInvalidBoVWStorageType
	}
	if b.vocabularySize == 0 {
		return BoVW{}, ErrInvalidVocabularySize
	}
	if b.maxKeypoints == 0 {
		return BoVW{}, ErrInvalidKeypoints
	}
	if b.minHashSize == 0 {
		return BoVW{}, ErrInvalidSignatureSize
	}
	if b.simHashBits == 0 {
		return BoVW{}, ErrInvalidSignatureSize
	}
	return b, nil
}

// Calculate returns a BoVW representation hash.
func (b BoVW) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(b.width, b.height, img, b.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	gray, w, h := bovwGrayToFlat(g)

	var keypoints []bovwKeypoint
	switch b.featureType {
	case BoVWORB:
		keypoints = bovwDetectORB(gray, w, h, int(b.maxKeypoints))
	case BoVWAKAZE:
		keypoints = bovwDetectAKAZE(g, w, h, int(b.maxKeypoints))
	default:
		return nil, ErrInvalidBoVWFeatureType
	}

	descriptors := make([][]byte, 0, len(keypoints))
	switch b.featureType {
	case BoVWORB:
		for i := range keypoints {
			descriptors = append(descriptors, bovwDescriptorORB(gray, w, h, keypoints[i]))
		}
	case BoVWAKAZE:
		blurred := imgproc.GaussianBlur(g, 0, 1.2)
		blurredGray, berr := imgproc.Grayscale(blurred)
		if berr != nil {
			return nil, berr
		}
		bg, _, _ := bovwGrayToFlat(blurredGray)
		for i := range keypoints {
			descriptors = append(descriptors, bovwDescriptorAKAZE(bg, w, h, keypoints[i]))
		}
	}

	wordHist := b.bovwHistogram(descriptors)

	switch b.storageType {
	case BoVWHistogram:
		return hashtype.Float64(bovwNormalizeL2(wordHist)), nil
	case BoVWMinHash:
		return hashtype.Float64(b.bovwMinHash(wordHist)), nil
	case BoVWSimHash:
		return b.bovwSimHash(wordHist), nil
	default:
		return nil, ErrInvalidBoVWStorageType
	}
}

// Compare computes distance using cosine for histogram storage and Jaccard
// for MinHash/SimHash storage.
func (b BoVW) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	switch b.storageType {
	case BoVWHistogram:
		if err := validateFloat64CompareInputs(h1, h2); err != nil {
			return 0, err
		}
		if b.distFunc != nil {
			return b.distFunc(h1, h2)
		}
		return similarity.Cosine(h1, h2)
	case BoVWMinHash:
		if err := validateFloat64CompareInputs(h1, h2); err != nil {
			return 0, err
		}
		if b.distFunc != nil {
			return b.distFunc(h1, h2)
		}
		return similarity.Jaccard(h1, h2)
	case BoVWSimHash:
		if err := validateBinaryCompareInputs(h1, h2); err != nil {
			return 0, err
		}
		if b.distFunc != nil {
			return b.distFunc(h1, h2)
		}
		return similarity.Jaccard(h1, h2)
	default:
		return 0, ErrInvalidBoVWStorageType
	}
}

func bovwDetectORB(gray []uint8, w, h, maxKeypoints int) []bovwKeypoint {
	if w < 7 || h < 7 || maxKeypoints <= 0 {
		return nil
	}
	const threshold = 20
	responses := make([]float64, w*h)

	for y := 3; y < h-3; y++ {
		for x := 3; x < w-3; x++ {
			score, ok := bovwFASTScore(gray, w, x, y, threshold)
			if ok {
				responses[y*w+x] = score
			}
		}
	}
	return bovwSelectKeypoints(gray, w, h, responses, maxKeypoints, 8)
}

func bovwDetectAKAZE(grayImg *image.Gray, w, h, maxKeypoints int) []bovwKeypoint {
	if w < 7 || h < 7 || maxKeypoints <= 0 {
		return nil
	}
	gray, _, _ := bovwGrayToFlat(grayImg)
	scales := []float64{1.0, 1.6, 2.2}
	scaleResponses := make([][]float64, len(scales))
	for i, sigma := range scales {
		b := imgproc.GaussianBlur(grayImg, 0, sigma)
		bg, err := imgproc.Grayscale(b)
		if err != nil {
			return nil
		}
		flat, _, _ := bovwGrayToFlat(bg)
		resp := make([]float64, w*h)
		for y := 1; y < h-1; y++ {
			for x := 1; x < w-1; x++ {
				c := float64(flat[y*w+x])
				dxx := float64(flat[y*w+x+1]) + float64(flat[y*w+x-1]) - 2*c
				dyy := float64(flat[(y+1)*w+x]) + float64(flat[(y-1)*w+x]) - 2*c
				dxy := (float64(flat[(y+1)*w+x+1]) - float64(flat[(y+1)*w+x-1]) -
					float64(flat[(y-1)*w+x+1]) + float64(flat[(y-1)*w+x-1])) * 0.25
				det := math.Abs(dxx*dyy-dxy*dxy) * sigma * sigma
				resp[y*w+x] = det
			}
		}
		scaleResponses[i] = resp
	}

	responses := make([]float64, w*h)
	for y := 3; y < h-3; y++ {
		for x := 3; x < w-3; x++ {
			best := 0.0
			for s := range scaleResponses {
				r := scaleResponses[s][y*w+x]
				if r == 0 {
					continue
				}
				if s > 0 && scaleResponses[s-1][y*w+x] > r {
					continue
				}
				if s+1 < len(scaleResponses) && scaleResponses[s+1][y*w+x] > r {
					continue
				}
				if r > best {
					best = r
				}
			}
			responses[y*w+x] = best
		}
	}
	return bovwSelectKeypoints(gray, w, h, responses, maxKeypoints, 10)
}

func bovwSelectKeypoints(gray []uint8, w, h int, responses []float64, maxKeypoints, radius int) []bovwKeypoint {
	points := make([]bovwKeypoint, 0, maxKeypoints)
	for y := radius; y < h-radius; y++ {
		for x := radius; x < w-radius; x++ {
			r := responses[y*w+x]
			if r <= 0 {
				continue
			}
			if !bovwIsLocalMaximum(responses, w, h, x, y) {
				continue
			}
			points = append(points, bovwKeypoint{x: x, y: y, response: r})
		}
	}

	sort.Slice(points, func(i, j int) bool {
		if points[i].response != points[j].response {
			return points[i].response > points[j].response
		}
		if points[i].y != points[j].y {
			return points[i].y < points[j].y
		}
		return points[i].x < points[j].x
	})
	if len(points) > maxKeypoints {
		points = points[:maxKeypoints]
	}
	for i := range points {
		points[i].angle = bovwOrientation(gray, w, h, points[i].x, points[i].y, radius)
	}
	return points
}

func bovwIsLocalMaximum(responses []float64, w, h, x, y int) bool {
	center := responses[y*w+x]
	for yy := y - 1; yy <= y+1; yy++ {
		if yy < 0 || yy >= h {
			continue
		}
		for xx := x - 1; xx <= x+1; xx++ {
			if xx < 0 || xx >= w || (xx == x && yy == y) {
				continue
			}
			if responses[yy*w+xx] >= center {
				return false
			}
		}
	}
	return true
}

func bovwOrientation(gray []uint8, w, h, x, y, radius int) float64 {
	var m10, m01 float64
	for yy := y - radius; yy <= y+radius; yy++ {
		sy := bovwReflect101(yy, h)
		for xx := x - radius; xx <= x+radius; xx++ {
			sx := bovwReflect101(xx, w)
			v := float64(gray[sy*w+sx])
			m10 += float64(xx-x) * v
			m01 += float64(yy-y) * v
		}
	}
	return math.Atan2(m01, m10)
}

func bovwFASTScore(gray []uint8, w, x, y, threshold int) (float64, bool) {
	center := int(gray[y*w+x])
	var bright, dark [16]bool
	score := 0
	for i, off := range bovwFASTCircle {
		v := int(gray[(y+off[1])*w+x+off[0]])
		if v >= center+threshold {
			bright[i] = true
		}
		if v <= center-threshold {
			dark[i] = true
		}
		score += absInt(v - center)
	}
	if bovwLongestCircularRun(bright) < 9 && bovwLongestCircularRun(dark) < 9 {
		return 0, false
	}
	return float64(score), true
}

func bovwLongestCircularRun(mask [16]bool) int {
	maxRun := 0
	run := 0
	for i := 0; i < 32; i++ {
		if mask[i%16] {
			run++
			if run > maxRun {
				maxRun = run
			}
			continue
		}
		run = 0
	}
	if maxRun > 16 {
		maxRun = 16
	}
	return maxRun
}

func bovwDescriptorORB(gray []uint8, w, h int, kp bovwKeypoint) []byte {
	return bovwBinaryDescriptor(gray, w, h, kp, bovwORBPairs, false)
}

func bovwDescriptorAKAZE(gray []uint8, w, h int, kp bovwKeypoint) []byte {
	return bovwBinaryDescriptor(gray, w, h, kp, bovwAKAZEPairs, true)
}

func bovwBinaryDescriptor(gray []uint8, w, h int, kp bovwKeypoint, pairs []bovwPair, usePatchMean bool) []byte {
	desc := make([]byte, 32)
	sinA := math.Sin(kp.angle)
	cosA := math.Cos(kp.angle)

	for i := range pairs {
		p := pairs[i]
		ax := kp.x + int(math.Round(cosA*float64(p.x1)-sinA*float64(p.y1)))
		ay := kp.y + int(math.Round(sinA*float64(p.x1)+cosA*float64(p.y1)))
		bx := kp.x + int(math.Round(cosA*float64(p.x2)-sinA*float64(p.y2)))
		by := kp.y + int(math.Round(sinA*float64(p.x2)+cosA*float64(p.y2)))

		var va, vb float64
		if usePatchMean {
			va = bovwPatchMean(gray, w, h, ax, ay, 1)
			vb = bovwPatchMean(gray, w, h, bx, by, 1)
		} else {
			va = float64(bovwSample(gray, w, h, ax, ay))
			vb = float64(bovwSample(gray, w, h, bx, by))
		}

		if va < vb {
			desc[i/8] |= 1 << uint(i%8)
		}
	}
	return desc
}

func (b BoVW) bovwHistogram(descriptors [][]byte) []float64 {
	hist := make([]float64, b.vocabularySize)
	for i := range descriptors {
		idx := bovwHashBytes(descriptors[i]) % uint32(b.vocabularySize)
		hist[idx]++
	}
	return hist
}

func (b BoVW) bovwMinHash(hist []float64) []float64 {
	words := make([]uint32, 0, len(hist))
	for i := range hist {
		if hist[i] > 0 {
			words = append(words, uint32(i))
		}
	}
	signature := make([]float64, b.minHashSize)
	if len(words) == 0 {
		return signature
	}

	const prime uint64 = 4294967291
	for i := uint(0); i < b.minHashSize; i++ {
		s1 := bovwMix64(0x9E3779B97F4A7C15 + uint64(i)*0xD1B54A32D192ED03)
		s2 := bovwMix64(0x94D049BB133111EB + uint64(i)*0xBF58476D1CE4E5B9)
		a := uint64(uint32(s1)|1) % prime
		c := uint64(uint32(s2)) % prime
		minV := prime - 1
		for j := range words {
			v := (a*uint64(words[j]) + c) % prime
			if v < minV {
				minV = v
			}
		}
		signature[i] = float64(minV)
	}
	return signature
}

func (b BoVW) bovwSimHash(hist []float64) hashtype.Binary {
	out := hashtype.NewBinary(b.simHashBits)
	acc := make([]float64, b.simHashBits)
	for i := range hist {
		weight := hist[i]
		if weight == 0 {
			continue
		}
		base := bovwMix64(uint64(i) + 0xA0761D6478BD642F)
		for bit := uint(0); bit < b.simHashBits; bit++ {
			h := bovwMix64(base + uint64(bit)*0xE7037ED1A0B428DB)
			if h&1 == 1 {
				acc[bit] += weight
			} else {
				acc[bit] -= weight
			}
		}
	}
	for bit := uint(0); bit < b.simHashBits; bit++ {
		if acc[bit] > 0 {
			out[bit/8] |= 1 << (bit % 8)
		}
	}
	return out
}

func bovwNormalizeL2(v []float64) []float64 {
	out := make([]float64, len(v))
	copy(out, v)
	var norm float64
	for i := range out {
		norm += out[i] * out[i]
	}
	if norm == 0 {
		return out
	}
	norm = math.Sqrt(norm)
	for i := range out {
		out[i] /= norm
	}
	return out
}

func bovwGrayToFlat(img *image.Gray) ([]uint8, int, int) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	flat := make([]uint8, w*h)
	ox, oy := bounds.Min.X, bounds.Min.Y
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			flat[y*w+x] = img.GrayAt(x+ox, y+oy).Y
		}
	}
	return flat, w, h
}

func bovwPatchMean(gray []uint8, w, h, x, y, radius int) float64 {
	var sum float64
	count := 0
	for yy := y - radius; yy <= y+radius; yy++ {
		sy := bovwReflect101(yy, h)
		for xx := x - radius; xx <= x+radius; xx++ {
			sx := bovwReflect101(xx, w)
			sum += float64(gray[sy*w+sx])
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

func bovwSample(gray []uint8, w, h, x, y int) uint8 {
	sx := bovwReflect101(x, w)
	sy := bovwReflect101(y, h)
	return gray[sy*w+sx]
}

func bovwReflect101(idx, size int) int {
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

func bovwHashBytes(data []byte) uint32 {
	var h uint32 = 2166136261
	for i := range data {
		h ^= uint32(data[i])
		h *= 16777619
	}
	return h
}

func bovwMix64(x uint64) uint64 {
	x += 0x9E3779B97F4A7C15
	x = (x ^ (x >> 30)) * 0xBF58476D1CE4E5B9
	x = (x ^ (x >> 27)) * 0x94D049BB133111EB
	return x ^ (x >> 31)
}

func bovwGeneratePairs(count, radius int, seed uint64) []bovwPair {
	pairs := make([]bovwPair, count)
	s := seed
	for i := 0; i < count; i++ {
		s = bovwMix64(s)
		x1 := int(s%uint64(2*radius+1)) - radius
		s = bovwMix64(s)
		y1 := int(s%uint64(2*radius+1)) - radius
		s = bovwMix64(s)
		x2 := int(s%uint64(2*radius+1)) - radius
		s = bovwMix64(s)
		y2 := int(s%uint64(2*radius+1)) - radius
		pairs[i] = bovwPair{x1: x1, y1: y1, x2: x2, y2: y2}
	}
	return pairs
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

var bovwFASTCircle = [16][2]int{
	{0, -3}, {1, -3}, {2, -2}, {3, -1},
	{3, 0}, {3, 1}, {2, 2}, {1, 3},
	{0, 3}, {-1, 3}, {-2, 2}, {-3, 1},
	{-3, 0}, {-3, -1}, {-2, -2}, {-1, -3},
}
