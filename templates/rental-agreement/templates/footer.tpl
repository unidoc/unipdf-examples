<table columns="3"  margin="30 30 0 30" column-widths="0.1 0.7 0.2">
    <table-cell vertical-align="bottom">
        <division>
            <image src="path('templates/res/house.png')" width="50" height="50"></image>
        </division>
    </table-cell>
    <table-cell vertical-align="bottom">
        <table columns="2" margin = "0 0 0 30">
            <table-cell vertical-align="bottom" padding="0 20 0 20">
                <division margin="0 20 0 0">
                    <paragraph>
                        <text-chunk font="times" font-size="9">Initials</text-chunk>
                    </paragraph>
                        <line fit-mode="fill-width" position="relative" thickness= "0.5" margin="0 0 0 30"></line>
                </division>
            </table-cell>
        <table-cell vertical-align="bottom" margin = "0 0 0 30">
            <paragraph>
            <text-chunk outline-color= "faf9e4" font="times" color="#0000FF" font-size="10" underline="true" underline-thickness="0.5" underline-color="#0000FF"> http://www.bestlandlords.com/billing </text-chunk>
            </paragraph>
        </table-cell>
        </table>
    </table-cell>
    <table-cell vertical-align="bottom">
        <table columns="2">
            <table-cell vertical-align="bottom">
                <image src="path('templates/res/qr.png')" width="40" height="50"></image>
            </table-cell>
            <table-cell vertical-align="bottom">
            <paragraph>
                <text-chunk font="times">  Page {{.PageNum}} </text-chunk>
            </paragraph>
            </table-cell>
        </table>
    </table-cell>
</table>