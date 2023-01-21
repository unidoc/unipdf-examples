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

{{$props := dict "Colspan" 1 "Rowspan" 1 "Align" "left" "VerticalAlign" "top" "BackgroundColor" "#ffffff" 
    "BorderColor" "#000000" "BorderTopSize" 0 "BorderBottomSize" 0 "Indent" 0 
    "Margin" "3 0 5 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "TextColor" "#000000"}}

<table columns="3" column-widths="0.3 0.35 0.35">
    <table-cell align="center" vertical-align="middle">
        <image src="path('templates/res/hospital_logo.png')" fit-mode="fill-width" margin="0 5"></image>
    </table-cell>
    <table-cell align="center" vertical-align="middle">
        <image src="path('templates/res/clinic_care.png')" fit-mode="fill-width" margin="0 5"></image>
    </table-cell>
    <table-cell align="center" vertical-align="middle">
        <image src="path('templates/res/health_care.png')" fit-mode="fill-width" margin="0 5"></image>
    </table-cell>

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Text" "Sample Medical Center\n123 Main Street\nAnywhere, NY 12345 - 6789") }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "center" "Text" "To Contact Us Call: 123 - 456 - 7890\n\nPhone representatives are available:\n8am to 8pm Monday - Thursday\nand 8am to 4:30pm Friday") }}
    
    <table-cell>
        <table columns="2" column-widths="0.5 0.5">
            {{template "table-cell-paragraph" (extendDict $props "Margin" "0" "Align" "left" "Text" "Guarantor Number:") }}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlgin" "right" "Text" "2nnnnn") }}
            {{template "table-cell-paragraph" (extendDict $props "Margin" "-2 0 0 0" "Align" "left" "Text" "Guarantor Name:") }}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlgin" "right" "Text" "Sample Guarantor") }}
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Text" "Statement Date:") }}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlgin" "right" "Text" "01/01/2023") }}
            {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Text" "Due Date:") }}
            {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlgin" "right" "Text" "Upon Receipt") }}
        </table>
    </table-cell>
</table>

<table columns="5" column-widths="0.15 0.4 0.15 0.15 0.15">
    {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#0261AB" "Margin" "5" "Align" "center" "VerticalAlign" "middle" "TextAlign" "center" "TextColor" "#ffffff" "Font" "helvetica-bold" "Text" "Date of Service" ) }}
    {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#0261AB" "Align" "center" "VerticalAlign" "middle" "TextAlign" "center" "TextColor" "#ffffff" "Text" "Description" ) }}
    {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#0261AB" "Align" "center" "VerticalAlign" "middle" "TextAlign" "center" "TextColor" "#ffffff" "Text" "Charges" ) }}
    {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#0261AB" "Align" "center" "VerticalAlign" "middle" "TextAlign" "center" "TextColor" "#ffffff" "Text" "Payment/Adjustments" ) }}
    {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#0261AB" "Align" "center" "VerticalAlign" "middle" "TextAlign" "center" "TextColor" "#ffffff" "Text" "Patient Balance" ) }}

    {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#ffffff" "Margin" "2" "BorderBottomSize" 1 "BorderColor" "#000000" "Align" "left" "VerticalAlign" "middle" "TextAlign" "left" "Font" "helvetica" "TextColor" "#000000" "Text" "07/01/2020 to 07/01/2020" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Text" "Visit #123 Sample Patient" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "Pharmacy" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "60.53" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "Treatment or Observation Room" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "588.00" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "Insurance Payment" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "-598.53" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "Total Hospital Charge" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "638.53" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "Total Payments" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "-598.53" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "Total Adjustments" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "0.0" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}

    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "left" "TextAlign" "left" "Font" "helvetica-bold" "Text" "Patient Due" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Align" "right" "TextAlign" "right" "Text" "40.00" ) }}
</table>

<table columns="2" column-width="0.5 0.5">
    {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 0 "Align" "left" "TextAlign" "left" "Font" "helvetica" 
        "Text" "MESSAGES:\nWe have filed the medical claims with your insurance.They have indicated the balance is your responsibility. To pay your DIN online, please visit www.ourwebsite.com.\n\nIf you have questions regarding your bill, or for payment arrangements, please call 123 - 456 - 78 or send an email inquiry to aboutmybill@ourwebsite.com"
    )}}

    <table-cell>
        <table columns="2" column-widths="0.8 0.2">
            {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#0261AB" "Margin" "2" "TextColor" "#ffffff" "Font" "helvetica-bold" "Text" "Current Balance" ) }}
            {{template "table-cell-paragraph" (extendDict $props "BackgroundColor" "#0261AB" "Margin" "2" "Align" "right" "TextAlign" "right" "TextColor" "#ffffff" "Font" "helvetica-bold" "Text" "$40.00" ) }}

            {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Margin" "0" "BackgroundColor" "#ffffff" "TextColor" "#0261AB" "Align" "left" "TextAlign" "left"
                "Text" "This is your first notice for the visit above, which includes a list of itemized services rendered.")}}

            {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Margin" "0" "TextColor" "#000000" "Font" "helvetica"
                "Text" "We offer a Financial Aid program for qualified applicants. For more information, please call 123-456-7890 or visit our website at www.ourwebsite.com for more information.")}}
        </table>
    </table-cell>
</table>

<table columns="2" column-widths="0.5 0.5" margin="15 0 0 0">
    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Align" "center" "TextAlign" "center" "Font" "helvetica-bold" "FontSize" 12 "Text" "Please retain statement for your records" ) }}

    <table-cell colspan="2">
        <line position="relative" fit-mode="fill-width" thickness="1" style="dashed" dash-array="10"></line>
    </table-cell>

    {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "Align" "left" "TextAlign" "left" "Font" "helvetica" 
        "FontSize" 9 "Indent" 5 
        "Text" "Please check box if address is incorrect or insurance information has changed, and indicate change(s) on reverse side.") }}

    {{template "table-cell-paragraph" (extendDict $props "Rowspan" 2 "Align" "left" "TextAlign" "left" "VerticalAlign" "top" "Font" "helvetica-bold" 
        "FontSize" 9 "Text" "IF PAYING BY VISA, MASTERCARD, DISCOVER OR AMEX, FILL OUT BELOW") }}

    <table-cell>
        <line position="relative" fit-mode="fill-width" thickness="1" margin="5 0"></line>
    </table-cell>

    <table-cell>
        <division>
            {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 10 "TextColor" "#000000" "Text" "MAKE CHECKS PAYABLE TO"}}
            {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 10 "TextColor" "#0261AB" "Text" "Sample Medical Center\n123 Main Street\nAnywhere, NY 12345 - 6789"}}
        </division>
    </table-cell>

    <table-cell rowspan="2">
        <table columns="4" column-widths="0.25 0.25 0.25 0.25">
            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 1 "Margin" "0 0 10 0" "Font" "helvetica-bold" "FontSize" 9 "Text" "Visa")}}
            {{template "table-cell-paragraph" (extendDict $props "Text" "MasterCard")}}
            {{template "table-cell-paragraph" (extendDict $props "Text" "Discover")}}
            {{template "table-cell-paragraph" (extendDict $props "Text" "Amex")}}

            {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Font" "helvetica" "Text" "Card Number")}}
            {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "Text" "Exp. Date")}}
            {{template "table-cell-paragraph" (extendDict $props "Text" "Amount")}}

            {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Text" "Signature")}}
            {{template "table-cell-paragraph" (extendDict $props "Text" "SVV")}}

            <table-cell colspan="4" border-width-bottom="1" border-color="#000000">
                <table columns="3" column-widths="0.3 0.4 0.3">
                    {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 0 "Margin" "0" "Colspan" 1 "Rowspan" 1 "Text" "Statement Date")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Guarantor number")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Pay the Amount")}}

                    {{template "table-cell-paragraph" (extendDict $props "FontSize" 12 "Text" "07/10/2020")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "2nnnnn")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "$40.00")}}
                </table>
            </table-cell>

            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 1 "Colspan" 2 "Margin" "0 0 20 0" "FontSize" 9 "Text" "Visit # to apply payment")}}
            {{template "table-cell-paragraph" (extendDict $props "Text" "Show amount paid here")}}
        </table>
    </table-cell>

    <table-cell>
        <division>
            {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica-bold" "FontSize" 10 "TextColor" "#000000" "Text" "CHANGE SERVICE REQUESTED"}}
            {{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1 "Font" "helvetica" "FontSize" 10 "TextColor" "#000000" "Text" "For Billing inquries: 123 - 456 - 7890\nPatent Name: Sample Patent"}}
        </division>
    </table-cell>


    {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "Rowspan" 1 "BorderBottomSize" 0 "Margin" "0" "LineHeight" 1.5 "Font" "helvetica-bold" "FontSize" 12 "Text" "SAMPLE GUARANTOR\n123 MAIN STREET\nANYWHERE, NY 12345 - 6789" ) }}
    {{template "table-cell-paragraph" (extendDict $props "Margin" "0" "LineHeight" "1.5" "Font" "helvetica-bold" "FontSize" 12 "Text" "SAMPLE MEDICAL CENTER\n123 MAIN STREET\nANYWHERE, NY 12345 - 6789" ) }}
</table>

<page-break></page-break>

{{template "simple-paragraph" dict "Margin" "0 0 15 0" "LineHeight" 1.5 "Font" "helvetica-bold" "FontSize" 16 "TextColor" "#000000" "Text" "The Sample Medical Center financial assistance policy plain language summary"}}
{{template "simple-paragraph" dict "Margin" "0 0 10 0" "LineHeight" 1.5 "Font" "helvetica" "FontSize" 12 "TextColor" "#000000" "Text" "Sample Medical Center offers financial assistance to eligible patients who are uninsured, underinsured, and ineligible for a government health care program, or who are otherwise unable to pay for medically necessary care based on their individual financial situation.\nPatients seeking financial assistance must apply for the program, which is summarized below."}}

{{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1.5 "Font" "helvetica-bold" "FontSize" 14 "TextColor" "#000000" "Text" "Eligible Services"}}
{{template "simple-paragraph" dict "Margin" "0 0 10 0" "LineHeight" 1.5 "Font" "helvetica" "FontSize" 12 "TextColor" "#000000" "Text" "Eligible services include emergent or medically necessary services provided by the Hospital. Eligible patients include all patients who submit a financial assistance application (including requested documentation) and are determined to be eligible for financial assistance by the Patient Financial Services Department."}}

{{template "simple-paragraph" dict "Margin" "0" "LineHeight" 1.5 "Font" "helvetica-bold" "FontSize" 14 "TextColor" "#000000" "Text" "How to Apply"}}
{{template "simple-paragraph" dict "Margin" "0 0 10 0" "LineHeight" 1.5 "Font" "helvetica" "FontSize" 12 "TextColor" "#000000" "Text" "Financial Assistance applications may be obtained/completed/submitted as follows:"}}

