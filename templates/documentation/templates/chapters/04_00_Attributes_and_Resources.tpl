<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Attributes and Resources"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">This chapter includes information regarding methods of loading resources (</text-chunk>
        <text-chunk font="helvetica-bold" color="text">fonts</text-chunk>
        <text-chunk color="text">, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">images</text-chunk>
        <text-chunk color="text">, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">colors</text-chunk>
        <text-chunk color="text">) and the usage of attributes which have shorthand forms (</text-chunk>
        <text-chunk font="helvetica-bold" color="text">margin</text-chunk>
        <text-chunk color="text">, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">border-radius</text-chunk>
        <text-chunk color="text">).</text-chunk>
    </paragraph>

    <!-- Fonts -->
    {{template "04_01_Fonts" .}}

    <!-- Images -->
    {{template "04_02_Images" .}}

    <!-- Colors -->
    {{template "04_03_Colors" .}}

    <!-- Margin and Padding -->
    {{template "04_04_Margin_and_Padding" .}}

    <!-- Border Radius -->
    {{template "04_05_Border_Radius" .}}
</chapter>
