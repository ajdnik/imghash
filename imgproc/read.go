package imgproc

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func Read(file string) (image.Image, error) {
	iFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer iFile.Close()
	img, _, err := image.Decode(iFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}
