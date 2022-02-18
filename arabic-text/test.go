package main

import (
	//"fmt"
	"math"
	"os"
	//"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/textshaping"
	//"golang.org/x/text/unicode/bidi"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

// Commonly used colors.
var (
	ColorBlack  = ColorRGBFromArithmetic(0, 0, 0)
	ColorWhite  = ColorRGBFromArithmetic(1, 1, 1)
	ColorRed    = ColorRGBFromArithmetic(1, 0, 0)
	ColorGreen  = ColorRGBFromArithmetic(0, 1, 0)
	ColorBlue   = ColorRGBFromArithmetic(0, 0, 1)
	ColorYellow = ColorRGBFromArithmetic(1, 1, 0)
)

func main() {
	//c := creator.New()
	//c.NewPage()

	//// Load font as composite PDF font.
	//arabicFont, err := model.NewCompositePdfFontFromTTFFile("./Amiri-Regular.ttf")
	//if err != nil {
	//	panic(err)
	//}

	//// Set font size.
	//fontSize := float64(14)

	//textArabic := `هنالك العديد من الأنواع المتوفرة لنصوص لوريم إيبسوم، ولكن الغالبية تم تعديلها بشكل ما عبر إدخال بعض النوادر أو الكلمات العشوائية إلى النص. إن كنت تريد أن تستخدم نص لوريم إيبسوم ما، عليك أن تتحقق أولاً أن ليس هناك أي كلمات أو عبارات محرجة أو غير لائقة مخبأة في هذا النص. بينما تعمل جميع مولّدات نصوص لوريم إيبسوم على الإنترنت على إعادة تكرار مقاطع من نص لوريم إيبسوم نفسه عدة مرات بما تتطلبه الحاجة، يقوم مولّدنا هذا باستخدام كلمات من قاموس يحوي على أكثر من 200 كلمة لا تينية، مضاف إليها مجموعة من الجمل النموذجية، لتكوين نص لوريم إيبسوم ذو شكل منطقي قريب إلى النص الحقيقي. وبالتالي يكون النص الناتح خالي من التكرار، أو أي كلمات أو عبارات غير لائقة أو ما شابه. وهذا ما يجعله أول مولّد نص لوريم إيبسوم حقيقي على الإنترنت.`

	////shaped, err := textshaping.ArabicShape(textArabic)
	////if err != nil {
	////	panic(err)
	////}
	//shaped := textArabic

	//par := c.NewParagraph(shaped)
	//par.SetFont(arabicFont)
	//par.SetFontSize(fontSize)
	//par.SetTextAlignment(creator.TextAlignmentRight)
	//err = c.Draw(par)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////fmt.Println("org: ", textArabic)
	////fmt.Println("shp: ", bidi.ReverseString(shaped))
	//bidiP := bidi.Paragraph{}
	//bidiP.SetString(shaped)

	//o, _ := bidiP.Order()
	//for i := 0; i < o.NumRuns(); i++ {
	//	run := o.Run(i)
	//	fmt.Printf("i: %v, str: %v\n", i, run.String())
	//}
	//fmt.Println("org: ", textArabic)
	//fmt.Println("shp: ", shaped)

	//textArabic = `هناك حقيقة مثبتة منذ زمن طويل وهي أن المحتوى المقروء لصفحة ما سيلهي القارئ عن التركيز على الشكل الخارجي للنص أو شكل توضع الفقرات في الصفحة التي يقرأها. ولذلك يتم استخدام طريقة لوريم إيبسوم لأنها تعطي توزيعاَ طبيعياَ -إلى حد ما- للأحرف عوضاً عن استخدام "هنا يوجد محتوى نصي، هنا يوجد محتوى نصي" فتجعلها تبدو (أي الأحرف) وكأنها نص مقروء. العديد من برامح النشر المكتبي وبرامح تحرير صفحات الويب تستخدم لوريم إيبسوم بشكل إفتراضي كنموذج عن النص، وإذا قمت بإدخال "lorem ipsum" في أي محرك ث ستظهر العديد من المواقع الحديثة العهد في نتائج البحث. على مدى السنين ظهرت نسخ جديدة ومختلفة من نص لوريم إيبسوم، أحياناً عن طريق الصدفة، وأحياناً عن عمد كإدخال بعض العبارات الفكاهية إليها.`
	//shaped, err = textshaping.ArabicShape(textArabic)
	//if err != nil {
	//	panic(err)
	//}

	//par = c.NewParagraph(shaped)
	//par.SetFont(arabicFont)
	//par.SetFontSize(fontSize)
	//par.SetTextAlignment(creator.TextAlignmentRight)

	//err = c.Draw(par)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////fmt.Println("org: ", textArabic)
	////// Make sure to check error.
	//err = c.WriteToFile("arabic-text.pdf")
	//if err != nil {
	//	panic(err)
	//}

	//testLines()
	//testNormal()
	//testNormal1()
	//testNormal15()
	//TestStyledParagraphTextRise()
	TestUnderline()
}

func testNormal() {
	c := creator.New()
	font, err := model.NewCompositePdfFontFromTTFFile("./Amiri-Regular.ttf")
	if err != nil {
		panic(err)
	}

	text := `What is Lorem Ipsum?`

	p := c.NewStyledParagraph()
	chunk := p.Append(text)
	chunk.Style.FontSize = 18
	chunk.Style.Font = font
	p.SetLineHeight(1)

	err = c.Draw(p)
	if err != nil {
		panic(err)
	}

	p2 := c.NewStyledParagraph()

	chunk = p2.Append("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.")
	chunk.Style.FontSize = 12
	chunk.Style.Font = font
	p2.SetLineHeight(2)

	err = c.Draw(p2)
	if err != nil {
		panic(err)
	}

	err = c.WriteToFile("test_lines_latin-2.pdf")
	if err != nil {
		panic(err)
	}

}

func testNormal1() {
	c := creator.New()
	font, err := model.NewCompositePdfFontFromTTFFile("./Amiri-Regular.ttf")
	if err != nil {
		panic(err)
	}

	text := `What is Lorem Ipsum?`

	p := c.NewStyledParagraph()
	p.SetLineHeight(1)
	chunk := p.Append(text)
	chunk.Style.FontSize = 18
	chunk.Style.Font = font

	err = c.Draw(p)
	if err != nil {
		panic(err)
	}

	p2 := c.NewStyledParagraph()
	p2.SetLineHeight(1)

	chunk = p2.Append("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.")
	chunk.Style.FontSize = 12
	chunk.Style.Font = font

	err = c.Draw(p2)
	if err != nil {
		panic(err)
	}

	err = c.WriteToFile("test_lines_latin-1.pdf")
	if err != nil {
		panic(err)
	}

}

func testNormal15() {
	c := creator.New()
	font, err := model.NewCompositePdfFontFromTTFFile("./Amiri-Regular.ttf")
	if err != nil {
		panic(err)
	}

	text := `What is Lorem Ipsum?`

	p := c.NewStyledParagraph()
	p.SetLineHeight(1)
	chunk := p.Append(text)
	chunk.Style.FontSize = 18
	chunk.Style.Font = font

	err = c.Draw(p)
	if err != nil {
		panic(err)
	}

	p2 := c.NewStyledParagraph()
	p2.SetLineHeight(1.5)

	chunk = p2.Append("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.")
	chunk.Style.FontSize = 12
	chunk.Style.Font = font

	err = c.Draw(p2)
	if err != nil {
		panic(err)
	}

	err = c.WriteToFile("test_lines_latin-1.5.pdf")
	if err != nil {
		panic(err)
	}

}

func testLines() {
	c := creator.New()

	lines := []string{
		"The names of these states in Arabic are مصر, البحري and الكويت respectively",
		"عن النص، وإذا قمت بإدخال lorem ipsum في أي محرك بحث ستظهر العديد من المواقع",
		"Hello World!",
		"كتابة مثال عربي",
		"يبدو ، نضيف المزيد من الفقرة هنا",
		"Another simple word",
		"ماذا لو أضفنا المزيد من الفقرات",
		"الكلاسيكي منذ العام 45 قبل الميلاد، مما يجعله أكثر من 2000 عام في القدم",
		`فلقد اتضح أن كلمات نص لوريم إيبسوم تأتي من الأقسام 1.10.32 و 1.10.33 من كتاب "حول`,
		`أقاصي الخير والشر" (de Finibus Bonorum et Malorum) للمفكر شيشيرون`,
		` ipsum dolor sit amet.." يأتي من سطر في القسم 1.20.32 من هذا الكتاب.`,
		`عن النص، وإذا قمت بإدخال "lorem ipsum" في أي محرك بحث ستظهر العديد من المواقع`,
		`با تشکر از  "Asuni Nicola" و محمد علی گل کار برای پشتیبانی زبان فارسی.`,
	}

	arabicFont, err := model.NewCompositePdfFontFromTTFFile("./Amiri-Regular.ttf")
	if err != nil {
		panic(err)
	}

	// Enable font subsetting for the composite font - embed only needed glyphs
	// (much smaller file size for large fonts).
	//c.EnableFontSubsetting(arabicFont)

	for _, line := range lines {
		shaped, err := textshaping.ArabicShape(line)
		if err != nil {
			panic(err)
		}

		p := c.NewParagraph(shaped)
		p.SetFont(arabicFont)
		p.SetTextAlignment(creator.TextAlignmentRight)

		marginLeft, marginRight, marginTop, marginBottom := p.GetMargins()
		// Since arabic character having more height than latin character, add more margin.
		marginTop += 4
		marginBottom += 10

		p.SetMargins(marginLeft, marginRight, marginTop, marginBottom)

		//err = c.Draw(p)
		//if err != nil {
		//	panic(err)
		//}

		//	opt := bidi.DefaultDirection(bidi.LeftToRight)
		//	bidiP := bidi.Paragraph{}
		//	bidiP.SetString(shaped, opt)

		//	o, _ := bidiP.Order()
		//	for i := 0; i < o.NumRuns(); i++ {
		//		run := o.Run(i)
		//		fmt.Printf("i: %v, str: %v\n", i, run.String())
		//		fmt.Println("direction:", run.Direction())
		//	}

		//	split := strings.Split(line, " ")
		//	for i, str := range split {
		//		fmt.Printf("i: %v, str: %v\n", i, str)
		//	}

		//	runes := []rune(line)
		//	for i, r := range runes {
		//		fmt.Printf("i: %v, char: %v \n", i, string(r))
		//	}
		//	fmt.Println("\n")
		//	fmt.Println("bidi order LTR:", bidiP.Direction() == bidi.LeftToRight)
		//	fmt.Println("bidi order RTL:", bidiP.Direction() == bidi.RightToLeft)

	}

	black := creator.ColorBlack
	green := creator.ColorRGBFromArithmetic(0.0, 1.0, 0.0)

	ch := c.NewChapter("Text Color RGB")
	ch.GetHeading().SetColor(black)

	// Styled text.
	shaped, err := textshaping.ArabicShape(lines[0])
	if err != nil {
		panic(err)
	}

	textArabic := `هناك حقيقة مثبتة منذ زمن طويل وهي أن المحتوى المقروء لصفحة ما سيلهي القارئ عن التركيز على الشكل الخارجي للنص أو شكل توضع الفقرات في الصفحة التي يقرأها. ولذلك يتم استخدام طريقة لوريم إيبسوم لأنها تعطي توزيعاَ طبيعياَ -إلى حد ما- للأحرف عوضاً عن استخدام "هنا يوجد محتوى نصي، هنا يوجد محتوى نصي" فتجعلها تبدو (أي الأحرف) وكأنها نص مقروء. العديد من برامح النشر المكتبي وبرامح تحرير صفحات الويب تستخدم لوريم إيبسوم بشكل إفتراضي كنموذج عن النص، وإذا قمت بإدخال "lorem ipsum" في أي محرك ث ستظهر العديد من المواقع الحديثة العهد في نتائج البحث. على مدى السنين ظهرت نسخ جديدة ومختلفة من نص لوريم إيبسوم، أحياناً عن طريق الصدفة، وأحياناً عن عمد كإدخال بعض العبارات الفكاهية إليها.`

	shaped2, err := textshaping.ArabicShape(textArabic)
	if err != nil {
		panic(err)
	}

	textArabic2 := `هنالك العديد من الأنواع المتوفرة لنصوص لوريم إيبسوم، ولكن الغالبية تم تعديلها بشكل ما عبر إدخال بعض النوادر أو الكلمات العشوائية إلى النص. إن كنت تريد أن تستخدم نص لوريم إيبسوم ما، عليك أن تتحقق أولاً أن ليس هناك أي كلمات أو عبارات محرجة أو غير لائقة مخبأة في هذا النص. بينما تعمل جميع مولّدات نصوص لوريم إيبسوم على الإنترنت على إعادة تكرار مقاطع من نص لوريم إيبسوم نفسه عدة مرات بما تتطلبه الحاجة، يقوم مولّدنا هذا باستخدام كلمات من قاموس يحوي على أكثر من 200 كلمة لا تينية، مضاف إليها مجموعة من الجمل النموذجية، لتكوين نص لوريم إيبسوم ذو شكل منطقي قريب إلى النص الحقيقي. وبالتالي يكون النص الناتح خالي من التكرار، أو أي كلمات أو عبارات غير لائقة أو ما شابه. وهذا ما يجعله أول مولّد نص لوريم إيبسوم حقيقي على الإنترنت.`
	shaped3, err := textshaping.ArabicShape(textArabic2)
	if err != nil {
		panic(err)
	}

	p := c.NewStyledParagraph()
	p.SetLineHeight(1)
	tc := p.SetText(shaped)
	tc.Style.Color = green
	tc.Style.FontSize = 12
	tc.Style.Font = arabicFont

	chunk := p.Append("\n" + shaped)
	chunk.Style.Color = green
	chunk.Style.FontSize = 12
	chunk.Style.Font = arabicFont

	chunk = p.Append("\n" + shaped)
	chunk.Style.Color = green
	chunk.Style.FontSize = 12
	chunk.Style.Font = arabicFont

	chunk = p.Append("\n" + shaped2)
	chunk.Style.Font = arabicFont

	chunk = p.Append("\n" + shaped3)
	chunk.Style.Font = arabicFont

	// Styled text.
	shaped, err = textshaping.ArabicShape(lines[1])
	if err != nil {
		panic(err)
	}

	chunk = p.Append("\n" + shaped)
	chunk.Style.Font = arabicFont

	// Styled text.
	shaped, err = textshaping.ArabicShape(lines[7])
	if err != nil {
		panic(err)
	}

	chunk = p.Append("\n" + shaped)
	chunk.Style.Font = arabicFont

	p.SetMargins(20, 0, 10, 0)

	ch.Add(p)
	err = c.Draw(ch)
	if err != nil {
		panic(err)
	}

	err = c.WriteToFile("test_lines_fixes.pdf")
	if err != nil {
		panic(err)
	}
}

func TestStyledParagraphTextRise() {
	c := creator.New()

	// Basic usage.
	p := c.NewStyledParagraph()
	p.SetMargins(0, 0, 0, 20)
	p.SetLineHeight(1.2)
	p.Append("Styled paragraphs allow drawing ")

	style := &p.Append("subscript").Style
	style.TextRise = -8
	style.FontSize = 7
	style.Color = ColorRed

	p.Append(" and ")

	style = &p.Append("superscript").Style
	style.TextRise = 9
	style.FontSize = 7
	style.Color = ColorBlue

	p.Append(" text chunks.")

	// Draw paragraph.
	err := c.Draw(p)
	if err != nil {
		panic(err)
	}

	// Text rise combined with underlined and annotated text.
	p = c.NewStyledParagraph()
	p.SetMargins(0, 0, 0, 20)
	p.SetLineHeight(1.2)
	p.Append("Subscript and superscript text can be ")

	style = &p.Append("underlined").Style
	style.TextRise = 5
	style.FontSize = 7
	style.Color = ColorGreen
	style.Underline = true

	p.Append(" or turned into ")

	style = &p.AddExternalLink("link", "https://google.com").Style
	style.TextRise = -5
	style.FontSize = 7

	p.Append(" annotations.")

	// Draw paragraph.
	err = c.Draw(p)
	if err != nil {
		panic(err)
	}
	// Text rise and row wrapping.
	p = c.NewStyledParagraph()
	p.SetMargins(0, 0, 0, 20)
	p.SetLineHeight(1.2)

	p.Append("Let's see if subscript and superscript text is properly wrapped into lines. ")

	style = &p.Append("This is some superscript text which should not fit on the current line.").Style
	style.TextRise = 5
	style.FontSize = 7
	style.Color = ColorGreen
	style.Underline = true
	style.UnderlineStyle.Offset = 2

	p.Append(" And then some regular text again to also test subscripts. ")

	style = &p.AddExternalLink("This is some subscript text which should not fit on the current line.", "https://google.com").Style
	style.TextRise = -2
	style.FontSize = 7

	p.Append(" And then some regular text again.")

	// Draw paragraph.
	err = c.Draw(p)
	if err != nil {
		panic(err)
	}

	err = c.WriteToFile("test_subscript.pdf")
	if err != nil {
		panic(err)
	}
}

func TestUnderline() {
	c := creator.New()

	p := c.NewStyledParagraph()
	p.SetLineHeight(1.2)
	p.Append("Text chunks can be ")

	// Default underline style.
	chunk := p.Append("underlined")
	chunk.Style.Underline = true

	p.Append(" using the default style.\n")
	p.Append("By default, the ")

	// Default underline style based on the color of the text.
	chunk = p.Append("underline color")
	chunk.Style.Underline = true
	chunk.Style.Color = ColorBlue

	p.Append(" is the color of the text chunk. We can also add ")

	// Custom underline style color.
	chunk = p.Append("some long underlined text chunk which wraps")
	chunk.Style.FontSize = 15
	chunk.Style.Underline = true
	chunk.Style.UnderlineStyle.Color = ColorRed
	chunk.Style.UnderlineStyle.Offset = 1

	p.Append(" and then some more regular text.\n")

	// Custom underline thickness and offset.
	p.Append("Finally, we can customize the offset and the thickness of the ")

	chunk = p.Append("underlined text")
	chunk.Style.Underline = true
	chunk.Style.UnderlineStyle.Thickness = 2
	chunk.Style.UnderlineStyle.Offset = 2

	p.Append(".")

	// Draw paragraph.
	err := c.Draw(p)
	if err != nil {
		panic(err)
	}

	// Write output file.
	err = c.WriteToFile("styled_paragraph_underline.pdf")
	if err != nil {
		panic(err)
	}
}

// Color interface represents colors in the PDF creator.
type Color interface {
	ToRGB() (float64, float64, float64)
}

// rgbColor represents a color in the RGB color model.
type rgbColor struct {
	// Arithmetic representation of r,g,b (range 0-1).
	r, g, b float64
}

func (col rgbColor) ToRGB() (float64, float64, float64) {
	return col.r, col.g, col.b
}

// ColorRGBFromArithmetic creates a Color from arithmetic color values (0-1).
// Example:
//   green := ColorRGBFromArithmetic(0.0, 1.0, 0.0)
func ColorRGBFromArithmetic(r, g, b float64) Color {
	return rgbColor{
		r: math.Max(math.Min(r, 1.0), 0.0),
		g: math.Max(math.Min(g, 1.0), 0.0),
		b: math.Max(math.Min(b, 1.0), 0.0),
	}
}
