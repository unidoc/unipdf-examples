{{define "test-result-header"}}
    <table-cell background-color="#1B7A98" align="center" vertical-align="middle">
        <paragraph text-align="center">
            <text-chunk color="#ffffff">{{ .Text }}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

{{define "test-result-row"}}
    <table-cell align="center" border-width-bottom="1" border-color="#1B7A98">
        <paragraph text-align="center" line-height="1.1">
            <text-chunk>{{ .Text }}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

{{define "result-field"}}
    <table-cell border-width-left="1">
        <paragraph margin="0 10">
            <text-chunk color="#1B7A98" font="helvetica-bold">{{ .Title }}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell>
        <paragraph margin="0 10 0 0">
            <text-chunk>{{ .Value }}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

{{range $idx, $result := .Results}}
    <table columns="6" column-widths="0.15 0.15 0.17 0.13 0.13 0.17">
        <table-cell border-width-left="1" colspan="2">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold">Patient Details</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1" colspan="2">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold">Specimen Details</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1" colspan="2">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold">Physician Details</text-chunk>
            </paragraph>
        </table-cell>

        {{template "result-field" (dict "Title" "DOB(y/m/d):" "Value" $.Patient.Birthdate ) }}
        {{template "result-field" (dict "Title" "Date Collected:" "Value" $.Specimen.Collected ) }}
        {{template "result-field" (dict "Title" "Ordering:" "Value" $.Physician.Ordering ) }}
        {{template "result-field" (dict "Title" "Age:" "Value" $.Patient.Age ) }}
        {{template "result-field" (dict "Title" "Date received:" "Value" $.Specimen.Received ) }}
        {{template "result-field" (dict "Title" "Referring:" "Value" $.Physician.Referring ) }}
        {{template "result-field" (dict "Title" "Gender:" "Value" $.Patient.Gender ) }}
        {{template "result-field" (dict "Title" "Date entered:" "Value" $.Specimen.Entered ) }}
        {{template "result-field" (dict "Title" "ID:" "Value" $.Physician.Id ) }}
        {{template "result-field" (dict "Title" "Patient ID:" "Value" $.Patient.Id ) }}
        {{template "result-field" (dict "Title" "Date reported:" "Value" $.Specimen.Reported ) }}
        {{template "result-field" (dict "Title" "NPI:" "Value" $.Physician.Npi ) }}
    </table>

    <line position="relative" fit-mode="fill-width" thickness="1" color="#1B7A98" margin="10 0"></line>

    <table columns="2" indent="0">
        <table-cell>
            <paragraph>
                <text-chunk font="helvetica-bold" color="#1B7A98">General Comments and Additional Information</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font="helvetica-bold" color="#1B7A98">Ordered Items</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font="helvetica-bold">Clinical Info: </text-chunk>
                <text-chunk font="helvetica-bold" color="{{ (infoColor $result.Level) }}">{{ $result.Info }}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk>{{ $result.Test }}</text-chunk>
            </paragraph>
        </table-cell>
    </table>

    <table columns="6" column-widths="0.2 0.2 0.1 0.1 0.3 0.1" margin="10 0">
        {{template "test-result-header" (dict "Text" "TESTS") }}
        {{template "test-result-header" (dict "Text" "RESULT") }}
        {{template "test-result-header" (dict "Text" "FLAG") }}
        {{template "test-result-header" (dict "Text" "UNIT") }}
        {{template "test-result-header" (dict "Text" "REFERENCE INTERVAL") }}
        {{template "test-result-header" (dict "Text" "LAB") }}

        {{template "test-result-row" (dict "Text" $result.Test) }}
        {{template "test-result-row" (dict "Text" $result.Result) }}
        {{template "test-result-row" (dict "Text" $result.Flag) }}
        {{template "test-result-row" (dict "Text" $result.Units) }}
        {{template "test-result-row" (dict "Text" $result.RefInterval) }}
        {{template "test-result-row" (dict "Text" $result.Lab) }}
    </table>

    {{range $p := $result.Description}}
        <paragraph margin="0 0 10 0" line-height="1.1">
            <text-chunk>{{ $p }}</text-chunk>
        </paragraph>
    {{end}}

    <table columns="4" column-widths="0.1 0.1 0.4 0.4" margin="20 0 0 0">
        <table-cell border-width-bottom="1" border-width-top="1">
            <paragraph>
                <text-chunk>{{ $result.Lab }}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="1" border-width-top="1">
            <paragraph>
                <text-chunk>BN</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="1" border-width-top="1">
            <division margin="0 0 5 0">
                <paragraph>
                    <text-chunk>Medical Lab</text-chunk>
                </paragraph>
                <paragraph>
                    <text-chunk>1234 Main Street, New York, NY 12345-6789</text-chunk>
                </paragraph>
            </division>
        </table-cell>
        <table-cell border-width-bottom="1" border-width-top="1">
            <paragraph>
                <text-chunk>Dir: John Doe, MD</text-chunk>
            </paragraph>
        </table-cell>
    </table>

    {{/* Make sure each test result is rendered on a new page */}}
    {{if lt $idx (len (slice $.Results 1)) }}
        <page-break></page-break>
    {{end}}
{{end}}