/*
 * This example showcases the creation of styled paragraph.
 * The output is saved as styled_paragraph.pdf which illustrates some of the features
 * of the creator.
 */

package main

import (
	"log"

	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

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
	stChap := c.NewChapter("Styled Paragraphs")
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
	c.Draw(stChap)

	// Write output file.
	err = c.WriteToFile("styled_paragraph.pdf")
	if err != nil {
		log.Fatalf("Error %s", err)
	}
}
