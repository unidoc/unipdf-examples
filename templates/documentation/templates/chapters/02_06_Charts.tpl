<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Charts"}}

    <paragraph margin="10 0 5 0" line-height="1.1">
        <text-chunk color="text">Charts can be rendered using the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">chart </text-chunk>
        <text-chunk color="text">component. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<chart>"}} </text-chunk>
        <text-chunk color="text">tag. The charts are created using </text-chunk>
        <text-chunk color="secondary">unichart</text-chunk>
        <text-chunk color="text">, our in-house charting library. For more information, please visit </text-chunk>
        <text-chunk link="url('https://github.com/unidoc/unichart')">github.com/unidoc/unichart</text-chunk>
        <text-chunk color="text">.{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Text" "Supported chart types:"}}
    <list indent="10" margin="5 0 15 0">
        <list-item>{{template "paragraph" dict "Text" "Line chart"}}</list-item>
        <list-item>{{template "paragraph" dict "Text" "Bar chart"}}</list-item>
        <list-item>{{template "paragraph" dict "Text" "Stacked bar chart"}}</list-item>
        <list-item>{{template "paragraph" dict "Text" "Pie chart"}}</list-item>
        <list-item>{{template "paragraph" dict "Text" "Donut chart"}}</list-item>
        <list-item>{{template "paragraph" dict "Text" "Scatter chart"}}</list-item>
    </list>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of charts:"}}

    {{createLineChart "line-chart-1" 40 0.0 25.0}}
    {{$codeBlock :=
`<chart src="line-chart-1" height="175"></chart>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{createPieChart "pie-chart-1" (dict "Coffee" 0.44 "Tea" 0.30 "Orange Juice" 0.26) false}}
    {{$codeBlock :=
`<chart src="pie-chart-1" height="175"></chart>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    <page-break></page-break>
    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "10 0 10 0" "Text" "Supported attributes:"}}

    <!-- src -->
    {{$codeBlocks := array
`<chart src="line-chart-1"></chart>`
    }}
    {{template "attr-showcase" dict "AttrName" "src"
        "AttrDescr" "The `src` attribute is used to specify the name of the chart to render. The referenced source chart is loaded from the chart map of the options used to draw the template."
        "HideResult" true
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "5.0. Golang Interface"
        "AttrLink" "page(56, 0, 60)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "empty string"
    }}

    <!-- x -->
    {{template "attr-showcase" dict "AttrName" "x"
        "AttrDescr" "The `x` attribute sets the X coordinate of the top left corner of the component. The attribute is used only when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- y -->
    {{template "attr-showcase" dict "AttrName" "y"
        "AttrDescr" "The `y` attribute sets the Y coordinate of the top left corner of the component. The attribute is used only when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- width -->
    {{$codeBlocks := array
`<chart src="pie-chart-1" width="1024"></chart>`
    }}
    {{template "attr-showcase" dict "AttrName" "width"
        "AttrDescr" "The `width` attribute sets the width of the rendered chart. The attribute is ignored if the chart is rendered in relative position mode (i.e. it is not manually positioned by the user using the `x` and `y` attributes)."
        "HideResult" true
        "CodeBlocks" $codeBlocks
        "AttrDefault" "all available space"
    }}

    <!-- height -->
    {{$codeBlocks := array
`<chart src="donut-chart-1" height="250"></chart>`
    }}
    {{template "attr-showcase" dict "AttrName" "height"
        "AttrDescr" "The `height` attribute sets the height of the rendered chart."
        "HideResult" true
        "CodeBlocks" $codeBlocks
        "AttrDefault" "400"
    }}

    <!-- margin -->
    {{$codeBlocks := array
`<chart src="bar-chart-1" margin="10 5 20 5"></chart>`
    }}
    {{template "attr-showcase" dict "AttrName" "margin"
        "AttrDescr" "The `margin` attribute allows setting a configurable amount of space around the component."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "HideResult" true
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <page-break></page-break>
    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "10 0 20 0" "Text" "Chart types showcase:"}}

    <division enable-page-wrap="false">
        {{template "paragraph" dict "TextAlign" "center" "Font" "helvetica-bold" "Text" "Line Chart"}}
        <chart margin="10 0 60 0" height="250" src="{{createLineChart "line-chart-2" 100.0 0.0 50.0}}"></chart>
    </division>

    <table columns="2" enable-page-wrap="false">
        <table-cell align="center">
            <chart margin="0 0 5 0" height="200" src="{{createPieChart "pie-chart-2" (dict "Science" 0.36 "Literature" 0.15 "History" 0.22 "Math" 0.27) false}}"></chart>
        </table-cell>
        <table-cell align="center">
            <chart margin="-55 0 5 0" height="300" src="{{createPieChart "donut-chart-1" (dict "Science" 0.36 "Literature" 0.15 "History" 0.22 "Math" 0.27) true}}"></chart>
        </table-cell>
        {{template "table-cell-paragraph" dict "Margin" "-25 0 0 0" "TextAlign" "center" "Font" "helvetica-bold" "Text" "Pie Chart"}}
        {{template "table-cell-paragraph" dict "Margin" "-25 0 0 0" "TextAlign" "center" "Font" "helvetica-bold" "Text" "Donut Chart"}}
    </table>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "TextAlign" "center" "Font" "helvetica-bold" "Text" "Bar Chart"}}
        {{$src := createBarChart "bar-chart-1"
            (dict "SSH Server" 15.0 "Web Server" 9.0 "DNS Server" 5.0 "FTP Server" 4.0 "DHCP Server" 4.0)}}
        <chart height="300" src="{{$src}}"></chart>
    </division>

    <division margin="-50 0 0 0" enable-page-wrap="false">
        {{template "paragraph" dict "TextAlign" "center" "Font" "helvetica-bold" "Text" "Stacked Bar Chart"}}
        {{$src := createStackedBarChart "stacked-bar-chart-1"
            (createStackedBar "Q1" (dict "42K" 42. "48K" 48.0 "46K" 46.0 "32K" 32.0))
            (createStackedBar "Q2" (dict "53K" 53. "62K" 62.0 "60K" 60.0 "45K" 45.0))
            (createStackedBar "Q3" (dict "47K" 47. "55K" 55.0 "58K" 58.0 "54K" 54.0))
            (createStackedBar "Q4" (dict "60K" 60. "74K" 74.0 "70K" 70.0 "46K" 46.0))
        }}
        <chart height="200" src="{{$src}}" margin="-15 20 0 0"></chart>
    </division>

    <paragraph margin="-30 0 20 0">
        <text-chunk color="text">Full documentation and more examples of charts can be found at </text-chunk>
        <text-chunk link="url('https://github.com/unidoc/unichart')">github.com/unidoc/unichart</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>
</chapter>

