<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Ellipses"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">Ellipses can be rendered using the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">ellipse </text-chunk>
        <text-chunk color="text">component. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<ellipse>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of ellipses:"}}

    {{$codeBlock :=
`<ellipse position="relative" width="150" height="50"
    fill-color="#ff80ab" border-color="#f50057"></ellipse>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "15 0 10 0" "Text" "Supported attributes:"}}

    <!-- width -->
    {{$codeBlocks := array
`<ellipse position="relative" width="100" height="50"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "width"
        "AttrDescr" "The `width` attribute sets the width of the ellipse."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- height -->
    {{$codeBlocks := array
`<ellipse position="relative" width="50" height="75"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "height"
        "AttrDescr" "The `height` attribute sets the height of the ellipse."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- fit-mode -->
    {{$codeBlocks := array
`<ellipse position="relative" width="1" height="0.2" fit-mode="fill-width"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "fit-mode"
        "AttrDescr" "The `fit-mode` attribute controls the sizing of the component relative to the available space. When the attribute value is set to `fill-width`, the component is scaled so that it occupies the entire available width, preserving the original aspect ratio. The `width` and `height` of the ellipse must be specified in order to calculate the original aspect ratio of the ellipse."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "none" "fill-width")
        "AttrDefault" "none"
    }}

<!-- margin -->
    {{$codeBlocks := array
`<ellipse position="relative" width="50" height="50"
    margin="10 0 10 100"></ellipse>`
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
        "AttrDescr" "The `position` attribute controls whether the component uses relative or absolute position. In absolute mode, users position the component manually using the `x` and `y` attributes. In relative position mode, the component is positioned relative to the already rendered components."
        "CodeBlocks" ""
        "AttrDefault" "absolute"
        "AttrValues" (array "relative" "absolute")
    }}

    <!-- cx -->
    {{template "attr-showcase" dict "AttrName" "cx"
        "AttrDescr" "The `cx` attribute sets the X coordinate of the center of the component. The attribute is used only when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- cy -->
    {{template "attr-showcase" dict "AttrName" "cy"
        "AttrDescr" "The `cy` attribute sets the Y coordinate of the center of the component. The attribute is used only when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- fill-color -->
    {{$codeBlocks := array
`<ellipse position="relative" width="100" height="50"
    fill-color="#81d4fa"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "fill-color"
        "AttrDescr" "The `fill-color` attribute sets the fill color of the ellipse."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.3. Colors"
        "AttrLink" "page(49, 0, 50)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "#ffffff"
    }}

    <!-- fill-opacity -->
    {{$codeBlocks := array
`<ellipse position="relative" width="100" height="50"
    fill-color="#81d4fa" fill-opacity="0.5"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "fill-opacity"
        "AttrDescr" "The `fill-opacity` attribute allows setting the transparency of the ellipse background."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0"
        "AttrValues" (array "any value between 0 and 1")
    }}

    <!-- border-width -->
    {{$codeBlocks := array
`<ellipse position="relative" width="175" height="75"
    border-width="5"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-width"
        "AttrDescr" "The `border-width` attribute sets the width of the ellipse border."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1"
    }}

    <!-- border-color -->
    {{$codeBlocks := array
`<ellipse position="relative" width="100" height="75" border-width="10"
    border-color="#e64a19"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-color"
        "AttrDescr" "The `border-color` attribute sets the border color of the ellipse."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.3. Colors"
        "AttrLink" "page(49, 0, 50)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "#000000"
    }}

    <!-- border-opacity -->
    {{$codeBlocks := array
`<ellipse position="relative" width="100" height="75" border-width="10"
    border-color="#e64a19" border-opacity="0.5"></ellipse>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-opacity"
        "AttrDescr" "The `border-opacity` attribute allows setting the transparency of the ellipse border."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0"
        "AttrValues" (array "any value between 0 and 1")
    }}
</chapter>
