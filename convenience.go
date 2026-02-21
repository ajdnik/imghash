package imghash

import (
	"image"
	_ "image/gif"  // register GIF decoder
	_ "image/jpeg" // register JPEG decoder
	_ "image/png"  // register PNG decoder
	"io"
	"os"

	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

// OpenImage reads and decodes an image from the given file path.
// It supports JPEG, PNG, and GIF formats.
func OpenImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	img, _, err := image.Decode(f)
	return img, err
}

// DecodeImage decodes an image from the given reader.
// It supports JPEG, PNG, and GIF formats.
func DecodeImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	return img, err
}

// HashFile is a convenience that opens an image file and computes its hash.
func HashFile(hasher Hasher, path string) (hashtype.Hash, error) {
	img, err := OpenImage(path)
	if err != nil {
		return nil, err
	}
	return hasher.Calculate(img)
}

// HashReader is a convenience that decodes an image from a reader and computes its hash.
func HashReader(hasher Hasher, r io.Reader) (hashtype.Hash, error) {
	img, err := DecodeImage(r)
	if err != nil {
		return nil, err
	}
	return hasher.Calculate(img)
}

// Compare computes the distance between two hashes.
// By default it uses the natural metric for their type: Hamming distance
// for Binary hashes, L2 (Euclidean) distance for UInt8 and Float64 hashes.
// Pass an optional DistanceFunc to override the metric, e.g.:
//
//	Compare(h1, h2, similarity.Cosine)
func Compare(h1, h2 hashtype.Hash, fn ...DistanceFunc) (similarity.Distance, error) {
	if len(fn) > 0 && fn[0] != nil {
		return fn[0](h1, h2)
	}
	if _, ok := h1.(hashtype.Binary); ok {
		return similarity.Hamming(h1, h2)
	}
	return similarity.L2(h1, h2)
}
