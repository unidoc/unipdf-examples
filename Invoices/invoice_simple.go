/*
 * This example showcases the creation of a simple invoice
 * The output is saved as invoice_simple.pdf
 * For a more advance experience, check out the invoice_advanced.go example
 */

package main

import (
	"fmt"
	"log"

	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	c := creator.New()
	c.NewPage()

	logo, err := c.NewImageFromFile("unidoc-logo.png")
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}

	invoice := c.NewInvoice()

	// Set invoice logo.
	invoice.SetLogo(logo)

	// Set invoice information.
	invoice.SetNumber("0001")
	invoice.SetDate("28/07/2016")
	invoice.SetDueDate("28/07/2016")
	invoice.AddInfo("Payment terms", "Due on receipt")
	invoice.AddInfo("Paid", "No")

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

	// Add invoice line items.
	for i := 0; i < 75; i++ {
		invoice.AddLine(
			fmt.Sprintf("Test product #%d", i+1),
			"1",
			"$10",
			"$10",
		)
	}

	// Set invoice totals.
	invoice.SetSubtotal("$100.00")
	invoice.AddTotalLine("Tax (10%)", "$10.00")
	invoice.AddTotalLine("Shipping", "$5.00")
	invoice.SetTotal("$115.00")

	// Set invoice content sections.
	invoice.SetNotes("Notes", "Thank you for your business.")
	invoice.SetTerms("Terms and conditions", "Full refund for 60 days after purchase.")

	if err = c.Draw(invoice); err != nil {
		log.Fatalf("Error drawing: %v", err)
	}

	// Write output file.
	err = c.WriteToFile("invoice_simple.pdf")
	if err != nil {
		log.Fatalf("Fail: %v\n", err)
	}
}
