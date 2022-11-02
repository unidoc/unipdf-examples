# Bullet list extractor examples

This example demonstrates how to extract  bullet point lists from a pdf document.
bullet point lists are lists that are marked with bullet point characters, with sequential numbering system i.e roman numerals, english alphabet, or hindu arabic numerals (1,2,3).
The extracted content only contains the texts that are part of the bullet list and related sub lists. 
The rest of contents is ignored.

To extract bullet point lists from a given pdf page:
- First construct an `extractor` from a given pdf page.
- Then get the `PageText` object by calling `ExtractPageText()` method.
- Finally get the slice of bullet point lists from `PageText.List()`. </br>
This returns a slice of lists that have a tree structure.
```go
ex, err := extractor.New(pdfPage)
// handle error
pageText, _, _, err := ex.ExtractPageText()
// handle another error
bulletLists := pageText.List()
```

The textual representation is found by calling the `Text()` method on the bulletLists.

By default the bullet list in a given pdf page are extracted in two ways. 
1. Using document accessibility tags if the file is tagged
2. If tags are not available then simple regex matching is performed on the raw text of that page to detect the possible bullet point lists.

However this behavior can be changed using the `DisableDocumentTags` attribute of `Options` object in extractor package.

This can be done as follows.
```go
    options := &extractor.Options{
		DisableDocumentTags: true,
	}
	ex, err := extractor.NewWithOptions(pdfPage, options)
```
If the pdf document was tagged poorly, especially if the bullet lists were not tagged with appropriate tag structures, then the second way of extracting may give a better result. This way we can take advantage of both extraction methods. 