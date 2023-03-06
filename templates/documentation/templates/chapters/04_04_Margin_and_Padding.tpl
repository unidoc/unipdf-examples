<chapter margin="10 0 0 0">
    {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Margin and Padding"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk color="secondary">margin </text-chunk>
        <text-chunk color="text">and </text-chunk>
        <text-chunk color="secondary">padding </text-chunk>
        <text-chunk color="text">attributes are used generate configurable amounts of space for a component. The difference between them is that </text-chunk>
        <text-chunk font="helvetica-bold" color="text">margin </text-chunk>
        <text-chunk color="text">generates space around the exterior of components, whereas </text-chunk>
        <text-chunk font="helvetica-bold" color="text">padding </text-chunk>
        <text-chunk color="text">generates space inside components, around their content. Almost all components have a </text-chunk>
        <text-chunk color="secondary">margin </text-chunk>
        <text-chunk color="text">attribute whereas  </text-chunk>
        <text-chunk color="secondary">padding </text-chunk>
        <text-chunk color="text">is currently available only for </text-chunk>
        <text-chunk font="helvetica-bold" color="text">divisions</text-chunk>
        <text-chunk color="text">. The syntax of both attributes is identical.</text-chunk>
    </paragraph>

    <division>
        {{template "paragraph" dict "Margin" "10 0 10 0" "Font" "helvetica-bold" "Text" "Margin vs Padding"}}
        {{$codeBlocks := array
`<division margin="20">
    <background border-size="5" border-color="#ffc400"></background>

    <image src="path('templates/res/images/sample-image.jpg')"
        margin="2" fit-mode="fill-width"></image>
</division>`
`<division padding="20">
    <background border-size="5" border-color="#ffc400"></background>

    <image margin="2" src="path('templates/res/images/sample-image.jpg')"
        margin="2" fit-mode="fill-width"></image>
</division>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Syntax"}}

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Margin: top 10, right 10, bottom 10, left 10."}}
        {{$codeBlock :=
`<image src="path('templates/res/images/sample-image-5.jpg')"
    fit-mode="fill-width" margin="10"></image>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Margin: top 5, right 20, bottom 5, left 20."}}
        {{$codeBlock :=
`<image src="path('templates/res/images/sample-image-5.jpg')"
    fit-mode="fill-width" margin="5 20"></image>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Margin: top 0, right 20, bottom 5, left 20."}}
        {{$codeBlock :=
`<image src="path('templates/res/images/sample-image-5.jpg')"
    fit-mode="fill-width" margin="0 20 5"></image>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "10 0 5 0" "Text" "Margin: top 10, right 5, bottom 20, left 15."}}
        {{$codeBlock :=
`<image src="path('templates/res/images/sample-image-5.jpg')"
    fit-mode="fill-width" margin="10 5 20 15"></image>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
    </division>
</chapter>
