<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Rectangles"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">Rectangles can be rendered using the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">rectangle </text-chunk>
        <text-chunk color="text">component. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{.newline}}{{xmlEscape "<rectangle>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of rectangles:"}}

    {{$codeBlock :=
`<rectangle position="relative" width="200" height="75"
    fill-color="#98ee99" border-color="#003300"></rectangle>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "15 0 10 0" "Text" "Supported attributes:"}}

    <!-- width -->
    {{$codeBlocks := array
`<rectangle position="relative" width="75" height="50"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "width"
        "AttrDescr" "The `width` attribute sets the width of the rectangle."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- height -->
    {{$codeBlocks := array
`<rectangle position="relative" width="50" height="75"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "height"
        "AttrDescr" "The `height` attribute sets the height of the rectangle."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- fit-mode -->
    {{$codeBlocks := array
`<rectangle position="relative" width="1" height="0.15"
    fit-mode="fill-width"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "fit-mode"
        "AttrDescr" "The `fit-mode` attribute controls the sizing of the component relative to the available space. When the attribute value is set to `fill-width`, the component is scaled so that it occupies the entire available width, preserving the original aspect ratio. The `width` and `height` of the rectangle must be specified in order to calculate the original aspect ratio of the rectangle."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "none" "fill-width")
        "AttrDefault" "none"
    }}

    <!-- margin -->
    {{$codeBlocks := array
`<rectangle position="relative" width="50" height="50"
    margin="10 0 10 100"></rectangle>`
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
        "AttrDescr" "The `position` attribute controls whether the component uses relative or absolute position. In absolute mode, user position the component manually using the `x` and `y` attributes. In relative position mode, the component is positioned relative to the already rendered components."
        "CodeBlocks" ""
        "AttrDefault" "absolute"
        "AttrValues" (array "relative" "absolute")
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

    <!-- fill-color -->
    {{$codeBlocks := array
`<rectangle position="relative" width="100" height="50"
    fill-color="#e64a19"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "fill-color"
        "AttrDescr" "The `fill-color` attribute sets the fill color of the rectangle."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.3. Colors"
        "AttrLink" "page(49, 0, 50)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "#ffffff"
    }}

    <!-- fill-opacity -->
    {{$codeBlocks := array
`<rectangle position="relative" width="100" height="50"
    fill-color="#e64a19" fill-opacity="0.5"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "fill-opacity"
        "AttrDescr" "The `fill-opacity` attribute allows setting the transparency of the rectangle background."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0"
        "AttrValues" (array "any value between 0 and 1")
    }}

    <!-- border-width -->
    {{$codeBlocks := array
`<rectangle position="relative" width="200" height="75"
    border-width="5"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-width"
        "AttrDescr" "The `border-width` attribute sets the width of the rectangle border."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1"
    }}

    <!-- border-color -->
    {{$codeBlocks := array
`<rectangle position="relative" width="100" height="75" border-width="10"
    border-color="#5c007a"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-color"
        "AttrDescr" "The `border-color` attribute sets the border color of the rectangle."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.3. Colors"
        "AttrLink" "page(49, 0, 50)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "#000000"
    }}

    <!-- border-opacity -->
    {{$codeBlocks := array
`<rectangle position="relative" width="100" height="75" border-width="10"
    border-color="#5c007a" border-opacity="0.5"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-opacity"
        "AttrDescr" "The `border-opacity` attribute allows setting the transparency of the rectangle border."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0"
        "AttrValues" (array "any value between 0 and 1")
    }}

    <!-- border-radius -->
    {{$codeBlocks := array
`<rectangle position="relative" width="150" height="75"
    border-radius="10"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-radius"
        "AttrDescr" "The `border-radius` attribute sets the radius used when rendering the corners of the rectangle."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.5. Border Radius"
        "AttrLink" "page(54, 0, 310)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- border-top-left-radius -->
    {{$codeBlocks := array
`<rectangle position="relative" width="150" height="50" border-width="3"
    border-color="#283593" border-top-left-radius="10"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-top-left-radius"
        "AttrDescr" "The `border-top-left-radius` attribute sets the radius of the top left corner of the rectangle."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- border-top-right-radius -->
    {{$codeBlocks := array
`<rectangle position="relative" width="150" height="50" border-width="3"
    border-color="#009688" border-top-right-radius="10"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-top-right-radius"
        "AttrDescr" "The `border-top-right-radius` attribute sets the radius of the top right corner of the rectangle."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- border-bottom-left-radius -->
    {{$codeBlocks := array
`<rectangle position="relative" width="150" height="50" border-width="3"
    border-color="#ffeb3b" border-bottom-left-radius="10"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-bottom-left-radius"
        "AttrDescr" "The `border-bottom-left-radius` attribute sets the radius of bottom left corner of the rectangle."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- border-bottom-right-radius -->
    {{$codeBlocks := array
`<rectangle position="relative" width="150" height="50" border-width="3"
    border-color="#f44336" border-bottom-right-radius="10"></rectangle>`
    }}
    {{template "attr-showcase" dict "AttrName" "border-bottom-right-radius"
        "AttrDescr" "The `border-bottom-right-radius` attribute sets the radius of the bottom right corner of the rectangle."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}
</chapter>
