<division margin="10">
    <line position="relative" fit-mode="fill-width" thickness="2" margin="0 0 10 0"></line>

    <table columns="3" column-widths="0.3 0.3 0.4">
        <table-cell>
            <paragraph>
                <text-chunk font="helvetica-bold" font-size="14">FINAL REPORT</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell align="right">
            <paragraph>
                <text-chunk font-size="12">Date Issued: {{ .Date }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell align="right">
            <paragraph>
                <text-chunk font-size="12">Page {{ .PageNum }} of {{ .TotalPages }}</text-chunk>
            </paragraph>
        </table-cell>

        <table-cell colspan="2">
            <paragraph>
                <text-chunk font-size="12">This document contains private and confidential health information protected by state and federal law. If you have received this document in error, please call {{ .SupportPhone }}</text-chunk>
            </paragraph>
        </table-cell>
        
        <table-cell align="right" indent="0">
            <paragraph text-align="right" margin="0 10">
                <text-chunk font-size="12">© 1995-2020 Sample Corporation of America® Holdings All Rights Reserved - Enterprise Report Version: 1.00</text-chunk>
            </paragraph>
        </table-cell>
    </table>
</division>