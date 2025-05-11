/*
 * An example of writing Hebrew text
 *
 * Run as: go run pdf_write_hebrew.go
 */

package main

import (
	"os"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
	"github.com/unidoc/unipdf/v4/model"
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
	font, err := model.NewCompositePdfFontFromTTFFile("./OpenSans-Regular.ttf")
	if err != nil {
		panic(err)
	}

	c := creator.New()
	c.SetPageSize(creator.PageSizeA5)

	hebrew := `שונה תוכל חפש מה, דת חפש היום כדור מיזמים, היא מה גרמנית מרצועת. לעריכת האנציקלופדיה או קרן, ספינות מדויקים או לוח. של סדר ניהול ויקיפדיה. של המלחמה פיסיקה קישורים עזה. אם טבלאות ממונרכיה כתב, שתפו ישראל ננקטת אם אתה, אם חפש בהשחתה פסיכולוגיה. בדף על ברית ריקוד חבריכם, אחד אל ריקוד להפוך אווירונאוטיקה.`
	p := c.NewStyledParagraph()
	style := &p.Append(hebrew).Style
	style.Font = font

	err = c.Draw(p)
	if err != nil {
		panic(err)
	}

	hebrew = `ב אחד כלליים חופשית, עוד אירועים ותשובות האנציקלופדיה על. אל אתה מושגי הבהרה ויקימדיה. ב עזה חשמל בלשנות, מה החלה.`
	p = c.NewStyledParagraph()
	style = &p.Append(hebrew).Style
	style.Font = font

	err = c.Draw(p)
	if err != nil {
		panic(err)
	}

	if err := c.WriteToFile("hebrew-text.pdf"); err != nil {
		panic(err)
	}
}
