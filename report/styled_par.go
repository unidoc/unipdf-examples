/*
 * This example showcases the creation of styled paragraph
 * The output is saved as styled_paragraph.pdf which illustrates some of the features
 * of the creator.
 */
/*
 * NOTE: This example depends on github.com/unidoc/unipdf/v3/creator, MIT licensed,
 *       and github.com/unidoc/unipdf/v3/model,
 *       Apache-2 licensed.
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

	fontRegular, _ := model.NewStandard14Font(model.HelveticaName)
	fontBold, _ := model.NewStandard14Font(model.HelveticaBoldName)

	//Styled paragraphs ***********************************************
	//Setting heading style
	stChap := c.NewChapter("Styled Paragraphs")
	stChap.GetHeading().SetMargins(0, 0, 20, 0)
	stChap.GetHeading().SetFont(fontBold)
	stChap.GetHeading().SetFontSize(18)
	stChap.GetHeading().SetColor(creator.ColorRed)

	//Creatign new styled paragraph
	stylPar := c.NewStyledParagraph()
	stylPar.SetLineHeight(3) //setting line height

	boldStyle := c.NewTextStyle() //creating new style

	boldStyle.Font = fontBold                                                                                       //setting chunk's font style
	boldStyle.Color = creator.ColorGreen                                                                            //setting chunk's color
	chunk := stylPar.Append("This text is bolded and is in green color. We're showing how styled paragraphs work.") //setting chunk's text
	chunk.Style = boldStyle                                                                                         //setting chunk's style

	normStyle := c.NewTextStyle()                                                                                                                                               //creating new style
	normStyle.Font = fontRegular                                                                                                                                                //setting chunk's font style
	normStyle.Color = creator.ColorBlue                                                                                                                                         //setting chunk's color
	chunktwo := stylPar.Append("You can change the size, color and almost anything of the font using the StyledParagraph command. This font is in blue color and is not bold.") //setting chunk's text
	chunktwo.Style = normStyle                                                                                                                                                  //setting chunk's style

	hugeStyle := c.NewTextStyle()
	hugeStyle.Font = fontRegular
	hugeStyle.FontSize = 25
	chunkthree := stylPar.Append("This is HUGE and black.")
	chunkthree.Style = hugeStyle

	stChap.Add(stylPar) //Adding styled paragraph into the chapter

	c.Draw(stChap) //Drawing the chapter using the creator in the new page we creator

	// Write output file.
	err := c.WriteToFile("styled_paragraph.pdf")
	if err != nil {
		log.Fatalf("Error %s", err) //Writing error if any
	}
}
