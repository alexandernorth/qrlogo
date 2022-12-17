package imagetools

import (
	"image"
	"image/color"
)

// DrawFilledCircle returns a new image of a filled circle with radius r and fill color c
func DrawFilledCircle(r int, c color.Color) *image.NRGBA {
	dest := image.NewNRGBA(image.Rect(0, 0, r*2, r*2))
	x, y := r, r
	r2 := r * r
	area := r2 << 2
	rr := r << 1

	for i := 0; i < area; i++ {
		tx := (i % rr) - r
		ty := (i / rr) - r
		if tx*tx+ty*ty <= r2 {
			dest.Set(x+tx, y+ty, c)
		}
	}
	return dest
}
