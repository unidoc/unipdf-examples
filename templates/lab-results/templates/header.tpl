{{define "info-field"}}
    <table-cell>
        <paragraph>
            <text-chunk color="#1B7A98" font-size="12">{{ .Title }}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell>
        <paragraph>
            <text-chunk font-size="12">{{ .Value }}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

<division margin="10">
    <table columns="2">
        <table-cell border-width-bottom="2" border-color="#000000">
            <image src="path('templates/res/logo.png')" margin="0 0 5 5"></image>
        </table-cell>

        <table-cell vertical-align="bottom" border-width-bottom="2" border-color="#000000">
            <paragraph margin="5 0">
                <text-chunk font="helvetica-bold" font-size="14">SAMPLE REPORT</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell indent="0">
            <table columns="2" column-widths="0.3 0.7">
                {{template "info-field" (dict "Title" "Specimen ID:" "Value" .SpecimenID) }}
                {{template "info-field" (dict "Title" "Control ID:" "Value" .ControlID) }}
                {{template "info-field" (dict "Title" "Acct #:" "Value" .AcctNum) }}
                {{template "info-field" (dict "Title" "Phone:" "Value" .Phone) }}
                {{template "info-field" (dict "Title" "Rte:" "Value" .Rte) }}
            </table>
        </table-cell>

        <table-cell>
            <paragraph>
                <text-chunk font-size="12">SampleLab Test Master
Test Account
1234 Millstream Road
CLEANSVILLE NC 12345</text-chunk>
            </paragraph>
        </table-cell>
    </table>
</division>