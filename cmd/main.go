package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexandernorth/qrlogo"
	"golang.org/x/image/colornames"
)

var (
	backgroundColor = flag.String("background-color", "white", "set background color based on SVG 1.1 spec color names")
	codeColor       = flag.String("code-color", "black", "set code color based on SVG 1.1 spec color names")
	disableBorder   = flag.Bool("disable-border", false, "set to remove border around the whole QRCode")
	logo            = flag.String("logo", "", "the logo to place in the centre of the code")
	logoCoverage    = flag.Float64("logo-coverage", 0.2, "set percentage [0,1] of code to cover with logo, this value should not exceed 0.3")
	outputDir       = flag.String("output-dir", "qrcodes", "the output image directory")
	paddingType     = flag.String("padding", "dilate", "the type of padding around the logo")
	paddingWeight   = flag.Int("padding-weight", 20, "adjust the padding around the logo")
	qrSize          = flag.Int("qr-size", 2048, "size of the qr code output")
)

func main() {

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	if *logo == "" {
		log.Fatal("path to logo needs to be given using --logo")
	}

	if *logoCoverage > 0.3 {
		log.Print("--logo-coverage should be no more than 0.3 to ensure QRCode functions")
	}

	logo, err := qrlogo.OpenImage(*logo)
	if err != nil {
		panic(err)
	}
	qlg := qrlogo.NewQRLogo(
		logo,
		qrlogo.BackgroundColor(colornames.Map[*backgroundColor]),
		qrlogo.CodeColour(colornames.Map[*codeColor]),
		qrlogo.DisableBorder(*disableBorder),
		qrlogo.LogoCoverage(*logoCoverage),
		qrlogo.PaddingType(qrlogo.LogoPaddingType(*paddingType)),
		qrlogo.PaddingWeight(*paddingWeight),
		qrlogo.QRSize(*qrSize),
	)

	for i := 0; i < flag.NArg(); i++ {
		url := flag.Arg(i)
		qrLogo, err := qlg.QRCodeForURL(url)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(*outputDir, 0755)
		if err != nil {
			panic(err)
		}
		err = qrlogo.SaveToPNGFile(qrLogo, path.Join(*outputDir, fmt.Sprintf("qr-%03d.png", i+1)))
		if err != nil {
			panic(err)
		}

	}

}
