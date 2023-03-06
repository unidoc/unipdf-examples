<chapter margin="10 0 0 0">
    {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Fonts"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">There are multiple ways of loading and using </text-chunk>
        <text-chunk font="helvetica-bold" color="text">fonts </text-chunk>
        <text-chunk color="text">in templates. This section presents all of them, with code examples. In most cases, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">fonts </text-chunk>
        <text-chunk color="text">are specified using the </text-chunk>
        <text-chunk color="secondary">font </text-chunk>
        <text-chunk color="text">attribute.</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Margin" "10 0 5 0" "Font" "helvetica-bold" "Text" "Using standard fonts"}}
    <paragraph margin="0 0 5 0" line-height="1.2" enable-word-wrap="true">
        <text-chunk color="text">A list of </text-chunk>
        <text-chunk font="helvetica-bold" color="text">standard fonts </text-chunk>
        <text-chunk color="text">is provided by the library. These fonts are always available and can be specified by their names. If no font specified, components default to using </text-chunk>
        <text-chunk font="helvetica-bold" color="text">helvetica</text-chunk>
        <text-chunk color="text">.{{.newline}}{{.newline}}Here is the list of available standard fonts:</text-chunk>
    </paragraph>

    <list indent="15" margin="0 0 10 0">
        <list-item>{{template "paragraph" dict "Font" "helvetica" "Text" "helvetica"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "helvetica-bold" "Text" "helvetica-bold"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "helvetica-oblique" "Text" "helvetica-oblique"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "helvetica-bold-oblique" "Text" "helvetica-bold-oblique"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "courier" "Text" "courier"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "courier-bold" "Text" "courier-bold"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "courier-oblique" "Text" "courier-oblique"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "courier-bold-oblique" "Text" "courier-bold-oblique"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "times" "Text" "times"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "times-bold" "Text" "times-bold"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "times-italic" "Text" "times-italic"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "times-bold-italic" "Text" "times-bold-italic"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "helvetica" "Text" "symbol"}}</list-item>
        <list-item>{{template "paragraph" dict "Font" "helvetica" "Text" "zapf-dingbats"}}</list-item>
    </list>

    <division>
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<paragraph>
    <text-chunk color="#388e3c" font="times">Times</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk color="#4dd0e1"
        font="helvetica-bold">Helvetica Bold</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk color="#ff5722"
        font="courier-oblique">Courier Oblique</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk color="#880e4f">Zapf Dingbats: </text-chunk>
    <text-chunk color="#880e4f" font="zapf-dingbats">☞ ⑩✩</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk color="#ffc400">Symbol: </text-chunk>
    <text-chunk color="#ffc400" font="symbol">α ≥ π ∗ β</text-chunk>
</paragraph>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Margin" "5 0" "Font" "helvetica-bold" "Text" "Loading fonts from paths"}}
    <paragraph margin="0 0 10 0" line-height="1.2" enable-word-wrap="true">
        <text-chunk color="text">Fonts can be loaded at runtime, directly from templates, by specifying their paths. Optionally, there is an option to specify the type of font (regular, composite). The provided font paths can be absolute or relative to the application binary.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<paragraph>
    <text-chunk font="path('templates/res/fonts/DejaVuSansMono.ttf')"
        color="#00838f">DejaVu Sans Mono</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk font="path('templates/res/fonts/Roboto-Regular.ttf') type('composite')"
        color="#2e7d32">Roboto Regular</text-chunk>
</paragraph>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Columns" 1 "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "5 0" "Text" "Using preloaded fonts"}}
    <paragraph line-height="1.2" enable-word-wrap="true">
        <text-chunk color="text">Fonts can be made available to a template by loading and adding them to the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">font map </text-chunk>
        <text-chunk color="text">of the options used to render that template. The fonts are accessed in the templates using the names assigned to them in the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">font map</text-chunk>
        <text-chunk color="text">. Preloading fonts is useful when they are used more than once throughout templates, in order to avoid loading them multiple times.</text-chunk>
    </paragraph>

    <paragraph margin="10 0">
        <text-chunk color="text">Please see chapter </text-chunk>
        <text-chunk link="page(56)">5. Golang Interface </text-chunk>
        <text-chunk color="text">for examples of how to use the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">font map</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<paragraph>
    <text-chunk font="deja-vu-sans-mono"
        color="#dd2c00">deja-vu-sans-mono was preloaded for this example.</text-chunk>
</paragraph>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Columns" 1 "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "5 0" "Text" "Using font lists"}}
    <paragraph margin="0 0 10 0" line-height="1.2" enable-word-wrap="true">
        <text-chunk font="helvetica-bold" color="text">Font lists </text-chunk>
        <text-chunk color="text">are used in order to provide acceptable fallback fonts, in case some of the requested fonts may not be available for any reason. A list of fonts can be specified by using any of the previously described methods of loading/referencing fonts. It is basically an enumaration of font preferences, separated by </text-chunk>
        <text-chunk color="secondary">commas</text-chunk>
        <text-chunk color="text">. Each of the specified fonts in the list is attempted from left to right and the first available font is selected. If none of fonts in the list is valid, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">helvetica </text-chunk>
        <text-chunk>is used a fallback.</text-chunk>
    </paragraph>

    <division>
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<paragraph>
    <text-chunk font="path('templates/res/fonts/DejaVuSansMono.ttf'), courier"
        color="#2962ff">Font order: DejaVuSansMono (load from path) > courier (standard).</text-chunk>
</paragraph>`
`<paragraph>
    <text-chunk font="path('templates/res/fonts/Roboto-Regular.ttf') type('composite'), deja-vu-sans-mono, times"
        color="#558b2f">Font order: Roboto-Regular (load from path, composite) > deja-vu-sans-mono (preloaded) > times (standard).</text-chunk>
</paragraph>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Columns" 1 "Code" $codeBlock}}
        {{end}}
    </division>
</chapter>
