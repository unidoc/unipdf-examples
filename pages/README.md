# PDF Page Manipulation

The example explains how to work with and manipulate PDF pages using UniPDF library. You can perform a number of manipulations such as crop, merge, rotate and many more. 

## Examples

- [pdf_4up.go](pdf_4up.go) The example outputs multiple pages (4) per page to an output PDF from an input PDF. Showcases page templating by loading pages as Blocks and manipulating with the creator package.
- [pdf_crop.go](pdf_crop.go) The example Crop pages in a PDF file. Crops the view to a certain percentage of the original. The percentage specifies the trim-off percentage, both widthwise and heightwise.
- [pdf_merge.go](pdf_merge.go) The example highlights basic merging of PDF files. Simply loads all pages for each file and writes to the output file.
- [pdf_merge_advanced.go](pdf_merge_advanced.go) The example merges PDF files, including form field data (AcroForms). For a more basic merging of PDF page contents, see pdf_merge.go.
- [pdf_page_info.go](pdf_page_info.go) The example prints PDF page info: Mediabox size and other parameters. If [page num] is not specified prints out info for all pages.
- [pdf_rotate_flatten.go](pdf_rotate_flatten.go) The example rotates the contents of a PDF file in accordance with each page's Rotate entry and then sets Rotate to 0. That is, flattens the rotation. Will look the same in viewer, but when working with the PDF, the upper left corner will be the origin (in unidoc coordinate system).
- [pdf_rotate.go](pdf_rotate.go) The example rotate pages in a PDF file. Degrees needs to be a multiple of 90. Example of how to manipulate pages with the pdf creator.
- [pdf_split.go](pdf_split.go) The example highlights basic PDF split example: Splitting by page range.
- [pdf_split_advanced.go](pdf_split_advanced.go) The example highlights advanced PDF split example: Takes into account optional content - OCProperties (rarely used).
