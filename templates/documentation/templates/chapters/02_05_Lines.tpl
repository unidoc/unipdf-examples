<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Lines"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">Lines can be rendered using the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">line </text-chunk>
        <text-chunk color="text">component. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<line>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of lines:"}}

    {{$codeBlock :=
`<line position="relative" x1="10" y1="10" x2="280" y2="30"
    color="#dd2c00"></line>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{$codeBlock :=
`<line position="relative" fit-mode="fill-width"
    color="#1b5e20"></line>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "15 0 10 0" "Text" "Supported attributes:"}}

    <!-- fit-mode -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width"></line>`
`<line position="relative" fit-mode="fill-width"
    x1="0" y1="20" x2="0" y2="0"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "fit-mode"
        "AttrDescr" "The `fit-mode` attribute controls the sizing of the component relative to the available space. When the attribute value is set to `fill-width`, the component is scaled so that it occupies the entire available width, preserving the original orientation of the line. In `fill-width` mode, if (`x1`, `y1`) and (`x2`, `y2`) are specified, they are used only to determine orientation of the line."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "none" "fill-width")
        "AttrDefault" "none"
    }}

    <!-- color -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width"
    color="#283593"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "color"
        "AttrDescr" "The `color` attribute sets the color of the line."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.3. Colors"
        "AttrLink" "page(49, 0, 50)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "#000000"
    }}

    <!-- opacity -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width"
    color="#283593" opacity="0.5"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "opacity"
        "AttrDescr" "The `opacity` attribute allows setting the transparency of the line."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0"
        "AttrValues" (array "any value between 0 and 1")
    }}

    <!-- thickness -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width" color="#ad1457"
    thickness="5"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "thickness"
        "AttrDescr" "The `thickness` attribute sets the thickness of the line."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1"
    }}

    <!-- style -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width"
    color="#7f0000"></line>`
`<line position="relative" fit-mode="fill-width"
    color="#7f0000" style="dashed"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "style"
        "AttrDescr" "The `style` attribute sets the style of the line."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "solid" "dashed")
        "AttrDefault" "solid"
    }}

    <!-- dash-array -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width" color="#00b8d4"
    style="dashed" dash-array="4 1"></line>`
`<line position="relative" fit-mode="fill-width" color="#00b8d4"
    style="dashed" dash-array="1 4"></line>`
`<line position="relative" fit-mode="fill-width" color="#00b8d4"
    style="dashed" dash-array="4 1 2"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "dash-array"
        "AttrDescr" "The `dash-array` attribute controls the pattern of dashes and gaps used to paint the line. It specifies both the pattern and the lengths of the dashes and gaps. The attribute has no effect if the `style` attribute is not set to `dashed`."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "[1]"
    }}

    <!-- dash-phase -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width" color="#8b6b61"
    style="dashed" dash-array="1 3 5 7"></line>`
`<line position="relative" fit-mode="fill-width" color="#8b6b61"
    style="dashed" dash-array="1 3 5 7" dash-phase="3"></line>`
`<line position="relative" fit-mode="fill-width" color="#8b6b61"
    style="dashed" dash-array="1 3 5 7" dash-phase="6"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "dash-phase"
        "AttrDescr" "The `dash-phase` attribute controls the distance into the dash pattern at which to start the dash. The attribute has no effect if the `style` attribute is not set to `dashed`."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- margin -->
    {{$codeBlocks := array
`<line position="relative" fit-mode="fill-width" color="#1565c0"
    margin="5 25"></line>`
    }}
    {{template "attr-showcase" dict "AttrName" "margin"
        "AttrDescr" "The `margin` attribute allows setting a configurable amount of space around the component."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- position -->
    {{template "attr-showcase" dict "AttrName" "position"
        "AttrDescr" "The `position` attribute controls whether the component uses relative or absolute position. In absolute mode, users position the component manually using the (`x1`, `y1`) and (`x2`, `y2`) pairs of attributes. In relative position mode, the component is positioned relative to the already rendered components."
        "CodeBlocks" ""
        "AttrDefault" "absolute"
        "AttrValues" (array "relative" "absolute")
    }}

    <!-- x1 -->
    {{template "attr-showcase" dict "AttrName" "x1"
        "AttrDescr" "The `x1` attribute sets the X coordinate of the starting point of the line. The attribute is mainly used when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- y1 -->
    {{template "attr-showcase" dict "AttrName" "y1"
        "AttrDescr" "The `y1` attribute sets the Y coordinate of the starting point of the line. The attribute is mainly used when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- x2 -->
    {{template "attr-showcase" dict "AttrName" "x2"
        "AttrDescr" "The `x2` attribute sets the X coordinate of the ending point of the line. The attribute is mainly used when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- y2 -->
    {{template "attr-showcase" dict "AttrName" "y2"
        "AttrDescr" "The `y2` attribute sets the Y coordinate of the ending point of the line. The attribute is mainly used when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}
</chapter>
