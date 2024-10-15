/*
 * An example of adding svg to a PDF file
 *
 * Run as: go run add_svg.go
 */

package main

import (
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
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
	file := "./svgs/different_shapes.svg"
	graphicSvg, err := creator.NewGraphicSVGFromFile(file)
	if err != nil {
		panic(err)
	}

	err = c.Draw(graphicSvg)
	if err != nil {
		panic(err)
	}
	c.WriteToFile("pdf_shapes_and_text_svg.pdf")
}
