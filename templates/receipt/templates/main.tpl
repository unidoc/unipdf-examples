{{define "table-cell-paragraph"}}
    <table-cell align="left" border-style="none" border-width="0">
    <paragraph>
        <text-chunk font="times" font-size="9">{{.}}</text-chunk>
    </paragraph>
    </table-cell>
{{end}}
<division margin="0 10 10 10" padding="5">
    <image src="path('./templates/res/unidoc-logo.png')" width="55.87" height="18" align="center"></image>
    <paragraph text-align="center" margin="10 0 0 0">
        <text-chunk font="times" font-size="18"> {{.Title}} </text-chunk>    
    </paragraph>
    <line fit-mode="fill-width" position="relative" thickness= "1.2" margin="10 0 0 0"></line>
</division>
<division margin="0 0 10 10" padding="5">
    <paragraph text-align="left">
        <text-chunk font="times" font-size="11">Membership fees are billed at the beginning of each period</text-chunk>
        <text-chunk font="times" font-size="11"> and may take a few days after the billing date to appear on your account. Sales tax may apply.</text-chunk>
        <text-chunk font="times" font-size="11"> Thanks for staying with us!</text-chunk>
        <text-chunk font="times" font-size="11"> If you have any questions, please contact</text-chunk>
        <text-chunk font="times" color="#0000ff" font-size="12"> support@unidocprovider.com</text-chunk>
        <text-chunk font="times" font-size="11" >.</text-chunk>
    </paragraph>
</division>
<table columns="2" margin="10" column-widths="0.4 0.6">
    {{range  .Fields}}
    {{template "table-cell-paragraph" .FieldName}}
    {{template "table-cell-paragraph" .FieldValue}}
    {{end}}
</table>
