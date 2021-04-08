/*
 * An example of using CMYK color to set text color.
 *
 * Run as: go run pdf_cmyk_color.go
 */

package main

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/contentstream/draw"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	outputFile := "output.pdf"

	err := genPdfFile(outputFile)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func genPdfFile(outputFile string) error {
	fontRegular, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		return err
	}

	c := creator.New()
	c.SetPageMargins(50, 50, 100, 70)

	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		writeContent(c, fontRegular)
	})

	return c.WriteToFile(outputFile)
}

func writeContent(c *creator.Creator, font *model.PdfFont) {
	redColor := creator.ColorCMYKFromArithmetic(0.0, 1.0, 1.0, 0.0)
	blueColor := creator.ColorCMYKFrom8bit(100, 40, 0, 0)

	ch := c.NewChapter("CMYK color model")
	ch.GetHeading().SetColor(redColor)
	ch.GetHeading().SetFontSize(20)

	p := c.NewStyledParagraph()
	p.SetMargins(20, 10, 20, 0)
	text := p.SetText(`The CMYK color model (also known as process color, or four color) 
		is a subtractive color model, based on the CMY color model, used in color printing, 
		and is also used to describe the printing process itself. 
		CMYK refers to the four ink plates used in some color printing: 
		cyan, magenta, yellow, and key (black).`)
	text.Style.Color = blueColor
	text.Style.FontSize = 14

	ch.Add(p)

	p = c.NewStyledParagraph()
	p.SetMargins(20, 10, 20, 0)
	text = p.SetText(`The CMYK model works by partially or entirely masking colors on a lighter, 
		usually white, background. The ink reduces the light that would otherwise be reflected. 
		Such a model is called subtractive because inks "subtract" the colors 
		red, green and blue from white light. White light minus red leaves cyan, white light 
		minus green leaves magenta, and white light minus blue leaves yellow.`)
	text.Style.Color = blueColor
	text.Style.FontSize = 14

	ch.Add(p)

	p = c.NewStyledParagraph()
	p.SetMargins(20, 10, 20, 0)
	text = p.SetText(`In additive color models, such as RGB, white is the "additive" 
		combination of all primary colored lights, black is the absence of light. 
		In the CMYK model, it is the opposite: white is the natural color of the paper 
		or other background, black results from a full combination of colored inks. 
		To save cost on ink, and to produce deeper black tones, unsaturated and 
		dark colors are produced by using black ink instead of the combination of 
		cyan, magenta, and yellow. `)
	text.Style.Color = blueColor
	text.Style.FontSize = 14

	ch.Add(p)

	curve := c.NewPolyBezierCurve([]draw.CubicBezierCurve{
		draw.NewCubicBezierCurve(250, 600, 278, 584, 305, 610, 300, 640), // top right
		draw.NewCubicBezierCurve(300, 640, 300, 680, 279, 720, 250, 700), // bottom right
		draw.NewCubicBezierCurve(250, 700, 221, 720, 200, 680, 200, 640), // bottom left
		draw.NewCubicBezierCurve(200, 640, 195, 610, 222, 584, 250, 600), // top left
		draw.NewCubicBezierCurve(250, 600, 246, 588, 242, 560, 250, 550), // leaf
	})

	curve.SetBorderColor(creator.ColorCMYKFromArithmetic(1.0, 1.0, 0.0, 0.0))
	curve.SetFillColor(creator.ColorCMYKFromArithmetic(1.0, 0.0, 1.0, 0.0))
	curve.SetBorderWidth(2)

	c.Draw(ch)
	c.Draw(curve)
}
