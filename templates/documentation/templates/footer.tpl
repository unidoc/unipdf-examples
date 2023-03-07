<division margin="5 25 0 25">
    <division>
        <background fill-color="secondary-bg-gradient"></background>
        <paragraph><text-chunk color="white" font-size="2"> </text-chunk></paragraph>
    </division>

    <table columns="3" margin="2 0 0 0">
        <table-cell indent="0">
            <paragraph>
                <text-chunk font="helvetica-bold" color="text">© {{.Date.Year}} UniDoc ehf.</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell indent="0">
            <paragraph text-align="center" margin="-5 0 0 0">
                <text-chunk link="page(1)" underline="true" underline-offset="2" color="primary">Table of contents</text-chunk>
                <text-chunk link="page(1)" font="symbol" underline="true" underline-offset="2" color="primary"> ⇑</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell indent="0" align="right">
            <paragraph>
                <text-chunk color="text">Page </text-chunk>
                <text-chunk font="helvetica-bold" color="text">{{.PageNum}} </text-chunk>
                <text-chunk color="text">of </text-chunk>
                <text-chunk font="helvetica-bold" color="text">{{.TotalPages}}</text-chunk>
            </paragraph>
        </table-cell>
    </table>
</division>
