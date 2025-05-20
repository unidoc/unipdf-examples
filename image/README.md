# PDF Images

UniPDF allows you to add, extract, list images in your PDF documents. Having the ability to play around with images allows for the creation of attractive PDF reports.

## Examples

- [pdf_add_image_to_page.go](pdf_add_image_to_page.go) explains how to add an image in a PDF document
- [pdf_images_to_pdf.go](pdf_images_to_pdf.go) explains how to add multiple images in a PDF document, one image per page. 
- [pdf_list_images.go](pdf_list_images.go) explains how to list images in a PDF file. Passes through each page, goes through the content stream and finds instances of both. XObject Images and inline images. Also handles images referred within XObject Form content streams.
- Image extraction example can be found in [../extract](extract) example subfolder.
