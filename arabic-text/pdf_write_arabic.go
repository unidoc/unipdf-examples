/*
 * An example of writing arabic text
 * need to use external library for shaping arabic text.
 *
 * Run as: go run pdf_write_arabic.go
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"

	arabic "github.com/abdullahdiaa/garabic"
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
	c.NewPage()

	// Load font as composite PDF font.
	arabicFont, err := model.NewCompositePdfFontFromTTFFile("./Amiri-Regular.ttf")
	if err != nil {
		panic(err)
	}

	// Set font size.
	fontSize := float64(14)
	// For arabic text to be rendered correctly,
	// we need to use external library for shaping arabic text
	// in this example, we use: https://github.com/AbdullahDiaa/garabic.

	textBismillah := `بِسْمِ اللّٰهِ الرَّحْمٰنِ الرَّحِيْمِ`
	par := c.NewParagraph(arabic.Shape(textBismillah))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentCenter)
	// Get paragraph margins.
	marginLeft, marginRight, marginTop, marginBottom := par.GetMargins()
	// Since arabic character having more height than latin character, add more margin.
	marginTop += 4
	marginBottom += 10

	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)

	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textSurahAlIkhlas := `قُلْ هُوَ اللّٰهُ اَحَدٌۚ - ١`
	par = c.NewParagraph(arabic.Shape(textSurahAlIkhlas))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("1. Say, ˹O Prophet,˺ “He is Allah—One ˹and Indivisible˺;")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textSurahAlIkhlas = `اَللّٰهُ الصَّمَدُۚ - ٢`
	par = c.NewParagraph(arabic.Shape(textSurahAlIkhlas))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("2. Allah—the Sustainer ˹needed by all˺.")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textSurahAlIkhlas = `لَمْ يَلِدْ وَلَمْ يُوْلَدْۙ - ٣`
	par = c.NewParagraph(arabic.Shape(textSurahAlIkhlas))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("3. He has never had offspring, nor was He born.")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textSurahAlIkhlas = `وَلَمْ يَكُنْ لَّهٗ كُفُوًا اَحَدٌ ࣖ - ٤`
	par = c.NewParagraph(arabic.Shape(textSurahAlIkhlas))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("4. And there is none comparable to Him.”")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check error.
	err = c.WriteToFile("arabic-surah-al-ikhlas.pdf")
	if err != nil {
		panic(err)
	}
}
