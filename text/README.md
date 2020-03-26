# PDF Text Examples

The example explains how you can use UniPDF to perform a number of text manipulations such as detect signatures, extract and insert text. You can also search and replace text and find locations for specific text. Using UniPDF, you can also convert PDF to .csv files.

## Examples

- [pdf_detect_signature.go](pdf_detect_signature.go) The example highlights the basic functionality for text searching: Retrieving position of a signature line in PDF where the signature line is given by "__________________" text. And positioned with a Tm operation above.
- [pdf_extract_text.go](pdf_extract_text.go) The example showcases how to extract all text for each page of a PDF file.
- [pdf_insert_text.go](pdf_insert_text.go) The example showcases how to insert text to a specific page, location in a PDF file. If unsure about position, try getting the dimensions of a PDF with pdf/pages/pdf_page_info.go first or start with 0,0 (upper left corner) and increase to move right, down.
- [pdf_search_replace.go](pdf_search_replace.go) The example highlights a basic example of find and replace with UniPDF.
- [pdf_text_locations.go](pdf_text_locations.go) The example highlights how to find mark up locations of substrings of extracted text in a PDF file.
- [pdf_to_csv.go](pdf_to_csv.go) The example is illustrating capability to extract TextMarks from PDF, and grouping together into words, rows and columns for CSV data extraction. The example includes debugging capabilities such as outputting a marked-up PDF showing bounding boxes of marks, words, lines and columns.