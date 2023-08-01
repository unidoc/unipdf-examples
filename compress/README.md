# PDF compression (optimization)

Optimization of PDF output is implemented in the PDF writer of UniPDF and contains multiple options (optimize.Options)
```go
// Options describes PDF optimization parameters.
type Options struct {
	CombineDuplicateStreams         bool
	CombineDuplicateDirectObjects   bool
	ImageUpperPPI                   float64
	ImageQuality                    int
	UseObjectStreams                bool
	CombineIdenticalIndirectObjects bool
	CompressStreams                 bool
	CleanFonts                      bool
	SubsetFonts                     bool
	CleanContentstream              bool
	CleanUnusedResources            bool
}
```

From the available filters listed above, all of them except `ImageQuality` and `ImageUpperPPI` enable lossless compression.

## Examples

- [pdf_optimize.go](pdf_optimize.go) compresses a PDF file with some typical options.
- [pdf_font_subsetting.go](pdf_font_subsetting.go) illustrates how to reduce a PDF file size by subsetting all fonts used in the document using `SubsetFonts` Optimizer option.
- [pdf_remove_unused_resources.go](pdf_remove_unused_resources.go) reduces file size by removing unused resources such as Images, Xforms, fonts and external graphics state dictionaries.
