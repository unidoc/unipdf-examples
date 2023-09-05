{{define "simple-paragraph"}}
    <paragraph margin="{{.Margin}}" line-height="{{.LineHeight}}">
        <text-chunk font="{{.Font}}" font-size="{{.FontSize}}">{{.Text}}</text-chunk>
    </paragraph>
{{end}}

{{define "table-cell-paragraph"}}
    <table-cell colspan="{{.Colspan}}" background-color="{{.BackgroundColor}}" align="{{.Align}}" vertical-align="{{.VerticalAlign}}" border-color="{{.BorderColor}}" border-width="{{.BorderSize}}" indent="{{.Indent}}">
        {{template "simple-paragraph" .}}
    </table-cell>
{{end}}

{{define "gradeLabel"}}
    <division padding="7 0 3 0" margin="{{.Margin}}">
        <background fill-color="{{.FillColor}}" border-radius="5"></background>
        <paragraph text-align="center" vertical-text-align="center">
            <text-chunk font="helvetica-bold" font-size="{{.FontSize}}" color="#fff">{{.Text}}</text-chunk>
        </paragraph>
    </division>
{{end}}

{{define "gradeBarSection"}}
    <table-cell background-color="{{.BackgroundColor}}" vertical-align="middle">
        {{$label := .Label}}
        {{if not .ShowLabel}}
            {{$label = ""}}
        {{end}}

        <paragraph margin="{{.Margin}}" vertical-text-align="center">
            <text-chunk font="helvetica-bold" font-size="{{.FontSize}}" color="#ffffff">{{$label}}</text-chunk>
        </paragraph>
    </table-cell>
{{end}}

{{define "gradeBar"}}
    <table columns="5" column-widths="0.3 0.2 0.15 0.15 0.2" margin="{{.Margin}}">
        {{template "gradeBarSection" dict "FontSize" .SectionFontSize "Margin" .SectionMargin "ShowLabel" .ShowLabels "Label" "Very Poor" "BackgroundColor" "linear-gradient(#F7000B, #F95A0C)"}}
        {{template "gradeBarSection" dict "FontSize" .SectionFontSize "Margin" .SectionMargin "ShowLabel" .ShowLabels "Label" "Poor" "BackgroundColor" "linear-gradient(#F95A0C, #F9C20E)"}}
        {{template "gradeBarSection" dict "FontSize" .SectionFontSize "Margin" .SectionMargin "ShowLabel" .ShowLabels "Label" "Fair" "BackgroundColor" "linear-gradient(#F9C20E, #98ED1B)"}}
        {{template "gradeBarSection" dict "FontSize" .SectionFontSize "Margin" .SectionMargin "ShowLabel" .ShowLabels "Label" "Good" "BackgroundColor" "linear-gradient(#98ED1B, #27C729)"}}
        {{template "gradeBarSection" dict "FontSize" .SectionFontSize "Margin" .SectionMargin "ShowLabel" .ShowLabels "Label" "Excellent" "BackgroundColor" "#27C729"}}
    </table>
{{end}}

{{define "sectionTitle"}}
    <paragraph margin="10 0 10 0">
        <text-chunk color="#848484" font="helvetica-bold" font-size="11">{{.Text}}</text-chunk>
    </paragraph>
{{end}}

{{define "sectionContent"}}
    {{$columns := 3}}
    {{$columnWidths := "0.075 0.5 0.425"}}
    {{if .HideBar }}
        {{$columns = 2}}
        {{$columnWidths = "0.2 0.8"}}
    {{end}}

    <division padding="10">
        <background fill-color="{{.BackgroundColor}}"></background>
        <table columns="{{$columns}}" column-widths="{{$columnWidths}}" margin="0 0 3 0">
            <table-cell vertical-align="middle">
                {{template "gradeLabel" dict "FillColor" .GradeColor "Text" .GradeText "FontSize" 20 "Margin" "-3 0 0 0"}}
            </table-cell>
            <table-cell vertical-align="middle">
                <paragraph line-height="1.1">
                    <text-chunk font="helvetica-bold">{{.Title}}&#xA;</text-chunk>
                    <text-chunk font="helvetica">{{.Subtitle}}</text-chunk>
                </paragraph>
            </table-cell>
            {{if not .HideBar}}
                <table-cell vertical-align="middle">
                    {{template "gradeBar" dict "ShowLabels" false "Margin" "0" "SectionFontSize" 10 "SectionMargin" "5 0"}}
                </table-cell>
            {{end}}
        </table>
    </division>
{{end}}

{{define "status-circle-cell"}}
<table-cell align="center" vertical-align="middle" border-color="{{.BorderColor}}" border-width="1">
    <division margin="0 21">
        <background fill-color="{{.Color}}" border-radius="6.5" border-color="#e2ecf0" border-size="1"></background>
        {{template "simple-paragraph" dict "Margin" "5 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" " "}}
    </division>
</table-cell>
{{end}}

{{define "status-legend-cell"}}
    <table-cell indent="0">
        <table columns="2" column-widths="0.4 0.6" margin="0 5 0 0">
            {{template "status-circle-cell" dict "Color" .Color "BorderColor" "#ffffff"}}
            {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0 0 0 -20" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" .Text}}
        </table>
    </table-cell>
{{end}}

{{define "smart-alert-summary-row"}}
    {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 5 "Margin" "3 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" .Type}}
    {{template "status-circle-cell" dict "Color" .StatusColor "BorderColor" "#e2ecf0"}}
    {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "3 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" .NoActors}}
    {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "3 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" .FirstSeen}}
    {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "3 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" .LastTriggered}}
{{end}}

<chapter show-numbering="false">
    <chapter-heading font="helvetica-bold" font-size="16" margin="0 0 5 0">Threat Assessment for Acme</chapter-heading>
    {{template "simple-paragraph" dict "Margin" "0 0 5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 12 "Text" "Report generated on November 13, 2019"}}
    {{template "simple-paragraph" dict "Margin" "0 0 10 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "Text" "Analysis based on 19 days of retained data from October 25, 2019 to November 13, 2019"}}
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 10 0"></line>

    <table columns="3" column-widths="0.3 0.3 0.4">
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "0 0 5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 10 "Text" "Your threat score is"}}
                <table columns="2" column-widths="0.3 0.7">
                    <table-cell indent="0">
                        {{template "gradeLabel" dict "Text" "D" "FillColor" "#ff9721" "FontSize" "24" "Margin" "10 5 10 0"}}
                    </table-cell>
                    {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "bottom" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "10 5 5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 25 "Text" "650"}}
                </table>
            </division>
        </table-cell>
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "0 0 5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 10 "Text" "Target* threat score is"}}
                <table columns="2" column-widths="0.3 0.7">
                    <table-cell indent="0">
                        {{template "gradeLabel" dict "Text" "B" "FillColor" "#0dc63c" "FontSize" 24 "Margin" "10 5 10 0"}}
                    </table-cell>
                    {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "bottom" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "10 5 5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 25 "Text" "763"}}
                </table>
                {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica" "FontSize" 7 "Text" "* Based on threat scores from UniPDF users"}}
            </division>
        </table-cell>
        <table-cell indent="0" align="right">
            <division>
                {{template "simple-paragraph" dict "Margin" "-1 0 2 0" "LineHeight" 1.25 "Font" "helvetica-bold" "FontSize" 7 "Text" "Continued use of UniPDF is aimed at improving your threat score and securing your critical IT assets."}}
                {{template "simple-paragraph" dict "Margin" "-1 0 2 0" "LineHeight" 1.25 "Font" "helvetica" "FontSize" 7 "Text" "UniPDF identifies, detects, and responds to threats to your network without requiring any additional hardware, software or people. The UniPDF Cloud continuously analyzes the billions of conversations happening on your network, learns what is normal, and alerts when suspicious behaviors that users risk the security of your critical IT assets are detected."}}
            </division>
        </table-cell>
    </table>

    {{template "gradeBar" dict "ShowLabels" true "Margin" "15 0 10 0" "SectionFontSize" 10 "SectionMargin" "5"}}

    {{template "sectionTitle" dict "Text" "THREAT DETECTION"}}
    {{template "sectionContent" dict "BackgroundColor" "#e2ecf0" "GradeColor" "#159635" "GradeText" "A" "Title" "Open Smart Alerts" "Subtitle" "3 Currently open"}}
    {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#ff9721" "GradeText" "D" "Title" "Average Time to Close Smart Alerts" "Subtitle" "3.2 Days, (Using a trailing 7-day average)"}}
    {{template "sectionContent" dict "BackgroundColor" "#e2ecf0" "GradeColor" "#ffC700" "GradeText" "C" "Title" "Manual Effort Saved" "Subtitle" "1.8 Person Days per Week"}}

    {{template "sectionTitle" dict "Text" "NETWORK VISIBILITY"}}
    {{template "sectionContent" dict "BackgroundColor" "#e2ecf0" "GradeColor" "#ff0000" "GradeText" "F" "Title" "Unidentified Assets" "Subtitle" "68.2%"}}
    {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#159635" "GradeText" "A" "Title" "High Risk Assets" "Subtitle" "0.0%"}}
    {{template "sectionContent" dict "BackgroundColor" "#e2ecf0" "GradeColor" "#ff0000" "GradeText" "F" "Title" "Unidentified Subnets or IP Ranges" "Subtitle" "100.0%"}}

    {{template "sectionTitle" dict "Text" "POLICY ASSURANCE"}}
    {{template "sectionContent" dict "BackgroundColor" "#e2ecf0" "GradeColor" "#ffC700" "GradeText" "C" "Title" "Policy Alerts" "Subtitle" "3.0 Per day (Average of past 7 days)"}}
    {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#159635" "GradeText" "A" "Title" "Time Saved" "Subtitle" "4.2 Days"}}
    {{template "sectionContent" dict "BackgroundColor" "#e2ecf0" "GradeColor" "#8bc650" "GradeText" "B" "Title" "Policy Violations" "Subtitle" "1 per day"}}
</chapter>

<chapter show-numbering="false">
    <chapter-heading font="helvetica-bold" font-size="16" margin="20 0 5 0">Threat Detection</chapter-heading>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>
    <table columns="2" column-widths="0.5 0.5">
        <table-cell indent="0" align="center">
            <division>
                {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#ffC700" "GradeText" "C" "Title" "Open Smart Alerts" "Subtitle" "3 Currently open" "HideBar" true}}
                <chart height="80" src="{{CreatePieChart "pie-chart-1" (dict "Normal" 3.0 "Very High" 4.0 "Low" 5.0 "High" 5.0)}}" margin="-10 0 0 0"></chart>
            </division>
        </table-cell>
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "10 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "Having Less than 5 open alerts at any given time is a good indicator that you are addressing detected threats in a timely manner."}}
                <chart height="80" src="{{CreateLineChart "line-chart-1" 50 0.0 50.0}}"></chart>
            </division>
        </table-cell>
    </table>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>

    <table columns="2" column-widths="0.5 0.5">
        <table-cell indent="0" align="center">
            <division>
                {{$src := CreateStackedBarChart "stacked-bar-chart-1"
                    (CreateStackedBar "" (dict "Abnormal and Unauthorized" 40.0 "Abnormal but Authorized" 60.0))}}
                {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#ff9721" "GradeText" "D" "Title" "Average Time to Close Smart Alerts" "Subtitle" "3.2 Days (Using a trailing 7-day average)" "HideBar" true}}
                <chart height="50" src="{{$src}}" margin="0 0 0 0"></chart>
            </division>
        </table-cell>
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "10 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "An average time to close of less than 2 days indicates that you are taking a proactive approach to assessing and remediating threats and vulnerabilities."}}
                <chart height="80" src="{{CreateLineChart "line-chart-2" 20 0.0 20.0}}"></chart>
            </division>
        </table-cell>
    </table>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>

    <table columns="2" column-widths="0.5 0.5">
        <table-cell indent="0" align="center">
            <division>
                {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#ffC700" "GradeText" "C" "Title" "Manual Effort Saved" "Subtitle" "1.8 Person Days per Week" "HideBar" true}}
                <table columns="2" column-widths="0.2 0.8" margin="0 20">
                    <table-cell indent="0"></table-cell>
                    <table-cell indent="0">
                        <division>
                            {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Occurrences over past 1 week"}}
                            <table columns="2" column-widths="0.8 0.2" margin="10 0 0 0">
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Policy Alerts"}}
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "21"}}
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "New Smart Alerts"}}
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "3"}}
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Orphaned Behaviors"}}
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "15"}}
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Unconfirmed Smart Alerts"}}
                                {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "0"}}
                            </table>
                        </division>
                    </table-cell>
                </table>
            </division>
        </table-cell>
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "10 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "In networks of 100 to 200 unique internal IPs, you should see a savings of 5 work days per week. Larger networks should see more. Your target time saved is proportional to the size of your network."}}
                <chart height="80" src="{{CreateLineChart "line-chart-3" 10 0.0 10.0}}"></chart>
            </division>
        </table-cell>
    </table>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>

    <table columns="2" margin="0 0 10 0">
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 12 "Text" "Smart Alert Summary"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "right" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "2 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "Summary of Smart Alerts detected in your network during the report period."}}
    </table>

    <table columns="5" column-widths="0.35 0.1 0.15 0.2 0.2">
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 5 "Margin" "3 0" "LineHeight" "1" "Font" "helvetica-bold" "FontSize" 8 "Text" "Smart Alert Type"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 5 "Margin" "3 0" "LineHeight" "1" "Font" "helvetica-bold" "FontSize" 8 "Text" "Status"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 5 "Margin" "3 0" "LineHeight" "1" "Font" "helvetica-bold" "FontSize" 8 "Text" "# Major Actors"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 5 "Margin" "3 0" "LineHeight" "1" "Font" "helvetica-bold" "FontSize" 8 "Text" "Time First Seen"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 5 "Margin" "3 0" "LineHeight" "1" "Font" "helvetica-bold" "FontSize" 8 "Text" "Time Last Triggered"}}

        {{template "smart-alert-summary-row" dict "Type" "Suspicious Activity On an Asset" "StatusColor" "#ff0000" "NoActors" "1" "FirstSeen" "11/07/2019 04:00:00 UTC" "LastTriggered" "11/07/2019 05:00:00 UTC"}}
        {{template "smart-alert-summary-row" dict "Type" "Suspicious Tunneling Plus Data Exfiltration" "StatusColor" "#ffff00" "NoActors" "1" "FirstSeen" "11/07/2019 10:30:00 UTC" "LastTriggered" "11/08/2019 10:42:09 UTC"}}
        {{template "smart-alert-summary-row" dict "Type" "Internal to External Probing or Reconnaissance Activity" "StatusColor" "#159635" "NoActors" "" "FirstSeen" "" "LastTriggered" ""}}
        {{template "smart-alert-summary-row" dict "Type" "Probing or Reconnaissance Activity" "StatusColor" "#159635" "NoActors" "" "FirstSeen" "" "LastTriggered" ""}}
        {{template "smart-alert-summary-row" dict "Type" "Suspicious Activity On an Untrusted Private IP" "StatusColor" "#159635" "NoActors" "" "FirstSeen" "" "LastTriggered" ""}}
    </table>

    <table columns="3" margin="10 70 0 50">
        {{template "status-legend-cell" dict "Color" "#ff0000" "Text" " - High Threat"}}
        {{template "status-legend-cell" dict "Color" "#ffff00" "Text" " - Medium Threat"}}
        {{template "status-legend-cell" dict "Color" "#159635" "Text" " - Low Threat"}}
    </table>
</chapter>

<chapter show-numbering="false">
    <chapter-heading font="helvetica-bold" font-size="16" margin="0 0 5 0">Network Visibility</chapter-heading>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>
    {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Your Network over the previous 7 days"}}

    <table columns="5" margin="10 0">
        {{template "table-cell-paragraph" dict "Colspan" 2 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Internal IPs"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "External IPs"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Network Flows"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#e2ecf0" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Traffic in Bytes"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#cee5b5" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "Trusted: 96"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#fbc4ab" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "Untrusted: 94"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "16,879"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "577,458"}}
        {{template "table-cell-paragraph" dict "Colspan" 1 "Align" "center" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#e2ecf0" "BorderSize" 1 "Indent" 0 "Margin" "5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "16,112,333,568"}}
    </table>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>

    <table columns="2" column-widths="0.5 0.5">
        <table-cell indent="0" align="center">
            <division>
                {{$src := CreateBarChart "bar-chart-1"
                    (dict "SSH Server" 15.0 "Web Server" 9.0 "DNS Server" 5.0 "FTP Server" 4.0 "DHCP Server" 4.0)}}
                {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#ff0000" "GradeText" "F" "Title" "Unidentified Assets" "Subtitle" "68.2%" "HideBar" true}}
                <chart height="100" src="{{$src}}" margin="0 20 20 0"></chart>
            </division>
        </table-cell>
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 7 "Text" "Unidentified Assets are those that UniPDF sees that you have not labeled and rated. By applying labels and importance ratings, you provide important context to UniPDF in better understanding what threats are most critical to you."}}
                {{template "simple-paragraph" dict "Margin" "5 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 7 "Text" "Optimally, there should be no unidentified assets on your network. You may, however, have an assets or 2 pop up that needs to be identified. Address them quickly by labeling them or remediating any rouge assets. Don't let them accumulate."}}
                <chart height="80" src="{{CreateLineChart "line-chart-4" 5 0.0 5.0}}"></chart>
            </division>
        </table-cell>
    </table>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>

    <table columns="2" column-widths="0.5 0.5">
        <table-cell indent="0" align="center">
            <division>
                {{$src := CreateStackedBarChart "stacked-bar-chart-2"
                    (CreateStackedBar "SQL Database Server" (dict "" 1.0))
                    (CreateStackedBar "Exchange" (dict "" 1.0))
                    (CreateStackedBar "DHCP Server" (dict "" 1.0))
                }}
                {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#159635" "GradeText" "A" "Title" "High Risk Assets " "Subtitle" "0.0%" "HideBar" true}}
                <chart height="120" src="{{$src}}" margin="0 20 0 0"></chart>
            </division>
        </table-cell>
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 7 "Text" "You know which assets are important to your business. UniPDF knows which assets are most likely the target of threatening behavior. That’s how we rate risk."}}
                {{template "simple-paragraph" dict "Margin" "5 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 7 "Text" "Work to reduce the number of high risk assets to no more than a few. Do this by addressing Smart Alerts promptly and protecting your systems against attack."}}
                <chart height="80" src="{{CreateLineChart "line-chart-5" 15 0.0 15.0}}"></chart>
            </division>
        </table-cell>
    </table>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>

    <table columns="2" column-widths="0.5 0.5">
        <table-cell indent="0" align="center">
            <division>
                {{template "sectionContent" dict "BackgroundColor" "#ffffff" "GradeColor" "#ff0000" "GradeText" "F" "Title" "Unidentified Subnets or IP Ranges " "Subtitle" "100.0%" "HideBar" true}}
                <chart height="80" src="{{CreatePieChart "pie-chart-2" (dict "Defined" 60.0 "Undefined" 40.0)}}" margin="-10 0 0 0"></chart>
            </division>
        </table-cell>
        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 7 "Text" "Unidentified Subnets are those that UniPDF sees that you have not labeled. By applying labels, you provide important context to UniPDF in better understanding what threats are most critical to you."}}
                <chart height="80" src="{{CreateLineChart "line-chart-6" 80 0.0 50.0}}"></chart>
            </division>
        </table-cell>
    </table>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>
</chapter>

<page-break></page-break>

<chapter show-numbering="false">
    <chapter-heading font="helvetica-bold" font-size="16" margin="0 0 5 0">How to Use this Report</chapter-heading>
    <line fit-mode="fill-width" position="relative" color="#555" margin="0 0 5 0"></line>
    {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "Whether you are evaluating UniPDF for use or actively protecting your network with it, this report provides you with a quick and easy assessment of your network, enabling you to see where key threats and vulnerabilities are."}}

    <table columns="2" column-widths="0.5 0.5" margin="20 0 0 0">
        <table-cell indent="0" align="left">
            <image src="path('templates/res/page1.png')" fit-mode="fill-width" margin="0 50 0 0"></image>
        </table-cell>
        <table-cell indent="0" align="left">
            <division>
                {{template "simple-paragraph" dict "Margin" "5 0 0 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Threat Score"}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "Your threat score provides you an overall measure of threats and vulnerabilities that UniPDF detects. The score enables you to track your progress over time and compare your network to that of other UniPDF customers."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "The score is calculated like a credit score, on a scale of 300 to 850. Your letter grade reflects your performance compared to others. Most get a B. But we all strive for an A."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Metrics"}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "UniPDF tracks 7 key metrics across 3 key areas: Network Visibility, Policy Assurance and Threat Detection."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "These metrics allow you to see your progress in each area so you can work on increasing your score."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "The additional pages of the report provide more detail about each one of these three areas."}}
            </division>
        </table-cell>
        <table-cell indent="0" align="left">
            <image src="path('templates/res/pages.png')" fit-mode="fill-width" margin="20 20 0 0"></image>
        </table-cell>
        <table-cell indent="0" align="left">
            <division>
                {{template "simple-paragraph" dict "Margin" "30 0 0 0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 8 "Text" "Metric Detail"}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "Pages two, three and four provide more details into the key metrics displayed on page one. Each metric includes a fourteen day trend chart showing how the metrics has varied over the preceding 14 days."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "It also contains additional charts that show specific information about the metrics."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "The Manual Effort Saved metric is Not Applicable when there were no Smart Alerts or Behaviors generated during the reporting period."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "The Unidentified Assets metric is Not Applicable when there are no assets defined in the system nor any detected undefined assets."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "The High Risk Assets metric is Not Applicable when there are no assets defined in the system."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "The Unidentified Subnets or IP Ranges metric is Not Applicable when there are no subnet defined in the system nor any detected undefined subnets."}}
                {{template "simple-paragraph" dict "Margin" "15 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 8 "Text" "The Policy Alerts metric is Not Applicable when there are no active policies."}}
            </division>
        </table-cell>
    </table>
</chapter>
