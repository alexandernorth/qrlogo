package imagetools

import (
	"image"

	"golang.org/x/image/draw"
)

type Number interface {
	int64 | int32 | int16 | int8 | int | uint64 | uint32 | uint16 | uint8 | uint
}

const maxAlpha = 0xffff
const alphaCutoff = maxAlpha / 2

// min calculates the min of a,b using generics
func min[N Number](a N, b N) N {
	if a < b {
		return a
	}
	return b
}

// ImageWithPadding returns a new image, padded by the given padding
func ImageWithPadding(img image.Image, padding int) image.Image {
	out := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx()+2*padding, img.Bounds().Dy()+2*padding))
	draw.Draw(out, image.Rect(padding, padding, img.Bounds().Max.X+padding, img.Bounds().Max.Y+padding), img, image.Point{}, draw.Over)
	return out
}
