<chapter margin="10 0 0 0">
    {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Border Radius"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk color="secondary">border-radius </text-chunk>
        <text-chunk color="text">attribute sets the radius used when rendering border corners. The attribute can be used to set the radius for each of the border corners and is supported by most components which have border support. The </text-chunk>
        <text-chunk color="secondary">border-top-left-radius</text-chunk>
        <text-chunk color="text">, </text-chunk>
        <text-chunk color="secondary">border-top-right-radius</text-chunk>
        <text-chunk color="text">, </text-chunk>
        <text-chunk color="secondary">border-bottom-left-radius </text-chunk>
        <text-chunk color="text">and </text-chunk>
        <text-chunk color="secondary">border-bottom-right-radius </text-chunk>
        <text-chunk color="text">attributes can be used for setting the border radius of individual corners.</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Syntax"}}

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Border radius: top-left 10, top-right 10, bottom-left 10, bottom-right 10."}}
        {{$codeBlock :=
`<rectangle position="relative" width="250" height="100" border-width="5"
    border-color="#40c4ff" border-radius="10"></rectangle>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Border radius: top-left 25, top-right 5, bottom-left 5, bottom-right 25."}}
        {{$codeBlock :=
`<rectangle position="relative" width="250" height="100" border-width="5"
    border-color="#ffb300" border-radius="25 5"></rectangle>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Border radius: top-left 10, top-right 0, bottom-left 0, bottom-right 25."}}
        {{$codeBlock :=
`<rectangle position="relative" width="250" height="100" border-width="5"
    border-color="#e53935" border-radius="10 0 25"></rectangle>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Border radius: top-left 20, top-right 5, bottom-left 30, bottom-right 0."}}
        {{$codeBlock :=
`<rectangle position="relative" width="250" height="100" border-width="5"
    border-color="#388e3c" border-radius="20 5 30 0"></rectangle>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>
</chapter>
