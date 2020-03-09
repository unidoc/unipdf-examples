# PDF Outlines

The example explains how to work with outlines (bookmarks) using UniPDF. 

## Examples

- [pdf_get_outlines.go](pdf_get_outlines.go) explains how to retrieve outlines (bookmarks) from a PDF file and prints them out in JSON format. Note: The JSON output can be used with the related pdf_set_outlines.go example to apply outlines to a PDF file.
- [pdf_set_outlines.go](pdf_set_outlines.go) explains how to apply outlines to a PDF file. The files are read from a JSON formatted file, which can be created via pdf_get_outlines which outputs outlines for an input PDF file in the JSON format.   