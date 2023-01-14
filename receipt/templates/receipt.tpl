{{define "simple-paragraph"}}
    <paragraph>
        <text-chunk font="times" font-size="11">{{.}}</text-chunk>
    </paragraph>
{{end}}
{{define "table-cell-paragraph"}}
    <table-cell align="left" border-style="none" border-width="0">
        {{template "simple-paragraph" .}}
    </table-cell>
{{end}}
<division margin="0 0 0 0" padding="5, 5, 5, 5">
    <image src="path('./templates/res/unidoc-logo.png')" width="77.611" height="25" margin="0 0 0 0" align="center"></image>
    <paragraph text-align="center" margin="0 0 15 15">
        <text-chunk font="times" font-size="28"> {{.Title}} </text-chunk>    
    </paragraph>
    <line fit-mode="fill-width" position="relative" thickness= "1.2" margin="0 0 0 0"></line>
</division>
<division margin="10 0 0 0" padding="0, 0, 0, 5">
    <paragraph text-align="left">
        <text-chunk font="times" font-size="14">Membership fees are billed at the beginning of each period</text-chunk>
        <text-chunk font="times" font-size="14">and may take a few days after the billing date to appear on your account. Sales tax may apply.</text-chunk>
        <text-chunk font="times" font-size="14">Thanks for staying with us!</text-chunk>
        <text-chunk font="times" font-size="14">If you have any questions, please contact</text-chunk>
        <text-chunk font="times" color="#0000ff" font-size="14"> support@unidocprovider.com</text-chunk>
        <text-chunk font="times" font-size="14" >.</text-chunk>
    </paragraph>
</division>
<table columns="2" margin="10 0 0 0">
    {{range  .Fields}}
    {{template "table-cell-paragraph" .FieldName}}
    {{template "table-cell-paragraph" .FieldValue}}
    {{end}}
</table>
