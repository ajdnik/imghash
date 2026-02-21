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

// Compare computes the distance between two hashes using the natural
// similarity metric for their type: Hamming distance for Binary hashes,
// L2 (Euclidean) distance for UInt8 and Float64 hashes.
// For finer control (e.g. PCC), use the similarity package directly.
func Compare(h1, h2 hashtype.Hash) (similarity.Distance, error) {
	d, err := h1.Distance(h2)
	return similarity.Distance(d), err
}
