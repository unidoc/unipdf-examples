/*
 * Using RGB and CMYK color to colorize text.
 *
 * Run as: go run pdf_text_color.go
 */

package main

import (
	"fmt"
	"math"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
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
	c := creator.New()
	c.SetPageMargins(50, 50, 100, 70)

	colorRGB(c)
	colorCMYK(c)

	if err := c.WriteToFile("color.pdf"); err != nil {
		fmt.Printf("Error %v", err)
	}
}

func drawTextRGB(c *creator.Creator, chapter *creator.Chapter,
	colorName string, color creator.Color) {
	r, g, b := color.ToRGB()

	p := c.NewStyledParagraph()
	tc := p.SetText(fmt.Sprintf("%s color => R: %.1f; G: %.1f; B: %.1f", colorName, r, g, b))
	tc.Style.Color = color
	tc.Style.FontSize = 20

	p.SetMargins(20, 0, 10, 0)

	chapter.Add(p)
}

func colorRGB(c *creator.Creator) {
	// Color from one of predefined colors.
	black := creator.ColorBlack

	// Define RGB color
	red := creator.ColorRGBFrom8bit(255, 0, 0)
	green := creator.ColorRGBFromArithmetic(0.0, 1.0, 0.0)
	blue := creator.ColorRGBFromHex("#0000ff")

	ch := c.NewChapter("Text Color RGB")
	ch.GetHeading().SetColor(black)

	drawTextRGB(c, ch, "Red", red)
	drawTextRGB(c, ch, "Green", green)
	drawTextRGB(c, ch, "Blue", blue)

	c.Draw(ch)
}

func colorToCMYK(color creator.Color) (c, m, y, k float64) {
	r, g, b := color.ToRGB()

	k = 1 - math.Max(math.Max(r, g), b)
	c = (1 - r - k) / (1 - k)
	m = (1 - g - k) / (1 - k)
	y = (1 - b - k) / (1 - k)

	return
}

func drawTextCMYK(cr *creator.Creator, chapter *creator.Chapter,
	colorName string, color creator.Color) {
	c, m, y, k := colorToCMYK(color)

	p := cr.NewStyledParagraph()
	tc := p.SetText(fmt.Sprintf("%s color => C: %.1f; M: %.1f; Y: %.1f; K: %.1f", colorName, c, m, y, k))
	tc.Style.Color = color
	tc.Style.FontSize = 20

	p.SetMargins(20, 0, 10, 0)

	chapter.Add(p)
}

func colorCMYK(c *creator.Creator) {
	// Color from one of predefined colors.
	black := creator.ColorBlack

	// Define CMYK color
	red := creator.ColorCMYKFrom8bit(0, 100, 100, 0)
	green := creator.ColorCMYKFromArithmetic(1.0, 0.0, 1.0, 0.0)
	blue := creator.ColorCMYKFromArithmetic(1.0, 1.0, 0.0, 0.0)

	ch := c.NewChapter("Text Color CMYK")
	ch.GetHeading().SetColor(black)
	ch.SetMargins(0, 0, 50, 0)

	drawTextCMYK(c, ch, "Red", red)
	drawTextCMYK(c, ch, "Green", green)
	drawTextCMYK(c, ch, "Blue", blue)

	c.Draw(ch)
}
