## Advanced UniPDF Examples

- pdf_grayscale_transform.go converts an entire PDF file to grayscale in a vectorized fashion,
including images and all content.

This advanced example demonstrates some of the more complex capabilities of UniPDF, showing the capability to process
and transform objects and contents.

In case you are interested only in convert to grayscale functionality itself, use [https://github.com/unidoc/unipdf-examples/blob/master/grayscale/pdf_grayscale.go](the simple example)

The conversion applies to：
1. XObject images - and transforming through each color component to corresponding DeviceGray component.
2. Inline images.
3. Shapes such as lines as curves in the content stream.
4. Shadings and patterns.




