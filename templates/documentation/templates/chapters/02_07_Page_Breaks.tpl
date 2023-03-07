<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Page Breaks"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">page break </text-chunk>
        <text-chunk color="text">component has a single purpose and that is to end the current page and start a new one. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<page-break>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of page breaks:"}}

    {{$codeBlock :=
`<page-break></page-break>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "Code" $codeBlock}}
</chapter>
