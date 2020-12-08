/*
 * This example showcases row wrapping across pages in creator tables.
 * The output is saved as unipdf-tables-row-wrap.pdf.
 */

package main

import (
	"fmt"
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
	c.SetPageMargins(50, 50, 50, 50)

	headingFont, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Fatal(err)
	}

	if err := rowWrapDisabled(c, headingFont); err != nil {
		log.Fatal(err)
	}

	c.NewPage()

	if err := rowWrapEnabled(c, headingFont); err != nil {
		log.Fatal(err)
	}

	// Write to output file.
	if err := c.WriteToFile("unipdf-tables-row-wrap.pdf"); err != nil {
		log.Fatal(err)
	}
}

func rowWrapDisabled(c *creator.Creator, headingFont *model.PdfFont) error {
	heading := c.NewStyledParagraph()
	chunk := heading.Append("1. Table row wrap disabled")
	chunk.Style.Font = headingFont
	chunk.Style.FontSize = 20

	if err := c.Draw(heading); err != nil {
		return err
	}

	description := c.NewStyledParagraph()
	description.SetMargins(0, 0, 10, 20)
	chunk = description.Append("When table row wrapping is disabled, if one of the cells of a row does not fit in the available space of the current page, the whole row will be moved on the next one.")
	chunk.Style.FontSize = 14
	chunk.Style.Color = creator.ColorRGBFromHex("#777")

	if err := c.Draw(description); err != nil {
		return err
	}

	return fillTable(c, 22, false, func(table *creator.Table) error {
		sp1 := c.NewStyledParagraph()
		sp1.Append("This is a styled paragraph which will not fit on the current page. All its content should be moved on the next page, along with the entire row, as row wrapping is disabled.").Style.FontSize = 14

		p1 := c.NewParagraph("This is a regular paragraph which will fit on the current page. However, it will be moved on the next page.")
		p1.SetFontSize(14)

		sp2 := c.NewStyledParagraph()
		sp2.Append("This is a styled paragraph which will fit on the current page. However, it will be moved on the next page.").Style.FontSize = 14

		p2 := c.NewParagraph("This is a regular paragraph which will not fit on the current page. All its content should be moved on the next page, along with the entire row.")
		p2.SetFontSize(14)

		// Draw table row.
		for _, d := range []creator.VectorDrawable{sp1, p1, sp2, p2} {
			if err := drawCell(table, d); err != nil {
				return err
			}
		}

		return nil
	})
}

func rowWrapEnabled(c *creator.Creator, headingFont *model.PdfFont) error {
	heading := c.NewStyledParagraph()
	chunk := heading.Append("2. Table row wrap enabled")
	chunk.Style.Font = headingFont
	chunk.Style.FontSize = 20

	if err := c.Draw(heading); err != nil {
		return err
	}

	description := c.NewStyledParagraph()
	description.SetMargins(0, 0, 10, 20)
	description.SetLineHeight(1.1)
	chunk = description.Append("When table row wrapping is enabled, cells which contain styled paragraphs that don't fit in the available space of the current page will wrap across pages.\n")
	chunk.Style.FontSize = 14
	chunk.Style.Color = creator.ColorRGBFromHex("#777")
	chunk = description.Append("Other components behave as usual. If they fit in the available space, they remain there. Otherwise, they are placed on the next page.\n")
	chunk.Style.FontSize = 14
	chunk.Style.Color = creator.ColorRGBFromHex("#777")

	if err := c.Draw(description); err != nil {
		return err
	}

	return fillTable(c, 22, true, func(table *creator.Table) error {
		sp1 := c.NewStyledParagraph()
		sp1.Append("This is a styled paragraph. When table row wrapping is enabled, the content that fits on the current page stays on the current page. The rest of the content will be placed on the next page.").Style.FontSize = 14

		p1 := c.NewParagraph("This is a regular paragraph which will fit on the current page. All its content should remain on the current page.")
		p1.SetFontSize(14)

		sp2 := c.NewStyledParagraph()
		sp2.Append("This is a styled paragraph which will fit on the current page. All its content should remain on the current page.").Style.FontSize = 14

		p2 := c.NewParagraph("This is a regular paragraph which will not fit on the current page. All its content should be moved on the next page, in the wrapped row, leaving the current cell empty.")
		p2.SetFontSize(14)

		// Draw table row.
		for _, d := range []creator.VectorDrawable{sp1, p1, sp2, p2} {
			if err := drawCell(table, d); err != nil {
				return err
			}
		}

		return nil
	})
}

func drawCell(table *creator.Table, content creator.VectorDrawable) error {
	cell := table.NewCell()
	cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
	return cell.SetContent(content)
}

func fillTable(c *creator.Creator, lineCount int, enableRowWrap bool,
	addExtraRows func(table *creator.Table) error) error {
	// Create table.
	table := c.NewTable(4)
	table.SetMargins(0, 0, 10, 0)
	table.EnableRowWrap(enableRowWrap)

	for i := 0; i < lineCount; i++ {
		for j := 0; j < 4; j++ {
			sp := c.NewStyledParagraph()
			chunk := sp.Append(fmt.Sprintf("Row %d - Cell %d", i+1, j+1))
			chunk.Style.FontSize = 14

			if err := drawCell(table, sp); err != nil {
				return err
			}
		}
	}

	if err := addExtraRows(table); err != nil {
		return err
	}

	// Draw table.
	if err := c.Draw(table); err != nil {
		return err
	}

	return nil
}
