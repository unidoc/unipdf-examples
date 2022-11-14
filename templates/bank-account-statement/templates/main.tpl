{{define "simple-paragraph"}}
    {{$margin := "0"}}
    {{if .Margin}} {{$margin = .Margin}} {{end}}

    {{$lineHeight := 1}}
    {{if .LineHeight}} {{$lineHeight = .LineHeight}} {{end}}

    {{$align := "left"}}
    {{if .TextAlign}} {{$align = .TextAlign}} {{end}}

    {{$font := "helvetica"}}
    {{if .Font}} {{$font = .Font}} {{end}}

    {{$fontSize := 10}}
    {{if .FontSize}} {{$fontSize = .FontSize}} {{end}}

    {{$textColor := "#000000"}}
    {{if .TextColor}} {{$textColor = .TextColor}} {{end}}

    {{$text := ""}}
    {{if .Text}} {{$text = .Text}} {{end}}

    <paragraph margin="{{$margin}}" line-height="{{$lineHeight}}" text-align="{{$align}}">
        <text-chunk font="{{$font}}" font-size="{{$fontSize}}" color="{{$textColor}}">{{$text}}</text-chunk>
    </paragraph>
{{end}}

{{define "table-cell-paragraph"}}
    {{$colspan := 1}}
    {{if .Colspan}} {{$colspan = .Colspan}} {{end}}

    {{$backgroundColor := "#ffffff"}}
    {{if .BackgroundColor}} {{$backgroundColor = .BackgroundColor}} {{end}}

    {{$align := "left"}}
    {{if .Align}} {{$align = .Align}} {{end}}

    {{$verticalAlign := "top"}}
    {{if .VerticalAlign}} {{$verticalAlign = .VerticalAlign}} {{end}}

    {{$borderColor := "#ffffff"}}
    {{if .BorderColor}} {{$borderColor = .BorderColor}} {{end}}

    {{$borderLeftSize := 0}}
    {{if .BorderLeftSize}} {{$borderLeftSize = .BorderLeftSize}} {{end}}

    {{$borderRightSize := 0}}
    {{if .BorderRightSize}} {{$borderRightSize = .BorderRightSize}} {{end}}

    {{$borderTopSize := 0}}
    {{if .BorderTopSize}} {{$borderTopSize = .BorderTopSize}} {{end}}

    {{$borderBottomSize := 0}}
    {{if .BorderBottomSize}} {{$borderBottomSize = .BorderBottomSize}} {{end}}

    {{$indent := 0}}
    {{if .Indent}} {{$indent = .Indent}} {{end}}

    <table-cell colspan="{{$colspan}}" background-color="{{$backgroundColor}}" align="{{$align}}" vertical-align="{{$verticalAlign}}" border-color="{{$borderColor}}" border-width-left="{{$borderLeftSize}}" border-width-right="{{$borderRightSize}}" border-width-top="{{$borderTopSize}}" border-width-bottom="{{$borderBottomSize}}" indent="{{$indent}}">
        {{template "simple-paragraph" .}}
    </table-cell>
{{end}}

{{define "panel-title"}}
    <background border-radius="2" border-color="#00aff5" border-size="1"></background>
    <division padding="5">
        <background border-radius="2 2 0 0" fill-color="#00aff5"></background>
        {{template "simple-paragraph" dict "Margin" "0 0 5 0" "Font" "helvetica-bold" "FontSize" 14 "TextColor" "#ffffff" "Text" .Title}}
    </division>
{{end}}

<table columns="2" column-widths="0.4 0.6" margin="0">
    <table-cell indent="0">
        <division>
            {{template "panel-title" dict "Title" "Company Information"}}
            <division margin="10 5 108 10">
                <paragraph>
                    <text-chunk font="helvetica-bold" font-size="12">{{.CompanyName}}</text-chunk>
                </paragraph>
                <paragraph margin="10 0 0 0" line-height="1.1">
                    <text-chunk>{{.CompanyAddress}}</text-chunk>
                </paragraph>
            </division>
        </division>
    </table-cell>
    <table-cell indent="0">
        <division margin="0 0 0 7">
            {{template "panel-title" dict "Title" "Questions"}}
            <division margin="10 5 10 10">
                {{template "simple-paragraph" dict "Font" "helvetica-bold" "FontSize" 12 "Text" "Contact"}}
                {{template "simple-paragraph" dict "Margin" "10 0 0 0" "LineHeight" 1.1 "Text" "Available by phone 24 hours a day, 7 days a week. Telecommunications Relay Services calls are accepted."}}
                <paragraph margin="5 0 0 0">
                    <text-chunk font="helvetica-bold" font-size="12">{{.PhoneFree}} </text-chunk>
                    <text-chunk font-size="12">({{.Phone}})</text-chunk>
                </paragraph>
                <table columns="2" column-widths="0.3 0.7" margin="10 0 0 0">
                    {{$props := dict "LineHeight" 1.1}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "TTY")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" .Phone)}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Online")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" .Online)}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Write")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" .White)}}
                </table>
            </division>
        </division>
    </table-cell>
</table>

<division margin="7 0 0 0">
    {{template "panel-title" dict "Title" "Company Information"}}
    <division>
        <table columns="2" column-widths="0.6 0.4" margin="10">
            <table-cell indent="0">
                <division margin="0 5 0 0">
                    {{template "simple-paragraph" dict "Font" "helvetica-bold" "FontSize" 12 "Text" (printf "Your business and %s" .BankName)}}
                    {{template "simple-paragraph" dict "Margin" "10 0 0 0" "LineHeight" 1.1 "Text" "The plans you establish today will shape your business far into the future. The heart of the planning process is your business plan. Take your time to build a strong foundation."}}
                    <paragraph margin="5 0 0 0">
                        <text-chunk>Find out more at </text-chunk>
                        <text-chunk color="#0000ff">{{.BusinessPlanURL}}</text-chunk>
                        <text-chunk>.</text-chunk>
                    </paragraph>
                </division>
            </table-cell>
            <table-cell indent="0">
                <division margin="0 0 0 5">
                    {{template "simple-paragraph" dict "Font" "helvetica-bold" "FontSize" 12 "Text" "Account options"}}
                    <paragraph margin="10 0 0 0">
                        <text-chunk>A check mark in the box indicates you have these convenient services with your account(s). Go to </text-chunk>
                        <text-chunk color="#0000ff">{{.AccountOptionsURL}}</text-chunk>
                        <text-chunk> or call the number above if you have questions or if you would like to add new services.</text-chunk>
                    </paragraph>
                    <table columns="2" column-widths="0.8 0.2" margin="10 0 0 0">
                        {{range $label := .AccountOptionsLabels}}
                            <table-cell indent="0" vertical-align="middle">
                                {{template "simple-paragraph" dict "TextColor" "#000000" "Text" $label}}
                            </table-cell>
                            <table-cell indent="0" align="right" vertical-align="middle">
                                {{template "simple-paragraph" dict "Margin" "-4 0 0 0" "Font" "zapf-dingbats" "FontSize" 12 "Text" "❑"}}
                            </table-cell>
                        {{end}}
                    </table>
                </division>
            </table-cell>
        </table>
    </division>
</division>

<division margin="7 0 0 0">
    {{template "panel-title" dict "Title" "Account Summary"}}
    <division margin="10">
        {{template "simple-paragraph" dict "LineHeight" 1.1 "Text" .Advt}}
    </division>
</division>

<table columns="2" column-widths="0.5 0.5" margin="7 0 0 0">
    <table-cell indent="0">
        <division>
            {{template "panel-title" dict "Title" "Company Information"}}
            <table columns="2" column-widths="0.75 0.25" margin="10">
                {{template "table-cell-paragraph" (dict "Text" (printf "Beginning balance on %s" (formatTime .DateBegin "Jan 02 2006")))}}
                {{template "table-cell-paragraph" (dict "Align" "right" "Text" (printf " $%.2f" .BeginningBalance))}}
                {{template "table-cell-paragraph" (dict "Text" "Deposits/Credits")}}
                {{template "table-cell-paragraph" (dict "Align" "right" "Text" (printf " %.2f" .Deposits))}}
                {{template "table-cell-paragraph" (dict "Text" "Withdrawals/Debits")}}
                {{template "table-cell-paragraph" (dict "Align" "right" "Text" (printf "%.2f" .Withdrawals))}}
                {{template "table-cell-paragraph" (dict "Font" "helvetica-bold" "BorderColor" "#000000" "BorderTopSize" 1 "Text" (printf "Ending balance on %s" (formatTime .DateEnd "Jan 02 2006")))}}
                {{template "table-cell-paragraph" (dict "Align" "right" "Font" "helvetica-bold" "BorderColor" "#000000" "BorderTopSize" 1 "Text" (printf " $%.2f" .EndingBalance))}}
                {{template "table-cell-paragraph" (dict "Margin" "10 0 0 0" "Text" "Average ledger balance this period")}}
                {{template "table-cell-paragraph" (dict "Align" "right" "Margin" "10 0 1 0" "Text" (printf " %.2f" .AverageBalance))}}
            </table>
        </division>
    </table-cell>
    <table-cell indent="0">
        <division margin="0 0 0 7">
            {{template "panel-title" dict "Title" "Account Info"}}
            <division margin="10">
                {{template "simple-paragraph" dict "Font" "helvetica-bold" "Text" .CompanyName}}
                {{template "simple-paragraph" dict "Margin" "2 0 0 0" "Text" "New York account terms and conditions apply."}}
                {{template "simple-paragraph" dict "Margin" "10 0 0 0" "Text" "For Direct Deposit use:"}}
                {{template "simple-paragraph" dict "Margin" "2 0 0 0" "Text" (printf "Routing Number (RTN): %s" .DepositRTN)}}
                {{template "simple-paragraph" dict "Margin" "10 0 0 0" "Text" "For Wire Transfers use:"}}
                {{template "simple-paragraph" dict "Margin" "2 0 0 0" "Text" (printf "Routing Number (RTN): %s" .WireRTN)}}
            </division>
        </division>
    </table-cell>
</table>

<division margin="7 0 0 0">
    {{template "panel-title" dict "Title" "Overdraft Protection"}}
    <division margin="10 10 20 10">
        {{template "simple-paragraph" dict "Text" "Your account is linked to the following for overdraft protection:"}}
        {{template "simple-paragraph" dict "Text" "Savings - 000001234567890"}}
    </division>
</division>

<division>
    {{template "panel-title" dict "Title" "Transaction History"}}
    <division margin="5" padding="5">
        <table columns="6" column-widths="0.075 0.1 0.405 0.14 0.14 0.14">
            {{$props := dict "BorderColor" "#000000" "BorderTopSize" 1}}
            {{template "table-cell-paragraph" (extendDict $props "BorderTopSize" 0 "Text" "Date")}}
            {{template "table-cell-paragraph" (extendDict $props "BorderTopSize" 0 "TextAlign" "right" "Text" "Check&#xA;Number")}}
            {{template "table-cell-paragraph" (extendDict $props "BorderTopSize" 0 "Indent" 5 "Text" "Description")}}
            {{template "table-cell-paragraph" (extendDict $props "BorderTopSize" 0 "TextAlign" "right" "Text" "Deposits/&#xA;Credits")}}
            {{template "table-cell-paragraph" (extendDict $props "BorderTopSize" 0 "TextAlign" "right" "Text" "Withdrawals/&#xA;Debits")}}
            {{template "table-cell-paragraph" (extendDict $props "BorderTopSize" 0 "TextAlign" "right" "Text" "Ending Daily&#xA;Balance")}}

            {{range $transaction := .Transactions}}
                {{template "table-cell-paragraph" (extendDict $props "Text" (formatTime $transaction.Date "01/02"))}}
                {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" $transaction.Check)}}
                {{template "table-cell-paragraph" (extendDict $props "Indent" 5 "Text" $transaction.Details)}}
                {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf " $%.2f" $transaction.Deposits))}}
                {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf " $%.2f" $transaction.Withdrawals))}}
                {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf " $%.2f" $transaction.EndingDailyBalance))}}
            {{end}}

            {{template "table-cell-paragraph" (extendDict $props "Colspan" 5 "Font" "helvetica-bold" "Text" (printf "Ending balance on %s" (formatTime .DateEnd "Jan 02 2006")))}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Font" "helvetica-bold" "Indent" 3 "Text" (printf " $%.2f" .EndingBalance))}}

            {{template "table-cell-paragraph" (extendDict $props "Colspan" 3 "Font" "helvetica-bold" "BorderBottomSize" 0 "Text" "Totals")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Font" "helvetica-bold" "Indent" 3 "Text" (printf " $%.2f" .TransactionDeposits))}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Font" "helvetica-bold" "Text" (printf " $%.2f" .TransactionWithdrawals))}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Font" "helvetica-bold")}}
        </table>

        {{template "simple-paragraph" dict "Margin" "7 0 0 0" "LineHeight" 1.1 "Text" "The Ending Daily Balance does not reflect any pending withdrawals or holds on deposited funds that may have been outstanding on your account when your transaction posted. If you had insufficient available funds when a transaction posted, fees may have been assessed."}}
    </division>
</division>

<division margin="7 0 0 0">
    {{template "panel-title" dict "Title" "Service Fee Summary"}}
    <division margin="10">
        <paragraph line-height="1.1">
            <text-chunk>For a complete list of fees and detailed account information, please see the {{.BankName}} Fee and Information Schedule and Account Agreement applicable to your account or talk to a banker. Go to </text-chunk>
            <text-chunk color="#0000ff">{{.AccountOptionsURL}} </text-chunk>
            <text-chunk>to find the answers to common questions about the monthly service fee on your account.</text-chunk>
        </paragraph>

        <table columns="4" column-widths="0.45 0.26 0.26 0.03" margin="10 0 0 0">
            {{$props := dict "BorderColor" "#000000" "BorderBottomSize" 1 "Indent" 5}}
            {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" (printf "Fee period %s - %s" (formatTime .DateBegin "Jan 02 2006") (formatTime .DateEnd "Jan 02 2006")))}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf "Standard monthly fee %.2f" .StandardServiceFee))}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf "You paid %.2f" .ServiceFee))}}
            {{template "table-cell-paragraph" (extendDict $props)}}

            <table-cell indent="0" background-color="#ffffff">
                <paragraph>
                    <text-chunk font="helvetica-bold">How to reduce the monthly service fee by {{.ServiceDiscount}}&#xA;</text-chunk>
                    <text-chunk>Have </text-chunk>
                    <text-chunk font="helvetica-bold">ONE </text-chunk>
                    <text-chunk>of the following account requirements</text-chunk>
                </paragraph>
            </table-cell>
            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 0 "Text" "Minimum required")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "BorderBottomSize" 0 "Text" "This fee period")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "BorderBottomSize" 0)}}

            {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "Average ledger balance")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf "$%.2f" .MinimumRequired))}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf "$%.2f" .AverageBalance))}}
            {{template "table-cell-paragraph" (extendDict $props "VerticalAlign" "middle" "FontSize" 12 "Margin" "-4 0 0 0" "Font" "zapf-dingbats" "TextAlign" "right" "Text" "❑")}}
            {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Colspan" 4 "Font" "helvetica-bold" "BorderBottomSize" 0 "Text" "Monthly fee discount(s) (applied when box is checked)")}}

            {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "Online only statements")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf "$%.2f" .MinimumRequired))}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" (printf "$%.2f" .AverageBalance))}}
            {{template "table-cell-paragraph" (extendDict $props "VerticalAlign" "middle" "FontSize" 12 "Margin" "-4 0 0 0" "Font" "zapf-dingbats" "TextAlign" "right" "Text" "❑")}}
        </table>

        <table columns="6" column-widths="0.3 0.1 0.1 0.1 0.2 0.2" margin="20 0 0 0">
            {{$props := dict "BorderColor" "#000000" "BorderBottomSize" 1 "Indent" 5}}
            {{template "table-cell-paragraph" (extendDict $props "Indent" 0 "Text" "Service&#xA;charge description")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" "Units&#xA;used")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" "Units&#xA;included")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" "Excess&#xA;units")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" "Service charge per&#xA;excess units ($)")}}
            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Text" "Total service charge($)")}}

            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 3 "Indent" 0 "Text" "Transactions")}}
            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 3 "TextAlign" "right" "Text" .TransactionUnits)}}
            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 3 "TextAlign" "right" "Text" .TransactionUnitsIncluded)}}
            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 3 "TextAlign" "right" "Text" (printf "%.0f" .TransactionExcessUnits))}}
            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 3 "TextAlign" "right" "Text" (printf "%.2f" .ServiceCharge))}}
            {{template "table-cell-paragraph" (extendDict $props "BorderBottomSize" 3 "TextAlign" "right" "Text" (printf "%.2f" .TotalServiceCharge))}}

            {{template "table-cell-paragraph" (extendDict $props "Colspan" 5 "Font" "helvetica-bold" "BorderBottomSize" 0 "Indent" 0 "Text" "Total service charges")}}
            {{template "table-cell-paragraph" (extendDict $props "Font" "helvetica-bold" "BorderBottomSize" 0 "TextAlign" "right" "Text" (printf "%.2f" .TotalServiceCharge))}}
        </table>

        <table columns="2" column-widths="0.03 0.97" margin="20 0 0 0">
            <table-cell indent="0" border-color="#000000" border-width-top="1">
                {{template "simple-paragraph" dict "Margin" "0" "Font" "zapf-dingbats" "FontSize" 12 "Text" "❑"}}
            </table-cell>
            <table-cell indent="0" border-color="#000000" border-width-top="1">
                <paragraph margin="6 0 0 0">
                    <text-chunk font="helvetica-bold">Your feedback matters&#xA;</text-chunk>
                    <text-chunk>Share your compliments and complaints so we can better serve you.&#xA;</text-chunk>
                    <text-chunk>Call us at {{.FeedbackPhone}} ({{.Phone}}) or visit </text-chunk>
                    <text-chunk color="#0000ff">{{.BusinessPlanURL}}</text-chunk>
                    <text-chunk>.</text-chunk>
                </paragraph>
            </table-cell>
        </table>
    </division>
</division>

<division margin="7 0 0 0">
    {{template "panel-title" dict "Title" "Policies"}}
    <division>
        <table columns="2" column-widths="0.6 0.4" margin="10">
            <table-cell indent="0">
                <paragraph>
                    <text-chunk font="helvetica-bold">Notice</text-chunk>
                    <text-chunk>: {{.BankName}}, {{.BankNameState}} may furnish information about accounts belonging to individuals, including sole proprietorships, to consumer reporting agencies. If this applies to you, you have the right to dispute the accuracy of information that we have reporting by writing to us at: {{.ReportAddress}}</text-chunk>
                </paragraph>
            </table-cell>
            <table-cell indent="2">
                {{template "simple-paragraph" dict "LineHeight" 1.1 "Text" "You must describe the specific information that is innacurate or in dispute and the basis for any dispute with supporting documentation. In the case of information that relates to an identity theft, you will need to provide us with an identity theft report."}}
            </table-cell>
        </table>
    </division>
</division>

<table columns="2" margin="7 0 0 0">
    <table-cell indent="0">
        <division>
            {{template "panel-title" dict "Title" "Instruction"}}
            <division margin="10 5 120 10">
                <table columns="2" column-widths="0.05 0.95">
                    {{$props := dict "LineHeight" 1.1 "FontSize" 7 "Indent" 0}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "1.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Use the following worksheet to calculate your overall account balance.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "2.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Go through your register and mark each check, withdrawal, ATM transaction, payment, deposit or other credit listed on your statement. Be sure that your register shows any interest paid into your account and any service charges, automatic payments or ATM transactions withdrawn from your account during this statement period.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "3.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Use the chart to the right to list any deposits, transfers to your account, outstanding checks, ATM withdrawals, ATM payments or any other withdrawals (including any from previous months) which are listed in your register but not shown on your statement.")}}
                </table>

                {{template "simple-paragraph" dict "Font" "helvetica-bold" "FontSize" 7 "Margin" "15 0 0 0" "Text" "ENTER"}}
                <table columns="4" column-widths="0.05 0.5 0.15 0.3">
                    {{template "table-cell-paragraph" (extendDict $props "Font" "helvetica-bold" "Text" "A.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Text" "The ending balance shown on your statement.")}}
                    {{template "table-cell-paragraph" (extendDict $props)}}
                    {{template "table-cell-paragraph" (extendDict $props)}}
                    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Text" (strRepeat ". " 40))}}
                    <table-cell>
                        <table columns="2" column-widths="0.3 0.7">
                            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Margin" "-1 0 0 0" "Text" "$ ")}}
                            <table-cell indent="0">
                                <line position="relative" fit-mode="fill-width" color="#e9e9e9" margin="8 0 0 0"></line>
                            </table-cell>
                        </table>
                    </table-cell>
                </table>

                {{template "simple-paragraph" dict "Font" "helvetica-bold" "FontSize" 7 "Margin" "7 0 0 0" "Text" "ADD"}}
                <table columns="4" column-widths="0.05 0.5 0.15 0.3">
                    {{template "table-cell-paragraph" (extendDict $props "Font" "helvetica-bold" "Text" "B.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Text" "Any deposits listed in your register or transfers into account which are not shown on your statement.")}}

                    <table-cell>
                        <table columns="2" column-widths="0.3 0.7" margin="0 0 10 0">
                            {{range $i, $unused := (loop 4)}}
                                {{$text := "$ "}}
                                {{if eq $i 3}}
                                    {{$text = "+$ "}}
                                {{end}}

                                {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Margin" "-1 0 0 0" "Text" $text)}}
                                <table-cell indent="0">
                                    <line position="relative" fit-mode="fill-width" color="#e9e9e9" margin="8 0 0 0"></line>
                                </table-cell>
                            {{end}}
                        </table>
                    </table-cell>

                    {{template "table-cell-paragraph" (extendDict $props)}}
                    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Text" (strRepeat ". " 40))}}
                    <table-cell>
                        <table columns="2" column-widths="0.3 0.7">
                            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Margin" "-1 0 0 0" "Text" "$ ")}}
                            <table-cell indent="0">
                                <line position="relative" fit-mode="fill-width" color="#e9e9e9" margin="8 0 0 0"></line>
                            </table-cell>
                        </table>
                    </table-cell>
                </table>

                {{template "simple-paragraph" dict "Font" "helvetica-bold" "FontSize" 7 "Margin" "7 0 0 0" "Text" "SUBTRACT"}}
                <table columns="4" column-widths="0.05 0.5 0.15 0.3">
                    {{template "table-cell-paragraph" (extendDict $props "Font" "helvetica-bold" "Text" "C.")}}
                    {{template "table-cell-paragraph" (extendDict $props "Colspan" 2 "Text" (printf "The total outstanding checks and withdrawals from the chart above. %s" (strRepeat ". " 23)))}}
                    <table-cell vertical-align="bottom">
                        <table columns="2" column-widths="0.3 0.7">
                            {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "Margin" "-2 0 0 0" "Text" "-$ ")}}
                            <table-cell indent="0">
                                <line position="relative" fit-mode="fill-width" color="#e9e9e9" margin="7 0 0 0"></line>
                            </table-cell>
                        </table>
                    </table-cell>
                </table>

                <table columns="4" column-widths="0.05 0.5 0.15 0.3" margin="10 0 0 0">
                    {{template "table-cell-paragraph" (extendDict $props "Colspan" 4 "Font" "helvetica-bold" "Text" "CALCULATE THE ENDING BALANCE")}}
                    {{template "table-cell-paragraph" (extendDict $props "Colspan" 3 "Margin" "-3 0 0 0" "Text" (printf "(Part A + Part B + Part C)&#xA;This amount should be the same as the current balance shown in your check register. %s" (strRepeat ". " 13)))}}
                    <table-cell vertical-align="bottom">
                        <table columns="2" column-widths="0.3 0.7">
                            {{template "table-cell-paragraph" (extendDict $props "Font" "helvetica-bold" "TextAlign" "right" "Margin" "-2 0 0 0" "Text" "$ ")}}
                            <table-cell indent="0">
                                <line position="relative" fit-mode="fill-width" margin="7 0 0 0"></line>
                            </table-cell>
                        </table>
                    </table-cell>
                </table>
            </division>
        </division>
    </table-cell>
    <table-cell indent="0">
        <division margin="0 0 0 7">
            {{template "panel-title" dict "Title" "Balance calculation"}}
            <division margin="10">
                <table columns="3" column-widths="0.3 0.4 0.3">
                    {{$props = dict "TextAlign" "center" "Indent" 0 "Font" "helvetica-bold" "FontSize" 7}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Number")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Items Outstanding")}}
                    {{template "table-cell-paragraph" (extendDict $props "Text" "Amount")}}

                    {{range $unused := (loop 25)}}
                        {{$props = dict "BorderColor" "#000000" "BorderLeftSize" 1 "BorderRightSize" 1 "BorderTopSize" 1 "BorderBottomSize" 1 "Margin" "5 0" "Indent" 0}}
                        {{template "table-cell-paragraph" $props}}
                        {{template "table-cell-paragraph" $props}}
                        {{template "table-cell-paragraph" $props}}
                    {{end}}

                    {{template "table-cell-paragraph" (dict "BorderTopSize" 1 "BorderColor" "#000000")}}
                    {{template "table-cell-paragraph" (dict "VerticalAlign" "middle" "TextAlign" "right" "BorderTopSize" 1 "BorderColor" "#000000" "Font" "helvetica-bold" "FontSize" 7 "Text" "Total Amount $ ")}}
                    {{template "table-cell-paragraph" $props}}
                </table>
            </division>
        </division>
    </table-cell>
</table>
