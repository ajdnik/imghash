package dinov2

import (
	"image"

	"github.com/ajdnik/imghash/v2/internal/imgproc"
)

// Input image dimensions and channel count expected by the embedded ViT-S/14
// model. These are not user-configurable: the published weights are anchored
// to 224x224 RGB ImageNet-normalized inputs.
const (
	ImageSize  = 224
	Channels   = 3
	PixelCount = ImageSize * ImageSize
)

// ImageNet normalization constants used by the DINOv2 reference preprocessing.
var (
	imagenetMean = [3]float32{0.485, 0.456, 0.406}
	imagenetStd  = [3]float32{0.229, 0.224, 0.225}
)

// ImageToTensor resizes img to 224x224 with bilinear interpolation,
// normalizes each pixel with ImageNet statistics, and lays the result out in
// NCHW float32 order (3*224*224 = 150528 entries).
//
// The Python reference uses transforms.Resize((224, 224), interpolation=BILINEAR)
// followed by ToTensor()+Normalize(); golang.org/x/image's BilinearExact is
// the closest match available in this codebase.
func ImageToTensor(img image.Image) []float32 {
	resized := imgproc.Resize(ImageSize, ImageSize, img, imgproc.BilinearExact)
	out := make([]float32, Channels*PixelCount)

	switch r := resized.(type) {
	case *image.RGBA:
		for y := 0; y < ImageSize; y++ {
			for x := 0; x < ImageSize; x++ {
				i := r.PixOffset(x, y)
				rv := float32(r.Pix[i]) / 255
				gv := float32(r.Pix[i+1]) / 255
				bv := float32(r.Pix[i+2]) / 255
				idx := y*ImageSize + x
				out[0*PixelCount+idx] = (rv - imagenetMean[0]) / imagenetStd[0]
				out[1*PixelCount+idx] = (gv - imagenetMean[1]) / imagenetStd[1]
				out[2*PixelCount+idx] = (bv - imagenetMean[2]) / imagenetStd[2]
			}
		}
	case *image.Gray:
		for y := 0; y < ImageSize; y++ {
			for x := 0; x < ImageSize; x++ {
				v := float32(r.Pix[r.PixOffset(x, y)]) / 255
				idx := y*ImageSize + x
				out[0*PixelCount+idx] = (v - imagenetMean[0]) / imagenetStd[0]
				out[1*PixelCount+idx] = (v - imagenetMean[1]) / imagenetStd[1]
				out[2*PixelCount+idx] = (v - imagenetMean[2]) / imagenetStd[2]
			}
		}
	default:
		for y := 0; y < ImageSize; y++ {
			for x := 0; x < ImageSize; x++ {
				r16, g16, b16, _ := resized.At(x, y).RGBA()
				rv := float32(r16>>8) / 255
				gv := float32(g16>>8) / 255
				bv := float32(b16>>8) / 255
				idx := y*ImageSize + x
				out[0*PixelCount+idx] = (rv - imagenetMean[0]) / imagenetStd[0]
				out[1*PixelCount+idx] = (gv - imagenetMean[1]) / imagenetStd[1]
				out[2*PixelCount+idx] = (bv - imagenetMean[2]) / imagenetStd[2]
			}
		}
	}
	return out
}
