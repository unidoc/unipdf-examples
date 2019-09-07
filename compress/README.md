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
}
```


## Examples

- [pdf_optimize.go](pdf_optimize.go) compresses a PDF file with some typical options.

