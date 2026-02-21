package imgproc

import "math"

// saturateCastF32ToUI8 casts a float 32 bit value into an uint 8 bit
// value by cutting off values larger or smaller than the limits
func saturateCastF32ToUI8(val float32) uint8 {
	if val > 255 {
		return 255
	} else if val < 0 {
		return 0
	}
	return uint8(math.Round(float64(val)))
}

func saturateCastIToUI8(val int) uint8 {
	if val > 255 {
		return 255
	} else if val < 0 {
		return 0
	}
	return uint8(val)
}

func saturateCastUI8ToF32(val uint8) float32 {
	return float32(val)
}
