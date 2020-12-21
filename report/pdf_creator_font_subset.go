/*
 * This example showcases uses of EnableFontSubsetting flag in creator package.
 * This setting will embed the runes/glyphs that are used in the document
 * and in turn would reduce the document size by a lot.
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
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	c := creator.New()
	c.NewPage()

	fontRegular, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	fontBold, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	// `stChap` represents a styled paragraphs chapter
	stChap := c.NewChapter("Styled Paragraphs with Font Subsetting")
	stChap.GetHeading().SetMargins(0, 0, 20, 0)
	stChap.GetHeading().SetFont(fontBold)
	stChap.GetHeading().SetFontSize(18)
	stChap.GetHeading().SetColor(creator.ColorRed)

	// `stylPar` creates a new styled paragraph
	stylPar := c.NewStyledParagraph()
	stylPar.SetLineHeight(3)

	// `boldStyle` creates a new style
	boldStyle := c.NewTextStyle()
	boldStyle.Font = fontBold
	boldStyle.Color = creator.ColorGreen

	// Applying `boldStyle` to `chunk`
	chunk := stylPar.Append("This text is bolded and is in green color. We're showing how styled paragraphs work.")
	chunk.Style = boldStyle

	// Creating new style `normStyle`
	normStyle := c.NewTextStyle()
	normStyle.Font = fontRegular
	normStyle.Color = creator.ColorBlue

	// Applying `normStyle` to `chunkTwo`
	chunkTwo := stylPar.Append("You can change the size, color and almost anything of the font using the StyledParagraph command. This font is in blue color and is not bold.")
	chunkTwo.Style = normStyle

	// Creating new style `hugeStyle`
	hugeStyle := c.NewTextStyle()
	hugeStyle.Font = fontRegular
	hugeStyle.FontSize = 25

	// Applying `normStyle` to `chunkThree`
	chunkThree := stylPar.Append("This is HUGE and black.")
	chunkThree.Style = hugeStyle

	// Adding styled paragraph into the chapter `stChap` and drawing it using the creator
	stChap.Add(stylPar)

	// Enable font subsetting for both font.
	c.EnableFontSubsetting(fontRegular)
	c.EnableFontSubsetting(fontBold)

	c.Draw(stChap)

	// Write output file.
	err = c.WriteToFile("font_subset_example.pdf")
	if err != nil {
		log.Fatalf("Error %s", err)
	}
}
