/*
 * Filling text with a color gradient.
 *
 * A gradient color returned by Creator.NewLinearGradientColor or
 * Creator.NewRadialGradientColor can be assigned to a text chunk's style, in
 * which case the glyphs are filled with a shading pattern instead of a flat
 * color. The gradient's extent is derived from the chunks that use it: a color
 * set on a single chunk is centered on that chunk, while a color shared across
 * the paragraph spans the whole paragraph.
 *
 * NOTE: the bundled image renderer does not yet paint pattern-filled text and
 * falls back to solid black. The generated file itself is correct and displays
 * the gradient in PDF viewers.
 *
 * Run as: go run pdf_gradient_text.go
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v5/common/license"
	"github.com/unidoc/unipdf/v5/creator"
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
	c.SetPageMargins(50, 50, 50, 50)

	if err := linearGradientText(c); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := radialGradientText(c); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := perChunkGradientText(c); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := c.WriteToFile("gradient_text.pdf"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// linearGradientText fills a heading with a two-stop horizontal gradient.
func linearGradientText(c *creator.Creator) error {
	// Each color point pairs a color with its position along the gradient,
	// from 0.0 (start) to 1.0 (end).
	gradient := c.NewLinearGradientColor([]*creator.ColorPoint{
		creator.NewColorPoint(creator.ColorRGBFrom8bit(0xE5, 0x3E, 0x3E), 0.0),
		creator.NewColorPoint(creator.ColorRGBFrom8bit(0x2B, 0x6C, 0xB0), 1.0),
	})

	p := c.NewStyledParagraph()
	p.SetMargins(0, 0, 0, 20)
	p.SetTextAlignment(creator.TextAlignmentCenter)

	chunk := p.Append("Linear gradient heading")
	chunk.Style.FontSize = 32
	chunk.Style.Color = gradient

	return c.Draw(p)
}

// radialGradientText fills text with a gradient radiating from its center.
func radialGradientText(c *creator.Creator) error {
	// The first two arguments are an offset from the anchor point of the filled
	// area (the center of the text by default), not absolute page coordinates.
	// Passing 0, 0 centers the gradient on the text. An outer radius of -1 fits
	// the radius to the filled area; a fixed radius smaller than the text leaves
	// the glyphs outside the disc unpainted.
	gradient := c.NewRadialGradientColor(0, 0, 0, -1, []*creator.ColorPoint{
		creator.NewColorPoint(creator.ColorRGBFrom8bit(0xF6, 0xC3, 0x43), 0.0),
		creator.NewColorPoint(creator.ColorRGBFrom8bit(0xC0, 0x39, 0x2B), 1.0),
	})

	p := c.NewStyledParagraph()
	p.SetMargins(0, 0, 0, 20)
	p.SetTextAlignment(creator.TextAlignmentCenter)

	chunk := p.Append("Radial gradient heading")
	chunk.Style.FontSize = 32
	chunk.Style.Color = gradient

	return c.Draw(p)
}

// perChunkGradientText mixes gradient-filled and flat-colored chunks within a
// single paragraph. A gradient assigned to one chunk is scoped to that chunk.
func perChunkGradientText(c *creator.Creator) error {
	gradient := c.NewLinearGradientColor([]*creator.ColorPoint{
		creator.NewColorPoint(creator.ColorRGBFrom8bit(0x16, 0xA0, 0x85), 0.0),
		creator.NewColorPoint(creator.ColorRGBFrom8bit(0x8E, 0x44, 0xAD), 1.0),
	})

	p := c.NewStyledParagraph()
	p.SetTextAlignment(creator.TextAlignmentLeft)

	flat := p.Append("Flat color, then ")
	flat.Style.FontSize = 20
	flat.Style.Color = creator.ColorBlack

	graded := p.Append("a gradient chunk")
	graded.Style.FontSize = 20
	graded.Style.Color = gradient

	trailing := p.Append(", then flat again.")
	trailing.Style.FontSize = 20
	trailing.Style.Color = creator.ColorBlack

	return c.Draw(p)
}
