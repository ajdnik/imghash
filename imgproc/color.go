package imgproc

import (
	"errors"
	"image"
	"image/color"
)

// ErrImageIsNil ...
var ErrImageIsNil = errors.New("image argument is nil")

// Grayscale ...
func Grayscale(img image.Image) (*image.Gray, error) {
	if img == nil {
		return nil, ErrImageIsNil
	}
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray, nil
}

// YCrCb ...
func YCrCb(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, ErrImageIsNil
	}
	bounds := img.Bounds()
	ycrcb := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			yy, cb, cr := rgbToYCbCr(uint8(r/0x101), uint8(g/0x101), uint8(b/0x101))
			ycrcb.SetRGBA(x, y, color.RGBA{cb, cr, yy, uint8(a / 0x101)})
		}
	}
	return ycrcb, nil
}

// HSV ...
func HSV(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, ErrImageIsNil
	}
	bounds := img.Bounds()
	hsv := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			h, s, v := rgbToHSV(uint8(r/0x101), uint8(g/0x101), uint8(b/0x101))
			hsv.SetRGBA(x, y, color.RGBA{v, s, h, uint8(a / 0x101)})
		}
	}
	return hsv, nil
}

// rgbToYCbCr ...
func rgbToYCbCr(r, g, b uint8) (uint8, uint8, uint8) {
	yuvShift := uint8(14)
	var delta = 128 * (1 << yuvShift)
	y1 := int(r)*4899 + int(g)*9617 + int(b)*1868
	y2 := (((y1) + (1 << ((yuvShift) - 1))) >> (yuvShift))
	cr1 := (int(r)-y2)*11682 + delta
	cr2 := (((cr1) + (1 << ((yuvShift) - 1))) >> (yuvShift))
	cb1 := (int(b)-y2)*9241 + delta
	cb2 := (((cb1) + (1 << ((yuvShift) - 1))) >> (yuvShift))
	var Y, Cr, Cb uint8 = saturateCastIToUI8(y2), saturateCastIToUI8(cr2), saturateCastIToUI8(cb2)
	return Y, Cb, Cr
}

// rgbToHSV ...
func rgbToHSV(r, g, b uint8) (uint8, uint8, uint8) {
	hsvShift := uint8(12)
	var R, G, B int = int(r), int(g), int(b)
	var h, s int
	v := B
	vmin := B
	var vr, vg int = 0, 0

	if -256 <= (G-v) && (G-v) <= 512 {
		v += int(sat8U[(G-v)+256])
	}
	if -256 <= (R-v) && (R-v) <= 512 {
		v += int(sat8U[(R-v)+256])
	}
	if -256 <= (vmin-G) && (vmin-G) <= 512 {
		vmin -= int(sat8U[(vmin-G)+256])
	}
	if -256 <= (vmin-R) && (vmin-R) <= 512 {
		vmin -= int(sat8U[(vmin-R)+256])
	}

	diff := saturateCastIToUI8(v - vmin)

	if v == R {
		vr = -1
	}
	if v == G {
		vg = -1
	}

	s = (int(diff)*sdiv[v] + (1 << (hsvShift - 1))) >> hsvShift
	h = (vr & (G - B)) + (^vr & ((vg & (B - R + 2*int(diff))) + ((^vg) & (R - G + 4*int(diff)))))
	h = (h*hdiv180[diff] + (1 << (hsvShift - 1))) >> hsvShift
	if h < 0 {
		h += 180
	}
	H := saturateCastIToUI8(h)
	S := uint8(s)
	V := uint8(v)
	return H, S, V
}

var sdiv = [256]int{0, 1044480, 522240, 348160, 261120, 208896, 174080, 149211, 130560, 116053, 104448, 94953, 87040, 80345, 74606, 69632, 65280, 61440, 58027, 54973, 52224, 49737, 47476, 45412, 43520, 41779, 40172, 38684, 37303, 36017, 34816, 33693, 32640, 31651, 30720, 29842, 29013, 28229, 27486, 26782, 26112, 25475, 24869, 24290, 23738, 23211, 22706, 22223, 21760, 21316, 20890, 20480, 20086, 19707, 19342, 18991, 18651, 18324, 18008, 17703, 17408, 17123, 16846, 16579, 16320, 16069, 15825, 15589, 15360, 15137, 14921, 14711, 14507, 14308, 14115, 13926, 13743, 13565, 13391, 13221, 13056, 12895, 12738, 12584, 12434, 12288, 12145, 12006, 11869, 11736, 11605, 11478, 11353, 11231, 11111, 10995, 10880, 10768, 10658, 10550, 10445, 10341, 10240, 10141, 10043, 9947, 9854, 9761, 9671, 9582, 9495, 9410, 9326, 9243, 9162, 9082, 9004, 8927, 8852, 8777, 8704, 8632, 8561, 8492, 8423, 8356, 8290, 8224, 8160, 8097, 8034, 7973, 7913, 7853, 7795, 7737, 7680, 7624, 7569, 7514, 7461, 7408, 7355, 7304, 7253, 7203, 7154, 7105, 7057, 7010, 6963, 6917, 6872, 6827, 6782, 6739, 6695, 6653, 6611, 6569, 6528, 6487, 6447, 6408, 6369, 6330, 6292, 6254, 6217, 6180, 6144, 6108, 6073, 6037, 6003, 5968, 5935, 5901, 5868, 5835, 5803, 5771, 5739, 5708, 5677, 5646, 5615, 5585, 5556, 5526, 5497, 5468, 5440, 5412, 5384, 5356, 5329, 5302, 5275, 5249, 5222, 5196, 5171, 5145, 5120, 5095, 5070, 5046, 5022, 4998, 4974, 4950, 4927, 4904, 4881, 4858, 4836, 4813, 4791, 4769, 4748, 4726, 4705, 4684, 4663, 4642, 4622, 4601, 4581, 4561, 4541, 4522, 4502, 4483, 4464, 4445, 4426, 4407, 4389, 4370, 4352, 4334, 4316, 4298, 4281, 4263, 4246, 4229, 4212, 4195, 4178, 4161, 4145, 4128, 4112, 4096}

var hdiv180 = [256]int{0, 122880, 61440, 40960, 30720, 24576, 20480, 17554, 15360, 13653, 12288, 11171, 10240, 9452, 8777, 8192, 7680, 7228, 6827, 6467, 6144, 5851, 5585, 5343, 5120, 4915, 4726, 4551, 4389, 4237, 4096, 3964, 3840, 3724, 3614, 3511, 3413, 3321, 3234, 3151, 3072, 2997, 2926, 2858, 2793, 2731, 2671, 2614, 2560, 2508, 2458, 2409, 2363, 2318, 2276, 2234, 2194, 2156, 2119, 2083, 2048, 2014, 1982, 1950, 1920, 1890, 1862, 1834, 1807, 1781, 1755, 1731, 1707, 1683, 1661, 1638, 1617, 1596, 1575, 1555, 1536, 1517, 1499, 1480, 1463, 1446, 1429, 1412, 1396, 1381, 1365, 1350, 1336, 1321, 1307, 1293, 1280, 1267, 1254, 1241, 1229, 1217, 1205, 1193, 1182, 1170, 1159, 1148, 1138, 1127, 1117, 1107, 1097, 1087, 1078, 1069, 1059, 1050, 1041, 1033, 1024, 1016, 1007, 999, 991, 983, 975, 968, 960, 953, 945, 938, 931, 924, 917, 910, 904, 897, 890, 884, 878, 871, 865, 859, 853, 847, 842, 836, 830, 825, 819, 814, 808, 803, 798, 793, 788, 783, 778, 773, 768, 763, 759, 754, 749, 745, 740, 736, 731, 727, 723, 719, 714, 710, 706, 702, 698, 694, 690, 686, 683, 679, 675, 671, 668, 664, 661, 657, 654, 650, 647, 643, 640, 637, 633, 630, 627, 624, 621, 617, 614, 611, 608, 605, 602, 599, 597, 594, 591, 588, 585, 582, 580, 577, 574, 572, 569, 566, 564, 561, 559, 556, 554, 551, 549, 546, 544, 541, 539, 537, 534, 532, 530, 527, 525, 523, 521, 518, 516, 514, 512, 510, 508, 506, 504, 502, 500, 497, 495, 493, 492, 490, 488, 486, 484, 482}

var hdiv256 = [256]int{0, 174763, 87381, 58254, 43691, 34953, 29127, 24966, 21845, 19418, 17476, 15888, 14564, 13443, 12483, 11651, 10923, 10280, 9709, 9198, 8738, 8322, 7944, 7598, 7282, 6991, 6722, 6473, 6242, 6026, 5825, 5638, 5461, 5296, 5140, 4993, 4855, 4723, 4599, 4481, 4369, 4263, 4161, 4064, 3972, 3884, 3799, 3718, 3641, 3567, 3495, 3427, 3361, 3297, 3236, 3178, 3121, 3066, 3013, 2962, 2913, 2865, 2819, 2774, 2731, 2689, 2648, 2608, 2570, 2533, 2497, 2461, 2427, 2394, 2362, 2330, 2300, 2270, 2241, 2212, 2185, 2158, 2131, 2106, 2081, 2056, 2032, 2009, 1986, 1964, 1942, 1920, 1900, 1879, 1859, 1840, 1820, 1802, 1783, 1765, 1748, 1730, 1713, 1697, 1680, 1664, 1649, 1633, 1618, 1603, 1589, 1574, 1560, 1547, 1533, 1520, 1507, 1494, 1481, 1469, 1456, 1444, 1432, 1421, 1409, 1398, 1387, 1376, 1365, 1355, 1344, 1334, 1324, 1314, 1304, 1295, 1285, 1276, 1266, 1257, 1248, 1239, 1231, 1222, 1214, 1205, 1197, 1189, 1181, 1173, 1165, 1157, 1150, 1142, 1135, 1128, 1120, 1113, 1106, 1099, 1092, 1085, 1079, 1072, 1066, 1059, 1053, 1046, 1040, 1034, 1028, 1022, 1016, 1010, 1004, 999, 993, 987, 982, 976, 971, 966, 960, 955, 950, 945, 940, 935, 930, 925, 920, 915, 910, 906, 901, 896, 892, 887, 883, 878, 874, 869, 865, 861, 857, 853, 848, 844, 840, 836, 832, 828, 824, 820, 817, 813, 809, 805, 802, 798, 794, 791, 787, 784, 780, 777, 773, 770, 767, 763, 760, 757, 753, 750, 747, 744, 741, 737, 734, 731, 728, 725, 722, 719, 716, 713, 710, 708, 705, 702, 699, 696, 694, 691, 688, 685}

var sat8U = []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
