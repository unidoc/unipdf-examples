<division margin="20 25">
    <table position="relative" columns="3" column-widths="0.2 0.6 0.2">
        <table-cell indent="0">
            <image src="logo" fit-mode="fill-width" margin="0 40 0 0"></image>
        </table-cell>
        <table-cell indent="0" vertical-align="middle">
            <paragraph text-align="center">
                <text-chunk font="helvetica-bold" font-size="14" color="text">UniPDF templates documentation</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell indent="0" vertical-align="middle">
            <paragraph text-align="right" margin="5 0 0 0">
                <text-chunk color="text">Page </text-chunk>
                <text-chunk font="helvetica-bold" color="text">{{.PageNum}} </text-chunk>
                <text-chunk color="text">of </text-chunk>
                <text-chunk font="helvetica-bold" color="text">{{.TotalPages}}</text-chunk>
            </paragraph>
        </table-cell>
    </table>

    <division padding="0" margin="4 0 0 0">
        <background fill-color="secondary-bg-gradient"></background>
        <paragraph><text-chunk color="white" font-size="2"> </text-chunk></paragraph>
    </division>
</division>
