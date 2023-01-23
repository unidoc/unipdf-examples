{{define "header"}}
{{$fontSize := .FontSize}}
{{$subFont := .SubFont}}
    <division margin="0 0 0 0" padding="5, 5, 5, 5">
        <paragraph text-align="left" margin="0 0 0 0">
            <text-chunk font="helvetica-bold" font-size="{{$fontSize}}">Epic Rock Concert</text-chunk>
        </paragraph>
        <paragraph text-align="left" margin="5 5 0 0">
            <text-chunk font="helvetica" font-size="{{$subFont}}">25.05.2021  7:30PM</text-chunk> 
        </paragraph>
        <line fit-mode="fill-width" position="relative" thickness= "2.0" margin="5 0 0 0"></line>
    </division>
{{end}}
{{define "ticket-detail"}}
{{$fieled := .Name}}
{{$fontSize := .FontSize}}
{{$fontName := .FontName}}
    <table-cell indent="0">
        <division margin="10 0 0 0">
            <paragraph text-align="left" margin="0 0 0 0">
                <text-chunk font="{{$fontName}}" font-size="{{$fontSize}}">{{$fieled}}</text-chunk>
            </paragraph>
        </division>
    </table-cell>
{{end}}
<table columns="2" margin = "0 0 10 10" padding="0, 0, 30, 30" column-widths="0.20 0.80">
    <table-cell vertical-align="top">
        <image src="path('./res/red-guitar.png')" width="100" height="160" margin="5 0 0 0"></image>
    </table-cell>
    <table-cell vertical-align="top">
        <division>
            <division margin="0 0 0 15">
            {{template "header" dict "FontSize" 19 "SubFont" 11}}
            </division>
                <table columns="2" padding="0, 0, 0, 0" column-widths="0.3 0.7">
                <table-cell>
                    <table columns="1" margin="10 0 0 0">
                        <table-cell>
                            <division margin="0 0 0 0" padding="0, 0, 0, 0">
                                <paragraph text-align="left" margin="0 0 0 10" line-height="1.1">
                                    <text-chunk font="helvetica" font-size="10">E - ticket</text-chunk>
                                </paragraph>
                                <paragraph text-align="left" margin="2 0 0 10">
                                    <text-chunk font="helvetica" font-size="10">{{.TicketNumber}}</text-chunk>
                                </paragraph>
                                <image src="qr-code" height="70" width="70" margin="5 0 0 10"></image>
                            </division>
                        </table-cell>
                    </table>
                </table-cell>
                <table-cell>
                    <table columns="2">
                        {{range .Detail}}
                        {{template "ticket-detail" dict "FontName" "helvetica" "FontSize" 11 "Name" .FieldName }}
                        {{template "ticket-detail" dict "FontName" "helvetica" "FontSize" 11 "Name" .FieldValue }}
                        {{end}}
                    </table>
                </table-cell>
            </table>
        </division>
    </table-cell>
</table>
<table columns="2" margin="0 0 0 15">
 <table-cell border-width-top="0.5" border-width-bottom="0.5" border-width-left="0.5" vertical-align="top" indent="0.5" border-style="single">
    <division margin="0 0 15 15" padding="0, 0, 0, 0">
        <paragraph text-align="left" margin="10 0 10 0">
            <text-chunk font="helvetica-bold" font-size="12">Rules of attendance</text-chunk>
        </paragraph>
        {{range .RulesOfAttendance}}
            <paragraph text-align="left" margin="0 0 0 0" line-height="1.3">
                <text-chunk font="helvetica" font-size="9">{{.}}</text-chunk>
            </paragraph>
        {{end}}
    </division>
 </table-cell>
 <table-cell border-width-top="0.5" border-width-bottom="0.5" border-width-right="0.5" vertical-align="top" indent="0" border-style="single">
    <division margin="0 0 15 15" padding="0, 0, 0, 0">
        <paragraph text-align="left" margin="10 0 10 0">
            <text-chunk font="helvetica-bold" font-size="12">Rules of purchase</text-chunk>
        </paragraph>
        {{range .RulesOfPurchase}}
            <paragraph text-align="left" margin="0 0 0 0" line-height="1.3">
                <text-chunk font="helvetica" font-size="9">{{.}}</text-chunk>
            </paragraph>
        {{end}}
    </division>
 </table-cell>
</table>

<table columns="3" margin="10 0 0 10" column-widths="0.5 0.20 0.30">
    <table-cell vertical-align="top">
        <division margin="0 20 0 0" padding="0, 0, 0, 0">
            <paragraph text-align="top" margin="0 0 0 0">
                <text-chunk font="helvetica-bold" font-size="12">Program/Bands List/Event Description</text-chunk>
            </paragraph>
            <paragraph text-align="top" margin="5 0 0 0" line-height="1.3">
                <text-chunk font="helvetica" font-size="9">A City literally built on Rock n' Roll in Cadott WI, </text-chunk>
                <text-chunk font="helvetica" font-size="9"> Rock Fest is the true Rock experience you can't miss. In its 26th year, </text-chunk>
                <text-chunk font="helvetica" font-size="9"> it is the top venue for people of all ages to come together for one common purpose: </text-chunk>
                <text-chunk font="helvetica" font-size="9">to congregate with other rock fans from across the world, in a place where rock music still matters.</text-chunk>
                <text-chunk font="helvetica" font-size="9">Featuring the very best of active and classic rock and legendary names in Rock Music, entertainment and experience are the first priority.</text-chunk>
                <text-chunk font="helvetica" font-size="9">Aerosmith, Iron Maiden, Avenged Sevenfold, Kiss, Motley Crue, Fleetwood Mac, Tom Petty, Kid Rock, Shinedown,</text-chunk>
                <text-chunk font="helvetica" font-size="9">Five Finger Death Punch, Rob Zombie, Korn &amp; many more rock legends have graced this permanent Main Stage over the course of the last two and a half decades.</text-chunk>
            </paragraph>
        </division>
    </table-cell>
    <table-cell vertical-align="top">
        <image src="path('./res/map.png')" height="160" width="100" margin="5 0 0 0"></image>
    </table-cell>
    <table-cell vertical-align="top">
        <division>
            <paragraph text-align="top" margin="0 0 0 0">
                    <text-chunk font="helvetica-bold" font-size="12">How to find us</text-chunk>
            </paragraph>
            <paragraph text-align="top" margin="5 0 0 0" line-height="1.3">
                <text-chunk font="helvetica" font-size="9">Concert Hall is located three blocks west of the park and three blocks north the library.</text-chunk>
                <text-chunk font="helvetica" font-size="9">Free parking is available nearby, which fills up close to showtime, and on the streets around.</text-chunk>
                <text-chunk font="helvetica" font-size="9">Bike racks are located outside the main entrance tothe Hall.</text-chunk>
            </paragraph>
            <paragraph text-align="top" margin="10 0 0 0">
                <text-chunk font="helvetica-bold" font-size="12">Learn More At</text-chunk>
            </paragraph>
            <paragraph margin="5 0 0 0">
                <text-chunk font="helvetica" font-size="9">Facebook /ConcertHall</text-chunk>
            </paragraph>
            <paragraph margin="2 0 0 0">
                <text-chunk font="helvetica" font-size="9">Twitter @ConcertHall</text-chunk>
            </paragraph>
            <paragraph margin="2 0 0 0">
                <text-chunk font="helvetica" font-size="9">Instagram @ConcertHall</text-chunk>
            </paragraph>
        </division>
    </table-cell>
</table>

<table columns="2">
    <table-cell vertical-align="middle">
        <division margin="0 0 0 10">
            <paragraph text-align="top" margin="0 0 0 0">
                <text-chunk font="helvetica" font-size="9">This is your ticket. Print this Entire page, fold it and bring it with you to the event.</text-chunk>
                <text-chunk font="helvetica" font-size="9"> Please make sure the QR code is visible.</text-chunk>
            </paragraph>
            <image src="path('./res/ticket-img.png')" height="80" width="100" margin="5 0 0 0"></image>
        </division>
    </table-cell>
    <table-cell>
        <division margin="10 0 0 0">
            {{template "header" dict "FontSize" 19 "SubFont" 12}}
            <table columns="2">
            <table-cell>
                    <division margin="5 0 0 0">
                        <table columns="2">
                            {{range .Detail}}
                            {{if ne .FieldName "address"}}
                                {{template "ticket-detail" dict "FontName" "helvetica" "FontSize" 11 "Name" .FieldName }}
                                {{template "ticket-detail" dict "FontName" "helvetica-bold" "FontSize" 11 "Name" .FieldValue}}
                            {{end}}
                            {{end}}
                        </table>
                    </division>
            </table-cell>
            <table-cell>
                    <division>
                        <image src="qr-code" height="75" width="75" margin="15 0 0 10"></image>
                        <paragraph text-align="top" margin="0 0 0 15">
                            <text-chunk font="helvetica" font-size="10">{{.TicketNumber}}</text-chunk>
                        </paragraph>
                    </division>
            </table-cell> 
            </table>
        </division>
    </table-cell>
</table>
