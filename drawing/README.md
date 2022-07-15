# Drawing on PDF
Drawing in PDF is done through vectorized operands in the content stream.
UniPDF provides functions to do this for common shapes, otherwise custom
objects can always be applied to the content stream directly.

## Examples

- [pdf_draw_shapes.go](pdf_draw_shapes.go) draws multiple shapes in a new PDF file.
- [pdf_draw_lines.go](pdf_draw_lines.go) showcases the capabilities of creator lines.
