{{define "header"}}
    <division margin="0 0 0 0" padding="5, 5, 5, 5">
        <paragraph text-align="left" margin="0 0 0 0">
            <text-chunk font="times-bold" font-size="22">Epic Rock Concert</text-chunk>
        </paragraph>
        <paragraph text-align="left" margin="5 5 0 0">
            <text-chunk font="times" font-size="14">25.05.2021  7:30PM</text-chunk> 
        </paragraph>
        <line fit-mode="fill-width" position="relative" thickness= "1.2" margin="5 0 0 0"></line>
    </division>
{{end}}
{{define "ticket-detail"}}
    <table-cell vertical-align="middle" indent="0">
        <division margin="7 0 0 0">
            <paragraph text-align="left" margin="0 0 0 0">
                <text-chunk font="times" font-size="14">{{.}}</text-chunk>
            </paragraph>
        </division>
    </table-cell>
{{end}}
<table columns="2" margin = "0 0 10 10" padding="0, 0, 30, 30" column-widths="0.25 0.75">
<table-cell >
<image src="path('./res/0.png')" height="180" width="100" margin="5 0 0 0"></image>
</table-cell>
<table-cell>
<division>
    {{template "header"}}
<table columns="2" padding="0, 0, 0, 0" column-widths="0.3 0.7">
    <table-cell>
        <division margin="7 0 0 0" padding="5, 5, 0, 0">
            <paragraph text-align="left" margin="0 0 0 0">
                <text-chunk font="times" font-size="14">E - ticket</text-chunk>
            </paragraph>
            <paragraph text-align="left" margin="0 0 0 0">
                <text-chunk font="times" font-size="14">000385724</text-chunk>
            </paragraph>
            <image src="path('./res/1.png')" height="90" width="90" margin="5 0 0 0"></image>
        </division>
    </table-cell>
    <table-cell>
        <table columns="2">
            {{range .Detail}}
            {{template "ticket-detail" .FieldName}}
            {{template "ticket-detail" .FieldValue}}
            {{end}}
        </table>
    </table-cell>
</table>
</division>

</table-cell>
</table>
Rules of purchase
<table columns="2" margin="0 0 0 15">
 <table-cell border-width-top="0.5" border-width-bottom="0.5" border-width-left="0.5" vertical-align="top" indent="0.5" border-style="single">
    <division margin="0 0 15 15" padding="0, 0, 0, 0">
        <paragraph text-align="left" margin="0 0 10 0">
            <text-chunk font="times-bold" font-size="12">Rules of attendance</text-chunk>
        </paragraph>
        {{range .RulesOfAttendance}}
            <paragraph text-align="left" margin="0 0 0 0">
                <text-chunk font="times" font-size="11">{{.}}</text-chunk>
            </paragraph>
        {{end}}
    </division>
 </table-cell>
 <table-cell border-width-top="0.5" border-width-bottom="0.5" border-width-right="0.5" vertical-align="top" indent="0" border-style="single">
    <division margin="0 0 15 15" padding="0, 0, 0, 0">
        <paragraph text-align="left" margin="0 0 10 0">
            <text-chunk font="times-bold" font-size="12">Rules of purchase</text-chunk>
        </paragraph>
        {{range .RulesOfPurchase}}
            <paragraph text-align="left" margin="0 0 0 0">
                <text-chunk font="times" font-size="11">{{.}}</text-chunk>
            </paragraph>
        {{end}}
    </division>
 </table-cell>
</table>
        

