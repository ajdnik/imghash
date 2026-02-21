package imghash

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/ajdnik/imghash/hashtype"
)

// OpenImage reads and decodes an image from the given file path.
// It supports JPEG, PNG, and GIF formats.
func OpenImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
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
