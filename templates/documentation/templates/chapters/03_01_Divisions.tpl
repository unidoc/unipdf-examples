<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Divisions"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">division </text-chunk>
        <text-chunk color="text">component is used to group a collection of components, stacking them vertically. Optionally, divisions can have configurable backgrounds. Basically, a background can be assigned to all supported components by including them a division.</text-chunk>
        <text-chunk color="text">The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<division>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of divisions:"}}

    {{$codeBlock :=
`<division padding="10" margin="5">
    <background border-color="#0772cd" border-size="0.5"
        border-radius="5"></background>

    <image fit-mode="fill-width"
        src="path('templates/res/images/sample-image-5.jpg')"></image>

    <line position="relative" fit-mode="fill-width"
        thickness="0.5" color="#333333" margin="5 0"></line>

    <paragraph text-align="center">
        <text-chunk>A seemingly endless road...</text-chunk>
    </paragraph>
</division>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "15 0 10 0" "Text" "Supported attributes:"}}

    <!-- enable-page-wrap -->
    {{template "attr-showcase" dict "AttrName" "enable-page-wrap"
        "AttrDescr" "The `enable-page-wrap` attribute controls whether the division is wrapped across pages. Page wrapping is enabled by default. When the attribute is set to `false`, the division is moved in its entirety on a new page, if it does not fit in the available height. If the height of the division is larger than an entire page, wrapping is enabled automatically in order to avoid unwanted behavior."
        "CodeBlocks" ""
        "AttrDefault" "true"
        "AttrValues" (array "true" "false")
    }}

    <!-- margin -->
    {{$codeBlocks := array
`<division>
    <background border-color="#cccccc" border-size="0.5"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
`<division margin="10">
    <background border-color="#cccccc" border-size="0.5"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
    }}
    {{template "attr-showcase" dict "AttrName" "margin"
        "AttrDescr" "The `margin` attribute allows setting a configurable amount of space around the component."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- padding -->
    {{$codeBlocks := array
`<division>
    <background border-color="#cccccc" border-size="0.5"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
`<division padding="10">
    <background border-color="#cccccc" border-size="0.5"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
    }}
    {{template "attr-showcase" dict "AttrName" "padding"
        "AttrDescr" "The `padding` attribute allows setting a configurable amount of space inside the component, applied around its children. The background of the division is not affected by the padding attribute."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- backgrounds -->
    <chapter>
        {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Backgrounds" "AlternateText" "Division » Backgrounds"}}

        <paragraph margin="10 0" line-height="1.1" enable-word-wrap="true">
            <text-chunk color="text">A background can be added to a division by assigning it a </text-chunk>
            <text-chunk font="helvetica-bold" color="text">background </text-chunk>
            <text-chunk>component. The background component is represented in templates using the </text-chunk>
            <text-chunk color="secondary">{{xmlEscape "<background>"}} </text-chunk>
            <text-chunk color="text">tag.</text-chunk>
        </paragraph>

        {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Supported attributes: "}}

        <!-- fill-color -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background fill-color="#00695c"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk color="#ffffff">Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "fill-color"
            "AttrDescr" "The `fill-color` attribute sets the fill color of the division background."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#ffffff"
        }}

        <!-- border-size -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background border-size="1" border-color="#000000"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-size"
            "AttrDescr" "The `border-color` attribute sets the border size of the division background."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-color -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background border-size="1" border-color="#5c007a"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-color"
            "AttrDescr" "The `border-color` attribute sets the border color of the division background."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#ffffff"
        }}

        <!-- border-radius -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background border-size="3" border-color="#999999"
        border-radius="5"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-radius"
            "AttrDescr" "The `border-radius` attribute sets the radius used when rendering the corners of the division background."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.5. Border Radius"
            "AttrLink" "page(54, 0, 310)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-top-left-radius -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background border-size="3" border-color="#283593"
        border-top-left-radius="10"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-top-left-radius"
            "AttrDescr" "The `border-top-left-radius` attribute sets the radius of the top left corner of the division background."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-top-right-radius -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background border-size="3" border-color="#009688"
        border-top-right-radius="10"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-top-right-radius"
            "AttrDescr" "The `border-top-right-radius` attribute sets the radius of the top right corner of the division background."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-bottom-left-radius -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background border-size="3" border-color="#ffeb3b"
        border-bottom-left-radius="10"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-bottom-left-radius"
            "AttrDescr" "The `border-bottom-left-radius` attribute sets the radius of bottom left corner of the division background."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- border-bottom-right-radius -->
        {{$codeBlocks := array
`<division padding="5" margin="5">
    <background border-size="3" border-color="#f44336"
        border-bottom-right-radius="10"></background>

    <paragraph text-align="center" vertical-text-align="center">
        <text-chunk>Sample text</text-chunk>
    </paragraph>
</division>`
        }}
        {{template "attr-showcase" dict "AttrName" "border-bottom-right-radius"
            "AttrDescr" "The `border-bottom-right-radius` attribute sets the radius of bottom right corner of the division background."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}
    </chapter>
</chapter>
