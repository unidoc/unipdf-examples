# Drawing on PDF
Drawing in PDF is done through vectorized operands in the content stream.
UniPDF provides functions to do this for common shapes, otherwise custom
objects can always be applied to the content stream directly.

## Examples

- [pdf_draw_shapes.go](pdf_draw_shapes.go) draws multiple shapes in a new PDF file.
