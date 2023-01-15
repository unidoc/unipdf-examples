# PDF Text Examples

The example explains how you can use UniPDF to perform a number of text manipulations such as detect signatures, extract and insert text. You can also search and replace text and find locations for specific text. Using UniPDF, you can also convert PDF to .csv files.

## Examples

### Generation
- [pdf_cmyk_color.go](pdf_cmyk_color.go) The example showcase how to use CMYK color to colorize text.
- [pdf_formatted_text.go](pdf_formatted_text.go) The example showcases the usage of styled paragraphs. The output is saved as styled_paragraph.pdf which illustrates some of the features of the creator.
- [pdf_insert_text.go](pdf_insert_text.go) The example showcases how to insert text to a specific page, location in a PDF file. If unsure about position, try getting the dimensions of a PDF with pdf/pages/pdf_page_info.go first or start with 0,0 (upper left corner) and increase to move right, down.
- [pdf_text_color.go](pdf_text_color.go) The example showcase how to use RGB and CMYK color to colorize text.
- [pdf_using_unicode_font.go](pdf_using_unicode_font.go) The example illustrates how to use composite font (CJK font) file to render a text and subset the font to create a small output file.

### Extraction or modifying PDF
- [pdf_detect_signature.go](pdf_detect_signature.go) The example highlights the basic functionality for text searching: Retrieving position of a signature line in PDF where the signature line is given by "__________________" text. And positioned with a Tm operation above.
- [pdf_search_replace.go](pdf_search_replace.go) The example highlights a basic example of find and replace with UniPDF.
- [pdf_text_locations.go](pdf_text_locations.go) The example highlights how to find mark up locations of substrings of extracted text in a PDF file.
- [pdf_text_vertical_alignment.go](pdf_text_vertical_alignment.go) The example highlights an example of setting the vertical alignment of text chunks in a paragraph.
- [pdf_to_csv.go](pdf_to_csv.go) The example is illustrating capability to extract TextMarks from PDF, and grouping together into words, rows and columns for CSV data extraction. The example includes debugging capabilities such as outputting a marked-up PDF showing bounding boxes of marks, words, lines and columns.
