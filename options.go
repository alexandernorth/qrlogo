package qrlogo

import "image/color"

type LogoPaddingType string

const (
	// LogoPaddingCircle adds a circle border around the logo
	LogoPaddingCircle = "circle"
	// LogoPaddingDilate adds a dilated (traced) border around the logo
	LogoPaddingDilate = "dilate"
	// LogoPaddingNone adds no padding
	LogoPaddingNone = "none"
	// LogoPaddingSquare adds a square border around the logo
	LogoPaddingSquare = "square"
)

type qrLogoOptions struct {
	backgroundColor color.Color
	codeColour      color.Color
	disableBorder   bool
	logoCoverage    float64
	logoPaddingType LogoPaddingType
	paddingWeight   int
	qrSize          int
}

type Options func(*QRLogo)

// BackgroundColor sets the background color of the code
func BackgroundColor(backgroundColor color.Color) Options {
	return func(qrLogo *QRLogo) {
		qrLogo.backgroundColor = backgroundColor
	}
}

// CodeColour sets the colour of the QRCode bars
func CodeColour(codeColor color.Color) Options {
	return func(qrLogo *QRLogo) {
		qrLogo.codeColour = codeColor
	}
}

// DisableBorder sets whether there is a border around the whole QRCode
func DisableBorder(disableBorder bool) Options {
	return func(qrLogo *QRLogo) {
		qrLogo.disableBorder = disableBorder
	}
}

// LogoCoverage sets the percentage of the QRCode to be covered in range [0,1] with 1 being 100%
// The error correction ability of a QRCode is at max 30%, so this value should *not* be exceeded,
// but other factors may force this value down
func LogoCoverage(logoCoverage float64) Options {
	return func(qrLogo *QRLogo) {
		qrLogo.logoCoverage = logoCoverage
	}
}

// PaddingType sets the type of padding to use when rendering the logo
func PaddingType(logoPaddingType LogoPaddingType) Options {
	return func(qrLogo *QRLogo) {
		qrLogo.logoPaddingType = logoPaddingType
	}
}

// PaddingWeight sets the amount of padding around the logo
func PaddingWeight(paddingWeight int) Options {
	return func(qrLogo *QRLogo) {
		qrLogo.paddingWeight = paddingWeight
	}
}

// QRSize sets the size of the output image
func QRSize(qrSize int) Options {
	return func(qrLogo *QRLogo) {
		qrLogo.qrSize = qrSize
	}
}
