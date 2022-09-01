/*
 * This example showcases a masonry-like layout created using tables
 * and divisions.
 *
 * Run as: go run pdf_tables_division_layout.go
 */

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model/optimize"
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
	var (
		hPageMargin    = 50.0
		hContentMargin = 20.0
	)

	c := creator.New()
	c.SetPageMargins(hPageMargin, hPageMargin, 25, 25)
	c.NewPage()

	// Set optimizer.
	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   300,
	}))

	// Draw header section.
	headerTable := c.NewTable(1)
	headerTable.SetMargins(0, 0, 0, 10)

	p := c.NewStyledParagraph()
	p.SetMargins(0, 0, 5, 5)
	p.Append("Image gallery").Style.FontSize = 14

	cell := headerTable.NewCell()
	cell.SetBackgroundColor(creator.ColorRGBFromHex("#eeeeee"))
	if err := cell.SetContent(p); err != nil {
		log.Fatalf("failed to add paragraph to header cell: %v", err)
	}

	if err := c.Draw(headerTable); err != nil {
		log.Fatalf("failed to draw header table: %v", err)
	}

	// Create image layout.
	imagePaths := []string{
		"images/image1.jpg",
		"images/image2.jpg",
		"images/image3.jpg",
		"images/image4.jpg",
		"images/image5.jpg",
		"images/image6.jpg",
		"images/image7.jpg",
		"images/image8.jpg",
	}

	contentTable := c.NewTable(2)
	contentTable.EnableRowWrap(true)

	divLeft, divRight := c.NewDivision(), c.NewDivision()
	for i, imagePath := range imagePaths {
		img, err := c.NewImageFromFile(imagePath)
		if err != nil {
			log.Fatalf("failed to load img: %v", err)
		}

		// Scale image to width (account for horizontal margins as well).
		img.ScaleToWidth((c.Width() - 2*hPageMargin - 2*hContentMargin) / 2)

		p := c.NewStyledParagraph()
		p.Append(fmt.Sprintf("Photo %d", i+1))

		div := divLeft
		if i%2 != 0 {
			div = divRight
			img.SetMargins(hContentMargin, 0, 0, 0)
			p.SetMargins(hContentMargin, 0, 10, 15)
		} else {
			img.SetMargins(0, hContentMargin, 0, 0)
			p.SetMargins(0, hContentMargin, 10, 15)
		}

		if err = div.Add(img); err != nil {
			log.Fatalf("failed to add image to division: %v", err)
		}
		if err = div.Add(p); err != nil {
			log.Fatalf("failed to add paragraph to division: %v", err)
		}
	}

	// Create left column.
	colLeft := contentTable.NewCell()
	colLeft.SetIndent(0)
	colLeft.SetHorizontalAlignment(creator.CellHorizontalAlignmentLeft)

	if err := colLeft.SetContent(divLeft); err != nil {
		log.Fatalf("failed to add division to left cell: %v", err)
	}

	// Create right column.
	colRight := contentTable.NewCell()
	colRight.SetIndent(0)
	colRight.SetHorizontalAlignment(creator.CellHorizontalAlignmentRight)

	if err := colRight.SetContent(divRight); err != nil {
		log.Fatalf("failed to add division to right cell: %v", err)
	}

	// Draw table.
	if err := c.Draw(contentTable); err != nil {
		log.Fatalf("failed to draw content table: %v", err)
	}

	// Write output file.
	if err := c.WriteToFile("unipdf-tables-division-layout.pdf"); err != nil {
		log.Fatal(err)
	}
}
