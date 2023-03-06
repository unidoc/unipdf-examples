{{define "table-header"}}
    <table-cell  border-width-bottom="1.5" vertical-align="bottom">
        <paragraph margin="0 0 10 0">
            <text-chunk font="arial-bold">Drug &amp; Usage</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell  border-width-bottom="1.5" vertical-align="bottom" indent="0">
        <paragraph margin="0 0 10 0">
            <text-chunk font="arial-bold">Time</text-chunk>
        </paragraph>
    </table-cell>
    {{range .Days}}
        <table-cell border-width-bottom="1.5" vertical-align="bottom">
            <paragraph margin="0 0 10 0">
                <text-chunk font="arial-bold">{{.}}</text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
{{end}}

{{define "drug-schedule"}}   
    <table-cell border-width-left="0.5" border-width-bottom="0.5" rowspan="{{len .TimesOfTheDay}}">
        <division>
            <paragraph margin="10 0 0 0">
                <text-chunk font="arial-bold" font-size="11">{{.Name}}</text-chunk>
            </paragraph>
            <paragraph margin="5 5 5 0">
                <text-chunk font="arial" font-size="9">{{.Description}}</text-chunk>
            </paragraph>
        </division>
    </table-cell>
    {{$daysTaken := .DaysTaken}}
    {{range .TimesOfTheDay}}
        <table-cell border-width-right="0.5" border-width-bottom="0.5" vertical-align="bottom" indent="0">
            <paragraph margin="0 0 10 0">
                <text-chunk font="arial" font-size="9">{{.}}</text-chunk>
            </paragraph>
        </table-cell>
        {{range $daysTaken}}
            {{$bg:="#fffffe"}}
            {{if eq . "F"}}
                {{$bg ="#cfcfcb"}}
            {{end}}
            <table-cell border-width="0.5" vertical-align="bottom" background-color="{{$bg}}">
                <paragraph margin="0 0 10 0">
                    <text-chunk font="arial" font-size="9"></text-chunk>
                </paragraph>
            </table-cell>
        {{end}}
    {{end}}
{{end}}

{{define "sign-row"}}
    <table-cell border-width-bottom="{{.BorderWidth}}" indent="0">
        <paragraph>
            <text-chunk font-size="8" font="arial" indent="0">{{.Col1}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="{{.BorderWidth}}" indent="0">
        <paragraph>
            <text-chunk font-size="8" font="arial" indent="0">{{.Col2}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width-bottom="{{.BorderWidth}}" indent="0">
        <paragraph>
            <text-chunk font-size="8" font="arial">{{.Col3}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell>     
        <paragraph>
            <text-chunk font-size="8" font="arial"></text-chunk>
        </paragraph>
    </table-cell>
{{end}}

{{define "info-row"}}
    {{if .Colspan}}
        <table-cell colspan="2" margin="{{.Margin}}" indent="0">
            <paragraph>
                {{if .MultiLine}}
                <text-chunk font-size="8" font="{{.Font1}}">{{.Heading}}</text-chunk>
                {{end}}
                <text-chunk font-size="8" font="{{.Font}}">{{.Col1}}</text-chunk>
            </paragraph>
        </table-cell>
    {{ else }}
        <table-cell indent="0">
            <paragraph margin="{{.Margin}}">
                <text-chunk font-size="8" font="{{.Font}}">{{.Col1}}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell indent="0">
            <paragraph margin="{{.Margin}}">
                <text-chunk font-size="8" font="{{.Font}}">{{.Col2}}</text-chunk>
            </paragraph>
        </table-cell>
    {{end}}
{{end}}

{{define "patient-row"}}
    <table-cell vertical-align="{{.Valign}}" border-width-bottom="{{.BorderWidth}}" indent="0">
        <paragraph>
            <text-chunk font="arial" font-size="{{.FontSize}}">{{.Col1}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell vertical-align="{{.Valign}}" border-width-bottom="{{.BorderWidth}}" indent="0">
        <paragraph>
            <text-chunk font="arial" font-size="{{.FontSize}}">{{.Col2}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell vertical-align="{{.Valign}}" border-width-bottom="{{.BorderWidth}}" indent="0">
        <paragraph>
            <text-chunk font="arial" font-size="{{.FontSize}}">{{.Col3}}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

<table columns="5" indent="0" column-widths="0.15 0.55 0.03 0.2 0.07">
    <table-cell>
        <image src="path('templates/res/logo.png')" width="105.5" height="61.5" ></image>
    </table-cell>
    <table-cell>
        <table columns="3" column-widths="0.5 0.3 0.2">
            {{template "patient-row" dict "Col1" "Medication Reconcilliation/BPMH* for:" "Col2" "Social Security Number" "Col3" "DOB" "Valign" "middle" "BorderWidth" "0" "FontSize" "8"}}
            {{template "patient-row" dict "Col1" .Patient.Name "Col2" .Patient.SocialSecurityNumber "Col3" .Patient.Dob "Valign" "bottom" "BorderWidth" "0.5" "FontSize" "11"}}
            {{template "patient-row" dict "Col1" "Two Week Period From:" "Col2" "To:" "Col3" "" "Valign" "middle" "BorderWidth" "0" "FontSize" "8"}}
            {{template "patient-row" dict "Col1" .StartDate "Col2" .EndDate "Col3" "" "Valign" "bottom" "BorderWidth" "0.5" "FontSize" "11"}}
        </table>
    </table-cell>
    <table-cell>
        <paragraph>
            <text-chunk font-size="8" font="arial"></text-chunk>
        </paragraph>
    </table-cell>
    <table-cell rowspan="4" border-width-bottom="2.5" border-width-left="1">
        <table columns="2">
            {{template "info-row" dict "Col1" "K00" "Font" "arial-bold" "Margin" "5 0 0 0" "Colspan" 2}}
            {{template "info-row" dict "Heading" "PLEASE NOTE: " "Col1" "completed calendars MUST be returned to SHC as part of the patient's Medical Record" "Font" "arial" "Font1" "arial-bold" "Margin" "5 0 0 0" "Colspan" 2 "MultiLine" true}}
            {{template "info-row" dict "Col1" "Information:" "Col2" .InformationLine "Font" "arial" "Margin" "5 0 0 0"}}
            {{template "info-row" dict "Col1" "Emergency:" "Col2" .EmergencyLine "Font" "arial-bold" "Margin" "0 0 0 0"}}
            {{template "info-row" dict "Col1" "Website:" "Col2" .Website "Font" "arial" "Margin" "5 0 0 0"}}
        </table>
    </table-cell>
    <table-cell rowspan="4" border-width-bottom="2.5" vertical-align="middle">
        <division>
            <image src="path('templates/res/bar-code.png')" width="25" height="120"></image>
        </division>
    </table-cell>
    <table-cell colspan="3" rowspan="3" vertical-align="bottom">
        <table columns="4" column-widths="0.20 0.35 0.42 0.03">
            {{template "sign-row" dict "Col1" "Date (mm/dd/yyyy)" "Col2" "Prepared by (Signature/Printed Name)" "Col3" "Verified by PhC (Signature/Printed Name)" "BorderWidth" "0"}}
            {{template "sign-row" dict "Col1" "" "Col2" "" "Col3" "" "BorderWidth" "0.5"}}
            {{template "sign-row" dict "Col1" "Date (mm/dd/yyyy)" "Col2" "Verified by RN (Signature/Printed Name)**" "Col3" "Counselled by (Signature/Printed Name)" "BorderWidth" "0"}}
            {{template "sign-row" dict "Col1" "" "Col2" "" "Col3" "" "BorderWidth" "0.5"}}
            {{template "sign-row" dict "Col1" "Date (mm/dd/yyyy)" "Col2" "Parent/Legal Guardian (Signature/Printed Name)" "Col3" "" "BorderWidth" "0"}}
            {{template "sign-row" dict "Col1" "" "Col2" "" "Col3" "" "BorderWidth" "0.5"}}
        </table>
    </table-cell>
</table>

<division margin="5 0 0 0">
    <paragraph>
        <text-chunk font-size="8">*  Best Possible Medication History</text-chunk>
    </paragraph>
    <paragraph>
        <text-chunk font-size="8">** Verification of steroids medication that are part of the patients therapy treatment</text-chunk>
    </paragraph>
</division>

<table columns="16" margin="10 0 0 0" column-widths="0.25 0.08 {{getColumnWidths (len .ListOfDays) 0.67}}" enable-page-wrap="true" enable-row-wrap="true">
    {{template "table-header" dict "Days" .ListOfDays}}
        {{range $idx, $v := .Drugs}}
            {{template "drug-schedule" .}}
        {{end}}
</table>

<table columns="5" margin="0 0 0 200" column-widths="0.6 0.05 0.15 0.05 0.15" indent="0">
    <table-cell>
        <paragraph margin="10 0 0 0">
            <text-chunk>Mark each box with a checkmark after you have taken a dose of medicine. If you skipped a dose, please consult your physician or pharmacist. Do not take medicine on the days and times not clearly indicated on this schedule.</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell align="right">
        <division margin="15 0 0 0">
            <image src="path('templates/res/checkmark.png')"></image>
        </division>
    </table-cell>
    <table-cell>
        <division margin="10 0 0 0">
            <paragraph>
                <text-chunk>Take a medication</text-chunk>
            </paragraph>
        </division>
    </table-cell>
    <table-cell>
        <division margin="15 0 0 0">
            <image src="path('templates/res/checkmark-empty.png')"></image>
        </division>
    </table-cell>
    <table-cell align="right">
        <division margin="15 0 0 0">
            <paragraph>
                <text-chunk>Skip this day</text-chunk>
            </paragraph>
        </division>
    </table-cell>
</table>