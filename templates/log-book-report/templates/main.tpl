{{define "received-header"}}
    {{$receivedColumns := "#,Manufacturer,Model,VIN,DateReceived,Source"}}
    {{range $item := (getSlice $receivedColumns)}}
        <table-cell background-color="#cfcfcb" border-width-bottom="0.5" border-width-top="0.5">
            <paragraph>
                <text-chunk font="exo-bold" font-size="11">{{$item}}</text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
    {{range $i, $item := (getSlice $receivedColumns)}}
        <table-cell border-width-bottom="0.5" border-width-top="0.5">
            <paragraph>
                {{if eq $i 1}}
                    <text-chunk font="exo-regular" font-size="11">Description of model</text-chunk>
                {{else}}
                    <text-chunk font="exo-regular" font-size="11"></text-chunk>
                {{end}}
            </paragraph>
        </table-cell>
    {{end}}
{{end}}

{{define "sent-header"}}
    {{$sentColumns := "#,Date Sent,Buyer Name,Buyer Address"}}
    {{range $columName := (getSlice $sentColumns)}}
        <table-cell background-color="#cfcfcb" border-width-bottom="0.5" border-width-top="0.5">
            <paragraph>
                <text-chunk font="exo-bold" font-size="11">{{$columName}}</text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
    {{range $i, $columName := (getSlice $sentColumns)}}
        <table-cell border-width-bottom="0.5" border-width-top="0.5">
            <paragraph>
                {{if eq $i 1}}
                    <text-chunk font="exo-regular" font-size="11">Receipt</text-chunk>
                {{else}}
                    <text-chunk font="exo-regular" font-size="11"></text-chunk>
                {{end}}
            </paragraph>
        </table-cell>
    {{end}}
{{end}}

{{define "sent-row"}}
    {{$i := add 1 .Num}}
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular">{{$i}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular">{{.Item.Sent}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular">{{.Item.BuyerName}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular">{{.Item.BuyerAddress}}</text-chunk>
        </paragraph>
    </table-cell>
    {{if eq .Item.Discarded "true"}}
        <table-cell border-width-bottom="0.5" border-width-top="0.5" colspan="4" background-color="#fff8e2">
            <paragraph>
                <text-chunk font="exo-regular"></text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
{{end}}

{{define "received-row"}}
    {{$i := add 1 .Num}}
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular" underline-offset="-3" underline-color="#ff4933" underline="{{.Item.Discarded}}">{{$i}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular" underline-offset="-3" underline-color="#ff4933" underline="{{.Item.Discarded}}">{{.Item.Manufacturer | htmlescaper}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular" underline-offset="-3" underline-color="#ff4933" underline="{{.Item.Discarded}}">{{.Item.Model | htmlescaper}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular" underline-offset="-3" underline-color="#ff4933" underline="{{.Item.Discarded}}">{{.Item.VIN | htmlescaper}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular" underline-offset="-3" underline-color="#ff4933" underline="{{.Item.Discarded}}">{{.Item.Received | htmlescaper}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="0.5" border-width-top="0.5">
        <paragraph>
            <text-chunk font="exo-regular" underline-offset="-3" underline-color="#ff4933" underline="{{.Item.Discarded}}">{{.Item.Source | htmlescaper}}</text-chunk>
        </paragraph>
    </table-cell>

    {{if eq .Item.Discarded "true"}}
        <table-cell border-width-bottom="0.5" border-width-top="0.5" background-color="#fff8e2">
            <paragraph>
                <text-chunk font="exo-regular"></text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="0.5" border-width-top="0.5" colspan="2" background-color="#fff8e2">
            <paragraph>
                <text-chunk font="exo-regular">{{.Item.DiscardReason}}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="0.5" border-width-top="0.5" background-color="#fff8e2">
            <paragraph>
                <text-chunk font="exo-regular"></text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="0.5" border-width-top="0.5" background-color="#fff8e2">
            <paragraph>
                <text-chunk font="exo-regular"></text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="0.5" border-width-top="0.5" background-color="#fff8e2">
            <paragraph>
                <text-chunk font="exo-regular"></text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
{{end}}

{{define "sent-page"}}
    <table columns="4" column-widths = "0.1 0.3 0.3 0.3">
        {{template "sent-header"}}
        {{$sNum := .StartingNum}}
        {{range $i, $item := .Items}}
            {{$num := add $i $sNum}}
            {{template "sent-row" dict "Num" $num "Item" $item }}
        {{end}}
    </table>
{{end}}

{{define "received-page"}}
    <table columns="6" column-widths = "0.05 0.19 0.19 0.19 0.19 0.19">
        {{template "received-header"}}
        {{$sNum := .StartingNum}}
        {{range $i, $item := .Items}}
            {{$num := add $i $sNum}}
            {{template "received-row" dict "Item" $item "Num" $num}}
        {{end}}
    </table>
{{end}}

{{$currentPos := 0}}
{{$pageContent := .PageToItems}}
{{range $key, $items := $pageContent}}
    {{template "received-page" dict "Items" $items "StartingNum" $currentPos}}
        <page-break></page-break>
    {{template "sent-page" dict "Items" $items "StartingNum" $currentPos}}
        <page-break></page-break>
    {{$currentPos = add $currentPos (len $items)}}
{{end}}

<division margin="80 50 0 50 ">
    <paragraph text-align="center" margin="20 0 5 0">
        <text-chunk font-size="16" font="exo-regular">Book Name</text-chunk>
    </paragraph>
    <paragraph text-align="center">
        <text-chunk font-size="16" font="exo-bold">Operations Log Book</text-chunk>
    </paragraph>
    <paragraph text-align="center" margin="20 0 10 0">
        <text-chunk font-size="16" font="exo-regular">Date of Print</text-chunk>
    </paragraph>
    <paragraph text-align="center">
        <text-chunk font-size="16" font="exo-bold">{{.DateOfPrint}}</text-chunk>
    </paragraph>
     <paragraph text-align="center" margin="20 0 10 0">
        <text-chunk font-size="16" font="exo-regular">Date range</text-chunk>
    </paragraph>
    <paragraph text-align="center">
        <text-chunk font-size="16" font="exo-bold">{{.DateRange}}</text-chunk>
    </paragraph>
    <paragraph text-align="center" margin="20 0 10 0">
        <text-chunk font-size="16" font="exo-regular">Number of Records</text-chunk>
    </paragraph>
    <paragraph text-align="center">
        <text-chunk font-size="16" font="exo-bold">{{.NumOfRecords}}</text-chunk>
    </paragraph>
</division>