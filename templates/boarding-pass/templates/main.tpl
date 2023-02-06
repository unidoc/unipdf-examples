{{define "simple-paragraph"}}
    <paragraph margin="{{.Margin}}" line-height="{{.LineHeight}}">
        <text-chunk font="{{.Font}}" font-size="{{.FontSize}}" color="{{.TextColor}}">{{.Text}}</text-chunk>
    </paragraph>
{{end}}

{{define "table-cell-paragraph"}}
    <table-cell colspan="{{.Colspan}}" rowspan="{{.Rowspan}}" background-color="{{.BackgroundColor}}" align="{{.Align}}" vertical-align="{{.VerticalAlign}}" border-color="{{.BorderColor}}" border-width-top="{{.BorderTopSize}}" border-width-bottom="{{.BorderBottomSize}}" indent="{{.Indent}}">
        {{template "simple-paragraph" .}}
    </table-cell>
{{end}}

{{define "boarding-pass-form"}}
{{$props := dict "Colspan" 1 "Rowspan" 1 "Align" "left" "VerticalAlign" "top" "BackgroundColor" "#ffffff" 
    "BorderColor" "#000000" "BorderTopSize" 0 "BorderBottomSize" 0 "Indent" 0 
    "Margin" "3 0 5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "TextColor" "#000000"}}

<table columns="4" column-widths="0.1 0.4 0.4 0.1">
    <table-cell indent="0" rowspan="3" background-color="#ffffff">
        <image src="path('templates/res/logo.png')" fit-mode="fill-width" margin="10 10 0 0"></image>
    </table-cell>

    {{template "table-cell-paragraph" (extendDict $props "Rowspan" 1 "Margin" "10 0 0 0" "Font" "helvetica-bold" "FontSize" 10 "Text" "UniDoc Air") }}
    {{template "table-cell-paragraph" (extendDict $props "Rowspan" 3 "Margin" "0 20" "Align" "right" "Font" "helvetica" "FontSize" 10 "LineHeight" 1.5 "Text" (printf "ETK:\n%s\nReg No:%s" .Etk .RegNumber )) }}

    <table-cell indent="0" rowspan="3" background-color="#ffffff">
        <image src="qr-code-1" fit-mode="fill-width"></image>
    </table-cell>

    {{template "table-cell-paragraph" (extendDict $props "Rowspan" 1 "Margin" "0" "Align" "left" "Font" "helvetica-bold" "FontSize" 14 "LineHeight" 1.0 "Text" "Boarding pass") }}
    {{template "table-cell-paragraph" (extendDict $props "Rowspan" 1 "Align" "left" "Font" "helvetica" "FontSize" 10 "TextColor" "#777777" "Text" "Passenger's coupon") }}
</table>

<table columns="5" column-widths="0.15 0.2 0.2 0.2 0.25">
    {{template "table-cell-paragraph" (extendDict $props "Margin" "5 0 0 0" "Colspan" 5 "TextColor" "#000000"  "Text" "Passenger name") }}
    {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 1 "Margin" "0" "Font" "helvetica-bold" "FontSize" 12 "Text" .Name) }}

    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "BorderBottomSize" 0 "Margin" "5 0 0 0" "Font" "helvetica" "FontSize" 10 "Text" "From") }}
    {{template "table-cell-paragraph" (extendDict $props "Text" "To") }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "center" "Text" "Hand baggage allowed") }}

    {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 1 "Margin" "0" "Align" "left" "Font" "helvetica-bold" "FontSize" 12 "Text" (printf "%s / %s" .From.City .From.Code)) }}
    {{template "table-cell-paragraph" (extendDict $props "Text" (printf "%s / %s" .Destination.City .Destination.Code)) }}

    <table-cell rowspan="5" background-color="#ffffff">
        <image src="path('templates/res/baggage.png')" fit-mode="fill-width" margin="20"></image>
    </table-cell>

    {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "BorderBottomSize" 0 "Margin" "5 0 0 0" "Font" "helvetica" "FontSize" 10 "Text" "Flight") }}
    {{template "table-cell-paragraph" (extendDict $props "Text" "Gate") }}
    {{template "table-cell-paragraph" (extendDict $props "Text" "Class") }}
    {{template "table-cell-paragraph" (extendDict $props "Text" "Seat") }}

    {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 1 "Margin" "0" "Font" "helvetica-bold" "FontSize" 12 "TextColor" "#ff0000" "Text" .FlightNumber ) }}
    {{template "table-cell-paragraph" (extendDict $props "TextColor" "#000000" "Text" .Gate ) }}
    {{template "table-cell-paragraph" (extendDict $props "TextColor" "#000000" "Text" .Class ) }}

    <table-cell border-width-bottom="1" indent="0" background-color="#ffffff">
        <table columns="2" column-widths="0.1 0.9">
            <table-cell indent="0" vertical-align="middle">
                <image src="path('templates/res/seat_icon.png')" fit-mode="fill-width"></image>
            </table-cell>

            {{template "table-cell-paragraph" (extendDict $props "Margin" "0" "BorderBottomSize" 0 "TextColor" "#ff0000" "Text" .Seat ) }}
        </table>
    </table-cell>

    {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 0 "Margin" "5 0 0 0" "Font" "helvetica" "FontSize" 10 "TextColor" "#000000" "Text" "Date") }}
    {{template "table-cell-paragraph" (extendDict $props "Text" "Boarding time till") }}
    {{template "table-cell-paragraph" (extendDict $props "Text" "Departure time") }}
    {{template "table-cell-paragraph" (extendDict $props "Text" "Arrive") }}

    {{template "table-cell-paragraph" (extendDict $props "Margin" "0" "Font" "helvetica-bold" "FontSize" 12 "Text" .Date ) }}
    {{template "table-cell-paragraph" (extendDict $props "Text" .BoardingTime ) }}
    {{template "table-cell-paragraph" (extendDict $props "Text" .DepartureTime ) }}
    {{template "table-cell-paragraph" (extendDict $props "Text" .ArrivalTime ) }}
</table>
{{end}}

{{$props := dict "Colspan" 1 "Rowspan" 1 "Align" "left" "VerticalAlign" "top" "BackgroundColor" "#ffffff" 
    "BorderColor" "#000000" "BorderTopSize" 0 "BorderBottomSize" 0 "Indent" 0 
    "Margin" "3 0 5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "TextColor" "#000000"}}

{{template "boarding-pass-form" . }}

<table columns="2">
    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "BorderBottomSize" 2 "Margin" "10 0 0 0" "Font" "helvetica-bold" "FontSize" 10 "TextColor" "#000000" "Text" "What's next?" )}}

    <table-cell indent="0">
        <division>
            {{template "simple-paragraph" (extendDict $props "Font" "helvetica" "FontSize" 8 "Margin" "10 5 0 0" "Text" "Please, arrive at the airport in advance, taking into account the time required for baggage check-in, preflight screening, passport and customs control." )}}
            {{template "simple-paragraph" (extendDict $props "Text" "You may check-in your baggage at the check-in counter at the airport. Please note the baggage transportation rules. If the weight, width or height of your baggage exceeds the free baggage allowance, you have to pay for excess baggage." )}}
            {{template "simple-paragraph" (extendDict $props "Text" "If you carry cabin baggage only, please apply to the check-in counter to get a cabin baggage tag." )}}
            {{template "simple-paragraph" (extendDict $props "Text" "If after online check-in you decide to change or return your ticket, you have to cancel your check-in via the web-site at least 1 hours prior to the departure and apply to the place of purchase of your ticket." )}}
        </division>
    </table-cell>

    <table-cell indent="0">
        <division>
            {{template "simple-paragraph" (extendDict $props "Margin" "10 0 0 5" "Text" "To ensure flight safety the airline reserves right to change your seat onboard if required so by the pilot in command." )}}
            {{template "simple-paragraph" (extendDict $props "Text" "If you don't have an opportunity to print it out, you may apply to the check-in counters at the airport. For domestic flights departing from JFK airport you may show your boarding pass on the screen of an electronic device." )}}
            {{template "simple-paragraph" (extendDict $props "Text" "The current status of any flight is stated on the online timetable at sampleairlines.com" )}}
            {{template "simple-paragraph" (extendDict $props "Text" "The boarding closes on time stated on your boarding pass. Late passengers are not accepted for transportation." )}}
        </division>
    </table-cell>

    <table-cell colspan="2" align="center" border-line-style="dashed" border-width-bottom="1" border-color="#000000">
        {{template "simple-paragraph" (extendDict $props "Margin" "20 0" "Align" "center" "Font" "helvetica-bold" "FontSize" 12 "Text" "Have a good flight!") }}
    </table-cell>
</table>

{{template "boarding-pass-form" . }}