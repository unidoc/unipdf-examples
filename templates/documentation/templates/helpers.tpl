<!-- Reusable sub-templates -->

{{define "paragraph"}}
    {{$margin := "0"}}
    {{if .Margin}} {{$margin = .Margin}} {{end}}

    {{$lineHeight := 1}}
    {{if .LineHeight}} {{$lineHeight = .LineHeight}} {{end}}

    {{$align := "left"}}
    {{if .TextAlign}} {{$align = .TextAlign}} {{end}}

    {{$verticalAlign := "baseline"}}
    {{if .VerticalTextAlign}} {{$verticalAlign = .VerticalTextAlign}} {{end}}

    {{$font := "helvetica"}}
    {{if .Font}} {{$font = .Font}} {{end}}

    {{$fontSize := 10}}
    {{if .FontSize}} {{$fontSize = .FontSize}} {{end}}

    {{$textColor := "#000000"}}
    {{if .TextColor}} {{$textColor = .TextColor}} {{end}}

    {{$text := ""}}
    {{if .Text}} {{$text = .Text}} {{end}}

    <paragraph margin="{{$margin}}" line-height="{{$lineHeight}}" text-align="{{$align}}" vertical-text-align="{{$verticalAlign}}">
        <text-chunk font="{{$font}}" font-size="{{$fontSize}}" color="{{$textColor}}">{{$text}}</text-chunk>
    </paragraph>
{{end}}

{{define "table-cell-paragraph"}}
    {{$colspan := 1}}
    {{if .Colspan}} {{$colspan = .Colspan}} {{end}}

    {{$rowspan := 1}}
    {{if .Rowspan}} {{$rowspan = .Rowspan}} {{end}}

    {{$backgroundColor := "#ffffff"}}
    {{if .BackgroundColor}} {{$backgroundColor = .BackgroundColor}} {{end}}

    {{$align := "left"}}
    {{if .Align}} {{$align = .Align}} {{end}}

    {{$verticalAlign := "top"}}
    {{if .VerticalAlign}} {{$verticalAlign = .VerticalAlign}} {{end}}

    {{$borderColor := "#ffffff"}}
    {{if .BorderColor}} {{$borderColor = .BorderColor}} {{end}}

    {{$borderLeftSize := 0}}
    {{if .BorderLeftSize}} {{$borderLeftSize = .BorderLeftSize}} {{end}}

    {{$borderRightSize := 0}}
    {{if .BorderRightSize}} {{$borderRightSize = .BorderRightSize}} {{end}}

    {{$borderTopSize := 0}}
    {{if .BorderTopSize}} {{$borderTopSize = .BorderTopSize}} {{end}}

    {{$borderBottomSize := 0}}
    {{if .BorderBottomSize}} {{$borderBottomSize = .BorderBottomSize}} {{end}}

    {{$indent := 0}}
    {{if .Indent}} {{$indent = .Indent}} {{end}}

    <table-cell colspan="{{$colspan}}" rowspan="{{$rowspan}}" background-color="{{$backgroundColor}}" align="{{$align}}" vertical-align="{{$verticalAlign}}" border-color="{{$borderColor}}" border-width-left="{{$borderLeftSize}}" border-width-right="{{$borderRightSize}}" border-width-top="{{$borderTopSize}}" border-width-bottom="{{$borderBottomSize}}" indent="{{$indent}}">
        {{template "paragraph" .}}
    </table-cell>
{{end}}

{{define "chapter-title"}}
    {{$fillColor := "primary-bg-gradient"}}
    {{if .FillColor}} {{$fillColor = .FillColor}} {{end}}

    {{$textColor := "white"}}
    {{if .TextColor}} {{$textColor = .TextColor}} {{end}}

    {{$alternateText := .Text}}
    {{if .AlternateText}} {{$alternateText = .AlternateText}} {{end}}

    <chapter-heading color="white" font-size="10">{{.Text}}</chapter-heading>
    <division padding="4 0 8 0" margin="-10 0 0 0">
        <background fill-color="{{$fillColor}}"></background>
        {{template "paragraph" dict "TextAlign" "center" "Font" "helvetica-bold" "FontSize" 14 "TextColor" $textColor "Text" (strToUpper $alternateText)}}
    </division>
{{end}}

{{define "code-block"}}
    {{$columns := "2"}}
    {{if .Columns}} {{$columns = .Columns}} {{end}}

    {{$margin := "0"}}
    {{if .Margin}} {{$margin = .Margin}} {{end}}

    {{$codeLabel := "TEMPLATE"}}
    {{if .CodeLabel}} {{$codeLabel = .CodeLabel}} {{end}}

    {{$hideResult := false}}
    {{if .HideResult}}
        {{$hideResult = true}}
        {{$columns = 1}}
    {{end}}

    {{$result := .Code}}
    {{if .Result}} {{$result = .Result}} {{end}}

    <table columns="{{$columns}}" margin="{{$margin}}">
        <table-cell indent="0" border-color="medium-gray" border-width="0.5" background-color="light-gray">
            <division>
                <table columns="2" column-widths="0.2 0.8">
                    <table-cell indent="0">
                        <division padding="5 0 4 0">
                            <background fill-color="gray" border-color="medium-gray" border-size="0.5"></background>
                            {{template "paragraph" dict "TextAlign" "center" "VerticalTextAlign" "center" "Font" "helvetica-bold" "FontSize" 6 "TextColor" "white" "Text" $codeLabel}}
                        </division>
                    </table-cell>
                </table>
                <division margin="5 5 7 5">
                    {{template "paragraph" dict "Font" "deja-vu-sans-mono" "VerticalTextAlign" "center" "FontSize" 6 "TextColor" "text" "Text" (xmlEscape .Code)}}
                </division>
            </division>
        </table-cell>
        {{if not $hideResult}}
            <table-cell indent="0" border-color="medium-gray" border-width="0.5">
                <division>
                    <table columns="2" column-widths="0.2 0.8">
                        <table-cell indent="0">
                            <division padding="5 0 4 0">
                                <background fill-color="primary" border-color="medium-gray" border-size="1"></background>
                                {{template "paragraph" dict "TextAlign" "center" "VerticalTextAlign" "center" "Font" "helvetica-bold" "FontSize" 6 "TextColor" "white" "Text" "RESULT"}}
                            </division>
                        </table-cell>
                    </table>
                    <division margin="5 5 7 5">
                        {{print $result}}
                    </division>
                </division>
            </table-cell>
        {{end}}
    </table>
{{end}}

{{define "attr-showcase"}}
    {{$hideResult := false}}
    {{if .HideResult}}
        {{$hideResult = true}}
    {{end}}

    <division enable-page-wrap="false" margin="0 0 10 0">
        {{template "paragraph" dict "Font" "helvetica-bold" "TextColor" "secondary" "Margin" "0 0 3 0" "Text" (printf "» %s" .AttrName)}}
        <line position="relative" fit-mode="fill-width" color="secondary" thickness="0.5"></line>
        <table columns="2" column-widths="0.3 0.7" margin="2 0">
            <table-cell indent="0" vertical-align="middle" background-color="light-gray">
                {{if .AttrDefault}}
                    <paragraph margin="0 0 3 9">
                        <text-chunk font="deja-vu-sans-mono" font-size="8" color="text">Default: </text-chunk>
                        <text-chunk font="deja-vu-sans-mono" font-size="8" color="primary">{{.AttrDefault}}</text-chunk>
                        <text-chunk font="deja-vu-sans-mono" font-size="8" color="text">.</text-chunk>
                    </paragraph>
                {{end}}
            </table-cell>
            <table-cell indent="0" align="right" vertical-align="middle" background-color="light-gray">
                {{if .AttrValues}}
                    <paragraph margin="0 9 3 0">
                        <text-chunk font="deja-vu-sans-mono" font-size="8" color="text">Valid values: </text-chunk>
                        {{$lenAttrValues := len .AttrValues}}
                        {{range $i, $attrValue := .AttrValues}}
                            <text-chunk font="deja-vu-sans-mono" font-size="8" color="primary">{{$attrValue}}</text-chunk>
                            {{if not (eq $i (sum $lenAttrValues -1))}}
                                <text-chunk font="deja-vu-sans-mono" font-size="8" color="text">, </text-chunk>
                            {{end}}
                        {{end}}
                        <text-chunk font="deja-vu-sans-mono" font-size="8" color="text">.</text-chunk>
                    </paragraph>
                {{else}}
                    {{template "paragraph" dict "Font" "deja-vu-sans-mono" "FontSize" 8 "Text" " "}}
                {{end}}
            </table-cell>
        </table>

        {{if .AttrDescr}}
            {{template "paragraph" dict "Font" "deja-vu-sans-mono" "FontSize" 8 "TextColor" "text" "Margin" "0 0 2 0" "LineHeight" 1.2 "Text" .AttrDescr}}
        {{end}}
        {{if and .AttrLink .AttrLinkText}}
            <paragraph margin="0 0 2 0">
                {{if .AttrLinkDescr}}
                    <text-chunk font="deja-vu-sans-mono" font-size="8" color="text">{{.AttrLinkDescr}}</text-chunk>
                {{end}}
                <text-chunk font="deja-vu-sans-mono" font-size="8" link="{{.AttrLink}}">{{.AttrLinkText}}</text-chunk>
                <text-chunk font="deja-vu-sans-mono" font-size="8" color="text">.</text-chunk>
            </paragraph>
        {{end}}

        {{$lenCodeBlocks := len .CodeBlocks}}
        {{if gt $lenCodeBlocks 0}}
            <division margin="5 0">
                {{range $i, $codeBlock := .CodeBlocks}}
                    {{$margin := "5 0 0 0"}}
                    {{if eq $i 0}}
                        {{$margin = "0"}}
                    {{end}}

                    {{template "code-block" dict "Margin" $margin "HideResult" $hideResult "Code" $codeBlock}}
                {{end}}
            </division>
        {{end}}
    </division>
{{end}}
