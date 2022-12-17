package qrlogo

import (
	"image"
	"image/color"
	"log"
	"math"
	"sync"

	"github.com/alexandernorth/qrlogo/imagetools"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
)

type QRLogo struct {
	logo          image.Image
	processedLogo image.Image
	loadLogoOnce  sync.Once
	qrLogoOptions
}

// NewQRLogo creates a new QRLogo, taking the logo to be placed in the generated QRCode
func NewQRLogo(logo image.Image, opts ...Options) *QRLogo {
	qrLogo := &QRLogo{
		logo: logo,
		qrLogoOptions: qrLogoOptions{
			backgroundColor: color.White,
			codeColour:      color.Black,
			disableBorder:   false,
			logoCoverage:    .2,
			logoPaddingType: LogoPaddingDilate,
			paddingWeight:   20,
			qrSize:          2048,
		},
	}
	for _, opt := range opts {
		opt(qrLogo)
	}

	return qrLogo
}

// QRCodeForURL generates a new QRCode image.Image for the given URL, containing the logo
func (q *QRLogo) QRCodeForURL(url string) (image.Image, error) {
	qrCode, err := qrcode.New(url, qrcode.Highest)
	if err != nil {
		return nil, err
	}
	qrCode.BackgroundColor = q.backgroundColor
	qrCode.ForegroundColor = q.codeColour
	qrCode.DisableBorder = q.disableBorder

	// Draw QRCode onto canvas
	qrImage := qrCode.Image(q.qrSize)
	qrLogoCanvas := image.NewNRGBA(qrImage.Bounds())
	draw.Draw(qrLogoCanvas, qrLogoCanvas.Rect, qrImage, image.Point{}, draw.Src)

	// Get logo image to place
	logo := q.generateLogo()
	// Calculate the offset, to centre logo in QRCode
	offsetX := (qrLogoCanvas.Bounds().Max.X - logo.Bounds().Dx()) / 2
	offsetY := (qrLogoCanvas.Bounds().Max.Y - logo.Bounds().Dy()) / 2

	// Rectangle to draw logo into - centered using offsetX/Y
	logoRect := image.Rect(offsetX, offsetY, offsetX+logo.Bounds().Dx(), offsetY+logo.Bounds().Dy())

	// Draw logo onto QRCode
	draw.Draw(qrLogoCanvas, logoRect, logo, image.Point{}, draw.Over)

	return qrLogoCanvas, nil
}

// generateLogo generates the logo according to qrLogoOptions
// it caches the generated logo for future runs
func (q *QRLogo) generateLogo() image.Image {
	// Only generate logo once - performance boost for multiple calls
	q.loadLogoOnce.Do(func() {
		logoPaddedSize := math.Sqrt(float64(q.qrSize*q.qrSize) * q.logoCoverage)
		var logo image.Image
		switch q.logoPaddingType {
		case LogoPaddingNone:
			logo = q.scaledLogo(logoPaddedSize)
		case LogoPaddingSquare:
			logo = q.squareLogo(logoPaddedSize)
		case LogoPaddingCircle:
			logo = q.circleLogo(logoPaddedSize)
		case LogoPaddingDilate:
			logo = q.dilatedLogo(logoPaddedSize)
		default:
			log.Fatal("unknown padding type")
		}
		q.processedLogo = logo
	})
	return q.processedLogo
}

// scaledLogo resizes the logo image to fit to scaledSize
func (q *QRLogo) scaledLogo(scaledSize float64) image.Image {
	// calculate resize factor keeping aspect ratio
	ratio := math.Min(scaledSize/float64(q.logo.Bounds().Dx()), scaledSize/float64(q.logo.Bounds().Dy()))
	width := math.Min(math.Ceil(float64(q.logo.Bounds().Dx())*ratio), scaledSize)
	height := math.Min(math.Ceil(float64(q.logo.Bounds().Dy())*ratio), scaledSize)

	// Rectangle to draw logo into
	logoRect := image.Rect(0, 0, int(width), int(height))

	dest := image.NewNRGBA(logoRect)
	// Draw logo scaled into new Rectangle using Catmull-Rom kernel
	draw.CatmullRom.Scale(dest, logoRect, q.logo, q.logo.Bounds(), draw.Over, nil)
	return dest
}

// dilatedLogo generates a logo with padding type LogoPaddingDilate
func (q *QRLogo) dilatedLogo(paddedSize float64) image.Image {
	logoSize := paddedSize - (2 * float64(q.paddingWeight))
	scaledLogo := q.scaledLogo(logoSize)
	// get the mask from dilation algorithm
	paddingMask := imagetools.DilateMaskFromImage(scaledLogo, q.paddingWeight)

	dest := image.NewNRGBA(paddingMask.Bounds())
	// Draw q.backgroundColor into mask
	draw.DrawMask(dest, paddingMask.Bounds(), image.NewUniform(q.backgroundColor), image.Point{}, paddingMask, image.Point{}, draw.Src)

	// place logo
	logoRect := image.Rect(q.paddingWeight, q.paddingWeight, scaledLogo.Bounds().Max.X+q.paddingWeight, scaledLogo.Bounds().Max.Y+q.paddingWeight)
	draw.Draw(dest, logoRect, scaledLogo, image.Point{}, draw.Over)

	return dest
}

// squareLogo generates a logo with padding type LogoPaddingSquare
func (q *QRLogo) squareLogo(paddedSize float64) image.Image {
	logoSize := paddedSize - (2 * float64(q.paddingWeight))
	scaledLogo := q.scaledLogo(logoSize)

	// draw solid square
	dest := image.NewNRGBA(image.Rect(0, 0, int(math.Floor(paddedSize)), int(math.Floor(paddedSize))))
	draw.Draw(dest, dest.Bounds(), image.NewUniform(q.backgroundColor), image.Point{}, draw.Src)

	offsetX := (dest.Bounds().Max.X - scaledLogo.Bounds().Max.X) / 2
	offsetY := (dest.Bounds().Max.Y - scaledLogo.Bounds().Max.Y) / 2

	// place logo
	logoRect := image.Rect(offsetX, offsetY, offsetX+scaledLogo.Bounds().Max.X, offsetY+scaledLogo.Bounds().Max.Y)
	draw.Draw(dest, logoRect, scaledLogo, image.Point{}, draw.Over)
	return dest
}

// circleLogo generates a logo with padding type LogoPaddingCircle
func (q *QRLogo) circleLogo(paddedSize float64) image.Image {
	logoSize := paddedSize - (2 * float64(q.paddingWeight))
	scaledLogo := q.scaledLogo(logoSize)

	// Get a filled circle
	dest := imagetools.DrawFilledCircle(int(math.Floor(paddedSize/2)), q.backgroundColor)

	offsetX := (dest.Bounds().Max.X - scaledLogo.Bounds().Max.X) / 2
	offsetY := (dest.Bounds().Max.Y - scaledLogo.Bounds().Max.Y) / 2

	// place logo
	logoRect := image.Rect(offsetX, offsetY, offsetX+scaledLogo.Bounds().Max.X, offsetY+scaledLogo.Bounds().Max.Y)
	draw.Draw(dest, logoRect, scaledLogo, image.Point{}, draw.Over)
	return dest
}
