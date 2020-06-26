package imghash

import (
	"image"
	"image/color"
	"sort"

	"gocv.io/x/gocv"
)

func float64sMedian(slice []float64) float64 {
	sorted := make([]float64, len(slice))
	copy(sorted, slice)
	sort.Float64s(sorted)
	size := len(sorted)
	if size%2 == 0 {
		return (sorted[size/2-1] + sorted[size/2+1]) / 2
	} else {
		return sorted[size/2]
	}
}

// func dct(mat [][]float32) [][]float32 {
//	initialize coefficients
// c := make([]float64, len(mat))
// c[0] = 1 / math.Sqrt(2)
// for i := 1; i < N; i++ {
// c[i] = 1
// }
// res := make([][]float32, len(mat))
// for i := range mat {
// res[i] = make([]float32, len(mat))
// }
//
//
// }

func dctCV(mat [][]float32) [][]float32 {
	inputMat := gocv.NewMatWithSize(len(mat), len(mat[0]), gocv.MatTypeCV32F)
	for i := range mat {
		for j := range mat[i] {
			inputMat.SetFloatAt(i, j, mat[i][j])
		}
	}
	dst := gocv.NewMat()
	gocv.DCT(inputMat, &dst, 0)
	res := make([][]float32, dst.Rows())
	for i := range res {
		res[i] = make([]float32, dst.Cols())
		for j := range res[i] {
			res[i][j] = dst.GetFloatAt(i, j)
		}
	}
	return res
}

func filter2DCV(img image.Image, kernel [][]float32) [][]float32 {
	var imgMat gocv.Mat
	switch i := img.(type) {
	case *image.Gray:
		imgMat, _ = gocv.ImageGrayToMatGray(i)
	default:
		imgMat, _ = gocv.ImageToMatRGB(i)
	}
	dst := gocv.NewMat()
	kernelMat := gocv.NewMatWithSize(len(kernel), len(kernel), gocv.MatTypeCV32F)
	for i := range kernel {
		for j := range kernel[i] {
			kernelMat.SetFloatAt(i, j, kernel[i][j])
		}
	}
	gocv.Filter2D(imgMat, &dst, gocv.MatTypeCV32F, kernelMat, image.Point{-1, -1}, 0, gocv.BorderDefault)
	res := make([][]float32, dst.Rows())
	for i := range res {
		res[i] = make([]float32, dst.Cols())
		for j := range res[i] {
			res[i][j] = dst.GetFloatAt(i, j)
		}
	}
	return res
}

func resizeImageCV(width, height uint, img image.Image, resizeType ResizeType) image.Image {
	var interp gocv.InterpolationFlags
	switch resizeType {
	case NearestNeighbor:
		interp = gocv.InterpolationNearestNeighbor
	case Bilinear:
		interp = gocv.InterpolationLinear
	case Bicubic:
		interp = gocv.InterpolationCubic
	case MitchellNetravali:
		interp = gocv.InterpolationNearestNeighbor
	case Lanczos2:
		interp = gocv.InterpolationNearestNeighbor
	case Lanczos3:
		interp = gocv.InterpolationNearestNeighbor
	case BilinearExact:
		interp = 5
	default:
		interp = gocv.InterpolationNearestNeighbor
	}
	var imgMat gocv.Mat
	switch i := img.(type) {
	case *image.Gray:
		imgMat, _ = gocv.ImageGrayToMatGray(i)
	default:
		imgMat, _ = gocv.ImageToMatRGB(i)
	}
	dst := gocv.NewMat()
	gocv.Resize(imgMat, &dst, image.Point{int(width), int(height)}, 0, 0, interp)
	res, err := dst.ToImage()
	if err != nil {
		panic(err)
	}
	return res
}

func GaussianBlurCV(img image.Image, kernel int, sigma float64) image.Image {
	var imgMat gocv.Mat
	switch i := img.(type) {
	case *image.Gray:
		imgMat, _ = gocv.ImageGrayToMatGray(i)
	default:
		imgMat, _ = gocv.ImageToMatRGB(i)
	}
	dst := gocv.NewMat()
	gocv.GaussianBlur(imgMat, &dst, image.Point{kernel, kernel}, sigma, 0, gocv.BorderDefault)
	res, err := dst.ToImage()
	if err != nil {
		panic(err)
	}
	return res
}

func ImageSum(img image.Image) uint32 {
	bounds := img.Bounds()
	var sum uint32 = 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			sum += r / 0x101
		}
	}
	return sum
}

func ImageGraySum(img *image.Gray) uint32 {
	bounds := img.Bounds()
	var sum uint32 = 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixel := img.GrayAt(x, y).Y
			sum += uint32(pixel)
		}
	}
	return sum
}

func ReadImageCV(file string) (image.Image, error) {
	mat := gocv.IMRead(file, gocv.IMReadColor)
	return mat.ToImage()
}

type ResizeType int

const (
	NearestNeighbor ResizeType = iota
	Bilinear
	Bicubic
	MitchellNetravali
	Lanczos2
	Lanczos3
	BilinearExact
)

func convVertical(img image.Image, kernel []float64) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	mid := (len(kernel) - 1) / 2
	rem := len(kernel) - 1 - mid
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			var sumR, sumG, sumB float64 = float64(r/0x101) * kernel[mid], float64(g/0x101) * kernel[mid], float64(b/0x101) * kernel[mid]
			for i := 1; i <= rem; i++ {
				iY := y - i
				if iY < bounds.Min.Y {
					iY = bounds.Min.Y + i
				}
				rl, gl, bl, _ := img.At(x, iY).RGBA()
				sumR += float64(rl/0x101) * kernel[mid-i]
				sumG += float64(gl/0x101) * kernel[mid-i]
				sumB += float64(bl/0x101) * kernel[mid-i]
				iY = y + i
				if iY >= bounds.Max.Y {
					iY = bounds.Max.Y - 1 - i
				}
				rh, gh, bh, _ := img.At(x, iY).RGBA()
				sumR += float64(rh/0x101) * kernel[mid+i]
				sumG += float64(gh/0x101) * kernel[mid+i]
				sumB += float64(bh/0x101) * kernel[mid+i]
			}
			nR := int(sumR)
			if nR > 255 {
				nR = 255
			} else if nR < 0 {
				nR = 0
			}
			nG := int(sumG)
			if nG > 255 {
				nG = 255
			} else if nG < 0 {
				nG = 0
			}
			nB := int(sumB)
			if nB > 255 {
				nB = 255
			} else if nB < 0 {
				nB = 0
			}
			dst.Set(x, y, color.RGBA{uint8(nR), uint8(nG), uint8(nB), uint8(a / 0x101)})
		}
	}
	return dst
}

func convHorizontal(img image.Image, kernel []float64) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	mid := (len(kernel) - 1) / 2
	rem := len(kernel) - 1 - mid
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			var sumR, sumG, sumB float64 = float64(r/0x101) * kernel[mid], float64(g/0x101) * kernel[mid], float64(b/0x101) * kernel[mid]
			for i := 1; i <= rem; i++ {
				iX := x - i
				if iX < bounds.Min.X {
					iX = bounds.Min.X + i
				}
				rl, gl, bl, _ := img.At(iX, y).RGBA()
				sumR += float64(rl/0x101) * kernel[mid-i]
				sumG += float64(gl/0x101) * kernel[mid-i]
				sumB += float64(bl/0x101) * kernel[mid-i]
				iX = x + i
				if iX >= bounds.Max.X {
					iX = bounds.Max.X - 1 - i
				}
				rh, gh, bh, _ := img.At(iX, y).RGBA()
				sumR += float64(rh/0x101) * kernel[mid+i]
				sumG += float64(gh/0x101) * kernel[mid+i]
				sumB += float64(bh/0x101) * kernel[mid+i]
			}
			dst.Set(x, y, color.RGBA{uint8(sumR), uint8(sumG), uint8(sumB), uint8(a / 0x101)})
		}
	}
	return dst
}
