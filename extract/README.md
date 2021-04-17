Various extraction example code.

## Examples

- [extract_text_bound.go](extract_text_bound.go) The example showcases how to extract all text for each page along with it's boundary information.
- [pdf_extract_text.go](pdf_extract_text.go) The example showcases how to extract all text for each page of a PDF file.
- [reconstruct_text.go](reconstruct_text.go) Example that illustrates the accuracy of the text extraction, by first extracting all TextMarks and then reconstructing the text by writing out the text page-by-page to a new PDF with the creator package.
- [reconstruct_words.go](reconstruct_words.go) The example expands upon [reconstruct_text.go](reconstruct_text.go) to show word placements.