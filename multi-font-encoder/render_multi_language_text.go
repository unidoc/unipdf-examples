/*
 * renders a multi language text using a multi font encoder.
 *
 * Run as: go run render_multi_language_text.go
 */
package main

import (
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
	inputText := "Sample text 示例文本 טקסט לדוגמה ጽሑፍ ኣብነት 샘플 텍스트"
	fonts := getFonts([]string{
		"./fonts/PCSB Hebrew Regular.ttf",
		"./fonts/NotoSerifEthiopic_Condensed-Black.ttf",
		"./fonts/OpenSans-Regular.ttf",
		"./fonts/Batang.ttf",
	})

	c := creator.New()
	c.SetPageSize(creator.PageSizeA5)
	p := c.NewStyledParagraph()
	style := &p.Append(inputText).Style
	style.MultiFont = model.NewMultipleFontEncoder(fonts)

	err := c.Draw(p)
	if err != nil {
		panic(err)
	}
	if err := c.WriteToFile("multiple-language-text-multi-font.pdf"); err != nil {
		panic(err)
	}
}

// getFonts returns list of *model.PdfFont from list of paths to fonts files.
func getFonts(fontPaths []string) []*model.PdfFont {
	fonts := []*model.PdfFont{}
	for _, path := range fontPaths {
		font, err := model.NewCompositePdfFontFromTTFFile(path)
		if err != nil {
			panic(err)
		}
		fonts = append(fonts, font)
	}

	return fonts
}
