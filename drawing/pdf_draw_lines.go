/*
 * This example showcases the capabilities of creator lines.
 *
 * Run as: go run pdf_draw_lines.go
 */

package main

import (
	"log"
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/contentstream/draw"
	"github.com/unidoc/unipdf/v4/creator"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io.
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
	c := creator.New()

	// Create front page.
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		sp := c.NewStyledParagraph()
		sp.SetTextAlignment(creator.TextAlignmentCenter)
		sp.SetMargins(0, 0, 250, 0)
		style := &sp.Append("Quick guide to drawing lines \nusing the creator package.").Style
		style.FontSize = 24

		if err := c.Draw(sp); err != nil {
			log.Fatal(err)
		}
	})

	// Create customized table of contents.
	c.AddTOC = true
	c.CreateTableOfContents(func(toc *creator.TOC) error {
		// Set style of TOC heading just before render.
		style := c.NewTextStyle()
		style.FontSize = 18

		toc.SetHeading("Table of Contents", style)

		// Set style of TOC lines just before render.
		lines := toc.Lines()
		for _, line := range lines {
			line.Page.Style.FontSize = 15
		}

		return nil
	})

	// Showcase lines using absolute positioning.
	if err := drawLinesPositionAbsolute(c); err != nil {
		log.Fatal(err)
	}

	// Showcase lines using relative positioning.
	if err := drawLinesPositionRelative(c); err != nil {
		log.Fatal(err)
	}

	// Showcase lines inside divisions.
	if err := drawLinesInsideDivision(c); err != nil {
		log.Fatal(err)
	}

	// Showcase lines inside tables.
	if err := drawLinesInsideTable(c); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-draw-lines.pdf"); err != nil {
		log.Fatal(err)
	}
}

func drawLinesPositionAbsolute(c *creator.Creator) error {
	c.NewPage()

	// Draw chapter.
	subtitleText := "In absolute positioning mode, the absolute coordinates of the lines must be specified by the user. Let's draw some lines with different styles."

	chapter := c.NewChapter("Lines using absolute positioning")
	subtitle := newParagraph(c, subtitleText, 15, creator.TextAlignmentJustify, newMargins(0, 0, 8, 15))
	if err := chapter.Add(subtitle); err != nil {
		return err
	}
	if err := c.Draw(chapter); err != nil {
		return err
	}

	// Create page content.
	contents := []creator.VectorDrawable{
		newLine(c, 60, 150, 350, 150, true, false, 3, creator.ColorBlack, false, nil, 0, 1.0, newMargins(0, 0, 0, 0)),
		newLine(c, 500, 150, 500, 350, true, false, 10, creator.ColorRed, false, nil, 0, 1.0, newMargins(0, 0, 0, 0)),
		newLine(c, 60, 200, 400, 400, true, false, 2, creator.ColorBlue, true, []int64{6, 1}, 0, 1.0, newMargins(0, 0, 0, 0)),
		newLine(c, 60, 450, 550, 450, true, false, 15, creator.ColorGreen, false, nil, 0, 0.5, newMargins(0, 0, 0, 0)),
		newLine(c, 60, 650, 550, 480, true, false, 3, creator.ColorRed, false, nil, 0, 1.0, newMargins(0, 0, 0, 0)),
		newLine(c, 60, 700, 550, 700, true, false, 10, creator.ColorBlack, true, []int64{10, 2}, 0, 0.3, newMargins(0, 0, 0, 0)),
	}

	// Draw page content.
	for _, content := range contents {
		if err := c.Draw(content); err != nil {
			return err
		}
	}

	return nil
}

func drawLinesPositionRelative(c *creator.Creator) error {
	c.NewPage()

	// Draw chapter.
	subtitleText := "In relative positioning mode, lines are positioned relative to the current context. The specified coordinates are only used to determine the orientation and size of the line. Furthermore, relative lines can be made to fill the entire width available to them. They can also have margins."

	chapter := c.NewChapter("Lines using relative positioning")
	subtitle := newParagraph(c, subtitleText, 15, creator.TextAlignmentJustify, newMargins(0, 0, 8, 15))
	if err := chapter.Add(subtitle); err != nil {
		return err
	}
	if err := c.Draw(chapter); err != nil {
		return err
	}

	// Create page content.
	contents := []creator.VectorDrawable{
		newParagraph(c, "The line is always positioned relative to the current context. First, the size and orientation of the lines are determined from the provided coordinates. Then, the line is translated at the current context position. They can have any orientation.", 12, creator.TextAlignmentJustify, newMargins(0, 0, 10, 5)),
		newLine(c, 0, 0, 200, 0, false, false, 3, creator.ColorBlue, false, nil, 0, 1.0, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 50, false, false, 5, creator.ColorYellow, false, nil, 0, 1.0, newMargins(150, 0, 5, 5)),
		newLine(c, 0, 0, 400, 30, false, false, 5, creator.ColorRed, false, nil, 0, 0.5, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 50, 200, 0, false, false, 2, creator.ColorGreen, true, []int64{7, 2}, 0, 1.0, newMargins(0, 0, 5, 10)),

		newParagraph(c, "In full-width mode, the size of the line is calculated so that the line fills the entire available context width. The provided coordinates are only used to determine the orientation.", 12, creator.TextAlignmentJustify, newMargins(0, 0, 20, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, creator.ColorRed, false, nil, 0, 1.0, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 50, false, true, 5, creator.ColorBlue, true, []int64{10, 3}, 0, 1.0, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 50, 0, 0, false, true, 2, creator.ColorGreen, false, nil, 0, 0.5, newMargins(0, 0, 5, 5)),

		newParagraph(c, "Full-width horizontal lines are very useful as content separators. For this use case, all coordinates can be 0. Let's draw a couple of them, with and without horizontal margins.", 12, creator.TextAlignmentJustify, newMargins(0, 0, 20, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 1, creator.ColorBlue, false, nil, 0, 1.0, newMargins(0, 0, 5, 10)),
		newLine(c, 0, 0, 0, 0, false, true, 2, creator.ColorRed, true, []int64{5, 1}, 0, 1.0, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 10, creator.ColorGreen, false, nil, 0, 0.5, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 4, creator.ColorYellow, true, []int64{2, 4}, 0, 1.0, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, creator.ColorRed, false, nil, 0, 1.0, newMargins(150, 150, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 4, creator.ColorBlue, true, []int64{3}, 0, 1.0, newMargins(100, 100, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 5, creator.ColorBlack, false, nil, 0, 1.0, newMargins(50, 50, 5, 5)),
	}

	// Draw page content.
	for _, content := range contents {
		if err := c.Draw(content); err != nil {
			return err
		}
	}

	return nil
}

func drawLinesInsideDivision(c *creator.Creator) error {
	c.NewPage()

	var (
		lightGray    = creator.ColorRGBFromHex("#f3f3f3")
		darkGray     = creator.ColorRGBFromHex("#555555")
		sampleText   = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."
		subtitleText = "All styles of lines are supported inside divisions. They can be useful as section separators or to emphasize parts of the content."
	)

	// Draw chapter.
	chapter := c.NewChapter("Lines inside divisions")
	subtitle := newParagraph(c, subtitleText, 15, creator.TextAlignmentJustify, newMargins(0, 0, 8, 15))
	if err := chapter.Add(subtitle); err != nil {
		return err
	}
	if err := c.Draw(chapter); err != nil {
		return err
	}

	// Create division content.
	contents := []creator.VectorDrawable{
		newLine(c, 0, 0, 0, 0, false, true, 2, darkGray, false, nil, 0, 1.0, newMargins(125, 100, 10, 10)),
		newParagraph(c, "Sample section title", 16, creator.TextAlignmentCenter, newMargins(50, 50, 0, 0)),
		newLine(c, 0, 0, 0, 0, false, true, 2, darkGray, false, nil, 0, 1.0, newMargins(125, 100, 10, 20)),
		newParagraph(c, sampleText, 12, creator.TextAlignmentJustify, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 1, darkGray, false, nil, 0, 1.0, newMargins(0, 0, 5, 5)),
		newParagraph(c, sampleText, 12, creator.TextAlignmentJustify, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 1, darkGray, false, nil, 0, 1.0, newMargins(0, 0, 5, 5)),
		newParagraph(c, sampleText, 12, creator.TextAlignmentJustify, newMargins(0, 0, 5, 5)),
		newParagraph(c, sampleText, 12, creator.TextAlignmentJustify, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 1, creator.ColorRed, true, []int64{3}, 0, 1.0, newMargins(0, 0, 5, 5)),
		newParagraph(c, sampleText, 12, creator.TextAlignmentJustify, newMargins(0, 0, 5, 5)),
		newParagraph(c, sampleText, 12, creator.TextAlignmentJustify, newMargins(0, 0, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, creator.ColorBlue, false, nil, 0, 1.0, newMargins(25, 25, 10, 10)),
	}

	// Create division.
	div := c.NewDivision()
	div.SetPadding(10, 10, 10, 10)
	div.SetBackground(&creator.Background{
		FillColor:               lightGray,
		BorderRadiusTopLeft:     5,
		BorderRadiusTopRight:    5,
		BorderRadiusBottomLeft:  5,
		BorderRadiusBottomRight: 5,
	})

	// Draw division content.
	for _, content := range contents {
		if err := div.Add(content); err != nil {
			return err
		}
	}

	// Draw division.
	return c.Draw(div)
}

func drawLinesInsideTable(c *creator.Creator) error {
	c.NewPage()

	var (
		subtitleText = "Lines of all styles can be rendered inside table cells. Here is a quick reference of the customizable line attributes."
		blue         = creator.ColorBlue
		yellow       = creator.ColorYellow
		red          = creator.ColorRed
		black        = creator.ColorBlack
	)

	// Draw chapter.
	chapter := c.NewChapter("Lines inside tables")
	subtitle := newParagraph(c, subtitleText, 15, creator.TextAlignmentJustify, newMargins(0, 0, 8, 15))
	if err := chapter.Add(subtitle); err != nil {
		return err
	}
	if err := c.Draw(chapter); err != nil {
		return err
	}

	// Create table content.
	contents := []creator.VectorDrawable{
		// Solid lines.
		newParagraph(c, "Solid lines", 10, creator.TextAlignmentLeft, newMargins(5, 0, 0, 0)),
		newLine(c, 0, 0, 0, 0, false, true, 1, blue, false, nil, 0, 1.0, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, yellow, false, nil, 0, 1.0, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 5, red, false, nil, 0, 1.0, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 10, black, false, nil, 0, 1.0, newMargins(5, 5, 5, 5)),
		// Dashed lines.
		newParagraph(c, "Dashed lines", 10, creator.TextAlignmentLeft, newMargins(5, 0, 0, 0)),
		newLine(c, 0, 0, 0, 0, false, true, 1, blue, true, []int64{3}, 0, 1.0, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, yellow, true, []int64{3}, 0, 1.0, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 5, red, true, []int64{3}, 0, 1.0, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 10, black, true, []int64{3}, 0, 1.0, newMargins(5, 5, 5, 5)),

		// Transparent solid lines.
		newParagraph(c, "Transparent solid lines", 10, creator.TextAlignmentLeft, newMargins(5, 0, 0, 0)),
		newLine(c, 0, 0, 0, 0, false, true, 1, blue, false, nil, 0, 0.5, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, yellow, false, nil, 0, 0.5, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 5, red, false, nil, 0, 0.5, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 10, black, false, nil, 0, 0.5, newMargins(5, 5, 5, 5)),

		// Transparent dashed lines.
		newParagraph(c, "Transparent dashed lines", 9, creator.TextAlignmentLeft, newMargins(5, 0, 0, 0)),
		newLine(c, 0, 0, 0, 0, false, true, 1, blue, true, []int64{6, 3}, 0, 0.5, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, yellow, true, []int64{6, 3}, 0, 0.5, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 5, red, true, []int64{6, 3}, 0, 0.5, newMargins(5, 5, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 10, black, true, []int64{6, 3}, 0, 0.5, newMargins(5, 5, 5, 5)),

		// Solid lines with horizontal margins.
		newParagraph(c, "Solid lines - horizontal margins", 10, creator.TextAlignmentLeft, newMargins(5, 0, 0, 0)),
		newLine(c, 0, 0, 0, 0, false, true, 1, blue, false, nil, 0, 1.0, newMargins(20, 20, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 3, yellow, false, nil, 0, 1.0, newMargins(20, 20, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 5, red, false, nil, 0, 1.0, newMargins(20, 20, 5, 5)),
		newLine(c, 0, 0, 0, 0, false, true, 10, black, false, nil, 0, 1.0, newMargins(20, 20, 5, 5)),

		// Solid lines with vertical margins.
		newParagraph(c, "Solid lines - vertical margins", 10, creator.TextAlignmentLeft, newMargins(5, 0, 0, 0)),
		newLine(c, 0, 0, 0, 0, false, true, 1, blue, false, nil, 0, 1.0, newMargins(5, 5, 10, 10)),
		newLine(c, 0, 0, 0, 0, false, true, 3, yellow, false, nil, 0, 1.0, newMargins(5, 5, 10, 10)),
		newLine(c, 0, 0, 0, 0, false, true, 5, red, false, nil, 0, 1.0, newMargins(5, 5, 10, 10)),
		newLine(c, 0, 0, 0, 0, false, true, 10, black, false, nil, 0, 1.0, newMargins(5, 5, 10, 10)),
	}

	// Draw table content.
	table := c.NewTable(5)
	if err := table.SetColumnWidths(0.3, 0.175, 0.175, 0.175, 0.175); err != nil {
		return err
	}

	for _, content := range contents {
		cell := table.NewCell()
		cell.SetIndent(0)
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1.0)
		cell.SetVerticalAlignment(creator.CellVerticalAlignmentMiddle)

		if err := cell.SetContent(content); err != nil {
			return err
		}
	}

	// Draw table.
	return c.Draw(table)
}

func newParagraph(c *creator.Creator, title string, fontSize float64, textAlign creator.TextAlignment, margins creator.Margins) *creator.StyledParagraph {
	sp := c.NewStyledParagraph()
	sp.SetMargins(margins.Left, margins.Right, margins.Top, margins.Bottom)
	sp.Append(title).Style.FontSize = fontSize
	sp.SetTextAlignment(textAlign)
	sp.SetTextVerticalAlignment(creator.TextVerticalAlignmentCenter)
	return sp
}

func newLine(c *creator.Creator, x1, y1, x2, y2 float64, isAbsolute, fillWidth bool, lineWidth float64,
	color creator.Color, isDashed bool, dashArray []int64, dashPhase int64, opacity float64, margins creator.Margins) *creator.Line {
	positioning := creator.PositionRelative
	if isAbsolute {
		positioning = creator.PositionAbsolute
	}

	fitMode := creator.FitModeNone
	if fillWidth {
		fitMode = creator.FitModeFillWidth
	}

	style := draw.LineStyleSolid
	if isDashed {
		style = draw.LineStyleDashed
	}

	line := c.NewLine(x1, y1, x2, y2)
	line.SetLineWidth(lineWidth)
	line.SetMargins(margins.Left, margins.Right, margins.Top, margins.Bottom)
	line.SetPositioning(positioning)
	line.SetFitMode(fitMode)
	line.SetColor(color)
	line.SetStyle(style)
	line.SetDashPattern(dashArray, dashPhase)
	line.SetOpacity(opacity)
	return line
}

func newMargins(left, right, top, bottom float64) creator.Margins {
	return creator.Margins{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
}
