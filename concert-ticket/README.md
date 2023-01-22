Sample concert ticket
======================

This example shows how to create a concert ticket using creator templates. Templates are used to represent page components and their layout using XML based markup. This helps in creating pdf pages with less code and adds reusability to the already defined page layout just by updating the data and with fewer or no modifications to the template file.

Here is how the concert ticket is created using the creator templates.

- Define the page layout using template (.tpl) file [main.tpl](templates/main.tpl)
- Prepare the data that will be used on the ticket [concert-ticket.json](templates/concert-ticket.json)
- Prepare the images and other resources.
- Draw the template using the `Creator` object provided in the creator package. [pdf_concert_ticket.go](pdf_concert_ticket.go)
  
```go
    c := creator.New()
    // Modify margins and page size using the creator `c`
    // Read template data file.
    ticket, err := readTemplateData("./templates/concert-ticket.json")
    // Handle errors

    // Read template file.
	tpl, err := readTemplate("./templates/main.tpl")
    // Handle errors.
    // Create template options here if necessary

    // Draw template.
    if err := c.DrawTemplate(tpl, ticket, tplOpts)
    // Handle errors.
    
    // Write document to file.
    err := c.WriteToFile("unipdf_ticket.pdf")
    // Handle errors here.
```

