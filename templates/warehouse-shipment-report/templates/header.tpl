{{define "company-info"}}
    <paragraph line-height="1.3">
        <text-chunk font-size="12">{{ . }}</text-chunk>
    </paragraph>
{{end}}

<division margin="25 25 0 25">
    <table columns="3" column-widths="0.15 0.55 0.3">
        <table-cell rowspan="2" border-color="#994D00" border-width-bottom="2">
            <image src="path('templates/res/logo.png')" width="70" height="70"></image>
        </table-cell>

        <table-cell>
            <paragraph>
                <text-chunk color="#3C1C00" font-size="25" font="helvetica-bold">{{ .Brand }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell rowspan="2" border-color="#994D00" border-width-bottom="2">
            <division>
                {{template "company-info" .Name }}
                {{template "company-info" .Address }}
                {{template "company-info" (printf "%s, %s" .Area .City) }}
                {{template "company-info" .Country }}
            </division>
        </table-cell>

        <table-cell border-color="#994D00" border-width-bottom="2">
            <paragraph margin="0 0 10 0">
                <text-chunk color="#F5A623" font-size="25" font="helvetica-bold">Warehouse shipments</text-chunk>
            </paragraph>
        </table-cell>
    </table>
</division>