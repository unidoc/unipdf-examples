## Grayscale PDF example

- pdf_grayscale.go converts an entire PDF file to grayscale in a vectorized fashion,
including images and all content.

In case you are interested how exactly this is done, check [https://github.com/unidoc/unipdf-examples/blob/master/advanced/pdf_grayscale_transform.go](advanced example) demonstrating of the more complex capabilities of UniPDF, showing the capability to process and transform objects and contents.

The conversion applies to：
1. XObject images - and transforming through each color component to corresponding DeviceGray component.
2. Inline images.
3. Shapes such as lines as curves in the content stream.
4. Shadings and patterns.




