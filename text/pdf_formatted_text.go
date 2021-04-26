/*
 * This example showcases the creation of styled paragraph.
 * The output is saved as styled_paragraph.pdf which illustrates some of the features
 * of the creator.
 */

package main

import (
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
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
	chap.GetHeading().SetFont(fontBold)
	chap.GetHeading().SetFontSize(18)
	chap.GetHeading().SetColor(creator.ColorRed)

	// Generate styled paragraph text style subchapter.
	err = styledParagraphTextStyle(c, chap, fontRegular, fontBold)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Generate styled paragraph text wrapping and overflow subchapter.
	err = styledParagraphTextOverflow(c, chap, fontRegular, fontBold)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Generate styled paragraph text underline subchapter.
	err = styledParagraphUnderline(c, chap, fontRegular, fontBold)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Generate styled paragraph subscript and superscript subchapter.
	err = styledParagraphScript(c, chap, fontRegular, fontBold)
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

func styledParagraphTextOverflow(c *creator.Creator, ch *creator.Chapter,
	fontRegular, fontBold *model.PdfFont) error {
	// Create new subchapter.
	subchap := ch.NewSubchapter("Text wrapping and overflow")

	// Text overflow disabled and wrapped (default setting).
	p := c.NewStyledParagraph()
	p.SetMargins(0, 0, 20, 20)
	p.SetLineHeight(1.1)
	chunk := p.Append("Long styled paragraph will be wrapped by default if it doesn't fit in the available space, like this one for example.")
	chunk.Style.Font = fontRegular

	if err := subchap.Add(p); err != nil {
		return err
	}

	// Text wrapping disabled.
	p = c.NewStyledParagraph()
	p.SetMargins(0, 0, 0, 20)
	p.SetLineHeight(1.1)
	p.SetEnableWrap(false)
	chunk = p.Append("This long styled paragraph has it's wrap setting disabled, so it will NOT be wrapped even if it doesn't fit in the available space and will overflows out of the page.")
	chunk.Style.Font = fontRegular

	if err := subchap.Add(p); err != nil {
		return err
	}

	// Text wrapping disabled and text overflow set to hidden.
	p = c.NewStyledParagraph()
	p.SetMargins(0, 0, 0, 20)
	p.SetLineHeight(1.1)
	p.SetEnableWrap(false)
	p.SetTextOverflow(creator.TextOverflowHidden)
	chunk = p.Append("This long styled paragraph has it's wrap setting disabled and text overflow set to hidden, so it will be truncated if it doesn't fit in the available space.")
	chunk.Style.Font = fontRegular

	if err := subchap.Add(p); err != nil {
		return err
	}

	return nil
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

func styledParagraphScript(c *creator.Creator, ch *creator.Chapter,
	fontRegular, fontBold *model.PdfFont) error {
	// Create new subchapter.
	subchap := ch.NewSubchapter("Subscript & Superscript")
	subchap.GetHeading().SetMargins(0, 0, 0, 10)

	// Generate styled paragraph subscript and superscript subchapter.
	err := styledParagraphScriptBasic(c, subchap, fontRegular, fontBold)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Generate styled paragraph subscript and superscript subchapter
	// combined with underlined and annotated text.
	err = styledParagraphScriptCombination(c, subchap, fontRegular, fontBold)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return err
}

func styledParagraphScriptBasic(c *creator.Creator, ch *creator.Chapter,
	fontRegular, fontBold *model.PdfFont) error {
	// Create new subchapter.
	subchap := ch.NewSubchapter("Basic Example")
	subchap.GetHeading().SetMargins(10, 0, 0, 10)

	// Basic example.
	p := c.NewStyledParagraph()
	p.SetMargins(10, 0, 0, 20)
	p.SetLineHeight(1.2)
	p.Append("Styled paragraphs allow drawing ")

	style := &p.Append("subscript").Style
	style.TextRise = -8
	style.FontSize = 7
	style.Color = creator.ColorRed

	p.Append(" and ")

	style = &p.Append("superscript").Style
	style.TextRise = 9
	style.FontSize = 7
	style.Color = creator.ColorBlue

	p.Append(" text chunks.")

	// Add the styled paragraph to the created subchapter.
	return subchap.Add(p)
}

func styledParagraphScriptCombination(c *creator.Creator, ch *creator.Chapter,
	fontRegular, fontBold *model.PdfFont) error {
	// Create new subchapter.
	subchap := ch.NewSubchapter("Combined with Underlined and Annotated Text")
	subchap.GetHeading().SetMargins(10, 0, 0, 10)

	// Text rise combined with underlined and annotated text.
	p := c.NewStyledParagraph()
	p.SetMargins(10, 0, 0, 20)
	p.SetLineHeight(1.2)
	p.Append("Subscript and superscript text can be ")

	style := &p.Append("underlined").Style
	style.TextRise = 5
	style.FontSize = 7
	style.Color = creator.ColorGreen
	style.Underline = true

	p.Append(" or turned into ")

	style = &p.AddExternalLink("link", "https://google.com").Style
	style.TextRise = -5
	style.FontSize = 7

	p.Append(" annotations.")

	// Add the styled paragraph to the created subchapter.
	return subchap.Add(p)
}
