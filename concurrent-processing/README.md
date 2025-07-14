# Concurrent Processing

UniPDF supports document level concurrent processing (this means processing each document separately in its own go routine) for all operations. 
Page level concurrent processing is now supported by UniPDF for rendering to image and text extraction.


## Examples
- [concurrent_extraction.go](concurrent_extraction.go) Extracts text from multiple documents provided via the command line arguments concurrently and saves the result to a text file.
- [concurrent_extraction_page_level.go](concurrent_extraction_page_level.go) Extracts text from the document provided via the command line arguments concurrently on page level.