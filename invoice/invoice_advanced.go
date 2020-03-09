/*
 * This example showcases the creation of an advanced invoice with custom fields
 * The output is saved as invoice_advanced.pdf
 * The invoice contains improved formatting and styling as compared to simple invoice
 * You can easily create a personalized invoice following the code given below
 */

package main

import (
	"fmt"
	"log"

	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func main() {
	fontHelvetica := model.NewStandard14FontMustCompile(model.HelveticaName)

	c := creator.New()
	c.NewPage()

	logo, err := c.NewImageFromFile("unidoc-logo.png")
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	white := creator.ColorRGBFrom8bit(255, 255, 255)
	lightBlue := creator.ColorRGBFrom8bit(217, 240, 250)
	blue := creator.ColorRGBFrom8bit(2, 136, 209)

	invoice := c.NewInvoice()

	// Set invoice title.
	invoice.SetTitle("Unidoc Invoice")

	// Customize invoice title style.
	titleStyle := invoice.TitleStyle()
	titleStyle.Color = blue
	titleStyle.Font = fontHelvetica
	titleStyle.FontSize = 30

	invoice.SetTitleStyle(titleStyle)

	// Set invoice logo.
	invoice.SetLogo(logo)

	// Set invoice information.
	invoice.SetNumber("0001")
	invoice.SetDate("28/07/2016")
	invoice.SetDueDate("28/07/2016")
	invoice.AddInfo("Payment terms", "Due on receipt")
	invoice.AddInfo("Paid", "No")

	// Customize invoice information styles.
	for _, info := range invoice.InfoLines() {
		descCell, contentCell := info[0], info[1]
		descCell.BackgroundColor = lightBlue
		contentCell.TextStyle.Font = fontHelvetica
	}

	// Set invoice addresses.
	invoice.SetSellerAddress(&creator.InvoiceAddress{
		Name:    "John Doe",
		Street:  "8 Elm Street",
		City:    "Cambridge",
		Zip:     "CB14DH",
		Country: "United Kingdom",
		Phone:   "xxx-xxx-xxxx",
		Email:   "johndoe@email.com",
	})

	invoice.SetBuyerAddress(&creator.InvoiceAddress{
		Name:    "Jane Doe",
		Street:  "9 Elm Street",
		City:    "London",
		Zip:     "LB15FH",
		Country: "United Kingdom",
		Phone:   "xxx-xxx-xxxx",
		Email:   "janedoe@email.com",
	})

	// Customize address styles.
	addressStyle := invoice.AddressStyle()
	addressStyle.Font = fontHelvetica
	addressStyle.FontSize = 9

	addressHeadingStyle := invoice.AddressHeadingStyle()
	addressHeadingStyle.Color = blue
	addressHeadingStyle.Font = fontHelvetica
	addressHeadingStyle.FontSize = 16

	invoice.SetAddressStyle(addressStyle)
	invoice.SetAddressHeadingStyle(addressHeadingStyle)

	// Insert new column.
	col := invoice.InsertColumn(2, "Discount")
	col.Alignment = creator.CellHorizontalAlignmentRight

	// Customize column styles.
	for _, column := range invoice.Columns() {
		column.BackgroundColor = lightBlue
		column.BorderColor = lightBlue
		column.TextStyle.FontSize = 9
	}

	for i := 0; i < 7; i++ {
		cells := invoice.AddLine(
			fmt.Sprintf("Test product #%d", i+1),
			"1",
			"0%",
			"$10",
			"$10",
		)

		for _, cell := range cells {
			cell.BorderColor = white
			cell.TextStyle.FontSize = 9
		}
	}

	// Customize total line styles.
	titleCell, contentCell := invoice.Total()
	titleCell.BackgroundColor = lightBlue
	titleCell.BorderColor = lightBlue
	contentCell.BackgroundColor = lightBlue
	contentCell.BorderColor = lightBlue

	invoice.SetSubtotal("$100.00")
	invoice.AddTotalLine("Tax (10%)", "$10.00")
	invoice.AddTotalLine("Shipping", "$5.00")
	invoice.SetTotal("$85.00")

	// Set invoice content sections.
	invoice.SetNotes("Notes", "Thank you for your business.")
	invoice.SetTerms("Terms and conditions", "Full refund for 60 days after purchase.")
	invoice.AddSection("Custom section", "This is a custom section.")

	// Customize note styles.
	noteStyle := invoice.NoteStyle()
	noteStyle.Font = fontHelvetica
	noteStyle.FontSize = 12

	noteHeadingStyle := invoice.NoteHeadingStyle()
	noteHeadingStyle.Color = blue
	noteHeadingStyle.Font = fontHelvetica
	noteHeadingStyle.FontSize = 14

	invoice.SetNoteStyle(noteStyle)
	invoice.SetNoteHeadingStyle(noteHeadingStyle)

	if err = c.Draw(invoice); err != nil {
		log.Fatalf("Error drawing: %v", err)
	}

	// Write output file.
	err = c.WriteToFile("invoice_advanced.pdf")
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
}
