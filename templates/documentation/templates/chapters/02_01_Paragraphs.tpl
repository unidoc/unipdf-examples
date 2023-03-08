<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Paragraphs"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The paragraph is the primary component used to render blocks of text. A paragraph consists of one or more text chunks, which are the main building blocks of the component. A text chunk is a piece of text with a particular style, similar to a </text-chunk>
        <text-chunk font="helvetica-bold" color="text">HTML span</text-chunk>
        <text-chunk color="text">. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<paragraph>"}} </text-chunk>
        <text-chunk color="text">tag. Paragraphs start on a new line and any components rendered after a paragraph start on a new line as well.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of paragraphs:"}}

    {{$codeBlock :=
`<paragraph>
    <text-chunk>A single text chunk paragraph.</text-chunk>
</paragraph>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{$codeBlock =
`<paragraph>
    <text-chunk>A paragraph with </text-chunk>
    <text-chunk color="#0000ff">three </text-chunk>
    <text-chunk>text chunks.</text-chunk>
</paragraph>`
    }}
    {{template "code-block" dict "Margin" "0 0 10 0" "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Supported attributes: "}}

    <!-- text-align -->
    {{$codeBlocks := array
`<paragraph>
    <text-chunk>Sample text</text-chunk>
</paragraph>`
`<paragraph text-align="right">
    <text-chunk>Sample text</text-chunk>
</paragraph>`
`<paragraph text-align="center">
    <text-chunk>Sample text</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "text-align"
        "AttrDescr" "The `text-align` attribute aligns the paragraph text chunks horizontally based on the specified option."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "left" "right" "center" "justify")
        "AttrDefault" "left"
    }}

    <!-- vertical-text-align -->
    {{$codeBlocks := array
`<paragraph>
    <text-chunk>Sample text</text-chunk>
</paragraph>`
`<paragraph vertical-text-align="center">
    <text-chunk>Sample text</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "vertical-text-align"
        "AttrDescr" "The `vertical-text-align` attribute specifies the reference point used for vertically aligning the paragraph text."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "baseline" "center")
        "AttrDefault" "baseline"
    }}

    <!-- line-height -->
    {{$codeBlocks := array
`<paragraph>
    <text-chunk>Sample
text</text-chunk>
</paragraph>`
`<paragraph line-height="1.5">
    <text-chunk>Sample
text</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "line-height"
        "AttrDescr" "The `line-height` represents a scale factor for the vertical space a text line takes. Basically, the line height attribute controls the amount of vertical space between paragraph text lines."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0"
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

    <!-- margin -->
    {{$codeBlocks := array
`<paragraph>
    <text-chunk>Sample text</text-chunk>
</paragraph>`
`<paragraph margin="10 0 0 50">
    <text-chunk>Sample text</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "margin"
        "AttrDescr" "The `margin` attribute allows setting a configurable amount of space around the component."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- enable-wrap -->
    {{$codeBlocks := array
`<paragraph margin="0 20">
    <text-chunk>Sample text ............................</text-chunk>
    <text-chunk>........................................</text-chunk>
</paragraph>`
`<paragraph margin="0 20" enable-wrap="false">
    <text-chunk>Sample text ............................</text-chunk>
    <text-chunk>........................................</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "enable-wrap"
        "AttrDescr" "The `enable-wrap` attribute controls whether the content of the paragraph should be split into lines based on the available space."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "true"
    }}

    <!-- enable-word-wrap -->
    {{$codeBlocks := array
`<paragraph margin="0 40">
    <text-chunk>Far far away, behind the word mountains, </text-chunk>
    <text-chunk>far from the countries Vokalia and </text-chunk>
    <text-chunk>Consonantia, there live the blind texts.</text-chunk>
</paragraph>`
`<paragraph margin="0 40" enable-word-wrap="true">
    <text-chunk>Far far away, behind the word mountains, </text-chunk>
    <text-chunk>far from the countries Vokalia and </text-chunk>
    <text-chunk>Consonantia, there live the blind texts.</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "enable-word-wrap"
        "AttrDescr" "The `enable-word-wrap` attribute controls whether the content of the paragraph should be split into lines based on the available space."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "false"
    }}

    <!-- text-overflow -->
    {{$codeBlocks := array
`<paragraph margin="0 20" enable-wrap="false">
    <text-chunk>Sample text ............................</text-chunk>
    <text-chunk>........................................</text-chunk>
</paragraph>`
`<paragraph margin="0 20" enable-wrap="false" text-overflow="hidden">
    <text-chunk>Sample text ............................</text-chunk>
    <text-chunk>........................................</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "text-overflow"
        "AttrDescr" "The `text-overflow` attribute controls the bevahior of content which does not fit in the available space and overflows. The attribute has no effect if the `enable-wrap` attribute is set to `true` because in that case the content that does not fit is moved on a new line."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "visible" "hidden")
        "AttrDefault" "visible"
    }}

    <!-- angle -->
    {{$codeBlocks := array
`<paragraph margin="0 5" angle="90">
    <text-chunk>T</text-chunk>
</paragraph>`
`<paragraph margin="0 5" angle="180">
    <text-chunk>T</text-chunk>
</paragraph>`
    }}
    {{template "attr-showcase" dict "AttrName" "angle"
        "AttrDescr" "The `angle` attribute represents an angle specified in degrees to rotate the paragraph content by. The rotation is applied anti-clockwise."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- text chunks -->
    <page-break></page-break>
    <chapter>
        {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Text Chunks" "AlternateText" "Paragraph » Text Chunks"}}

        <paragraph margin="10 0" line-height="1.1" enable-word-wrap="true">
            <text-chunk color="text">The content of paragraphs is constructed using text chunks. Text chunks are bits of text with a particular style, configured using attributes. The text chunk component is represented in templates using the </text-chunk>
            <text-chunk color="secondary">{{xmlEscape "<text-chunk>"}} </text-chunk>
            <text-chunk color="text">tag.</text-chunk>
        </paragraph>

        {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Supported attributes: "}}

        <!-- font -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk font="courier">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "font"
            "AttrDescr" "The `font` attribute specifies the font used for rendering the text chunk."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.1. Fonts"
            "AttrLink" "page(46, 0, 130)"
            "AttrDefault" "helvetica"
            "CodeBlocks" $codeBlocks
        }}

        <!-- font-size -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk font-size="12">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "font-size"
            "AttrDescr" "The `font-size` attribute specifies the size of the font used by the text chunk."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "10"
        }}

        <!-- character-spacing -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk character-spacing="6">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "character-spacing"
            "AttrDescr" "The `character-spacing` attribute is used to change the distance between individual characters."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- horizontal-scaling -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk horizontal-scaling="200">Sample text</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk horizontal-scaling="75">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "horizontal-scaling"
            "AttrDescr" "The `horizontal-scaling` attribute is used to scale the text horizontally by the specified percent."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "100"
        }}

        <!-- color -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk color="#00695c">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "color"
            "AttrDescr" "The `color` attribute is used to specify the text color."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#000000"
        }}

        <!-- link -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk>Go to the start of </text-chunk>
    <text-chunk link="page(8)">page 8</text-chunk>
    <text-chunk>.</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk>Go to </text-chunk>
    <text-chunk link="page(8, 0, 50)">page 8</text-chunk>
    <text-chunk>, at coordinates (0, 50).</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk>Visit </text-chunk>
    <text-chunk link="url('https://unidoc.io')">unidoc.io</text-chunk>
    <text-chunk>.</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "link"
            "AttrDescr" "The `link` attribute is used to turn the text chunk into an internal or external link. The destination of internal links is a page or part of a page within the current document, whereas external links send users to external URLs."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "\"\""
        }}

        <!-- underline -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk underline="true">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "underline"
            "AttrDescr" "The `underline` attribute controls whether the text chunk text is underlined."
            "CodeBlocks" $codeBlocks
            "AttrValues" (array "true" "false")
            "AttrDefault" "false"
        }}

        <!-- underline-color -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk underline="true"
        underline-color="#ff00ff">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "underline-color"
            "AttrDescr" "The `underline-color` attribute configures the color of the decoration line."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "color of the text"
        }}

        <!-- underline-offset -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk underline="true"
        underline-offset="3">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "underline-offset"
            "AttrDescr" "The `underline-offset` attribute configures the offset of the decoration line."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- underline-thickness -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk underline="true" underline-offset="2"
        underline-thickness="2">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "underline-thickness"
            "AttrDescr" "The `underline-thickness` attribute configures the thickness of the decoration line."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "1"
        }}

        <!-- text-rise -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk >Sample </text-chunk>
    <text-chunk text-rise="5">text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "text-rise"
            "AttrDescr" "The `text-rise` attribute represents a vertical offset applied to the text chunk. Text rise values can be either positive or negative."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- rendering-mode -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk rendering-mode="fill">Sample text</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk rendering-mode="stroke"
        outline-size="0.1">Sample text</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk rendering-mode="fill-stroke"
        outline-size="0.1" color="#ffff00">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "rendering-mode"
            "AttrDescr" "The `rendering-mode` attribute determines whether showing text shall cause character outlines to be stroked, filled or a combination of the two. The attribute can also be used to render invisible text."
            "CodeBlocks" $codeBlocks
            "AttrValues" (array "fill" "stroke" "fill-stroke" "invisible")
            "AttrDefault" "fill"
        }}

        <!-- outline-size -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk rendering-mode="fill-stroke"
        outline-size="0.01" color="#ffffff">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "outline-size"
            "AttrDescr" "The `outline-size` attribute is used to specify the outline size of the text. By default, outlines are not rendered. That behavior can be changed by specifying an appropriate `rendering-mode` value such as `fill-stroke`."
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <!-- outline-color -->
        {{$codeBlocks := array
`<paragraph>
    <text-chunk rendering-mode="fill-stroke" outline-size="0.01"
        outline-color="#ff0000" color="#ffff00">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "attr-showcase" dict "AttrName" "outline-color"
            "AttrDescr" "The `outline-color` attribute is used to specify the color of the text outline. By default, outlines are not rendered. That behavior can be changed by specifying an appropriate `rendering-mode` value such as `fill-stroke`."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#ffffff"
        }}
    </chapter>
</chapter>
