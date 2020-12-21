/*
 * Font subsetting example by using optimizer option of `SubsetFonts`.
 *
 * Note: This option would subset all fonts used in the documents.
 * To subset only some font file, consider using `EnableFontSubsetting`
 * in creator package.
 *
 * Run as: go run font_subsetting.go <font.ttf>
 */

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

const usage = "Usage: %s FONT_TTF_FILE\n"
const outputOptimizesFile = "font_subsetting_optimizer.pdf"

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	fontFile := args[1]

	font, err := model.NewPdfFontFromTTFFile(fontFile)
	if err != nil {
		log.Fatalln(err)
	}

	subsetUsingOptimizer(font)
}

func subsetUsingOptimizer(font *model.PdfFont) {
	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)
	c.NewPage()

	// Set optimizer with `SubsetFonts` enabled.
	c.SetOptimizer(optimize.New(optimize.Options{
		SubsetFonts: true,
	}))

	p := c.NewStyledParagraph()
	p.SetPos(100, 100)
	text := p.SetText("This is an example of using Optimizer")
	text.Style.Font = font

	c.Draw(p)

	err := c.WriteToFile(outputOptimizesFile)
	if err != nil {
		log.Fatalln(err)
	}
}
