<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Chapters"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk font="helvetica-bold" color="text">Chapters </text-chunk>
        <text-chunk color="text">are used to group content into logical units, much like a book does. </text-chunk>
        <text-chunk color="text">They are top-level components, the only valid parent a chapter can have is another chapter, thus making it a subchapter. </text-chunk>
        <text-chunk font="helvetica-bold" color="text">Chapters </text-chunk>
        <text-chunk color="text">automatically appear in the document's table of contents. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<chapter>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of chapters:"}}

    {{$codeBlock :=
`<chapter>
    <chapter-heading>Chapter title</chapter-heading>
</chapter>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "Code" $codeBlock}}

    {{$codeBlock :=
`<chapter>
    <chapter-heading>Chapter title</chapter-heading>
    <chapter>
        <chapter-heading>Subchapter title</chapter-heading>
    </chapter>
</chapter>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "10 0" "Text" "Supported attributes:"}}

    <!-- show-numbering -->
    {{$codeBlocks := array
`<chapter>
    <chapter-heading>Chapter title</chapter-heading>
</chapter>`
`<chapter show-numbering="false">
    <chapter-heading>Chapter title</chapter-heading>
</chapter>`
    }}
    {{template "attr-showcase" dict "AttrName" "show-numbering"
        "AttrDescr" "The `show-numbering` attribute controls whether the chapter numbers are displayed."
        "HideResult" true
        "CodeBlocks" $codeBlocks
        "AttrDefault" "true"
        "AttrValues" (array "true" "false")
    }}

    <!-- include-in-toc -->
    {{$codeBlocks := array
`<chapter>
    <chapter-heading>Chapter title</chapter-heading>
</chapter>`
`<chapter include-in-toc="false">
    <chapter-heading>Chapter title</chapter-heading>
</chapter>`
    }}
    {{template "attr-showcase" dict "AttrName" "include-in-toc"
        "AttrDescr" "The `include-in-toc` attribute controls whether the chapter is included in the document's table of contents."
        "HideResult" true
        "CodeBlocks" $codeBlocks
        "AttrDefault" "true"
        "AttrValues" (array "true" "false")
    }}

    <!-- margin -->
    {{$codeBlocks := array
`<chapter margin="5 10">
    <chapter-heading>Chapter title</chapter-heading>
</chapter>`
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
    <chapter margin="10 0 0 0">
        {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Chapter Headings" "AlternateText" "Chapter » Headings"}}

        <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
            <text-chunk font="helvetica-bold" color="text">Chapter headings </text-chunk>
            <text-chunk color="text">are used to display the title and optionally the number of a </text-chunk>
            <text-chunk font="helvetica-bold" color="text">chapter</text-chunk>
            <text-chunk color="text">. The component is represented in templates using the </text-chunk>
            <text-chunk color="secondary">{{xmlEscape "<chapter-heading>"}} </text-chunk>
            <text-chunk color="text">tag.</text-chunk>
        </paragraph>

        {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "10 0" "Text" "Supported attributes:"}}

        <!-- font -->
        {{$codeBlocks := array
`<chapter>
    <chapter-heading font="courier">Chapter title</text-heading>
</chapter>`
        }}
        {{template "attr-showcase" dict "AttrName" "font"
            "AttrDescr" "The `font` attribute specifies the font used for rendering the chapter heading."
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.1. Fonts"
            "AttrLink" "page(46, 0, 130)"
            "AttrDefault" "helvetica"
            "HideResult" true
            "CodeBlocks" $codeBlocks
        }}

        <!-- font-size -->
        {{$codeBlocks := array
`<chapter>
    <chapter-heading font-size="12">Chapter title</chapter-heading>
</chapter>`
        }}
        {{template "attr-showcase" dict "AttrName" "font-size"
            "AttrDescr" "The `font-size` attribute specifies the size of the font used by the chapter heading."
            "CodeBlocks" $codeBlocks
            "HideResult" true
            "AttrDefault" "10"
        }}

        <!-- text-align -->
        {{$codeBlocks := array
`<chapter text-align="center">
    <chapter-heading>Chapter title</chapter-heading>
</chapter>`
        }}
        {{template "attr-showcase" dict "AttrName" "text-align"
            "AttrDescr" "The `text-align` attribute aligns the chapter heading horizontally based on the specified option."
            "CodeBlocks" $codeBlocks
            "HideResult" true
            "AttrValues" (array "left" "right" "center" "justify")
            "AttrDefault" "left"
        }}

        <!-- line-height -->
        {{$codeBlocks := array
`<chapter>
    <chapter-heading line-height="1.5">Chapter title</chapter-title>
</chapter>`
        }}
        {{template "attr-showcase" dict "AttrName" "line-height"
            "AttrDescr" "The `line-height` represents a scale factor for the vertical space a text line takes. Basically, the line height attribute controls amount of vertical space between the chapter heading text lines."
            "HideResult" true
            "CodeBlocks" $codeBlocks
            "AttrDefault" "1.0"
        }}

        <!-- enable-wrap -->
        {{$codeBlocks := array
`<chapter enable-wrap="false">
    <chapter-heading>Long chapter title</chapter-heading>
</chapter>`
        }}
        {{template "attr-showcase" dict "AttrName" "enable-wrap"
            "AttrDescr" "The `enable-wrap` attribute controls whether the content of the chapter heading should be split into lines based on the available space."
            "HideResult" true
            "CodeBlocks" $codeBlocks
            "AttrDefault" "true"
        }}

        <!-- color -->
        {{$codeBlocks := array
`<chapter>
    <chapter-heading color="#00695c">Chapter title</chapter-heading>
</chapter>`
        }}
        {{template "attr-showcase" dict "AttrName" "color"
            "AttrDescr" "The `color` attribute is used to specify the text color."
            "HideResult" true
            "AttrLinkDescr" "For more information see chapter "
            "AttrLinkText" "4.3. Colors"
            "AttrLink" "page(49, 0, 50)"
            "CodeBlocks" $codeBlocks
            "AttrDefault" "#000000"
        }}

        <!-- margin -->
        {{$codeBlocks := array
`<chapter>
    <chapter-heading margin="5 25">Chapter title</chapter-heading>
</chapter>`
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

        <!-- max-lines -->
        {{$codeBlocks := array
`<chapter>
    <chapter-heading max-lines="2">Chapter title</chapter-heading>
</chapter>`
        }}
        {{template "attr-showcase" dict "AttrName" "max-lines"
            "AttrDescr" "The `max-lines` attribute specifies the maximum number of lines before the chapter heading content is truncated."
            "HideResult" true
            "CodeBlocks" $codeBlocks
            "AttrDefault" "0"
        }}

        <paragraph>
            <text-chunk text="color">Please see chapter </text-chunk>
            <text-chunk link="page(4)">2.1. Paragraphs </text-chunk>
            <text-chunk text="color"> for visual examples of most attributes enumerated in this section.</text-chunk>
        </paragraph>
    </chapter>
</chapter>
