# PDF redactor Examples

The example explains how a text is redacted from the a pdf file using the `redactor` package. 

To redact a given content from a pdf file:
- First prepare the regular expressions that can match the target texts.
- Then build `redactor.RedactionTerm` using the `regex` as shown in the [example](redact_text.go).
- Finaly initialize a `redactor.Redactor` object using `model.PdfReader`, `redactor.RedactionOptions` and `redactor.RectangleProps` to apply redaction on the given pdf file.

## Example

- [redact_text.go](redact_text.go) The example shows redaction of credit card numbers and emails from a pdf file using regex patterns.