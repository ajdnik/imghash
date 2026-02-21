package imghash_test

import (
	"fmt"
	"os"

	. "github.com/ajdnik/imghash"
)

func ExampleOpenImage() {
	img, err := OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	fmt.Println(img.Bounds())
	// Output: (0,0)-(490,733)
}

func ExampleHashFile() {
	avg := NewAverage()
	hash, err := HashFile(avg, "assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
	// Output: [255 255 15 7 1 0 0 0]
}

func ExampleHashReader() {
	f, err := os.Open("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	avg := NewAverage()
	hash, err := HashReader(avg, f)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
	// Output: [255 255 15 7 1 0 0 0]
}

func ExampleCompare() {
	avg := NewAverage()
	h1, _ := HashFile(avg, "assets/lena.jpg")
	h2, _ := HashFile(avg, "assets/cat.jpg")

	dist, err := Compare(h1, h2)
	if err != nil {
		panic(err)
	}
	fmt.Println(dist)
	// Output: 30
}

func ExampleNewMarrHildreth_options() {
	mh := NewMarrHildreth(
		WithScale(1),
		WithAlpha(2),
		WithSize(512, 512),
		WithInterpolation(Bicubic),
		WithKernelSize(7),
		WithSigma(0),
	)
	hash, _ := HashFile(mh, "assets/cat.jpg")
	fmt.Println(hash)
	// Output: [92 190 42 111 87 107 101 164 184 24 75 41 185 54 178 162 26 236 155 150 108 98 233 112 56 235 124 177 139 159 148 66 89 38 229 47 195 44 158 180 85 115 79 165 92 131 225 252 54 148 218 61 99 92 82 141 141 96 112 186 185 208 174 112 252 150 153 164 173 206 43 130]
}

func ExampleInterpolation_String() {
	fmt.Println(Bilinear)
	fmt.Println(Bicubic)
	fmt.Println(Lanczos3)
	// Output:
	// Bilinear
	// Bicubic
	// Lanczos3
}
