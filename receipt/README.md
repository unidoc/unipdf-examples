Sample receipt document
======================
The creator package is used to generate pdf pages using components and simple interfaces. It also provides templates to easily use page components and manipulate their layout with fewer code and generate pdf documents fast. The templates use XML based markup to represent page components and they are converted to creator components at runtime.
This example shows creating a sample receipt document using the creator templates.  
Generally the flow of creating pdf documents using templates goes as follows.

- Define the components and their layouts using template file. 
 [template file](templates/main.tpl)
- Read the template file.
    ```go
    tpl, err := readTemplate("./templates/main.tpl")
    ```
- Instantiate new `Creator` object and draw the template.  
  [pdf_receipt.go](pdf_receipt.go) 
  ```go
    c := creator.New()
    // customize page size and margins of the page using `c` here.
    // prepare your data (either from json or any other source) and provide it to DrawTemplate()
    // prepare `creator.TemplateOptions` if necessary and provide it to DrawTemplate using `opts`
    if err := c.DrawTemplate(tpl, data, opts)
    // handle errors
    ```
- Write the document to file.
    ```go 
    err := c.WriteToFile("your-file-name.pdf")
    // handle errors
    ```
And the generated file can be seen below.
<p>
    <img src="./templates/res/screenshot-prev.png" alt="sample-receipt-preview" width="100%" />
</p>