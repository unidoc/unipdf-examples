{{if and (gt .PageNum 1) (lt .PageNum 74)}}
    <table columns="2" margin="20 10 35 10">
        <table-cell>
            <paragraph>
                {{if isEven .PageNum}}
                    <text-chunk font="exo-regular" font-size="22">RECEIVED</text-chunk>
                {{else}}
                    <text-chunk font="exo-regular" font-size="22">SENT</text-chunk>
                {{end}}
            </paragraph>
        </table-cell>
        <table-cell vertical-align="bottom" align="right">
            <paragraph margin="20 0 0 0">
                <text-chunk font="exo-regular" font-size="14">Operations Log Book Printout</text-chunk>
            </paragraph>
        </table-cell>
    </table>
{{end}}