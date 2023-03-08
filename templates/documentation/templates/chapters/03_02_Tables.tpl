<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Tables"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">Tables can be created using the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">table </text-chunk>
        <text-chunk color="text">component, which is somewhat similar to a </text-chunk>
        <text-chunk font="helvetica-bold" color="text">HTML table</text-chunk>
        <text-chunk color="text">. Tables are comprised of </text-chunk>
        <text-chunk font="helvetica-bold" color="text">cell </text-chunk>
        <text-chunk color="text">components, which are grouped into </text-chunk>
        <text-chunk font="helvetica-bold" color="text">rows </text-chunk>
        <text-chunk color="text">based on the user specified number of </text-chunk>
        <text-chunk font="helvetica-bold" color="text">columns</text-chunk>
        <text-chunk>. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<table>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of tables:"}}

    {{$codeBlock :=
`<table columns="2">
    <table-cell indent="0" border-width-right="1" border-width-bottom="1">
        <image src="path('templates/res/images/sample-image-2.jpg')"
            fit-mode="fill-width" margin="5"></image>
    </table-cell>
    <table-cell indent="0" border-width-left="1" border-width-bottom="1">
        <image src="path('templates/res/images/sample-image-4.jpg')"
            fit-mode="fill-width" margin="5"></image>
    </table-cell>
    <table-cell indent="0" border-width-right="1" border-width-top="1">
        <image src="path('templates/res/images/sample-image-3.jpg')"
            fit-mode="fill-width" margin="5"></image>
    </table-cell>
    <table-cell indent="0" border-width-left="1" border-width-top="1">
        <image src="path('templates/res/images/sample-image-5.jpg')"
            fit-mode="fill-width" margin="5"></image>
    </table-cell>
</table>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{$codeBlock :=
`<table columns="3">
    <table-cell border-width="1" align="center" colspan="3">
        <paragraph>
            <text-chunk font="helvetica-bold">Warm colors</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk color="#d32f2f">Red</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk color="#e64a19">Orange</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk color="#ffc107">Yellow</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center" colspan="3">
        <paragraph>
            <text-chunk font="helvetica-bold">Cool colors</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk color="#1565c0">Blue</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk color="#1b5e20">Green</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk color="#8e24aa">Purple</text-chunk>
        </paragraph>
    </table-cell>
</table>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    <page-break></page-break>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "10 0" "Text" "Supported attributes:"}}

    <!-- columns -->
    {{$codeBlocks := array
`<table>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>D</text-chunk></paragraph>
    </table-cell>
</table>`
`<table columns="2">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>D</text-chunk></paragraph>
    </table-cell>
</table>`
`<table columns="4">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>D</text-chunk></paragraph>
    </table-cell>
</table>`
    }}
    {{template "attr-showcase" dict "AttrName" "columns"
        "AttrDescr" "The `columns` attribute sets the number of columns the table has. The attribute basically controls how the cells are split into rows."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1"
    }}

    <!-- column-widths -->
    {{$codeBlocks := array
`<table columns="3">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
</table>`
`<table columns="3" column-widths="0.2 0.5 0.3">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
</table>`
    }}
    {{template "attr-showcase" dict "AttrName" "column-widths"
        "AttrDescr" "The `column-widths` attribute allows setting the sizes of the columns. The column sizes are specified as an array of ratios between `0.0` and `1.0` and the sum of the elements must equal `1.0`. If the attribute is not provided, all columns are considered equal."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0 / column count"
    }}

    <!-- margin -->
    {{$codeBlocks := array
`<table columns="3">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
</table>`
`<table columns="3" margin="30 20 0 20">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
</table>`
    }}
    {{template "attr-showcase" dict "AttrName" "margin"
        "AttrDescr" "The `margin` attribute allows setting a configurable amount of space around the component."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- header-start-row -->
    {{$codeBlocks := array
`<table columns="2" header-start-row="1" header-end-row="1">
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk font="helvetica-bold">Header column 1</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph>
            <text-chunk font="helvetica-bold">Header column 2</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
</table>`
    }}
    {{template "attr-showcase" dict "AttrName" "header-start-row"
        "AttrDescr" "The `header-start-row` attribute sets the starting row of the optional table header. The header rows are repeated on every page the table spans. If the `header-end-rows` attribute is not specified or its value is 0, the attribute has no effect."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- header-end-row -->
    {{$codeBlocks := array
`<table columns="2" header-start-row="1" header-end-row="2">
    <table-cell border-width="1" align="center" colspan="2">
        <paragraph>
            <text-chunk font="helvetica-bold">Header row 1</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center" colspan="2">
        <paragraph>
            <text-chunk font="helvetica-bold">Header row 2</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
</table>`
    }}
    {{template "attr-showcase" dict "AttrName" "header-end-row"
        "AttrDescr" "The `header-end-row` attribute sets the ending row of the optional table header. The header rows are repeated on every page the table spans. If the `header-start-rows` attribute is not specified or its value is 0, the attribute has no effect."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
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

    <!-- enable-page-wrap -->
    {{template "attr-showcase" dict "AttrName" "enable-page-wrap"
        "AttrDescr" "The `enable-page-wrap` attribute controls whether the table is wrapped across pages. Page wrapping is enabled by default. When the attribute is set to `false`, the table is moved in its entirety on a new page if it does not fit in the available height. If the height of the table is larger than an entire page, wrapping is enabled automatically in order to avoid unwanted behavior."
        "CodeBlocks" ""
        "AttrDefault" "true"
        "AttrValues" (array "true" "false")
    }}

    <!-- enable-row-wrap -->
    {{template "attr-showcase" dict "AttrName" "enable-row-wrap"
        "AttrDescr" "The `enable-page-wrap` attribute controls whether the individual table rows are wrapped across pages, basically splitting their content. Row wrapping is currently available for table cells consisting of paragraph and division components. The behavior is disabled by default."
        "CodeBlocks" ""
        "AttrDefault" "false"
        "AttrValues" (array "true" "false")
    }}

    <chapter margin="5 0 0 0">
        {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Table cells" "AlternateText" "Table » Cells"}}

        <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
            <text-chunk font="helvetica-bold" color="text">Table cells </text-chunk>
            <text-chunk color="text">are the main building blocks of the </text-chunk>
            <text-chunk font="helvetica-bold" color="text">table </text-chunk>
            <text-chunk color="text">component. The cells are arranged into </text-chunk>
            <text-chunk font="helvetica-bold" color="text">rows </text-chunk>
            <text-chunk color="text">based on the number of </text-chunk>
            <text-chunk font="helvetica-bold" color="text">columns </text-chunk>
            <text-chunk color="text">the parent </text-chunk>
            <text-chunk font="helvetica-bold" color="text">table </text-chunk>
            <text-chunk color="text">has. The component is represented in templates using the </text-chunk>
            <text-chunk color="secondary">{{xmlEscape "<table-cell>"}} </text-chunk>
            <text-chunk color="text">tag.</text-chunk>
        </paragraph>

        {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "10 0" "Text" "Supported attributes:"}}

        <!-- colspan -->
        {{$codeBlocks := array
`<table columns="3">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center" colspan="2">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center" colspan="3">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "colspan"
            "AttrDescr" "The `colspan` attribute sets the number of columns a cell spans."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "1"
        }}

        <!-- rowspan -->
        {{$codeBlocks := array
`<table columns="2">
    <table-cell border-width="1" align="center"
        vertical-align="middle" rowspan="2">
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "rowspan"
            "AttrDescr" "The `rowspan` attribute sets the number of rows a cell spans."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "1"
        }}

        <!-- indent -->
        {{$codeBlocks := array
`<table>
    <table-cell border-width="1">
        <paragraph>
            <text-chunk>A</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" indent="0">
        <paragraph>
            <text-chunk>B</text-chunk>
        </paragraph>
    </table-cell>
    <table-cell border-width="1" indent="15">
        <paragraph>
            <text-chunk>C</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "indent"
            "AttrDescr" "The `indent` attribute represents an horizontal offset applied on the left of the table cell content."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "5"
        }}

        <!-- align -->
        {{$codeBlocks := array
`<table columns="3">
    <table-cell border-width="1">
        <paragraph><text-chunk>Left</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>Center</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="right">
        <paragraph><text-chunk>Right</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "align"
            "AttrDescr" "The `align` attribute sets the horizontal alignment of the cell content."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "left"
            "AttrValues" (array "left" "center" "right")
        }}

        <!-- vertical-align -->
        {{$codeBlocks := array
`<table columns="4">
    <table-cell border-width="1" align="center">
        <paragraph><text-chunk>Top</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center" vertical-align="middle">
        <paragraph><text-chunk>Middle</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" align="center" vertical-align="bottom">
        <paragraph><text-chunk>Bottom</text-chunk></paragraph>
    </table-cell>
    <table-cell border-width="1" indent="0" align="center">
        <image src="path('templates/res/images/sample-image-2.jpg')"
            fit-mode="fill-width" margin="5 5 2 2"></image>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "vertical-align"
            "AttrDescr" "The `vertical-align` attribute sets the vertical alignment of the cell content."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "top"
            "AttrValues" (array "top" "middle" "bottom")
        }}

        <!-- border-width -->
        {{$codeBlocks := array
`<table>
    <table-cell>
        <paragraph><text-chunk>No border</text-chunk></paragraph>
    </table-cell>
</table>`
`<table>
    <table-cell border-width="2">
        <paragraph><text-chunk>Border width 2</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-width"
            "AttrDescr" "The `border-width` attribute sets the width of the cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-width-left -->
        {{$codeBlocks := array
`<table>
    <table-cell border-width="1" border-width-left="3">
        <paragraph><text-chunk>Left border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-width-left"
            "AttrDescr" "The `border-width-left` attribute sets the width of the left cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-width-right -->
        {{$codeBlocks := array
`<table>
    <table-cell align="right" border-width="1" border-width-right="3">
        <paragraph><text-chunk>Right border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-width-right"
            "AttrDescr" "The `border-width-right` attribute sets the width of the right cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-width-top -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="1" border-width-top="3">
        <paragraph><text-chunk>Top border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-width-top"
            "AttrDescr" "The `border-width-top` attribute sets the width of the top cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-width-bottom -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="1" border-width-bottom="3">
        <paragraph><text-chunk>Bottom border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-width-bottom"
            "AttrDescr" "The `border-width-bottom` attribute sets the width of the bottom cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-style -->
        {{$codeBlocks := array
`<table>
    <table-cell border-width="1">
        <paragraph><text-chunk>Single border</text-chunk></paragraph>
    </table-cell>
</table>`
`<table>
    <table-cell border-width="1" border-style="double">
        <paragraph><text-chunk>Double border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-style"
            "AttrDescr" "The `border-style` attribute sets the style of the cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "single"
            "AttrValues" (array "single" "double")
        }}

        <!-- border-style-left -->
        {{$codeBlocks := array
`<table>
    <table-cell border-width="1" border-style-left="double">
        <paragraph><text-chunk>Left border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-style-left"
            "AttrDescr" "The `border-style-left` attribute sets the style of the left cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "single"
            "AttrValues" (array "single" "double")
        }}

        <!-- border-style-right -->
        {{$codeBlocks := array
`<table>
    <table-cell align="right" border-width="1" border-style-right="double">
        <paragraph><text-chunk>Right border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-style-right"
            "AttrDescr" "The `border-style-right` attribute sets the style of the right cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "single"
            "AttrValues" (array "single" "double")
        }}

        <!-- border-style-top -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="1" border-style-top="double">
        <paragraph><text-chunk>Top border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-style-top"
            "AttrDescr" "The `border-style-top` attribute sets the style of the top cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "single"
            "AttrValues" (array "single" "double")
        }}

        <!-- border-style-bottom -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="1"
        border-style-bottom="double">
        <paragraph><text-chunk>Bottom border</text-chunk></paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-style-bottom"
            "AttrDescr" "The `border-style-bottom` attribute sets the style of the bottom cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "single"
            "AttrValues" (array "single" "double")
        }}

        <!-- border-color -->
        {{$codeBlocks := array
`<table>
    <table-cell border-width="3">
        <paragraph>
            <text-chunk>Default border color</text-chunk>
        </paragraph>
    </table-cell>
</table>`
`<table>
    <table-cell border-width="3" border-color="#ff6d00">
        <paragraph>
            <text-chunk>Custom border color</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-color"
            "AttrDescr" "The `border-color` attribute sets the color of the cell border."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#000000"
        }}

        <!-- border-color-left -->
        {{$codeBlocks := array
`<table>
    <table-cell border-width="3" border-color="#cccccc"
        border-color-left="#004d40">
        <paragraph>
            <text-chunk>Border color left</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-color-left"
            "AttrDescr" "The `border-color-left` attribute sets the color of the left cell border."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#000000"
        }}

        <!-- border-color-right -->
        {{$codeBlocks := array
`<table>
    <table-cell align="right" border-width="3" border-color="#cccccc"
        border-color-right="#fdd835">
        <paragraph>
            <text-chunk>Border color right</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-color-right"
            "AttrDescr" "The `border-color-right` attribute sets the color of the right cell border."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#000000"
        }}

        <!-- border-color-top -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="3" border-color="#cccccc"
        border-color-top="#283593">
        <paragraph>
            <text-chunk>Border color top</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-color-top"
            "AttrDescr" "The `border-color-top` attribute sets the color of the top cell border."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#000000"
        }}

        <!-- border-color-bottom -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="3" border-color="#cccccc"
        border-color-bottom="#c62828">
        <paragraph>
            <text-chunk>Border color bottom</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-color-bottom"
            "AttrDescr" "The `border-color-bottom` attribute sets the color of the bottom cell border."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#000000"
        }}

        <!-- border-line-style -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="3" border-color="#f50057">
        <paragraph>
            <text-chunk>Line style solid</text-chunk>
        </paragraph>
    </table-cell>
</table>`
`<table>
    <table-cell align="center" border-width="3" border-color="#f50057"
        border-line-style="dashed">
        <paragraph>
            <text-chunk>Line style dashed</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-line-style"
            "AttrDescr" "The `border-line-style` attribute sets the line style of the cell border."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "solid"
            "AttrValues" (array "solid" "dashed")
        }}

        <!-- background-color -->
        {{$codeBlocks := array
`<table>
    <table-cell align="center" border-width="1">
        <paragraph><text-chunk>Default background</text-chunk></paragraph>
    </table-cell>
</table>`
`<table>
    <table-cell align="center" border-width="1" border-color="#102027"
        background-color="#62727b">
        <paragraph>
            <text-chunk color="#ffffff">Custom background</text-chunk>
        </paragraph>
    </table-cell>
</table>`
        }}
        {{template "attr-showcase" dict "AttrName" "background-color"
            "AttrDescr" "The `background-color` attribute sets the background color of the cell."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#ffffff"
        }}
    </chapter>
</chapter>
