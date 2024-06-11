# concurrent processing

UniPDF supports document level concurrent processing. This means processing each document separately in its own go routine.   
The concurrency is supported on document level for now. Page level concurrency in UniPDF is not safe yet.

## Examples
- [concurrent_extraction.go](concurrent_extraction.go) Extracts text from multiple documents provided via the command line arguments concurrently and saves the result to a text file.