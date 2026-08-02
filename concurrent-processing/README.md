# Concurrent Processing

UniPDF supports document level concurrent processing (this means processing each document separately in its own go routine) for all operations.

Page level concurrent processing is supported as well. The parser and reader I/O layer takes an `io.ReaderAt`, which makes the parser internally thread-safe, so a document can be loaded once and then read from several goroutines. What is safe once the document structure is loaded:

- `*core.PdfParser` cached object lookups, backed by a concurrent object cache.
- `*model.PdfReader` page access (`GetPage`) and content traversal.

The supported pattern is **load once, then read in parallel**. `model.NewPdfReaderFromParser` resolves the document structure eagerly and lets one parser back a reader that is shared across goroutines.

Readers built with `NewPdfReaderFromParser` do not support signature verification or encrypted documents, and cannot be passed to `model.NewPdfAppender`. Use `model.NewPdfReader` for those. The single threaded pattern (one reader, one goroutine) remains the simplest and most efficient choice for most workloads.

## Examples
- [concurrent_extraction.go](concurrent_extraction.go) Extracts text from multiple documents provided via the command line arguments concurrently and saves the result to a text file.
- [concurrent_extraction_page_level.go](concurrent_extraction_page_level.go) Extracts text from the document provided via the command line arguments concurrently on page level.
- [concurrent_extraction_shared_parser.go](concurrent_extraction_shared_parser.go) Extracts text on page level with all goroutines sharing a single parser created via `core.NewParser` and `model.NewPdfReaderFromParser`.
