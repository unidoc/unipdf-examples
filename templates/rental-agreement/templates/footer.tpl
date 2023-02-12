<table columns="3"  margin="10 20 30 20" column-widths="0.085 0.415 0.5">
    <table-cell vertical-align="bottom">
        {{if gt 10 .PageNum}}
            <image src="path('templates/res/house.png')" width="45" height="48"></image>
        {{end}}
    </table-cell>
    <table-cell vertical-align="bottom">
        <table columns="1" margin = "0 0 0 30">
            <table-cell vertical-align="bottom">
                {{if gt 10 .PageNum}}
                {{$width := getWidth "Initials " "Times-Roman"}}
                <division>
                    <paragraph>
                        <text-chunk font="times" font-size="9">Initials</text-chunk>
                    </paragraph>
                        <line fit-mode="fill-width" position="relative" thickness= "0.5" margin="0 0 10 30"></line>
                </division>
                {{end}}
            </table-cell>
        </table>
    </table-cell>
    <table-cell vertical-align="bottom">
        <table columns="3" column-widths="0.58 0.22 0.3">
            <table-cell vertical-align="bottom">
                {{if gt 10 .PageNum}}
                <paragraph margin="0 0 8 0">
                    <text-chunk outline-color="faf9e4" font="arial-bold" font-size="9" underline="true" underline-thickness="0.5" link="url('http://www.bestlandlords.com/billing')">http://www.bestlandlords.com/billing</text-chunk>
                </paragraph>
                {{end}}
            </table-cell>
            <table-cell>
                {{if gt 10 .PageNum}}
                <division>
                <image src="path('templates/res/qr.png')" width="50" height="70" margin="0 0 10 0"></image>
                </division>
                {{end}}
            </table-cell>
            <table-cell vertical-align="bottom">
                <paragraph margin="0 0 10 0">
                    <text-chunk font="times">Page {{.PageNum}}</text-chunk>
                </paragraph>
            </table-cell>
        </table>
    </table-cell>
</table>