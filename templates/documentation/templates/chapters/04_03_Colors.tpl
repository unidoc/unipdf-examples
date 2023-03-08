<chapter margin="10 0 0 0">
    {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "Colors"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">In this section, we will cover all the ways of defining and using </text-chunk>
        <text-chunk font="helvetica-bold" color="text">colors </text-chunk>
        <text-chunk color="text">in templates. Color usage is identical with all relevant attributes (e.g. background, colors, border colors) unless noted otherwise.</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Margin" "5 0" "Font" "helvetica-bold" "Text" "Defining colors using hexadecimal notation"}}
    <paragraph margin="0 0 10 0" line-height="1.2" enable-word-wrap="true">
        <text-chunk color="text">Colors can be defined using hex values, directly from templates. Both the regular </text-chunk>
        <text-chunk font="helvetica-bold" color="text">six-digit notation </text-chunk>
        <text-chunk color="text">(e.g. </text-chunk>
        <text-chunk color="secondary">#00aaff</text-chunk>
        <text-chunk color="text">) and its shorter </text-chunk>
        <text-chunk font="helvetica-bold" color="text">three-digit notation </text-chunk>
        <text-chunk color="text">(e.g. </text-chunk>
        <text-chunk color="secondary">#0af</text-chunk>
        <text-chunk color="text">) are supported. Hex color values are preceded by </text-chunk>
        <text-chunk color="secondary"># </text-chunk>
        <text-chunk color="text">and they are case insensitive.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<rectangle width="6" height="1" position="relative" fit-mode="fill-width"
    border-width="0" fill-color="#ffc107"></rectangle>`
`<rectangle width="6" height="1" position="relative" fit-mode="fill-width"
    border-width="3" border-color="#0bf"></rectangle>`
`<rectangle width="6" height="1" position="relative" fit-mode="fill-width"
    border-width="3" border-color="#a0f" fill-color="#EA80FC"></rectangle>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Margin" "5 0" "Font" "helvetica-bold" "Text" "Defining linear gradient patterns"}}

    <paragraph margin="0 0 10 0" line-height="1.2" enable-word-wrap="true">
        <text-chunk font="helvetica-bold" color="text">Linear gradients </text-chunk>
        <text-chunk color="text">transition colors progressively along an imaginary line. They are defined by two or more </text-chunk>
        <text-chunk font="helvetica-bold" color="text">color stops </text-chunk>
        <text-chunk color="text">which consist of a color and an optional starting point for that color. The starting point is defined as a percentage (0-100) of the available space. Optionally, an </text-chunk>
        <text-chunk font="helvetica-bold" color="text">angle </text-chunk>
        <text-chunk color="text">specified in degrees (0-360) can be specified, which rotates the entire pattern.{{.newline}}{{.newline}}</text-chunk>
        <text-chunk font="helvetica-bold" color="text">Note</text-chunk>
        <text-chunk color="text">: currently, gradient patterns can only be applied to background colors of components. Border or outline gradient pattern colors are not supported yet.</text-chunk>
    </paragraph>

    <division>
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<division>
    <background fill-color="linear-gradient(#227cc3, #78ac77, #f1b025)"></background>

    <paragraph text-align="center" margin="45 0">
        <text-chunk color="#fff" font="helvetica-bold">No stops defined, colors are evenly distributed.</text-chunk>
    </paragraph>
</division>`
`<division>
    <background fill-color="linear-gradient(#227cc3, #78ac77 20%, #f1b025)"></background>

    <paragraph text-align="center" margin="45 0">
        <text-chunk color="#fff" font="helvetica-bold">First color stops at 20%, the rest are evenly distributed.</text-chunk>
    </paragraph>
</division>`
`<division>
    <background fill-color="linear-gradient(#227cc3, #78ac77 20%, #f1b025 35%)"></background>

    <paragraph text-align="center" margin="45 0">
        <text-chunk color="#fff" font="helvetica-bold">First color stops at 20%, the second at 35%, the third takes the rest of the space.</text-chunk>
    </paragraph>
</division>`
`<division>
    <background fill-color="linear-gradient(180deg, #227cc3, #78ac77 40%, #f1b025)"></background>

    <paragraph text-align="center" margin="45 0">
        <text-chunk color="#fff" font="helvetica-bold">Optionally, an angle specified in degrees can be provided.</text-chunk>
    </paragraph>
</division>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Columns" 1 "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Margin" "5 0" "Font" "helvetica-bold" "Text" "Defining radial gradient patterns"}}

    <paragraph margin="0 0 10 0" line-height="1.2" enable-word-wrap="true">
        <text-chunk font="helvetica-bold" color="text">Radial gradients </text-chunk>
        <text-chunk color="text"> transition colors progressively from the center of the assigned space. They are defined by two or more </text-chunk>
        <text-chunk font="helvetica-bold" color="text">color stops </text-chunk>
        <text-chunk color="text">which consist of a color and an optional starting point for that color. The starting point is defined as a percentage (0-100) of the available space.{{.newline}}{{.newline}}</text-chunk>
        <text-chunk font="helvetica-bold" color="text">Note</text-chunk>
        <text-chunk color="text">: currently, gradient patterns can only be applied to background colors of components. Border or outline gradient pattern colors are not supported yet.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<division>
    <background fill-color="radial-gradient(#eeaeca, #c3b4d9, #94bbe9)"></background>

    <paragraph text-align="center" margin="45 0">
        <text-chunk font="helvetica-bold">No stops defined, colors are evenly distributed.</text-chunk>
    </paragraph>
</division>`
`<division>
    <background fill-color="radial-gradient(#eeaeca, #c3b4d9 70%, #94bbe9)"></background>

    <paragraph text-align="center" margin="45 0">
        <text-chunk font="helvetica-bold">First color stops at 70%, the rest are evenly distributed.</text-chunk>
    </paragraph>
</division>`
`<division>
    <background fill-color="radial-gradient(#eeaeca, #c3b4d9 20%, #94bbe9 60%)"></background>

    <paragraph text-align="center" margin="45 0">
        <text-chunk font="helvetica-bold">First color stops at 20%, the second at 60%, the third takes the rest of the space.</text-chunk>
    </paragraph>
</division>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Columns" 1 "Code" $codeBlock}}
        {{end}}
    </division>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "5 0" "Text" "Using preloaded colors"}}
    <paragraph line-height="1.2" enable-word-wrap="true">
        <text-chunk color="text">Colors can be made available to a template by creating and adding them to the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">color map </text-chunk>
        <text-chunk color="text">of the options used to render that template. The colors are accessed in the templates using the names assigned to them in the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">color map</text-chunk>
        <text-chunk color="text">. Preloading colors is useful when they are used more than once throughout templates, in order to avoid defining them multiple times.</text-chunk>
    </paragraph>

    <paragraph margin="10 0">
        <text-chunk color="text">Please see chapter </text-chunk>
        <text-chunk link="page(56)">5. Golang Interface </text-chunk>
        <text-chunk color="text">for examples of how to use the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">color map</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{template "paragraph" dict "Margin" "0 0 10 0" "Text" "Usage"}}
        {{$codeBlocks := array
`<rectangle width="6" height="1" position="relative" fit-mode="fill-width"
    border-width="0" fill-color="secondary"></rectangle>`
`<rectangle width="6" height="1" position="relative" fit-mode="fill-width"
    border-width="0" fill-color="primary-bg-gradient"></rectangle>`
        }}

        {{range $i, $codeBlock := $codeBlocks}}
            {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
        {{end}}
    </division>
</chapter>
