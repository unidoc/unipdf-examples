    <table columns="7" indent="0" column-widths= "0.15 0.25 0.2 0.1 0.03 0.2 0.07">
        <table-cell rowspan="4">
            <image src="path('templates/res/logo.png')" width="95.5" height="41.5" ></image>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial">Medication Reconcilliation/BPMH* for:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial">Social Security Number</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial">DOB</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial"></text-chunk>
            </paragraph>
        </table-cell>
        <table-cell rowspan="5" border-width-bottom="2.5" border-width-left="1">
                    <table columns="2">
                        <table-cell colspan="2" margin="5 0 0 0">
                            <paragraph>
                                <text-chunk font-size="8" font="arial-bold">K00</text-chunk>
                            </paragraph>
                        </table-cell>
                        <table-cell colspan="3">
                            <paragraph>
                                <text-chunk font-size="8" font="arial-bold">PLEASE NOTE: </text-chunk>
                                <text-chunk font-size="8" font="arial">completed calendars MUST be returned to SHC as part of the patient's Medical Record</text-chunk>
                            </paragraph>
                        </table-cell>
                        <table-cell>
                            <paragraph margin="15 0 0 0">
                                <text-chunk font-size="8" font="arial">Information:</text-chunk>
                            </paragraph>
                        </table-cell>
                        <table-cell>
                            <paragraph margin="15 0 0 0">
                                <text-chunk font-size="8" font="arial">0-123-456-789</text-chunk>
                            </paragraph>
                        </table-cell>
                        <table-cell>
                            <paragraph>
                                <text-chunk font-size="8" font="arial">Emergenc:</text-chunk>
                            </paragraph>
                        </table-cell>
                        <table-cell>
                            <paragraph>
                                <text-chunk font-size="8" font="arial">0-123-456-789</text-chunk>
                            </paragraph>
                        </table-cell>
                        <table-cell>
                            <paragraph>
                                <text-chunk font-size="8" font="arial">Website:</text-chunk>
                            </paragraph>
                        </table-cell>
                        <table-cell>
                            <paragraph>
                                <text-chunk font-size="8" font="arial">www.samplehealthcare.com</text-chunk>
                            </paragraph>
                        </table-cell>
                    </table>
        </table-cell>
        <table-cell rowspan="5"  border-width-bottom="2.5" vertical-align="middle">
            <image src="path('templates/res/bar-code.png')" width="30" height="150"></image>
        </table-cell>

        {{/** ----------------------------------------- **/}}
        <table-cell vertical-align="bottom" border-width-bottom="0.5">
            <paragraph>
                <text-chunk font-size="12" font="arial-bold">{{.Patient.Name}}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell vertical-align="bottom" border-width-bottom="0.5">
            <paragraph>
                <text-chunk font-size="12" font="arial-bold">{{.Patient.SocialSecurityNumber}}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="0.5" vertical-align="bottom">
            <paragraph>
                <text-chunk font-size="12" font="arial-bold">{{.Patient.Dob}}</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial"></text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial">Two Week Period From:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial">To:</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk></text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk font-size="8" font="arial"></text-chunk>
            </paragraph>
        </table-cell>
        {{/** ----------------------------------------- **/}}
        <table-cell border-width-bottom="0.5" vertical-align="bottom">
            <paragraph>
                <text-chunk font-size="12" font="arial-bold">02/14/2021</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell border-width-bottom="0.5" vertical-align="bottom">
            <paragraph>
                <text-chunk font-size="12" font="arial-bold">02/28/2021</text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk></text-chunk>
            </paragraph>
        </table-cell>
        <table-cell>
            <paragraph>
                <text-chunk></text-chunk>
            </paragraph>
        </table-cell>
        {{/** ----------------------------------------- **/}}
           <table-cell colspan="5">
            <table columns="4" column-widths="0.20 0.35 0.42 0.03">
                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial">Date (mm/dd/yyyy)</text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                    <paragraph>
                    <text-chunk font-size="8" font="arial">Prepared by (Signature/Printed Name)</text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial">Verified by PhC (Signature/Printed Name)</text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell border-width-bottom="0.5">
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell border-width-bottom="0.5">
                    <paragraph>
                    <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell border-width-bottom="0.5">
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                 <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>

                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial">Date (mm/dd/yyyy)</text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                    <paragraph>
                    <text-chunk font-size="8" font="arial">Verified by RN (Signature/Printed Name)**</text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial">Counselled by (Signature/Printed Name)</text-chunk>
                    </paragraph>
                </table-cell>
                 <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>

                <table-cell border-width-bottom="0.5">
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell border-width-bottom="0.5">
                    <paragraph>
                    <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell border-width-bottom="0.5">
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                 <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>

                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial" >Date (mm/dd/yyyy)</text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial">Parent/Legal Guardian (Signature/Printed Name)</text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                </table-cell>
                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>

                <table-cell border-width-bottom="0.5">
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell border-width-bottom="0.5">
                    <paragraph>
                    <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell border-width-bottom="0.5">
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
                <table-cell>
                    <paragraph>
                        <text-chunk font-size="8" font="arial"></text-chunk>
                    </paragraph>
                </table-cell>
            </table>
        </table-cell>
    </table>

<division margin="20 0 0 0">
<paragraph>
<text-chunk>*  Best Possible Medication History</text-chunk>
</paragraph>
<paragraph>
    <text-chunk>** Verification of steroids medication that are part of the patients therapy treatment</text-chunk>
</paragraph>
</division>


{{define "table-header"}}
    <table-cell  border-width-bottom="2.0" vertical-align="bottom">
        <paragraph margin="0 0 10 0">
            <text-chunk font="arial-bold" font-size="11">Drug &amp; Usage</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell  border-width-bottom="2.0" vertical-align="bottom">
        <paragraph margin="0 0 10 0">
            <text-chunk font="arial-bold" font-size="11">Time</text-chunk>
        </paragraph>
    </table-cell>
    {{range .Days}}
    <table-cell border-width-bottom="2.0" vertical-align="bottom">
        <paragraph margin="0 0 10 0">
        <text-chunk font="arial-bold" font-size="11">{{.}}</text-chunk>
        </paragraph>
    </table-cell>
    {{end}}
{{end}}

{{define "drug-schedule"}}   
    <table-cell  border-width-left="0.5" border-width-top="0.5" border-width-bottom="0.5" vertical-align="top" rowspan="{{len .TimesOfTheDay}}">
        <division>
            <paragraph margin="10 0 0 0">
                <text-chunk font="arial-bold" font-size="11">{{.Name}}</text-chunk>
            </paragraph>
            <paragraph margin="5 0 5 0">
                <text-chunk font="arial" font-size="9"> {{.Description}}</text-chunk>
            </paragraph>
        </division>
    </table-cell>
    {{$daysTaken := .DaysTaken}}
    {{range .TimesOfTheDay}}
        <table-cell  border-width-right="0.5" border-width-top="0.5" border-width-bottom="0.5" vertical-align="bottom">
            <paragraph margin="0 0 10 0">
                <text-chunk font="arial" font-size="9">{{.}}</text-chunk>
            </paragraph>
        </table-cell>
        {{range $daysTaken}}
            {{$bg:="#fffffe"}}
            {{if eq . "T"}}
                {{$bg ="#cfcfcb"}}
            {{end}}
            <table-cell  border-width="0.5" vertical-align="bottom" background-color="{{$bg}}">
                <paragraph margin="0 0 10 0">
                    <text-chunk font="arial" font-size="9"></text-chunk>
                </paragraph>
            </table-cell>
        {{end}}
    {{end}}
{{end}}

<table columns="16" margin="20 0 0 0" column-widths="0.25 0.08 {{getColumnWidths (len .ListOfDays) 0.67}}">
    {{template "table-header" dict "Days" .ListOfDays}}
    {{range $idx, $v := .Drugs}}
        {{template "drug-schedule" .}}
    {{end}}
</table>

<table columns="5" margin="0 0 0 200" column-widths="0.6 0.05 0.15 0.05 0.15">
    <table-cell>
        <paragraph margin="10 0 0 0">
            <text-chunk>Mark each box with a checkmark after you have taken a dose of medicine. If you skipped a dose, please consult your physician or pharmacist. Do not take medicine on the days and times not clearly indicated on this schedule.</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell>
        <division margin="18 0 0 0">
            <ellipse position="relative" width="10" height="10" border-color="#000000" fill-color="#ffffff" border-width="0.4"></ellipse>
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
        <division margin="18 0 0 0">
            <ellipse position="relative" width="10" height="10" border-color="#cfcfcb" fill-color="#cfcfcb" border-width="0.4">Text</ellipse>
        </division>
    </table-cell>
    <table-cell>
        <division margin="18 0 0 0">
            <paragraph>
                <text-chunk>Skip this day</text-chunk>
            </paragraph>
        </division>
    </table-cell>
</table>

