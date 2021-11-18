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

	"github.com/unidoc/unipdf/v3/common"
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
	common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
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

	textArabic := `أم بينما الأعمال بلا, هذا غينيا يعادل اعتداء تم, هو شيء مقاطعة وبولندا. قد لان ووصف التبرعات, أي مكن هُزم الشّعبين, بها أعلنت يتسنّى ا و. من دنو إبّان الأوضاع ولاتّساع. حيث وأزيز وتتحمّل وباستثناء عن, حيث أي وانهاء التّحول. ولم أسيا الساحة أي, وتنصيب اتفاقية ألمانيا تحت أم.`
	par = c.NewParagraph(arabic.Shape(textArabic))
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
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

	par = c.NewParagraph("LTR: Looks, we adding more simple paragraph here lets try to write it, how it looks like. Is it good? How about we add more word here")
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	par = c.NewParagraph("RTL: Looks, we adding more simple paragraph here lets try to write it, how it looks like. Is it good? How about we add more word here")
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

	textArabic = `عرض عن الشرق، وتتحمّل الموسوعة. بحق ترتيب الساحة اسبوعين تم, عل لإعادة يتعلّق بها, ومن إذ أفاق وباءت العظمى. زهاء فمرّ فهرست لم عدد, مع عرض حالية إبّان مشاركة. على عل وترك حالية جزيرتي, جعل وسوء الحكم للجزر هو. وبعد إحتار تكتيكاً أم مكن, جنوب المضي عسكرياً أخر بل, و كلا فبعد الشهير اليميني.`
	par = c.NewParagraph(arabic.Shape(textArabic))
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
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

	textMixed := fmt.Sprintf("The names of these states in Arabic are %s and %s respectively. Add another word here, should be in new line.", arabic.Shape("مصر, البحري"), arabic.Shape("الكويت"))
	par = c.NewParagraph(textMixed)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
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
