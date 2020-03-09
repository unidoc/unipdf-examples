# PDF Testing Examples

The example explains how you can use UniPDF to test your PDF documents for any errors. The examples run through the documents and check for any errors produced. 

## Examples

- [pdf_count_color_pages_bench.go](pdf_count_color_pages_bench.go) The example detects the number of pages and the color pages (1-offset) all pages in a list of PDF files. Compares these results to running Ghostscript on the PDF files and reports an error if the results don't match.
- [pdf_grayscale_convert_bench.go](pdf_grayscale_convert_bench.go) The example showcases how to transform all content streams in all pages in a list of pdf files. This will transform all .pdf file in testdata and write the results to output.
- [pdf_passthrough_bench.go](pdf_passthrough_bench.go) The example showcases how to perform the pass through benchmark on all pdf files and write results to stdout.

