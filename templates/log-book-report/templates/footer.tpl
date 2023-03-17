{{if and (gt .PageNum 1) (lt .PageNum 74)}}
    <table columns="2" margin="10 10 35 10" indent="0">
        <table-cell>
            <paragraph>
                <text-chunk font="exo-regular" font-size="11">Page {{.PageNum}} | </text-chunk>
                <text-chunk font="exo-bold" font-size="11">Page Operations Log Book </text-chunk>
            </paragraph>
        </table-cell>
        <table-cell align="right">
            <paragraph>
                <text-chunk font="exo-regular" font-size="11">Date Of Print </text-chunk>
                <text-chunk font="exo-bold" font-size="11">{{.DateOfPrint}}</text-chunk>
                <text-chunk font="exo-regular" font-size="11"> | Date Range </text-chunk>
                <text-chunk font="exo-bold" font-size="11">{{.DateRange}}</text-chunk>
            </paragraph>
        </table-cell>
    </table>
{{end}}