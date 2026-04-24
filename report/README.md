# PDF Report Creation

These examples demonstrate how to generate PDF reports with UniPDF's creator package. They cover report layouts in both portrait and landscape orientation, building and styling tables, generating custom table of contents (including TOC pages with additional arbitrary content), and producing reports sourced from CSV data.

## Examples

- [pdf_custom_toc.go](pdf_custom_toc.go) The example showcases the capabilities of generating PDF custom table of contents layout. The output is saved as pdf-custom-toc.pdf.
- [pdf_custom_toc_with_content.go](pdf_custom_toc_with_content.go) The example extends the custom table of contents layout by adding extra arbitrary content on the TOC page. The output is saved as pdf-custom-toc_with_content.pdf.
- [pdf_report.go](pdf_report.go) The example showcases PDF report generation with UniPDF's creator package. The output is saved as unidoc-report.pdf which illustrates some of the features of the creator.
- [pdf_report_from_csv.go](pdf_report_from_csv.go) This example showcases how to prepare a report from CSV data. The output is saved as report_from_csv.pdf.
- [pdf_report_landscape.go](pdf_report_landscape.go) The example showcases PDF report generation in landscape mode with UniPDF's creator package. The output is saved as unidoc-report-landscape.pdf.
- [pdf_tables.go](pdf_tables.go) The example showcases PDF tables features using UniPDF's creator package. The output is saved as unipdf-tables.pdf which illustrates some of the features of the creator.