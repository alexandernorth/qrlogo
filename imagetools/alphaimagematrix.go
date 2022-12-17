package imagetools

import (
	"fmt"
	"image"
)

type AlphaImageMatrix struct {
	// backingImage stores the image represented by the matrix
	backingImage image.Image
	// pixels is a slice containing the alpha value of the pixels of the backing image
	pixels []uint32
	// stride is the stride (in bytes) between vertically adjacent pixels.
	stride int
}

// NewAlphaImageMatrix creates a new AlphaImageMatrix from the specified src image
func NewAlphaImageMatrix(img image.Image) *AlphaImageMatrix {
	matrix := &AlphaImageMatrix{
		backingImage: img,
		pixels:       make([]uint32, img.Bounds().Dx()*img.Bounds().Dy()),
		stride:       1 * img.Bounds().Dx(),
	}

	for i := img.Bounds().Min.X; i < img.Bounds().Max.X; i++ {
		for j := img.Bounds().Min.Y; j < img.Bounds().Max.Y; j++ {
			_, _, _, a := img.At(i, j).RGBA()
			matrix.Set(i, j, a)
		}
	}

	return matrix
}

// In returns whether an (x,y) point lies within the image bounds
func (m *AlphaImageMatrix) In(x int, y int) bool {
	return image.Point{X: x, Y: y}.In(m.backingImage.Bounds())
}

// pixOffset returns the index of a point in the matrix backing slice
func (m *AlphaImageMatrix) pixOffset(x int, y int) int {
	return (y-m.backingImage.Bounds().Min.Y)*m.stride + (x-m.backingImage.Bounds().Min.X)*1
}

// Get gets the alpha value at point (x,y)
func (m *AlphaImageMatrix) Get(x int, y int) uint32 {
	if !m.In(x, y) {
		panic(fmt.Errorf("out of bounds"))
	}
	return m.pixels[m.pixOffset(x, y)]
}

// Set sets the alpha value at point (x,y)
func (m *AlphaImageMatrix) Set(x int, y int, value uint32) {
	if !m.In(x, y) {
		panic(fmt.Errorf("out of bounds"))
	}
	m.pixels[m.pixOffset(x, y)] = value
}

// NumColumns returns the number of columns of the matrix
func (m *AlphaImageMatrix) NumColumns() int {
	return m.backingImage.Bounds().Dx()
}

// NumRows returns the number of rows of the matrix
func (m *AlphaImageMatrix) NumRows() int {
	return m.backingImage.Bounds().Dy()
}

// ImageBounds returns the bounding image.Rectangle of the image
func (m *AlphaImageMatrix) ImageBounds() image.Rectangle {
	return m.backingImage.Bounds()
}

// Clone returns a copy of the matrix
func (m *AlphaImageMatrix) Clone() *AlphaImageMatrix {
	clone := NewAlphaImageMatrix(m.backingImage)
	return clone
}
