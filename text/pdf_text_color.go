/*
 * Using RGB and CMYK color to colorize text.
 *
 * Run as: go run pdf_text_color.go
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
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

func colorRGB(c *creator.Creator) {
	// Color from one of predefined colors.
	black := creator.ColorBlack

	// Define RGB color
	red := creator.ColorRGBFrom8bit(255, 0, 0)
	green := creator.ColorRGBFromArithmetic(0.0, 1.0, 0.0)
	blue := creator.ColorRGBFromHex("#0000ff")
	outlineColor := creator.ColorRGBFromHex("#283A3F")

	ch := c.NewChapter("Text Color RGB")
	ch.GetHeading().SetFontColor(black)

	// Red colored text.
	p := c.NewStyledParagraph()
	tc := p.SetText("Red color")
	tc.Style.Color = red
	tc.Style.FontSize = 50

	p.SetMargins(20, 0, 10, 0)

	ch.Add(p)

	// Green colored text.
	p = c.NewStyledParagraph()
	tc = p.SetText("Green color with underline")
	tc.Style.Color = green
	tc.Style.FontSize = 50
	tc.Style.Underline = true
	tc.Style.UnderlineStyle.Color = outlineColor
	tc.Style.UnderlineStyle.Thickness = 2

	p.SetMargins(20, 0, 10, 0)

	ch.Add(p)

	// Blue colored text.
	p = c.NewStyledParagraph()
	tc = p.SetText("Blue color with outline")
	tc.Style.Color = blue
	tc.Style.FontSize = 50
	tc.Style.OutlineColor = outlineColor
	tc.Style.OutlineSize = 2
	tc.Style.RenderingMode = creator.TextRenderingModeFillStroke

	p.SetMargins(20, 0, 10, 0)

	ch.Add(p)

	c.Draw(ch)
}

func colorCMYK(c *creator.Creator) {
	// Color from one of predefined colors.
	black := creator.ColorBlack

	// Define CMYK color
	cyan := creator.ColorCMYKFrom8bit(100, 0, 0, 0)
	magenta := creator.ColorCMYKFromArithmetic(0.0, 1.0, 0.0, 0.0)
	outlineColor := creator.ColorCMYKFrom8bit(37, 8, 0, 75)

	ch := c.NewChapter("Text Color CMYK")
	ch.GetHeading().SetFontColor(black)
	ch.SetMargins(0, 0, 50, 0)

	p := c.NewStyledParagraph()
	tc := p.SetText("Cyan color")
	tc.Style.Color = cyan
	tc.Style.FontSize = 50

	p.SetMargins(20, 0, 10, 0)

	ch.Add(p)

	p = c.NewStyledParagraph()
	tc = p.SetText("Magenta color with outline")
	tc.Style.Color = magenta
	tc.Style.FontSize = 50
	tc.Style.OutlineColor = outlineColor
	tc.Style.OutlineSize = 2
	tc.Style.RenderingMode = creator.TextRenderingModeFillStroke

	p.SetMargins(20, 0, 10, 0)

	ch.Add(p)

	c.Draw(ch)
}
