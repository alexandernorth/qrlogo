package imagetools

import (
	"image"
	"image/color"
)

// Translated from https://blog.ostermiller.org/efficiently-implementing-dilate-and-erode-image-functions/

// GetManhattanMatrix generates a new matrix containing the manhattan distance away from the nearest
// pixel with an alpha value of at least alphaCutoff
func (m *AlphaImageMatrix) GetManhattanMatrix() *AlphaImageMatrix {
	matrix := m.Clone()
	// Start top left and work to bottom right
	imgBounds := matrix.ImageBounds()
	// Start top left, go bottom right
	for i := imgBounds.Min.X; i < imgBounds.Max.X; i++ {
		for j := imgBounds.Min.Y; j < imgBounds.Max.Y; j++ {
			// Get the alpha value of the pixel and check if it is set (> alphaCutoff) or not set
			a := matrix.Get(i, j)
			if a > alphaCutoff {
				// pixel set, set MM value zero
				matrix.Set(i, j, 0)
			} else {
				// pixel not set
				// it is at most the width+height distance away from the next pixel
				// away from a pixel it is on
				dist := uint32(matrix.NumColumns() + matrix.NumRows())

				// Are we closer to a pixel to the left?
				if i > imgBounds.Min.X {
					westPix := matrix.Get(i-1, j)
					dist = min(dist, westPix+1)
				}

				// Are we closer to a pixel above?
				if j > imgBounds.Min.Y {
					northPix := matrix.Get(i, j-1)
					dist = min(dist, northPix+1)
				}
				// Remember the closest pixel
				matrix.Set(i, j, dist)
			}
		}
	}

	// Then try again, bottom right to top left
	for i := imgBounds.Max.X - 1; i >= imgBounds.Min.X; i-- {
		for j := imgBounds.Max.Y - 1; j >= imgBounds.Min.Y; j-- {

			dist := matrix.Get(i, j)

			// Are we closer to a pixel to the right?
			if i+1 < imgBounds.Max.X {
				eastPix := matrix.Get(i+1, j)
				dist = min(dist, eastPix+1)
			}

			// Are we closer to a pixel below?
			if j+1 < imgBounds.Max.Y {
				southPix := matrix.Get(i, j+1)
				dist = min(dist, southPix+1)
			}
			// Remember the closest pixel
			matrix.Set(i, j, dist)
		}
	}
	return matrix
}

// DilateMask returns an Alpha mask with the specified padding around the source image.
// This does not modify the original matrix or image.
func (m *AlphaImageMatrix) DilateMask(padding uint32) *image.Alpha {
	// Calculate the manhattan distance to the (alpha) pixels in the AlphaImageMatrix
	manhattan := m.GetManhattanMatrix()
	// Create new Alpha image to define alpha mask
	output := image.NewAlpha(image.Rect(0, 0, m.NumColumns(), m.NumRows()))

	// Iterate through manhattan distance
	// if we are closer or equal to the padding, set mask to be opaque
	for i := 0; i < m.NumColumns(); i++ {
		for j := 0; j < m.NumRows(); j++ {
			if manhattan.Get(i, j) <= padding {
				output.Set(i, j, color.Opaque)
			}
		}
	}
	return output
}

// DilateMaskFromImage takes a src image, and runs the dilation algorithm with the specified padding
// It returns a new image and does not modify the original
func DilateMaskFromImage(img image.Image, padding int) image.Image {
	img = ImageWithPadding(img, padding)
	mat := NewAlphaImageMatrix(img)
	return mat.DilateMask(uint32(padding))
}
