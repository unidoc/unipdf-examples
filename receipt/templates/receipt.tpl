<division margin="0 0 0 0" padding="5, 5, 5, 5">
    <image src="path('./templates/res/unidoc-logo.png')" width="77.611" height="25" margin="0 0 0 0" align="center"></image>
    <paragraph text-align="center" margin="0 0 15 15">
        <text-chunk font="times-italic" font-size="28"> Receipt </text-chunk>    
    </paragraph>
    <line fit-mode="fill-width" position="relative" thickness= "1.2" margin="0 0 5 5"></line>
</division>
<division margin="5 0 0 0" padding="5, 5, 5, 5">
    <paragraph text-align="left">
        <text-chunk font="courier-oblique" font-size="14">We received payment for your subscription.</text-chunk>
        <text-chunk font-size="14" font="courier-oblique">Thanks for staying with us! Questions? </text-chunk>
        <text-chunk font-size="14" font="courier-oblique">Please contact </text-chunk>
        <text-chunk color="#0000ff" font-size="14" font="courier-oblique">support@yourhomeprovider.com</text-chunk>
         <text-chunk font-size="14" font="courier-oblique">.</text-chunk>
    </paragraph>
</division>

<table columns="2" margin="10 0 0 0">
    {{range  .Fields}}
        <table-cell align="left" border-style="none" border-width="0">
            <paragraph>
                <text-chunk font="courier-oblique" font-size="11">{{.FieldName}}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell align="left" border-style="none" border-width="0">
            <paragraph>
                <text-chunk font="courier-oblique" font-size="11">{{.FieldValue}}</text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
</table>
