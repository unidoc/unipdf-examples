<division padding="0" margin="0 25">
    <background border-color="primary" border-size="1" border-radius="3"></background>
    <division padding="25 0 20 0" margin = "0.5">
        <background fill-color="light-gray" border-radius="3"></background>
        <paragraph text-align="center" vertical-text-align="center">
            <text-chunk font="helvetica" font-size="40" color="secondary">Uni</text-chunk>
            <text-chunk font="helvetica-bold" font-size="40" color="primary">PDF</text-chunk>
        </paragraph>
    </division>
    <line position="relative" fit-mode="fill-width" color="medium-gray" margin="0.5"></line>

    <image src="path('templates/res/images/reports.png')" fit-mode="fill-width" margin="25 50"></image>

    <division padding="10 0 16 0" margin="0.5">
        <background fill-color="primary-bg-gradient"></background>
        <paragraph text-align="center">
            <text-chunk font="helvetica-bold" font-size="24" color="white">{{strToUpper "Templates Documentation"}}</text-chunk>
        </paragraph>
    </division>

    <table columns="4" margin="35 0 30 0">
        {{range $i, $unused := (loop 4)}}
            <table-cell indent="0">
                <image src="path('templates/res/images/hero{{$i}}.png')" fit-mode="fill-width" margin="0 10"></image>
            </table-cell>
        {{end}}
    </table>

    <line position="relative" fit-mode="fill-width" color="medium-gray" margin="0.5"></line>
    <division padding="10 0 12 0" margin = "0.5">
        <background fill-color="light-gray" border-radius="3"></background>
        <paragraph text-align="center" margin="0 0 8 0">
            <text-chunk font="helvetica-bold" font-size="12" color="text">POWERED BY</text-chunk>
        </paragraph>
        <image src="logo" fit-mode="fill-width" margin="0 200"></image>
    </division>
</division>
