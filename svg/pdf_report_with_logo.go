/*
 * This example code shows how to use an svg image when creating
 * PDF reports using unipdf
 *
 * Run as: go run add_svg.go
 */
package main

import (
	"fmt"
	"os"
	"time"

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
	err := generatePDFReport("unidoc-report_with_svg.pdf")
	if err != nil {
		panic(err)
	}
}

func generatePDFReport(outputPath string) error {
	robotoFontRegular, err := model.NewPdfFontFromTTFFile("./fonts/Roboto-Regular.ttf")
	if err != nil {
		return err
	}

	robotoFontPro, err := model.NewPdfFontFromTTFFile("./fonts/Roboto-Bold.ttf")
	if err != nil {
		return err
	}

	c := creator.New()
	c.SetPageMargins(50, 50, 100, 70)

	file := "./svgs/unidoc-logo.svg"
	graphicSvg, err := creator.NewGraphicSVGFromFile(file)
	if err != nil {
		panic(err)
	}

	graphicSvg.SetPos(58, 20)

	doFirstHeader(c, robotoFontRegular, robotoFontPro)

	// Setup a front page (always placed first).
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		createCoverPage(c, robotoFontRegular, robotoFontPro)
	})

	// Draw a header on each page.
	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		// Draw the header on a block. The block size is the size of the page's top margins.
		block.Draw(graphicSvg)
	})

	// Draw footer on each page.
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		// Draw the on a block for each page.
		p := c.NewParagraph("unidoc.io")
		p.SetFont(robotoFontRegular)
		p.SetFontSize(8)
		p.SetPos(50, 20)
		p.SetColor(creator.ColorRGBFrom8bit(63, 68, 76))
		block.Draw(p)

		strPage := fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages)
		p = c.NewParagraph(strPage)
		p.SetFont(robotoFontRegular)
		p.SetFontSize(8)
		p.SetPos(300, 20)
		p.SetColor(creator.ColorRGBFrom8bit(63, 68, 76))
		block.Draw(p)
	})

	err = c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}

// createCoverPage generates the front page.
func createCoverPage(c *creator.Creator, fontRegular *model.PdfFont, fontBold *model.PdfFont) {
	helvetica, _ := model.NewStandard14Font("Helvetica")
	helveticaBold, _ := model.NewStandard14Font("Helvetica-Bold")

	p := c.NewParagraph("UniDoc")
	p.SetFont(helvetica)
	p.SetFontSize(48)
	p.SetMargins(85, 0, 150, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)

	p = c.NewParagraph("Example Report")
	p.SetFont(helveticaBold)
	p.SetFontSize(30)
	p.SetMargins(85, 0, 0, 0)
	p.SetColor(creator.ColorRGBFrom8bit(45, 148, 215))
	c.Draw(p)

	t := time.Now().UTC()
	dateStr := t.Format("01 Jan, 2006 15:04")

	p = c.NewParagraph(dateStr)
	p.SetFont(helveticaBold)
	p.SetFontSize(12)
	p.SetMargins(90, 0, 5, 0)
	p.SetColor(creator.ColorRGBFrom8bit(56, 68, 77))
	c.Draw(p)
}

// doFirstHeader creates a Chapter with one header and and a sub header.
func doFirstHeader(c *creator.Creator, fontRegular *model.PdfFont, fontBold *model.PdfFont) {
	// Ensure that the chapter starts on a new page.
	c.NewPage()

	ch := c.NewChapter("Sample Header")

	chapterFont := fontRegular
	chapterFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	chapterFontSize := 18.0

	normalFont := fontRegular
	normalFontColor := creator.ColorRGBFrom8bit(72, 86, 95)
	normalFontSize := 10.0

	ch.GetHeading().SetFont(chapterFont)
	ch.GetHeading().SetFontSize(chapterFontSize)
	ch.GetHeading().SetColor(chapterFontColor)

	p := c.NewParagraph("This is an example sentence showcasing the content of the first header. It provides a brief introduction to the section, highlighting the key points that will be discussed under the header.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	ch.Add(p)

	// Paragraphs.
	sc := ch.NewSubchapter("Sub Header")
	sc.GetHeading().SetMargins(0, 0, 20, 0)
	sc.GetHeading().SetFont(chapterFont)
	sc.GetHeading().SetFontSize(chapterFontSize)
	sc.GetHeading().SetColor(chapterFontColor)

	p = c.NewParagraph("Paragraphs are used to represent text, as little as a single character, a word or " +
		"multiple words forming multiple sentences. UniDoc handles automatically wrapping those across lines and pages, making " +
		"it relatively easy to work with. They can also be left, center, right aligned or justified as illustrated below:")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	// Example paragraphs:
	p = c.NewParagraph("This paragraph shows the first paragraph of the subheader. It introduces the topic that will be " +
		"discussed in the subsequent sections. The purpose of this introduction is to set the context for " +
		"the reader. By doing so, it ensures a smooth transition into the more detailed points. The reader " +
		"can expect an overview of the important aspects related to the subheader.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	p = c.NewParagraph("This paragraph demonstrates the second paragraph of the subheader. It delves deeper into the " +
		"specific details of the topic introduced earlier. Here, more focus is given to explaining key " +
		"concepts or arguments that are central to the subject. The paragraph aims to elaborate on " +
		"important points while keeping the reader engaged. By offering clear explanations, it supports the " +
		"overall understanding of the subject matter.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	p = c.NewParagraph("This is the third paragraph of the subheader. It provides a conclusion or a summary of the key " +
		"takeaways discussed in the previous paragraphs. This section often reinforces the main message " +
		"or insights of the topic. By wrapping up the discussion, it offers closure and emphasizes the most " +
		"important points. The paragraph also paves the way for the reader to reflect on the subject and " +
		"prepare for what might come next.")
	p.SetFont(normalFont)
	p.SetFontSize(normalFontSize)
	p.SetColor(normalFontColor)
	p.SetMargins(0, 0, 5, 0)
	sc.Add(p)

	c.Draw(ch)
}
