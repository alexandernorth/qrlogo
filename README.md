# QRLogo
A command line tool and library to generate QR Codes with logos/images embedded. The logos can be padded in different ways to make them stand out.

## Contents
<!-- TOC -->
* [QRLogo](#qrlogo)
  * [Contents](#contents)
  * [Installation](#installation)
  * [Quickstart](#quickstart)
  * [Options](#options)
    * [BackgroundColor](#backgroundcolor)
    * [CodeColour](#codecolour)
    * [DisableBorder](#disableborder)
    * [LogoCoverage](#logocoverage)
    * [PaddingType](#paddingtype)
    * [PaddingWeight](#paddingweight)
    * [QRSize](#qrsize)
<!-- TOC -->

## Installation

1. Go needs to be installed then the following command can be used to add qrlogo to your project:
```shell
go get -u github.com/alexandernorth/qrlogo
```

2. Import `qrlogo` into your code:
```go
import "github.com/alexandernorth/qrlogo"
```

## Quickstart
Create a file with the following contents:
```go
package main

import (
	"github.com/alexandernorth/qrlogo"
	"os"
	"path"
)

func main() {

	logo, err := qrlogo.OpenImage("logo.png")
	if err != nil {
		panic(err)
	}
	
	qlg := qrlogo.NewQRLogo(logo)
    
	qrLogo, err := qlg.QRCodeForURL("https://github.com/alexandernorth/qrlogo")
    if err != nil {
        panic(err)
    }
	
	outputDir := "qrcodes"
	
    err = os.MkdirAll(outputDir, 0755)
    if err != nil {
        panic(err)
    }
    err = qrlogo.SaveToPNGFile(qrLogo, path.Join(outputDir, "qr.png"))
    if err != nil {
        panic(err)
    }
}
```

## Options
QRLogo can be configured with several options:

### BackgroundColor
Sets the background color of the code
```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.BackgroundColor(colornames.Map["white"]),
)
```

### CodeColour
Sets the colour of the QRCode bars
```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.CodeColour(colornames.Map["black"]),
)
```

### DisableBorder
Sets whether there is a border around the whole QRCode
```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.DisableBorder(false),
)
```

### LogoCoverage
Sets the percentage of the QRCode to be covered in range [0,1] with 1 being 100%. The error correction ability of 
a QRCode is at max 30%, so this value should *not* be exceeded, but other factors may force this value down
```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.LogoCoverage(0.2),
)
```

### PaddingType
Sets the type of padding to use when rendering the logo.

Valid options are: 
- dilate
- none
- square
- circle

```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.PaddingType(qrlogo.LogoPaddingDilate),
)
```
`LogoPaddingType` is a string, so it is also possible to specify the type through a string:
```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.PaddingType(qrlogo.LogoPaddingType("dilate")),
)
```

### PaddingWeight
Sets the amount of padding around the logo
```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.PaddingWeight(20),
)
```

### QRSize
Sets the size of the output image
```go
qlg := qrlogo.NewQRLogo(
  logo,
  qrlogo.QRSize(2048),
)
```
