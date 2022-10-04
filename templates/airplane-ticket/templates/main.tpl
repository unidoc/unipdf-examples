{{define "simple-paragraph"}}
    <paragraph margin="{{.Margin}}" line-height="{{.LineHeight}}">
        <text-chunk font="{{.Font}}" font-size="{{.FontSize}}" color="{{.TextColor}}">{{.Text}}</text-chunk>
    </paragraph>
{{end}}

{{define "table-cell-paragraph"}}
    <table-cell colspan="{{.Colspan}}" background-color="{{.BackgroundColor}}" align="{{.Align}}" vertical-align="{{.VerticalAlign}}" border-color="{{.BorderColor}}" border-width-top="{{.BorderTopSize}}" border-width-bottom="{{.BorderBottomSize}}" indent="{{.Indent}}">
        {{template "simple-paragraph" .}}
    </table-cell>
{{end}}

<table columns="3" column-widths="0.12 0.68 0.2">
    <table-cell vertical-align="middle" indent="0">
        <division padding="10 5">
            <background fill-color="#f7d351" border-radius="10"></background>
            <image src="path('templates/res/logo.png')" fit-mode="fill-width" margin="0 5"></image>
        </division>
    </table-cell>
    <table-cell vertical-align="middle" indent="10">
        <division margin="7 0 0 0">
            {{template "simple-paragraph" dict "Margin" "0 5" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 26 "TextColor" "#000000" "Text" .company}}
            {{template "simple-paragraph" dict "Margin" "0 5" "LineHeight" 1 "Font" "helvetica" "FontSize" 22 "TextColor" "#000000" "Text" "E-ticket"}}
        </division>
    </table-cell>
    <table-cell vertical-align="middle" indent="0">
        <image src="qr-code-1" fit-mode="fill-width" margin="5 0 0 30"></image>
    </table-cell>
</table>

<table columns="2" margin="20 0 0 0">
    {{$props := dict "Colspan" 1 "VerticalAlign" "top" "BackgroundColor" "#ffffff" "BorderColor" "#000000" "Indent" 0 "Margin" "3 0 5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "TextColor" "#000000"}}
    <table-cell indent="0">
        <table columns="2" margin="0 5 0 0" column-widths="0.4 0.6">
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "BorderTopSize" 0 "BorderBottomSize" 0 "Text" "Passenger:")}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "BorderTopSize" 0 "BorderBottomSize" 0 "Text" .ticket.Passenger)}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "BorderTopSize" 1 "BorderBottomSize" 0 "Text" "Document:")}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "BorderTopSize" 1 "BorderBottomSize" 0 "Text" .ticket.Document)}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "BorderTopSize" 1 "BorderBottomSize" 1 "Text" "Ticket No:")}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "BorderTopSize" 1 "BorderBottomSize" 1 "Text" .ticket.Number)}}
        </table>
    </table-cell>
    <table-cell indent="0">
        <table columns="2" margin="0 0 0 5" column-widths="0.4 0.6">
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "BorderTopSize" 0 "BorderBottomSize" 0 "Text" "Order:")}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "BorderTopSize" 0 "BorderBottomSize" 0 "Text" .ticket.Order)}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "BorderTopSize" 1 "BorderBottomSize" 0 "Text" "Issued:")}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "BorderTopSize" 1 "BorderBottomSize" 0 "Text" (formatTime .ticket.Issued "02 Jan 2006"))}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "BorderTopSize" 1 "BorderBottomSize" 1 "Text" "Status:")}}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "BorderTopSize" 1 "BorderBottomSize" 1 "Text" .ticket.Status)}}
        </table>
    </table-cell>
</table>

{{template "simple-paragraph" dict "Margin" "20 0 5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 12 "TextColor" "#000000" "Text" "Route"}}
<line position="relative" fit-mode="fill-width" thickness="2"></line>
<table columns="6" margin="5 0 0 0">
    {{$props := dict "Colspan" 1 "Align" "left" "VerticalAlign" "top" "BackgroundColor" "#ffffff" "BorderColor" "#000000" "BorderTopSize" 0 "BorderBottomSize" 0 "Indent" 0 "Margin" "0 0 10 0" "LineHeight" 1.1 "Font" "helvetica-bold" "FontSize" 10 "TextColor" "#808080"}}
    {{template "table-cell-paragraph" (extendDict $props "Text" "FLIGHT")}}
    {{template "table-cell-paragraph" (extendDict $props "Text" "DEPARTURE")}}
    {{template "table-cell-paragraph" (extendDict $props "Text" "ARRIVAL")}}
    {{template "table-cell-paragraph" (extendDict $props "Text" "CLASS")}}
    {{template "table-cell-paragraph" (extendDict $props "Text" "BAGGAGE")}}
    {{template "table-cell-paragraph" (extendDict $props  "Text" "CHECK-IN")}}

    {{range $route := .ticket.Routes}}
        {{$props := dict "Colspan" 1 "Align" "left" "VerticalAlign" "top" "BackgroundColor" "#ffffff" "BorderColor" "#000000" "BorderTopSize" 0 "BorderBottomSize" 0 "Indent" 0 "Margin" "2 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "TextColor" "#000000"}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.Flight)}}
        {{template "table-cell-paragraph" (extendDict $props "Text" (formatTime $route.Departure "02 Jan"))}}
        {{template "table-cell-paragraph" (extendDict $props "Text" (formatTime $route.Arrival "02 Jan"))}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.Class)}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.Baggage)}}
        {{template "table-cell-paragraph" (extendDict $props "Text" (formatTime $route.CheckIn "15:04"))}}

        {{$props = (extendDict $props "Margin" "5 0 0 0" "FontSize" 8 "TextColor" "#000000")}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.FlightCompany)}}
        {{template "table-cell-paragraph" (extendDict $props "Text" (formatTime $route.Departure "15:04"))}}
        {{template "table-cell-paragraph" (extendDict $props "Text" (formatTime $route.Arrival "15:04"))}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.ClassAdd)}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.BaggageAdd)}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.CheckInAirport)}}

        {{$props = (extendDict $props "Margin" "0")}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.FlightPlaner)}}
        {{template "table-cell-paragraph" (extendDict $props "Text" $route.DepartureAirport)}}
        {{template "table-cell-paragraph" (extendDict $props "Colspan" 4 "Text" $route.ArrivalAirport)}}

        <table-cell indent="0" colspan="6">
            <line position="relative" fit-mode="fill-width" thickness="1" margin="7 0 0 0"></line>
        </table-cell>
    {{end}}
</table>

{{template "simple-paragraph" dict "Margin" "15 0 5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 12 "TextColor" "#000000" "Text" "About your trip"}}
<line position="relative" fit-mode="fill-width" thickness="2"></line>
<table columns="2" column-widths="0.01 0.99" margin="5 0 0 0">
    {{$props := dict "Colspan" 1 "Align" "left" "VerticalAlign" "top" "BackgroundColor" "#ffffff" "BorderColor" "#000000" "BorderTopSize" 0 "BorderBottomSize" 0 "Margin" "0" "LineHeight" 1.1 "Font" "helvetica" "FontSize" 9 "TextColor" "#000000"}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "•")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 3 "Text" (printf "Use your Trip ID for all communication with %s about this booking." .company))}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "•")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 3 "Text" "Check-in counters at all Airports close 45 minutes before departure.")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "•")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 3 "Text" (printf "For %s flights the free check-in baggage allowance is 15 kgs in Economy class for domestic travel within US." .company))}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "•")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 3 "Text" "Your carry-on baggage shouldn't weigh more than 7 kgs.")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "•")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 3 "Text" "Carry photo identification, you will need it as proof of identity while checking in.")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "•")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 3 "Text" (printf "If Cancellation is done through the customer support executives assistance, we will levy 500.00 USD per passenger per flight. However, if you do it online using your %s account, we will only levy 250 USD per passenger per flight as %s Processing charges. This is over and above the airline cancellation charges." .company .company))}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "•")}}
    {{template "table-cell-paragraph" (extendDict $props "Indent" 3 "Text" (printf "For hassle free refund processing, cancel/amend your tickets with %s Customer Care instead of doing so directly with the airline." .company))}}
</table>

<table columns="2" margin="10 0 0 0">
    {{$props := dict "Colspan" 1 "VerticalAlign" "top" "BackgroundColor" "#ffffff" "BorderColor" "#000000" "BorderTopSize" 1 "BorderBottomSize" 1 "Indent" 0 "Margin" "3 0 5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "TextColor" "#000000"}}
    {{$propsHeader := dict "Colspan" 2 "Align" "left" "VerticalAlign" "top" "BackgroundColor" "#ffffff" "BorderColor" "#000000" "BorderTopSize" 0 "BorderBottomSize" 0 "Indent" 0 "Margin" "3 0 0 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 12 "TextColor" "#000000"}}
    <table-cell indent="0">
        <table columns="2" margin="0 5 0 0" column-widths="0.6 0.4">
            {{template "table-cell-paragraph" (extendDict $propsHeader "Text" "Fare breakdown")}}
            {{range $i, $fare := .ticket.Fares}}
                {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Text" $fare.Name)}}
                {{template "table-cell-paragraph" (extendDict $props "Align" "right" "Text" (printf "%.2f USD" $fare.Charge))}}
            {{end}}
        </table>
    </table-cell>
    <table-cell indent="0">
        <table columns="2" margin="0 0 0 5" column-widths="0.6 0.4">
            {{template "table-cell-paragraph" (extendDict $propsHeader "Text" "Need help?")}}
            {{range $phoneNumber := .ticket.PhoneNumbers}}
                {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Text" $phoneNumber.Name)}}
                {{template "table-cell-paragraph" (extendDict $props "Align" "right" "Text" $phoneNumber.Value)}}
            {{end}}
        </table>
    </table-cell>
</table>
