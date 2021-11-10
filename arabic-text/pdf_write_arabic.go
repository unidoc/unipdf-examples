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

	textTitle := `كتابة مثال عربي`
	par := c.NewParagraph(arabic.Shape(textTitle))
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

	textArabic := `هذه فقرة بسيطة جدا`
	par = c.NewParagraph(arabic.Shape(textArabic))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("This is a pretty simple paragraph")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabic = `يبدو ، نضيف المزيد من الفقرة هنا`
	par = c.NewParagraph(arabic.Shape(textArabic))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("Looks, we adding more paragraph here")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabic = `دعنا نحاول إضافة الرقم 300 و 500`
	par = c.NewParagraph(arabic.Shape(textArabic))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("Let's try to adding number 300 and 500")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabic = `ماذا لو أضفنا المزيد من الفقرات؟`
	par = c.NewParagraph(arabic.Shape(textArabic))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("How about we add more paragraphs?")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check error.
	err = c.WriteToFile("arabic-text.pdf")
	if err != nil {
		panic(err)
	}
}
