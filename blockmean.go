package imghash

import (
	"image"
	"math"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/internal/imgproc"
	"github.com/ajdnik/imghash/v2/similarity"
)

// BlockMean is a perceptual hash that uses the method described in
// Block Mean Value Based Image Perceptual Hashing; Yang et. al.
//
// See https://ieeexplore.ieee.org/document/4041692 for more information.
type BlockMean struct {
	baseConfig
	// Block width.
	bWidth uint
	// Block height.
	bHeight  uint
	distFunc DistanceFunc
	// Block mean computation method.
	method BlockMeanMethod
}

// BlockMeanMethod represents the method used when computing the mean of blocks.
type BlockMeanMethod int

const (
	// Direct method constructs blocks with no overlap or rotation.
	Direct BlockMeanMethod = iota
	// Overlap method constructs blocks by overlapping them, the degree of overlap is set to be half of a block.
	Overlap
	// Rotation method uses the same approach as Direct but also hashes 24 rotated
	// image variants in 15-degree steps.
	Rotation
	// RotationOverlap uses the same approach as Overlap but also hashes 24 rotated
	// image variants in 15-degree steps.
	RotationOverlap
)

const (
	blockMeanRotationStepDegrees = 15
	blockMeanRotationCount       = 360 / blockMeanRotationStepDegrees
)

// NewBlockMean creates a new BlockMean hash with the given options.
// Without options, sensible defaults are used.
func NewBlockMean(opts ...BlockMeanOption) (BlockMean, error) {
	b := BlockMean{
		baseConfig: baseConfig{width: 256, height: 256, interp: BilinearExact},
		bWidth:     16,
		bHeight:    16,
		method:     Direct,
	}
	for _, o := range opts {
		o.applyBlockMean(&b)
	}
	if err := b.validate(); err != nil {
		return BlockMean{}, err
	}
	if b.bWidth == 0 || b.bHeight == 0 {
		return BlockMean{}, ErrInvalidBlockSize
	}
	return b, nil
}

// Calculate returns a perceptual image hash.
func (bh BlockMean) Calculate(img image.Image) (hashtype.Hash, error) {
	r := imgproc.Resize(bh.width, bh.height, img, bh.interp.resizeType())
	g, err := imgproc.Grayscale(r)
	if err != nil {
		return nil, err
	}
	if bh.method == Rotation || bh.method == RotationOverlap {
		return bh.computeRotatedHash(g)
	}
	mm := bh.computeMean(g)
	med, err := imgproc.Mean(g)
	if err != nil {
		return nil, err
	}
	return bh.computeHash(mm, med), nil
}

func (bh BlockMean) computeRotatedHash(img *image.Gray) (hashtype.Binary, error) {
	meansPerRotation := len(bh.computeMean(img))
	totalBits := meansPerRotation * blockMeanRotationCount
	hash := hashtype.NewBinary(uint(totalBits))
	bitOffset := 0
	for d := 0; d < 360; d += blockMeanRotationStepDegrees {
		var rotated *image.Gray
		if d == 0 {
			rotated = img
		} else {
			rotated = rotateGrayCrop(img, float64(d)*math.Pi/180)
		}
		means := bh.computeMean(rotated)
		med, err := imgproc.Mean(rotated)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(means); i++ {
			if means[i] >= med {
				_ = hash.Set(uint(bitOffset + i))
			}
		}
		bitOffset += len(means)
	}
	return hash, nil
}

// rotateGrayCrop rotates a grayscale image around its center and keeps the
// original image size by cropping pixels that rotate outside the bounds.
func rotateGrayCrop(img *image.Gray, angle float64) *image.Gray {
	bounds := img.Bounds()
	rotated := image.NewGray(bounds)
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	cx := float64(bounds.Dx()-1) / 2
	cy := float64(bounds.Dy()-1) / 2
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		fy := float64(y-bounds.Min.Y) - cy
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			fx := float64(x-bounds.Min.X) - cx
			srcX := int(math.Round(cosA*fx+sinA*fy+cx)) + bounds.Min.X
			srcY := int(math.Round(-sinA*fx+cosA*fy+cy)) + bounds.Min.Y
			if image.Pt(srcX, srcY).In(bounds) {
				rotated.SetGray(x, y, img.GrayAt(srcX, srcY))
			}
		}
	}
	return rotated
}

// Computes mean values of constructed blocks.
func (bh BlockMean) computeMean(img *image.Gray) []float64 {
	blksInX := int(bh.width / bh.bWidth)
	blksInY := int(bh.height / bh.bHeight)
	numB := blksInX * blksInY
	xS := int(bh.bWidth)
	yS := int(bh.bHeight)
	if bh.method == Overlap || bh.method == RotationOverlap {
		blksInX = int(bh.width/bh.bWidth)*2 - 1
		blksInY = int(bh.height/bh.bHeight)*2 - 1
		numB = blksInX * blksInY
		xS /= 2
		yS /= 2
	}
	means := make([]float64, numB)
	pixPerBlk := float64(bh.bWidth * bh.bHeight)
	blkCnt := 0
	for i := 0; i < blksInY; i++ {
		for j := 0; j < blksInX; j++ {
			var sum float64
			for x := 0; x < int(bh.bWidth); x++ {
				for y := 0; y < int(bh.bHeight); y++ {
					pix := img.GrayAt(j*xS+x, i*yS+y).Y
					sum += float64(pix)
				}
			}
			means[blkCnt] = sum / pixPerBlk
			blkCnt++
		}
	}
	return means
}

// Computes binary hash value based on block means.
func (bh BlockMean) computeHash(means []float64, median float64) hashtype.Binary {
	mSize := len(means)
	hSize := (mSize + 7) / 8
	hash := make(hashtype.Binary, hSize)
	for i := 0; i < mSize; i++ {
		if means[i] >= median {
			_ = hash.Set(uint(i))
		}
	}
	return hash
}

// Compare computes the Hamming distance between two BlockMean hashes.
func (bh BlockMean) Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	if err := validateBinaryCompareInputs(h1, h2); err != nil {
		return 0, err
	}
	if bh.distFunc != nil {
		return bh.distFunc(h1, h2)
	}
	return similarity.Hamming(h1, h2)
}
