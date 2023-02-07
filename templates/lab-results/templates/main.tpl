{{define "test-result-header"}}
    <table-cell background-color="#1B7A98" align="center" vertical-align="middle">
        <paragraph text-align="center">
            <text-chunk font-size="10" color="#ffffff">{{ .Text }}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

{{define "test-result-row"}}
    <table-cell align="center" vertical-align="top" border-width-bottom="1" border-color="#1B7A98">
        <paragraph text-align="center" line-height="1.1">
            <text-chunk font-size="10">{{ .Text }}</text-chunk>
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

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">DOB(y/m/d):</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Patient.Birthdate }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Date Collected:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Specimen.Collected }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Ordering:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Physician.Ordering }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Age:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Patient.Age }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Date received:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Specimen.Received }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Referring:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Physician.Referring }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Gender:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Patient.Gender }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Date entered:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Specimen.Entered }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">ID:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Physician.Id }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Patient ID:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Patient.Id }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">Date reported:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Specimen.Reported }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell border-width-left="1">
            <paragraph margin="0 10">
                <text-chunk font="helvetica-bold" color="#1B7A98">NPI:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph margin="0 10 0 0">
                <text-chunk>{{ $.Physician.Npi }}</text-chunk>
            </paragraph>
        </table-cell>
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
                <text-chunk>Dir: Michael Tannenbaum, MD</text-chunk>
            </paragraph>
        </table-cell>
    </table>

    {{/* Make sure each test result is rendered on a new page */}}
    {{if lt $idx (len (slice $.Results 1)) }}
        <page-break></page-break>
    {{end}}
{{end}}