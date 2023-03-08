<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Data Components"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text" font="helvetica-bold">Data </text-chunk>
        <text-chunk color="text">components are used to store and render some form of content. Most creator </text-chunk>
        <text-chunk color="text" font="helvetica-bold">data </text-chunk>
        <text-chunk color="text">components can be used in templates, with support for more to come in the near future.{{.newline}}{{.newline}}</text-chunk>
        <text-chunk color="text">Here is the list of currently supported components:</text-chunk>
    </paragraph>

    {{$components := array
             (dict "name" "paragraph" "children" (array (dict "name" "text-chunk")))
             (dict "name" "image" "children" (array))
             (dict "name" "rectangle" "children" (array))
             (dict "name" "ellipse" "children" (array))
             (dict "name" "line" "children" (array))
             (dict "name" "chart" "children" (array))
             (dict "name" "page-break" "children" (array))
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

    <!-- Paragraphs -->
    <page-break></page-break>
    {{template "02_01_Paragraphs" .}}

    <!-- Images -->
    <page-break></page-break>
    {{template "02_02_Images" .}}

    <!-- Rectangles -->
    <page-break></page-break>
    {{template "02_03_Rectangles" .}}

    <!-- Ellipses -->
    <page-break></page-break>
    {{template "02_04_Ellipses" .}}

    <!-- Lines -->
    <page-break></page-break>
    {{template "02_05_Lines" .}}

    <!-- Charts -->
    <page-break></page-break>
    {{template "02_06_Charts" .}}

    <!-- Page Breaks -->
    {{template "02_07_Page_Breaks" .}}
</chapter>
