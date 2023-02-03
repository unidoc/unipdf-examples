<table columns="3"  margin="10 20 50 20" column-widths="0.1 0.4 0.5">
    <table-cell vertical-align="bottom">
        <division>
            <image src="path('templates/res/house.png')" width="45" height="48"></image>
        </division>
    </table-cell>
    <table-cell vertical-align="bottom">
        <table columns="1" margin = "0 0 0 30">
            <table-cell vertical-align="top">
                <division margin="0 20 10 0">
                    <paragraph>
                        <text-chunk font="times" font-size="9">Initials</text-chunk>
                    </paragraph>
                        <line fit-mode="fill-width" position="relative" thickness= "0.5" margin="0 0 0 30"></line>
                </division>
            </table-cell>
        </table>
    </table-cell>
    <table-cell vertical-align="bottom">
        <table columns="3" column-widths="0.6 0.2 0.3">
            <table-cell vertical-align="bottom">
                <paragraph margin="0 0 10 0">
                    <text-chunk outline-color= "faf9e4" font="times-bold" color="#0000FF" font-size="10" underline="true" underline-thickness="0.5" underline-color="#0000FF">http://www.bestlandlords.com/billing</text-chunk>
                </paragraph>
            </table-cell>
            <table-cell vertical-align="top">
                <division>
                <image src="path('templates/res/qr.png')" width="40" height="50" margin="0 0 12 0"></image>
                </division>
            </table-cell>
            <table-cell vertical-align="bottom">
                <paragraph margin="0 0 12 0">
                    <text-chunk font="times">Page {{.PageNum}}</text-chunk>
                </paragraph>
            </table-cell>
        </table>
    </table-cell>
</table>