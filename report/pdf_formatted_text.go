/*
 * This example showcases the creation of styled paragraph.
 * The output is saved as styled_paragraph.pdf which illustrates some of the features
 * of the creator.
 */

package main

import (
	"log"

	"github.com/unidoc/unipdf/v3/common/license"
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
	// common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	c := creator.New()

	fontRegular, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fontBold, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create styled paragraph chapter.
	chap := c.NewChapter("Styled Paragraphs")
	chap.GetHeading().SetMargins(0, 0, 0, 20)
	chap.GetHeading().SetFont(fontBold)
	chap.GetHeading().SetFontSize(18)
	chap.GetHeading().SetColor(creator.ColorRed)

	// Generate styled paragraph text style subchapter.
	err = styledParagraphTextStyle(c, chap, fontRegular, fontBold)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Generate styled paragraph text underline subchapter.
	err = styledParagraphUnderline(c, chap, fontRegular, fontBold)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Draw styled paragraph chapter.
	if err = c.Draw(chap); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	err = c.WriteToFile("styled_paragraph.pdf")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func styledParagraphTextStyle(c *creator.Creator, ch *creator.Chapter,
	fontRegular, fontBold *model.PdfFont) error {
	// Create new subchapter.
	subchap := ch.NewSubchapter("Text styles")

	// Create a new styled paragraph.
	p := c.NewStyledParagraph()
	p.SetMargins(0, 0, 20, 20)
	p.SetLineHeight(1.1)

	// Change individual text style properties.
	chunk := p.Append("Styled paragraphs are fully customizable.")
	chunk.Style.Font = fontRegular

	p.Append("You can change the ")
	chunk = p.Append("color")
	chunk.Style.Color = creator.ColorRed

	p.Append(", the ")
	chunk = p.Append("font ")
	chunk.Style.Font = fontBold

	p.Append("and the ")
	chunk = p.Append("font size ")
	chunk.Style.FontSize = 14
	p.Append("of text chunks. ")

	// Change text font size.
	chunk = p.Append("Let's draw some text with a larger font size. ")
	chunk.Style.FontSize = 16

	// Assign custom style to a text chunk.
	boldStyle := c.NewTextStyle()
	boldStyle.Font = fontBold
	boldStyle.Color = creator.ColorBlue

	chunk = p.Append("Now some blue bold text in order to showcase what styled paragraphs can do.")
	chunk.Style = boldStyle

	// Add the styled paragraph to the created subchapter.
	return subchap.Add(p)
}

func styledParagraphUnderline(c *creator.Creator, ch *creator.Chapter,
	fontRegular, fontBold *model.PdfFont) error {
	// Create new subchapter.
	subchap := ch.NewSubchapter("Text underline")

	// Create a new styled paragraph.
	p := c.NewStyledParagraph()
	p.SetMargins(0, 0, 20, 20)
	p.SetLineHeight(1.2)
	p.Append("Text chunks can be ")

	// Default underline style.
	chunk := p.Append("underlined")
	chunk.Style.Underline = true

	p.Append(" using the default style.\n")
	p.Append("By default, the ")

	// Default underline style based on the color of the text.
	chunk = p.Append("underline color")
	chunk.Style.Underline = true
	chunk.Style.Color = creator.ColorBlue

	p.Append(" is the color of the text chunk. We can also add ")

	// Custom underline style color.
	chunk = p.Append("some long underlined text chunk which wraps")
	chunk.Style.FontSize = 15
	chunk.Style.Underline = true
	chunk.Style.UnderlineStyle.Color = creator.ColorRed
	chunk.Style.UnderlineStyle.Offset = 1

	p.Append(" and then some more regular text.\n")

	// Custom underline thickness and offset.
	p.Append("Finally, we can customize the offset and the thickness of the ")

	chunk = p.Append("underlined text")
	chunk.Style.Underline = true
	chunk.Style.UnderlineStyle.Thickness = 2
	chunk.Style.UnderlineStyle.Offset = 2

	p.Append(".")

	// Add the styled paragraph to the created subchapter.
	return subchap.Add(p)
}
