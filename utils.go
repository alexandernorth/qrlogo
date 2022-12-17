package qrlogo

import (
	"image"
	"image/png"
	"os"
)

// OpenImage loads an image file from disk, and decodes it into an image.Image
func OpenImage(logoPath string) (image.Image, error) {
	file, err := os.Open(logoPath)
	if err != nil {
		return nil, err
	}
	defer func(img *os.File) {
		_ = img.Close()
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// SaveToPNGFile saves an image to disk at the given filepath
func SaveToPNGFile(img image.Image, outFile string) error {
	of, err := os.Create(outFile)
	if err != nil {
		return err
	}
	err = png.Encode(of, img)
	if err != nil {
		return err
	}
	return nil
}
