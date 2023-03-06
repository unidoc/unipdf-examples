<chapter margin="10 0 0 0">
    {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Images"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">There are multiple ways of loading and using </text-chunk>
        <text-chunk font="helvetica-bold" color="text">images </text-chunk>
        <text-chunk color="text">in templates. This section presents all of them, with code examples. Source images are specified in the </text-chunk>
        <text-chunk color="secondary">src </text-chunk>
        <text-chunk color="text">attribute of the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">image </text-chunk>
        <text-chunk color="text">component.</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Margin" "5 0" "Font" "helvetica-bold" "Text" "Loading images from paths"}}
    <paragraph margin="0 0 10 0" line-height="1.2" enable-word-wrap="true">
        <text-chunk color="text">Images can be loaded at runtime, directly from templates, by specifying their paths. The provided image paths can be absolute or relative to the application binary.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-3.jpg')"
    width="215" height="150" align="center"></image>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "5 0" "Text" "Using preloaded images"}}
    <paragraph line-height="1.2" enable-word-wrap="true">
        <text-chunk color="text">Images can be made available to a template by loading and adding them to the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">image map </text-chunk>
        <text-chunk color="text">of the options used to render that template. The images are accessed in the templates using the names assigned to them in the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">image map</text-chunk>
        <text-chunk color="text">. Preloading images is useful when they are used more than once throughout templates, in order to avoid loading them multiple times.</text-chunk>
    </paragraph>

    <paragraph margin="10 0">
        <text-chunk color="text">Please see chapter </text-chunk>
        <text-chunk link="page(56)">5. Golang Interface </text-chunk>
        <text-chunk color="text">for examples of how to use the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">image map</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<image src="sample"
    width="215" height="150" align="center"></image>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
        {{end}}
    </division>
</chapter>
