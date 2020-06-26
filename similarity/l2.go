package similarity

import (
	"math"

	"github.com/ajdnik/imghash/hashtype"
)

// L2Float64 calculates distance between two float64 hashes.
func L2Float64(h1, h2 hashtype.Float64) Distance {
	return l2Calculate(h1, h2)
}

// L2UInt8 calculates distance between two uint8 hashes.
func L2UInt8(h1, h2 hashtype.UInt8) Distance {
	return l2Calculate(h1, h2)
}

// Calculates L2 norm for a generic slice.
func l2Calculate(h1, h2 interface{}) Distance {
	len1, len2 := genLen(h1), genLen(h2)
	l := int(math.Min(float64(len1), float64(len2)))
	var s float64
	for i := 0; i < l; i++ {
		val1, val2 := genElem(h1, i), genElem(h2, i)
		s += math.Pow(val1-val2, 2)
	}
	return Distance(math.Sqrt(s))
}

// Get length of slice passed as generic interface.
func genLen(val interface{}) int {
	var l int
	switch vv := val.(type) {
	case hashtype.Float64:
		l = len(vv)
	case hashtype.UInt8:
		l = len(vv)
	}
	return l
}

// Get element of slice passed as generic interface.
func genElem(val interface{}, idx int) float64 {
	var v float64
	switch vv := val.(type) {
	case hashtype.Float64:
		v = vv[idx]
	case hashtype.UInt8:
		v = float64(vv[idx])
	}
	return v
}
