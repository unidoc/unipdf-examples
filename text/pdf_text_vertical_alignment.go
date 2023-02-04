/* Setting the vertical alignment of text chunks in a paragraph.
 *
 * Run as: go run pdf_text_vertical_alignment.go
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
	ctor := creator.New()

	textWithIndicator(ctor)
	basicTextRenderingResult(ctor)

	if err := ctor.WriteToFile("pdf_text_vertical_alignment.pdf"); err != nil {
		log.Fatalf("failed to write pdf: %v", err)
	}
}

func textWithIndicator(ctor *creator.Creator) error {
	ch := ctor.NewChapter("Text with Indicator")
	ch.SetMargins(0, 0, 10, 10)

	line := ctor.NewLine(50, 200, 550, 200)
	line.SetColor(creator.ColorRed)
	line.SetLineWidth(1)

	if err := ch.Add(line); err != nil {
		log.Fatalf("failed to add line: %v", err)
	}

	sp := ctor.NewStyledParagraph()
	sp.SetPos(50, 200-30)

	c := sp.SetText("Vertically aligned on baseline (default)")
	c.Style.FontSize = 30

	if err := ch.Add(sp); err != nil {
		log.Fatalf("failed to add sp: %v", err)
	}

	line = ctor.NewLine(50, 300, 550, 300)
	line.SetColor(creator.ColorRed)
	line.SetLineWidth(1)

	if err := ch.Add(line); err != nil {
		log.Fatalf("failed to add line: %v", err)
	}

	sp = ctor.NewStyledParagraph()
	sp.SetPos(50, 300-30)

	c = sp.SetText("Vertically aligned on center")
	c.Style.FontSize = 30
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	if err := ch.Add(sp); err != nil {
		log.Fatalf("failed to add sp: %v", err)
	}

	line = ctor.NewLine(50, 400, 550, 400)
	line.SetColor(creator.ColorRed)
	line.SetLineWidth(1)

	if err := ch.Add(line); err != nil {
		log.Fatalf("failed to add line: %v", err)
	}

	sp = ctor.NewStyledParagraph()
	sp.SetPos(50, 400-30)

	c = sp.SetText("Vertically aligned below baseline")
	c.Style.FontSize = 30
	c.VerticalAlignment = creator.TextVerticalAlignmentBottom

	if err := ch.Add(sp); err != nil {
		log.Fatalf("failed to add sp: %v", err)
	}

	line = ctor.NewLine(50, 500, 550, 500)
	line.SetColor(creator.ColorRed)
	line.SetLineWidth(1)

	if err := ch.Add(line); err != nil {
		log.Fatalf("failed to add line: %v", err)
	}

	sp = ctor.NewStyledParagraph()
	sp.SetPos(50, 500-60)

	c = sp.Append("Hello")
	c.Style.FontSize = 42

	c = sp.Append("World")
	c.Style.FontSize = 10
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	c = sp.Append("Hello")
	c.Style.FontSize = 42
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	c = sp.Append("World")
	c.Style.FontSize = 60

	c = sp.Append("Hello")
	c.Style.FontSize = 10
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	c = sp.Append("World")
	c.Style.FontSize = 42
	c.VerticalAlignment = creator.TextVerticalAlignmentBottom

	if err := ch.Add(sp); err != nil {
		log.Fatalf("failed to add sp: %v", err)
	}

	return ctor.Draw(ch)

}

func basicTextRenderingResult(ctor *creator.Creator) error {
	fontRegular, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		log.Fatalf("Failed initiating font: %v", err)
	}

	fontCourier, err := model.NewStandard14Font(model.CourierBoldName)
	if err != nil {
		log.Fatalf("Failed initiating font: %v", err)
	}

	fontTimes, err := model.NewStandard14Font(model.TimesBoldItalicName)
	if err != nil {
		log.Fatalf("Failed initiating font: %v", err)
	}

	ctor.NewPage()

	ch := ctor.NewChapter("Basic Text Rendering")
	ch.SetMargins(0, 0, 10, 10)

	sp := ctor.NewStyledParagraph()

	c := sp.Append("Hello")
	c.Style.FontSize = 40
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	c = sp.Append("World")
	c.Style.FontSize = 10
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	if err := ch.Add(sp); err != nil {
		log.Fatalf("failed to add sp: %v", err)
	}

	sp = ctor.NewStyledParagraph()
	sp.SetMargins(0, 0, 50, 0)

	c = sp.Append("Lorem ipsum dolor sit amet")
	c.Style.FontSize = 40

	c = sp.Append("consectetur adipiscing elit,")
	c.Style.Font = fontCourier
	c.Style.FontSize = 10
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	c = sp.Append("urna consequat felis vehicula")
	c.Style.FontSize = 20
	c.VerticalAlignment = creator.TextVerticalAlignmentBottom

	c = sp.Append("class ultricies mollis dictumst,")
	c.Style.Font = fontCourier
	c.Style.FontSize = 30
	c.VerticalAlignment = creator.TextVerticalAlignmentBaseline

	c = sp.Append("aenean non a in donec nulla.")
	c.Style.Font = fontRegular
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	c = sp.Append(" Phasellus ante pellentesque erat cum risus consequat imperdiet aliquam,")
	c.Style.FontSize = 14

	c = sp.Append(" integer placerat et turpis mi eros nec lobortis tacit")
	c.Style.Font = fontTimes
	c.Style.FontSize = 18
	c.VerticalAlignment = creator.TextVerticalAlignmentCenter

	if err := ch.Add(sp); err != nil {
		log.Fatalf("failed to add sp: %v", err)
	}

	return ctor.Draw(ch)
}
