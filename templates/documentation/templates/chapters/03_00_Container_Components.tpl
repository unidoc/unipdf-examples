<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Container Components"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text" font="helvetica-bold">Container </text-chunk>
        <text-chunk color="text">components are used to group other components, usually with the goal of presenting them in a structured way. Most creator </text-chunk>
        <text-chunk color="text" font="helvetica-bold">container </text-chunk>
        <text-chunk color="text">components can be used in templates, with support for more to come in the near future.{{.newline}}{{.newline}}</text-chunk>
        <text-chunk color="text">Here is the list of currently supported components:</text-chunk>
    </paragraph>

    {{$components := array
             (dict "name" "division" "children" (array (dict "name" "background")))
             (dict "name" "table" "children" (array (dict "name" "table-cell")))
             (dict "name" "list" "children" (array (dict "name" "list-item") (dict "name" "list-marker")))
             (dict "name" "chapter" "children" (array (dict "name" "chapter-heading")))
    }}

    {{range $i, $component := $components}}
        <paragraph>
            <text-chunk color="text" text-rise="-1">• </text-chunk>
            <text-chunk color="text">{{$component.name}}</text-chunk>
        </paragraph>
        {{range $i, $child := $component.children}}
            <paragraph>
                <text-chunk color="text" text-rise="-1">  • </text-chunk>
                <text-chunk color="text">{{$child.name}}</text-chunk>
            </paragraph>
        {{end}}
        <line position="relative" fit-mode="fill-width" color="medium-gray" thickness="0.5" margin="4 0 0 0"></line>
    {{end}}

    <!-- Divisions -->
    {{template "03_01_Divisions" .}}

    <!-- Tables -->
    <page-break></page-break>
    {{template "03_02_Tables" .}}

    <!-- Lists -->
    <page-break></page-break>
    {{template "03_03_Lists" .}}

    <!-- Chapters -->
    <page-break></page-break>
    {{template "03_04_Chapters" .}}
</chapter>
