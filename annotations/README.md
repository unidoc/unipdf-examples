# PDF Annotations

Annotations are different from normal PDF page contents and are stored separately.  They are intended to mark up
document contents, without changing the original document contents.

Support for creating annotations in UniPDF is through the unipdf/annotator package.
The annotator package creates the Annotation object and also an appearance stream, which is required to make
the annotation look the same in all viewers.
UniPDF's model package has support for all types of PDF annotations, whereas the annotator currently supports Square,
Circle, Line annotations.  Support for more annotation types will be added over time.  If you need support for an
unsupported annotation, please create an issue in this repository.

The examples in this folder illustrate a few capabilities for creating ellipse, lines, rectangles.

## Examples

- [pdf_annotate_add_ellipse.go](pdf_annotate_add_ellipse.go) adds a circle/ellipse annotation to a specified location on a page.
- [pdf_annotate_add_line.go](pdf_annotate_add_line.go) adds a line with arrowhead between two specified points on a page.
- [pdf_annotate_add_rectangle.go](pdf_annotate_add_rectangle.go) adds a rectangle annotation to a specified location on a page.
- [pdf_annotate_add_text.go](pdf_annotate_add_text.go) adds a text annotation with a user specified string to a fixed location on every page.
- [pdf_list_annotations.go](pdf_list_annotations.go) lists all annotations in a PDF file.

