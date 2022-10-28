/*
 * This example showcases the capabilities of creating a custom table of contents layout.
 *
 * Run as: go run pdf_custom_toc.go
 */

package main

import (
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
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

	// Create customized table of contents.
	c.AddTOC = true
	c.CustomTOC = true
	c.CreateTableOfContents(func(toc *creator.TOC) error {
		tocTitle := c.NewParagraph("Table of Contents")
		tocTitle.SetFontSize(20)
		tocTitle.SetMargins(0, 0, 0, 20)

		err := c.Draw(tocTitle)
		if err != nil {
			common.Log.Debug("Error drawing table of content's title: %v", err)
			return err
		}

		tocTable := c.NewTable(3)
		tocTable.SetColumnWidths(0.05, 0.85, 0.1)

		for _, tocLine := range toc.Lines() {
			tocLine.Page.Style.FontSize = 15
			pageLink, x, y := tocLine.Link()

			annotation := model.NewPdfAnnotationLink()

			// Set border style.
			bs := model.NewBorderStyle()
			bs.SetBorderWidth(0)
			annotation.BS = bs.ToPdfObject()

			annotation.Dest = core.MakeArray(
				core.MakeInteger(pageLink),
				core.MakeName("XYZ"),
				core.MakeFloat(x),
				core.MakeFloat(y),
				core.MakeFloat(0),
			)

			drawCell(c, tocTable, tocLine.Number.Text, annotation, creator.CellHorizontalAlignmentLeft, creator.CellVerticalAlignmentTop)
			drawTitleCell(c, tocTable, tocLine.Title.Text, annotation)
			drawCell(c, tocTable, tocLine.Page.Text, annotation, creator.CellHorizontalAlignmentRight, creator.CellVerticalAlignmentBottom)
		}

		err = c.Draw(tocTable)
		if err != nil {
			common.Log.Debug("Error drawing table of content's content: %v", err)
			return err
		}

		return nil
	})

	c.NewPage()
	chapter := c.NewChapter("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")
	err := c.Draw(chapter)
	if err != nil {
		common.Log.Error(err.Error())
	}

	c.NewPage()
	chapter = c.NewChapter("At iaculis dignissim curae fusce nam viverra libero, natoque tincidunt mi vitae neque fermentum maecenas, lobortis nisl duis sem montes turpis.")
	err = c.Draw(chapter)
	if err != nil {
		common.Log.Error(err.Error())
	}

	err = c.WriteToFile("pdf-custom-toc.pdf")
	if err != nil {
		log.Fatal(err)
	}
}

func drawCell(c *creator.Creator, tocTable *creator.Table, text string, annotation *model.PdfAnnotationLink, align creator.CellHorizontalAlignment, valign creator.CellVerticalAlignment) {
	p := c.NewStyledParagraph()
	chunk := p.Append(text)
	chunk.Style.Color = creator.ColorBlue
	chunk.SetAnnotation(annotation.PdfAnnotation)

	cell := tocTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 0)
	cell.SetHorizontalAlignment(align)
	cell.SetVerticalAlignment(valign)
	cell.SetContent(p)
}

func drawTitleCell(c *creator.Creator, tocTable *creator.Table, title string, annotation *model.PdfAnnotationLink) {
	p := c.NewStyledParagraph()
	chunk := p.Append(title)
	chunk.Style.Color = creator.ColorBlue
	chunk.SetAnnotation(annotation.PdfAnnotation)

	cell := tocTable.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 0)
	cell.SetVerticalAlignment(creator.CellVerticalAlignmentBottom)
	cell.SetContent(p)
}
