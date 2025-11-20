# PDF OCR Examples

UniPDF supports integration with HTTP-based OCR (Optical Character Recognition) services to extract text from images and scanned PDF documents. These examples demonstrate how to configure and use OCR services to process images and reconstruct searchable PDFs from scanned documents.

The OCR functionality works by sending images to a configured HTTP endpoint that performs text recognition and returns the results in various formats including plain text and HOCR (HTML-based OCR format).

## Examples

- [hocr_sample.go](hocr_sample.go) illustrates how to process HOCR formatted OCR output, parsing word-level information including bounding boxes and confidence scores.
- [ocr_batch.go](ocr_batch.go) shows how to perform batch OCR processing on multiple images concurrently, with error handling and summary reporting.
- [ocr_sample.go](ocr_sample.go) demonstrates basic OCR usage by sending a single image to an HTTP OCR service and extracting the text content.
- [reconstruct_pdf_from_hocr.go](reconstruct_pdf_from_hocr.go) demonstrates a complete workflow to extract images from a PDF, perform OCR with HOCR output, parse the structured results, and reconstruct a searchable PDF with properly positioned text.

## Requirements

These examples require an HTTP OCR service running on `http://localhost:8080/file`. The examples are created using [unidoc/ocrserver](https://github.com/unidoc/ocrserver) as the OCR service. However, UniPDF's OCR API is designed to be flexible and should support other OCR services that accept image uploads via multipart form data and return text or HOCR formatted results.
