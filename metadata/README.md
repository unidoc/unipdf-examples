# Metadata

## Examples

- pdf_metadata_get_docinfo.go outputs the document information dictionary information
- pdf_metadata_get_xml.go outputs metadata streams XML 
- pdf_metadata_set_docinfo.go showcase how to set a default and custom metadata information

## Background
According to section 14.3 Metadata (p. 556 in PDF32000_2008) metadata can be stores in two ways:
1. In metadata streams associated with the document or a component of the document (newer, preferred approach)
2. In a document information dictionary associated with the document (old way)

The document information dictionary has fixed field such as:
- Title  
  The document’s title.

- Author  
  The name of the person who created the document.

- Subject  
  The subject of the document.

- Keywords  
  Keywords associated with the document.

- Creator  
   If the document was converted to PDF from another format, the name of the conforming product that created the original document from which it was converted.

- Producer  
  If the document was converted to PDF from another format, the name of the conforming product that converted it to PDF.

- CreationDate  
  The date and time the document was created, in human- readable form.

- ModDate  
  The date and time the document was most recently modified, in human-readable form.

- Trapped  
  A name object indicating whether the document has been modified to include trapping information.
  

Aside from the default fields mentioned above, the metadata dictionary could contain other keys that hold custom string metadata.
