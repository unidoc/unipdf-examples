# PDF Forms

Forms and fields in PDF enables creating interactive forms as well as on the client side, filling in and submitting
forms.

## Examples

- [pdf_form_add.go](pdf_form_add.go) illustates adding a basic form to a document.
- [pdf_form_fill_fdf_merge.go](pdf_form_fill_fdf_merge.go) illustates FDF merging - merging FDF form data (values) with a template PDF, producing a flattened output PDF (with appearances streams generated).
- [pdf_form_fill_json.go](pdf_form_fill_json.go) supports exporting form data as JSON as well filling form and outputting a flattened PDF (see below).
- [pdf_form_flatten.go](pdf_form_flatten.go) flattens a form, making the fields part of the document and no longer editable.
- [fdf_fields_info.go](fdf_fields_info.go) outputs information about fields in a Field Data Format (FDF) file.
- [pdf_form_get_field_data.go](pdf_form_get_field_data.go) gets field data for a single field by field name.
- [pdf_form_list_fields.go](pdf_form_list_fields.go) lists form fields in a PDF.

## Use cases

1. Conveniently export form data as JSON to file:
```bash
$ ./bin/pdf_form_fill_json example.pdf > fields.json
[DEBUG]  parser.go:747 Pdf version 1.6
```
Contents of `fields.json`
```json
[
    {
        "name": "HIGH SCHOOL DIPLOMA",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "TRADE CERTIFICATE",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "COLLEGE NO DEGREE",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "PHD",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "OTHER DOCTORATE",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "ASSOCIATES DEGREE",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "MASTERS DEGREE",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "PROFESSIONAL DEGREE",
        "value": "Off",
        "options": [
            "Off",
            "On"
        ]
    },
    {
        "name": "STATE",
        "value": "WI"
    },
    {
        "name": "ZIP",
        "value": "30231"
    },
    {
        "name": "Name_Last",
        "value": "Johnsson"
    },
    {
        "name": "Name_First",
        "value": "John"
    },
    {
        "name": "Name_Middle",
        "value": "K."
    },
]
```

2. Edit fields data, simply by altering the values in the JSON file.


3. Import as JSON back and write out as flattened output file.

```bash
$ ./bin/pdf_form_fill_json ~/wh/Documents/UniDoc/bench/forms/interactiveform_filled.pdf fdata.json filled.pdf
```

The output filled.pdf is flattened so that it is no longer editable.


