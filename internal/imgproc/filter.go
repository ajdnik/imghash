package imgproc

import (
	"image"
	"image/color"
)

// Filter2DGray applies a 2-D convolution kernel to a grayscale image.
func Filter2DGray(img *image.Gray, kernel [][]float32) [][]float32 {
	bounds := img.Bounds()
	width, height := getSize(img)
	ret := make([][]float32, height)
	for i := 0; i < height; i++ {
		ret[i] = make([]float32, width)
	}
	midY, midX := (len(kernel)-1)/2, (len(kernel[0])-1)/2
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sum float64
			for yo, ky := y-midY, 0; yo <= y+midY; yo, ky = yo+1, ky+1 {
				for xo, kx := x-midX, 0; xo <= x+midX; xo, kx = xo+1, kx+1 {
					iX, iY := borderReflect101(xo, yo, bounds)
					sum += float64(kernel[ky][kx]) * float64(img.GrayAt(iX, iY).Y)
				}
			}
			ret[y][x] = float32(sum)
		}
	}
	return ret
}

func borderReflect101(x, y int, bounds image.Rectangle) (int, int) {
	rX, rY := x, y
	if rX < bounds.Min.X {
		rX = bounds.Min.X + (bounds.Min.X - x)
	}
	if rX >= bounds.Max.X {
		rX = bounds.Max.X - 1 - (x - bounds.Max.X + 1)
	}
	if rY < bounds.Min.Y {
		rY = bounds.Min.Y + (bounds.Min.Y - y)
	}
	if rY >= bounds.Max.Y {
		rY = bounds.Max.Y - 1 - (y - bounds.Max.Y + 1)
	}
	return rX, rY
}

func sepFilter2DGray(img *image.Gray, kernel []int) image.Image {
	bounds := img.Bounds()
	width, height := getSize(img)
	dst := image.NewGray(image.Rect(0, 0, width, height))
	buff := make([]int, width*height)
	mid := (len(kernel) - 1) / 2
	rem := len(kernel) - 1 - mid
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			g := img.GrayAt(x, y).Y
			sG := int(g) * kernel[mid]
			for i := 1; i <= rem; i++ {
				iX := x - i
				if iX < bounds.Min.X {
					iX = bounds.Min.X + i
				}
				gP := img.GrayAt(iX, y).Y
				iX = x + i
				if iX >= bounds.Max.X {
					iX = bounds.Max.X - 1 - i
				}
				gN := img.GrayAt(iX, y).Y
				sG += int(gP)*kernel[mid-i] + int(gN)*kernel[mid+i]
			}
			buff[y*width+x] = sG
		}
	}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			g := buff[y*width+x]
			sG := g*kernel[mid] + 65536
			for i := 1; i <= rem; i++ {
				iY := y - i
				if iY < bounds.Min.Y {
					iY = bounds.Min.Y + i
				}
				gP := buff[iY*width+x]
				iY = y + i
				if iY >= bounds.Max.Y {
					iY = bounds.Max.Y - 1 - i
				}
				gN := buff[iY*width+x]
				sG += gP*kernel[mid-i] + gN*kernel[mid+i]
			}
			sG = (sG + 1<<15) >> 16
			if sG > 255 {
				sG = 255
			} else if sG < 0 {
				sG = 0
			} else {
				sG--
			}
			dst.SetGray(x, y, color.Gray{uint8(sG)})
		}
	}
	return dst
}

func sepFilter2D(img image.Image, kernel []int) image.Image {
	bounds := img.Bounds()
	width, height := getSize(img)
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	buffR := make([]int, width*height)
	buffG := make([]int, width*height)
	buffB := make([]int, width*height)
	mid := (len(kernel) - 1) / 2
	rem := len(kernel) - 1 - mid
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			sR := int(r/0x101) * kernel[mid]
			sG := int(g/0x101) * kernel[mid]
			sB := int(b/0x101) * kernel[mid]
			for i := 1; i <= rem; i++ {
				iX := x - i
				if iX < bounds.Min.X {
					iX = bounds.Min.X + i
				}
				rP, gP, bP, _ := img.At(iX, y).RGBA()
				iX = x + i
				if iX >= bounds.Max.X {
					iX = bounds.Max.X - 1 - i
				}
				rN, gN, bN, _ := img.At(iX, y).RGBA()
				sR += int(rP/0x101)*kernel[mid-i] + int(rN/0x101)*kernel[mid+i]
				sG += int(gP/0x101)*kernel[mid-i] + int(gN/0x101)*kernel[mid+i]
				sB += int(bP/0x101)*kernel[mid-i] + int(bN/0x101)*kernel[mid+i]
			}
			buffR[y*width+x] = sR
			buffG[y*width+x] = sG
			buffB[y*width+x] = sB
		}
	}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b := buffR[y*width+x], buffG[y*width+x], buffB[y*width+x]
			sR := r*kernel[mid] + 65536
			sG := g*kernel[mid] + 65536
			sB := b*kernel[mid] + 65536
			for i := 1; i <= rem; i++ {
				iY := y - i
				if iY < bounds.Min.Y {
					iY = bounds.Min.Y + i
				}
				rP, gP, bP := buffR[iY*width+x], buffG[iY*width+x], buffB[iY*width+x]
				iY = y + i
				if iY >= bounds.Max.Y {
					iY = bounds.Max.Y - 1 - i
				}
				rN, gN, bN := buffR[iY*width+x], buffG[iY*width+x], buffB[iY*width+x]
				sR += rP*kernel[mid-i] + rN*kernel[mid+i]
				sG += gP*kernel[mid-i] + gN*kernel[mid+i]
				sB += bP*kernel[mid-i] + bN*kernel[mid+i]
			}
			sR = (sR + 1<<15) >> 16
			if sR > 255 {
				sR = 255
			} else if sR < 0 {
				sR = 0
			} else {
				sR--
			}
			sG = (sG + 1<<15) >> 16
			if sG > 255 {
				sG = 255
			} else if sG < 0 {
				sG = 0
			} else {
				sG--
			}
			sB = (sB + 1<<15) >> 16
			if sB > 255 {
				sB = 255
			} else if sB < 0 {
				sB = 0
			} else {
				sB--
			}
			_, _, _, a := img.At(x, y).RGBA()
			dst.Set(x, y, color.RGBA{uint8(sR), uint8(sG), uint8(sB), uint8(a / 0x101)})
		}
	}
	return dst
}
