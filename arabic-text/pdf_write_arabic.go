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

	"github.com/unidoc/unipdf/v4/common"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/unidoc/unipdf/v4/textshaping"
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

	textArabicShaped, err := textshaping.ArabicShape(`كتابة مثال عربي`)
	if err != nil {
		panic(err)
	}

	par := c.NewParagraph(textArabicShaped)
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

	textArabicShaped, err = textshaping.ArabicShape(`أم بينما الأعمال بلا, هذا غينيا يعادل اعتداء تم, هو شيء مقاطعة وبولندا. قد لان ووصف التبرعات, أي مكن هُزم الشّعبين, بها أعلنت يتسنّى ا و. من دنو إبّان الأوضاع ولاتّساع. حيث وأزيز وتتحمّل وباستثناء عن, حيث أي وانهاء التّحول. ولم أسيا الساحة أي, وتنصيب اتفاقية ألمانيا تحت أم.`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
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

	textArabicShaped, err = textshaping.ArabicShape(`يبدو ، نضيف المزيد من الفقرة هنا`)
	if err != nil {
		panic(err)
	}
	par = c.NewParagraph(textArabicShaped)
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

	textArabicShaped, err = textshaping.ArabicShape(`عرض عن الشرق، وتتحمّل الموسوعة. بحق ترتيب الساحة اسبوعين تم, عل لإعادة يتعلّق بها, ومن إذ أفاق وباءت العظمى. زهاء فمرّ فهرست لم عدد, مع عرض حالية إبّان مشاركة. على عل وترك حالية جزيرتي, جعل وسوء الحكم للجزر هو. وبعد إحتار تكتيكاً أم مكن, جنوب المضي عسكرياً أخر بل, و كلا فبعد الشهير اليميني.`)
	if err != nil {
		panic(err)
	}
	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`ماذا لو أضفنا المزيد من الفقرات؟`)
	if err != nil {
		panic(err)
	}
	par = c.NewParagraph(textArabicShaped)
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

	textMixed, err := textshaping.ArabicShape("The names of these states in Arabic are مصر, البحري and الكويت respectively. Add another word here, should be in new line.")
	par = c.NewParagraph(textMixed)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`ما فائدته ؟`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`هناك حقيقة مثبتة منذ زمن طويل وهي أن المحتوى المقروء لصفحة ما سيلهي القارئ عن التركيز على الشكل الخارجي للنص أو شكل توضع الفقرات في الصفحة التي يقرأها. ولذلك يتم استخدام طريقة لوريم إيبسوم لأنها تعطي توزيعاَ طبيعياَ -إلى حد ما- للأحرف عوضاً عن استخدام "هنا يوجد محتوى نصي، هنا يوجد محتوى نصي" فتجعلها تبدو (أي الأحرف) وكأنها نص مقروء. العديد من برامح النشر المكتبي وبرامح تحرير صفحات الويب تستخدم لوريم إيبسوم بشكل إفتراضي كنموذج عن النص، وإذا قمت بإدخال "lorem ipsum" في أي محرك بحث ستظهر العديد من المواقع الحديثة العهد في نتائج البحث. على مدى السنين ظهرت نسخ جديدة ومختلفة من نص لوريم إيبسوم، أحياناً عن طريق الصدفة، وأحياناً عن عمد كإدخال بعض العبارات الفكاهية إليها.؟`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`أين أجده ؟`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`هنالك العديد من الأنواع المتوفرة لنصوص لوريم إيبسوم، ولكن الغالبية تم تعديلها بشكل ما عبر إدخال بعض النوادر أو الكلمات العشوائية إلى النص. إن كنت تريد أن تستخدم نص لوريم إيبسوم ما، عليك أن تتحقق أولاً أن ليس هناك أي كلمات أو عبارات محرجة أو غير لائقة مخبأة في هذا النص. بينما تعمل جميع مولّدات نصوص لوريم إيبسوم على الإنترنت على إعادة تكرار مقاطع من نص لوريم إيبسوم نفسه عدة مرات بما تتطلبه الحاجة، يقوم مولّدنا هذا باستخدام كلمات من قاموس يحوي على أكثر من 200 كلمة لا تينية، مضاف إليها مجموعة من الجمل النموذجية، لتكوين نص لوريم إيبسوم ذو شكل منطقي قريب إلى النص الحقيقي. وبالتالي يكون النص الناتح خالي من التكرار، أو أي كلمات أو عبارات غير لائقة أو ما شابه. وهذا ما يجعله أول مولّد نص لوريم إيبسوم حقيقي على الإنترنت.`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`عرض عن الشرق، وتتحمّل الموسوعة. بحق ترتيب الساحة اسبوعين تم, عل لإعادة يتعلّق بها, ومن إذ أفاق وباءت العظمى. زهاء فمرّ فهرست لم عدد, مع عرض حالية إبّان مشاركة. على عل وترك حالية جزيرتي, جعل وسوء الحكم للجزر هو. وبعد إحتار تكتيكاً أم مكن, جنوب المضي عسكرياً أخر بل, و كلا فبعد الشهير اليميني.`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`أثره، كُلفة الحدود من هذا, مئات بمعارضة الإحتفاظ أسر هو. دول أي أملاً بتخصيص, بل وبدأت والفلبين البولندي بحث. وصل بحشد قِبل والديون كل, ان كلا وإيطالي ايطاليا، الأوربيين, أي الأرواح والكساد الخارجية على. حين الأمم ويعزى التاريخ، أي, أم قام أراضي الشرقي ليرتفع.`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`فقد ثم إعمار لبلجيكا،, وحرمان الأوروبية ثم على, ميناء الولايات بل قام. إذ كانت بداية النفط أما, اتّجة الخطّة تزامناً تم إيو. لم وتنصيب المنتصر أوراقهم هذه. أم مدن لدحر طوكيو الأسيوي. أم أما السبب العناد الإتفاقية, أملاً والقرى الشهيرة بعد أم. أي لكون الستار الولايات بال, تم بلاده قتيل، فعل. على الشتاء واتّجه الهادي و.`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
	par.SetTextAlignment(creator.TextAlignmentRight)
	par.SetMargins(marginLeft, marginRight, marginTop, marginBottom)
	err = c.Draw(par)
	if err != nil {
		fmt.Println(err)
	}

	textArabicShaped, err = textshaping.ArabicShape(`هنالك العديد من الأنواع المتوفرة لنصوص لوريم إيبسوم، ولكن الغالبية تم تعديلها بشكل ما عبر إدخال بعض النوادر أو الكلمات العشوائية إلى النص. إن كنت تريد أن تستخدم نص لوريم إيبسوم ما، عليك أن تتحقق أولاً أن ليس هناك أي كلمات أو عبارات محرجة أو غير لائقة مخبأة في هذا النص. بينما تعمل جميع مولّدات نصوص لوريم إيبسوم على الإنترنت على إعادة تكرار مقاطع من نص لوريم إيبسوم نفسه عدة مرات بما تتطلبه الحاجة، يقوم مولّدنا هذا باستخدام كلمات من قاموس يحوي على أكثر من 200 كلمة لا تينية، مضاف إليها مجموعة من الجمل النموذجية، لتكوين نص لوريم إيبسوم ذو شكل منطقي قريب إلى النص الحقيقي. وبالتالي يكون النص الناتح خالي من التكرار، أو أي كلمات أو عبارات غير لائقة أو ما شابه. وهذا ما يجعله أول مولّد نص لوريم إيبسوم حقيقي على الإنترنت.`)
	if err != nil {
		panic(err)
	}

	par = c.NewParagraph(textArabicShaped)
	par.SetFont(arabicFont)
	par.SetFontSize(fontSize)
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
