{{define "simple-paragraph"}}
    <paragraph margin="{{.Margin}}" line-height="{{.LineHeight}}" text-align="{{.TextAlign}}">
        <text-chunk font="{{.Font}}" font-size="{{.FontSize}}" color="{{.TextColor}}">{{.Text}}</text-chunk>
    </paragraph>
{{end}}

{{define "table-cell-paragraph"}}
    <table-cell colspan="{{.Colspan}}" rowspan="{{.Rowspan}}" background-color="{{.BackgroundColor}}" align="{{.Align}}" vertical-align="{{.VerticalAlign}}" border-color="{{.BorderColor}}" border-width="{{.BorderSize}}" indent="{{.Indent}}">
        {{template "simple-paragraph" .}}
    </table-cell>
{{end}}

{{define "calculation-entry"}}
    <table-cell align="left">
        <paragraph>
            <text-chunk font="helvetica" font-size="9" color="#000000">{{.CalcLabel}}</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell align="right">
        <paragraph>
            <text-chunk font="helvetica" font-size="9" color="#000000">{{.CalcValue}}</text-chunk>
        </paragraph>
    </table-cell>

    <table-cell colspan="2">
        <line position="relative" fit-mode="fill-width" thickness="{{.LineThickness}}" margin="0 0 0 0"></line>
    </table-cell>
{{end}}

<chapter show-numbering="false" margin="0 0 0 0">
    {{$props := dict "Colspan" 1 "Rowspan" 1 "Align" "left" "VerticalAlign" "middle" "BackgroundColor" "#ffffff" "BorderColor" "#ffffff" "BorderSize" 0 "Indent" 0 "Margin" "0 0 0 0" "LineHeight" 1 "Font" "helvetica" "FontSize" 9 "TextColor" "#000000" "TextAlign" "left"}}

    <table columns="2" column-widths="0.5 0.5">
        {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Colspan" 2 "Font" "helvetica-bold" "FontSize" 24 "Text" .firmName ) }}

        {{template "table-cell-paragraph" (extendDict $props "Align" "right" "Colspan" 2 "Font" "Helvetica" "FontSize" 9 "Text" .trade.Date ) }}
        {{template "table-cell-paragraph" (extendDict $props "Align" "right" "Colspan" 2 "Margin" "0 0 10 0" "Text" (printf "Trade Confirmation - Account # %s" .trade.AccountNumber )) }}

        {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "Margin" "0 0 0 0" "Align" "left" "TextAlign" "left" "Text" .trade.Name ) }}
        {{template "table-cell-paragraph" (extendDict $props "Margin" "0 0 0 0" "TextAlign" "right" "Text" .firmName ) }}

        {{template "table-cell-paragraph" (extendDict $props "Margin" "-3 0 0 0" "Align" "left" "TextAlign" "left" "VerticalAlign" "top" "Text" .trade.Address ) }}
        {{template "table-cell-paragraph" (extendDict $props "TextAlign" "right" "VerticalAlign" "top" "Text" .firmAddress ) }}

        {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Margin" "5 0 0 0" "TextAlign" "left" "Colspan" 2 "Font" "helvetica-bold" "Text" (printf "Trade Confirmation - Account # %s" .trade.AccountNumber )) }}
        {{template "table-cell-paragraph" (extendDict $props "Align" "left" "Margin" "-5 0 0 0" "TextAlign" "left" "Colspan" 2 "Font" "helvetica" "Text" "We are pleased to confirm the below transaction:") }}
    </table>

    <line position="relative" fit-mode="fill-width" thickness="1" margin="15 0 15 0"></line>

    <table columns="2" column-widths="0.2 0.8">
        {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "Rowspan" 3 "Margin" "30 0 30 0" "Align" "center" "BackgroundColor" "#00B91E" "Font" "helvetica-bold" "FontSize" 24 "TextColor" "#ffffff" "Text" .trade.Action ) }}

        {{template "table-cell-paragraph" (extendDict $props "Rowspan" 1 "Margin" "0 0 0 0" "Indent" 10 "Align" "left" "BackgroundColor" "#ffffff" 
            "Font" "helvetica" "FontSize" 9 "TextColor" "#000000" "TextAlign" "left" "Text" .trade.ProductDesc ) }}

        <table-cell rowspan="2" indent="0">
            <table columns="2" column-widths="0.5 0.5">
                {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "Text" "You Bought:" ) }}
                {{template "table-cell-paragraph" (extendDict $props "Colspan" 1 "Text" "Price:" ) }}

                {{template "table-cell-paragraph" (extendDict $props "Font" "helvetica-bold" "Text" .trade.BoughtUnit ) }}
                {{template "table-cell-paragraph" (extendDict $props "Font" "helvetica-bold" "Text" .trade.BoughtPrice ) }}
            </table>
        </table-cell>

        {{template "table-cell-paragraph" (extendDict $props "Align" "center" "Font" "helvetica" "Text" "Trade" ) }}

        <table-cell indent="0">
            <division>
                {{template "simple-paragraph" dict "Margin" "0 0 0 10" "LineHeight" 1 "TextAlign" "left" "Font" "helvetica" "FontSize" 9 "TextColor" "#000000" "Text" (printf "Order Number: %s" .trade.OrderNumber) }}
                {{template "simple-paragraph" dict "Margin" "5 0 5 10" "LineHeight" 1 "TextAlign" "left" "Font" "helvetica-bold" "FontSize" 12 "TextColor" "#000000" "Text" "Trade Calculation" }}

                <line position="relative" fit-mode="fill-width" thickness="2" margin="0 0 5 5"></line>

                <table columns="2" column-widths="0.5 0.5">
                    <table-cell>
                        <table columns="2" column-widths="0.5 0.5">
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Principal Amount*:" "CalcValue" .trade.Calculation.PrincipalAmount }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Accrued Interest:" "CalcValue" .trade.Calculation.AccruedInterest }}
                            {{template "calculation-entry" dict "LineThickness" 2 "CalcLabel" "Transaction Fee:" "CalcValue" .trade.Calculation.TransactionFee }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Total:" "CalcValue" .trade.Calculation.Total }}
                            {{template "calculation-entry" dict "LineThickness" 2 "CalcLabel" "Bank Qualified:" "CalcValue" .trade.Calculation.BankQualified }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "State:" "CalcValue" .trade.Calculation.State }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Bank Qualified:" "CalcValue" .trade.Calculation.BankQualified }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Dated to Date:" "CalcValue" .trade.Calculation.DatedToDate }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Yield to Maturity:" "CalcValue" .trade.Calculation.YieldToMaturity }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Yield to Call:" "CalcValue" .trade.Calculation.YieldToCall }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" (printf "Callable %s\nFederally Tax Exempt" .trade.Calculation.Callable) "CalcValue" .trade.Calculation.TaxExcempt }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Capacity" "CalcValue" .trade.Calculation.Capacity }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Bond Form" "CalcValue" .trade.Calculation.BondForm }}
                        </table>
                    </table-cell>

                    <table-cell>
                        <table columns="2" column-widths="0.5 0.5">
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Trade Date:" "CalcValue" .trade.Calculation.TradeDate }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Trade Time:" "CalcValue" .trade.Calculation.TradeTime }}
                            {{template "calculation-entry" dict "LineThickness" 1 "CalcLabel" "Settlement Date:" "CalcValue" .trade.Calculation.SettlementDate }}                
                        </table>
                    </table-cell>
                </table>

                {{template "simple-paragraph" dict "Margin" "0 0 10 10" "LineHeight" 1 "TextAlign" "left" 
                    "Font" "helvetica" "FontSize" 9 "TextColor" "#000000" "Text" 
                    (printf "*This Principal Amount includes a mark-up of %s (1.00%% of the prevailing market price of the security). A mark-up is the amount you paid to %s over and above the prevailing market price of the security. It typically includes compensation to your financial advisor and an additional amount that may account for %s's expenses in the transaction and/or risk taken by %s." .trade.MarkupValue .firmName .firmName .firmName) }}
                {{template "simple-paragraph" dict "Margin" "0 0 0 10" "LineHeight" 1 "TextAlign" "left" 
                    "Font" "helvetica" "FontSize" 9 "TextColor" "#000000" "Text" 
                    (printf "For more information about this security (including the official statement and trade and price history), visit %s." .trade.InfoUrl ) }}
            </division>
        </table-cell>
    </table>
</chapter>