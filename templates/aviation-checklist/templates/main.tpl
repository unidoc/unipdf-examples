{{define "check-table"}}
<table enable-page-wrap="false" margin="0 0 10 0">
    <table-cell border-width="1" background-color="#000C66">
        <paragraph>
            <text-chunk font="helvetica-bold" color="#FFFFFF">{{ .Title }}</text-chunk>
        </paragraph>
    </table-cell>

    {{range $idx, $item := .Items }}
        {{template "check-entry" dict "WithBottomBorder" (eq $idx (len (slice $.Items 1))) "Item" $item}}
    {{end}}
</table>
{{end}}

{{define "check-entry"}}
    {{$text := ""}}
    {{$align := "left"}}
    {{$font := "helvetica"}}

    {{if .Item.Action}}
        {{$text = html .Item.Action}}
        {{$align = "center"}}
        {{$font = "helvetica-bold"}}
    {{else}}
        {{$text = html .Item.DisplayText}}
    {{end}}

    {{$borderBottomWidth := 0}}
    {{if .WithBottomBorder}}
        {{$borderBottomWidth = 1}}
    {{end}}

    <table-cell border-width-left="1" border-width-right="1" border-width-bottom="{{$borderBottomWidth}}" align="{{$align}}">
        <paragraph>
            <text-chunk font="{{$font}}">{{$text}}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

{{$bottomMargin := 10.0}}
{{$contentHeight := 0.0}}
{{$column := 0}}

<table columns="2">
    <table-cell>
        <division margin="0 5 0 0">
            {{range .Checks}}
                {{$tableHeight := add (calcTableHeight .) $bottomMargin}}
                {{$newHeight := add $contentHeight $tableHeight}}

                {{if not (isFitInPageHeight $newHeight)}}
                    {{$newHeight = $tableHeight}}
                    {{$margin := "0 5 0 0"}}

                    {{/* Close current cell */}}
                    </division>
                    </table-cell>

                    {{if eq $column 0}}}
                        {{$column = 1}}
                        {{$margin := "0 0 0 5"}}
                    {{else}}
                        {{/* Close current table and move to new table in new page */}}
                        </table>
                        <table columns="2">

                        {{$column = 0}}
                    {{end}}

                    {{/* Start new cell */}}
                    <table-cell>
                    <division margin="{{$margin}}">
                {{end}}

                {{template "check-table" . }}

                {{$contentHeight = $newHeight}}
            {{end}}
        </division>
    </table-cell>
</table>