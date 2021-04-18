/*
 * Add side note on page's left and right margin area.
 *
 * Run as: go run pdf_page_side_note.go
 */

package main

import (
	"fmt"
	"log"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	c := creator.New()

	fontRegular, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fontBold, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Draw side note on each content page
	c.PageFinalize(func(args creator.PageFinalizeFunctionArgs) error {
		p := c.NewStyledParagraph()
		p.SetAngle(90)

		chunk := p.Append(fmt.Sprintf("Page %d/%d", args.PageNum, args.TotalPages))
		chunk.Style.FontSize = 14
		chunk.Style.Color = creator.ColorBlue

		if args.PageNum%2 != 0 {
			p.SetPos(args.PageWidth-p.Height()-10, (args.PageHeight-p.Width())/2)
		} else {
			p.SetPos(p.Height()+10, (args.PageHeight-p.Width())/2)
		}

		if err := c.Draw(p); err != nil {
			return err
		}

		return nil
	})

	// Create page one content.
	chap := c.NewChapter("Page One")
	chap.GetHeading().SetMargins(0, 0, 0, 20)
	chap.GetHeading().SetFont(fontBold)
	chap.GetHeading().SetFontSize(18)

	err = drawContent(c, chap, fontRegular, fontBold, 0, 70)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Draw chapter.
	if err = c.Draw(chap); err != nil {
		log.Fatal(err)
	}

	c.NewPage()

	// Create page two content.
	chap = c.NewChapter("Page Two")
	chap.GetHeading().SetMargins(70, 0, 0, 20)
	chap.GetHeading().SetFont(fontBold)
	chap.GetHeading().SetFontSize(18)

	err = drawContent(c, chap, fontRegular, fontBold, 70, 0)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Draw chapter.
	if err = c.Draw(chap); err != nil {
		log.Fatal(err)
	}

	// Write output file.
	err = c.WriteToFile("page_side_notes.pdf")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func drawContent(c *creator.Creator, ch *creator.Chapter,
	fontRegular, fontBold *model.PdfFont,
	leftMargin, rightMargin float64) error {

	p := c.NewStyledParagraph()
	p.SetMargins(leftMargin, rightMargin, 0, 20)
	p.SetLineHeight(1.1)
	p.SetText("Praesent pellentesque tincidunt odio sed rhoncus. Praesent finibus ligula sapien, quis semper quam porta a. Fusce pulvinar elit et augue hendrerit cursus. Duis nec vehicula magna. Pellentesque mattis nunc quis maximus convallis. Ut dictum malesuada pulvinar. Pellentesque et lectus mollis, semper nisi quis, eleifend odio. Suspendisse congue tellus sed convallis pellentesque. Aenean elementum, turpis sit amet posuere ultrices, risus nisl rutrum sapien, sed efficitur sem nisi ac ipsum. Quisque eu aliquet nisi. Aenean eu magna id ipsum venenatis malesuada at at lectus.")

	if err := ch.Add(p); err != nil {
		return err
	}

	p = c.NewStyledParagraph()
	p.SetMargins(leftMargin, rightMargin, 0, 20)
	p.SetLineHeight(1.1)
	p.SetText("Maecenas luctus mauris ac nisl auctor sagittis. Aenean sollicitudin pellentesque nibh, at convallis ex tincidunt vitae. Maecenas feugiat sit amet ex sit amet malesuada. Nunc mauris justo, posuere quis dolor euismod, ultrices ullamcorper lacus. Donec euismod, libero ut tincidunt consectetur, nulla nisi venenatis est, nec mattis felis erat a metus. In condimentum quam ut nibh malesuada, non dignissim elit dapibus. Integer justo arcu, tincidunt eu placerat eget, laoreet quis eros. Vivamus convallis mattis dolor eget dictum.")

	if err := ch.Add(p); err != nil {
		return err
	}

	p = c.NewStyledParagraph()
	p.SetMargins(leftMargin, rightMargin, 0, 20)
	p.SetLineHeight(1.1)
	p.SetText("Sed lacinia ex sit amet luctus sodales. Nulla eget molestie erat. Nullam quis nunc ornare, rutrum nunc volutpat, dapibus metus. Donec sit amet quam ornare, finibus libero ut, ullamcorper erat. Mauris mollis dolor non scelerisque consectetur. Integer euismod tincidunt vulputate. Nulla ante sapien, sodales quis urna a, pellentesque tincidunt erat. Vivamus euismod nisi eros, sed euismod justo faucibus scelerisque. Curabitur venenatis eros in dui finibus, vel dignissim lectus porta. Quisque pellentesque risus quis libero tempor, in eleifend tortor suscipit. Etiam quis ultrices justo.")

	if err := ch.Add(p); err != nil {
		return err
	}

	p = c.NewStyledParagraph()
	p.SetMargins(leftMargin, rightMargin, 0, 20)
	p.SetLineHeight(1.1)
	p.SetText("Etiam volutpat elit et sem pharetra, et hendrerit orci eleifend. Sed eu consequat erat. Donec in augue vitae felis scelerisque pulvinar. Sed hendrerit faucibus finibus. Curabitur risus est, suscipit eu varius vel, tempus eget enim. Morbi magna tellus, rhoncus ut eros ac, sagittis faucibus orci. Duis quis fringilla mi. Donec in dolor nec nisi commodo viverra.")

	if err := ch.Add(p); err != nil {
		return err
	}

	p = c.NewStyledParagraph()
	p.SetMargins(leftMargin, rightMargin, 0, 20)
	p.SetLineHeight(1.1)
	p.SetText("Nunc tincidunt mauris semper arcu vulputate semper. Donec in scelerisque velit. Phasellus vel enim eget leo dictum tempor at sit amet purus. Aliquam et nulla nulla. Morbi posuere commodo ex, nec dignissim elit laoreet vitae. Praesent fringilla tempor pharetra. Quisque ut imperdiet odio. In dictum dui vel lectus congue, feugiat eleifend mi convallis. Sed at massa et nunc auctor gravida. Maecenas varius nec nisi et accumsan. Donec fringilla vel nunc ac ultrices. Sed bibendum mauris tempor leo rhoncus fringilla. Maecenas vehicula lorem viverra dignissim elementum. Vivamus euismod lorem vel purus semper, vulputate finibus mauris auctor. Nam fermentum sapien nec dignissim convallis.")

	if err := ch.Add(p); err != nil {
		return err
	}

	return nil
}
