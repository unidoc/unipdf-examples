Various extraction example code.

## Examples

- [extract_text_bound.go](extract_text_bound.go) The example showcases how to extract all text for each page along with it's boundary information.
- [pdf_extract_location.go](pdf_extract_location.go) The example showcases how to extract text at certain location.
- [pdf_extract_text.go](pdf_extract_text.go) The example showcases how to extract all text for each page of a PDF file.
- [pdf_simple_extraction.go](pdf_simple_extraction.go) The example showcases how to use a simple extraction procedure in the text extraction process, which skips table processing and topological ordering. This procedure could be useful in some cases where full extraction process giving unwanted extraction result, especially when processing tables.
- [pdf_tables.go](pdf_tables.go) The example showcase how to extract all tables from the specified pages of one or more PDF files.
- [reconstruct_text.go](reconstruct_text.go) Example that illustrates the accuracy of the text extraction, by first extracting all TextMarks and then reconstructing the text by writing out the text page-by-page to a new PDF with the creator package.
- [reconstruct_words.go](reconstruct_words.go) The example expands upon [reconstruct_text.go](reconstruct_text.go) to show word placements.
- [pdf_extract_images.go](pdf_extract_images.go) explains how to extract images from an existing PDF. The code passes through each page, goes through the content stream and finds XObject Images and inline images. Also handles images referred within XObject Form content streams. The output files are saved as a zip archive.
