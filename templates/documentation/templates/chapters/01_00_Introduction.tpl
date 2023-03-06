<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Introduction"}}

    <paragraph margin="5 0" line-height="1.1">
        <text-chunk color="text" font="helvetica-bold">Templates </text-chunk>
        <text-chunk color="text">provide a new way of building PDF files using the components available in the </text-chunk>
        <text-chunk color="text" font="helvetica-bold">creator </text-chunk>
        <text-chunk color="text">package. The templates are defined using an XML based markup language and they are parsed and translated into creator components at runtime by the internal template processor.</text-chunk>
        <text-chunk>{{.newline}}{{.newline}}</text-chunk>
        <text-chunk>The templates have two processing phases. First, they are executed as </text-chunk>
        <text-chunk link="url('https://pkg.go.dev/text/template#Template')">text/template#Template </text-chunk>
        <text-chunk>instances, which allows actions to be executed and data to be injected. The second processing phase takes the output of the first and parses it into components which are then rendered.</text-chunk>
        <text-chunk>{{.newline}}{{.newline}}</text-chunk>
        <text-chunk color="text">Components are split into </text-chunk>
        <text-chunk font="helvetica-bold" color="text">data </text-chunk>
        <text-chunk>and </text-chunk>
        <text-chunk font="helvetica-bold" color="text">container </text-chunk>
        <text-chunk>components. Data components produce some form of content like text, images, geometric shapes. Container components hold or group other components (e.g. tables, lists).</text-chunk>
        <text-chunk>{{.newline}}{{.newline}}</text-chunk>
        <text-chunk>The aim of templates is to make creating complex layouts a lot easier, faster, and by writing less code. In fact, this document was written entirely using templates. This documentation and other examples of templates in action, including their source code, can be found in our </text-chunk>
        <text-chunk link="url('https://github.com/unidoc/unipdf-examples/tree/master/templates')">examples repository</text-chunk>
        <text-chunk>.</text-chunk>
    </paragraph>

    <paragraph margin="5 0" line-height="1.1">
        <text-chunk color="text">Let's dive right into it. Here's a simple example of a </text-chunk>
        <text-chunk font="helvetica-bold" color="text">division </text-chunk>
        <text-chunk color="text">containing a </text-chunk>
        <text-chunk font="helvetica-bold" color="text">paragraph</text-chunk>
        <text-chunk color="text">, a </text-chunk>
        <text-chunk font="helvetica-bold" color="text">line </text-chunk>
        <text-chunk color="text">and an </text-chunk>
        <text-chunk font="helvetica-bold" color="text">image</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>

    {{$codeBlock :=
`<division padding="15 10" margin="3">
    <background border-color="#333333" border-size="0.5"
        fill-color="#f9fafe"></background>

    <paragraph>
        <text-chunk>Far far away, behind the word mountains, </text-chunk>
        <text-chunk>far from the countries Vokalia and </text-chunk>
        <text-chunk>Consonantia, there live the blind texts. </text-chunk>
        <text-chunk>Separated they live in Bookmarksgrove </text-chunk>
        <text-chunk>right at the coast of the Semantics, </text-chunk>
        <text-chunk>a large language ocean.</text-chunk>
    </paragraph>

    <line position="relative" fit-mode="fill-width"
        thickness="0.5" color="#333333" margin="5 0"></line>

    <image fit-mode="fill-width"
        src="path('templates/res/images/sample-image.jpg')"></image>
</division>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}
</chapter>
