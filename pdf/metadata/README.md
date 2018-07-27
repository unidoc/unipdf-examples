# Metadata

## Examples

- pdf_metadata_get_docinfo.go outputs the document information dictionary information
- pdf_metadata_get_xml.go outputs metadata streams XML 

## Background
According to section 14.3 Metadata (p. 556 in PDF32000_2008) metadata can be stores in two ways:
1. In metadata streams associated with the document or a component of the document (newer, preferred approach)
2. In a document information dictionary associated with the document (old way)

The document information dictionary has fixed field such as:
- Title
- Author
- Creator
- Producer
- Subject

## References

https://github.com/unidoc/unidoc/issues/164
