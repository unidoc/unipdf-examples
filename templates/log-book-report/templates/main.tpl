{{define "received-header"}}
    {{range $item := (getSlice "#,Manufacturer,Model,VIN,DateReceived,Source")}}
        <table-cell background-color="#cfcfcb" border-width-bottom="0.5" border-width-top="0.5">
            <paragraph>
                <text-chunk font="exo-bold" font-size="11">{{$item}}</text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
      {{range $i, $item := (getSlice "#,Manufacturer,Model,VIN,DateReceived,Source")}}
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
    {{range $item := (getSlice "#,Date Sent,Buyer Name,Buyer Address")}}
        <table-cell background-color="#cfcfcb" border-width-bottom="0.5" border-width-top="0.5">
            <paragraph>
                <text-chunk font="exo-bold" font-size="11">{{$item}}</text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
    {{range $i, $item := (getSlice "#,Date Sent,Buyer Name,Buyer Address")}}
        <table-cell border-width-bottom="0.5" border-width-top="0.5">
            <paragraph>
                {{if eq $i 1}}
                <text-chunk font="exo-regular" font-size="11">Reciept</text-chunk>
                {{else}}
                <text-chunk font="exo-regular" font-size="11"></text-chunk>
                {{end}}
            </paragraph>
        </table-cell>
    {{end}}
{{end}}

{{define "recieved-row"}}
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
<table columns="6" column-widths = "0.05 0.19 0.19 0.19 0.19 0.19">
    {{template "received-header"}}
    {{range $i, $item := .}}
        {{template "recieved-row" dict "Item" $item "Num" $i}}
    {{end}}
</table>
<page-break></page-break>
<table columns="4" column-widths = "0.1 0.3 0.3 0.3">
    {{template "sent-header"}}
</table>