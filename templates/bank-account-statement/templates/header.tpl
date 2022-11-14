<division margin="20 50">
    <table position="relative" columns="2" column-widths="0.9 0.1">
        <table-cell indent="0" vertical-align="middle">
            <division>
                <paragraph>
                    <text-chunk font="helvetica-bold" font-size="16">{{.Statement.BankName}} Simple Business Checking</text-chunk>
                </paragraph>
                <table columns="5" column-widths="0.3 0.05 0.3 0.05 0.3" margin="8 0 0 0">
                    <table-cell indent="0">
                        <paragraph>
                            <text-chunk>Account Number: </text-chunk>
                            <text-chunk font="helvetica-bold">{{.Statement.AccountNumber}}</text-chunk>
                        </paragraph>
                    </table-cell>
                    <table-cell indent="0">
                        <paragraph text-align="center">
                            <text-chunk>•</text-chunk>
                        </paragraph>
                    </table-cell>
                    <table-cell indent="0">
                        <paragraph text-align="center">
                            <text-chunk>{{formatTime .Statement.DateBegin "Jan 02 2006"}} - {{formatTime .Statement.DateEnd "Jan 02 2006"}}</text-chunk>
                        </paragraph>
                    </table-cell>
                    <table-cell indent="0">
                        <paragraph text-align="center">
                            <text-chunk>•</text-chunk>
                        </paragraph>
                    </table-cell>
                    <table-cell indent="0">
                        <paragraph>
                            <text-chunk>Page {{.PageNum}} of {{.TotalPages}}</text-chunk>
                        </paragraph>
                    </table-cell>
                </table>
            </division>
        </table-cell>
        <table-cell indent="0" vertical-align="middle">
            <image src="path('templates/res/logo.png')" fit-mode="fill-width" margin="0 0 0 5"></image>
        </table-cell>
    </table>
    <line position="relative" fit-mode="fill-width" margin="5 0 0 0" color="#00aff5"></line>
</division>
