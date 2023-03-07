<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Golang Interface"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk font="helvetica-bold" color="text">Templates </text-chunk>
        <text-chunk color="text">can be rendered using the </text-chunk>
        <text-chunk color="secondary">DrawTemplate </text-chunk>
        <text-chunk color="text">method of the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">creator </text-chunk>
        <text-chunk color="text">package, which invokes the internal template processor. The processing of templates is done in two phases. First, templates are executed as </text-chunk>
        <text-chunk link="url('https://pkg.go.dev/text/template#Template')">text/template#Template </text-chunk>
        <text-chunk color="text">instances. This allows injecting data and executing actions. The second processing phase takes the output of the first one, parses and translates it into </text-chunk>
        <text-chunk font="helvetica-bold" color="text">components</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{$codeBlock :=
`func (c *Creator) DrawTemplate(r io.Reader, data interface{}, options *TemplateOptions) error`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "CodeLabel" "GO" "Code" $codeBlock}}
    </division>

    <paragraph margin="0 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The method renders the template specified through the specified reader </text-chunk>
        <text-chunk color="secondary">r</text-chunk>
        <text-chunk color="text">, using the provided </text-chunk>
        <text-chunk color="secondary">options</text-chunk>
        <text-chunk color="text">. The passed in </text-chunk>
        <text-chunk color="secondary">data </text-chunk>
        <text-chunk color="text">is available in the rendered template.</text-chunk>
    </paragraph>

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">Resources like </text-chunk>
        <text-chunk font="helvetica-bold" color="text">fonts</text-chunk>
        <text-chunk color="text">, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">images</text-chunk>
        <text-chunk color="text">, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">colors</text-chunk>
        <text-chunk color="text"> and </text-chunk>
        <text-chunk font="helvetica-bold" color="text">charts </text-chunk>
        <text-chunk color="text">can be shared with the templates through the resource maps in the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">template options</text-chunk>
        <text-chunk color="text">. In addition, </text-chunk>
        <text-chunk font="helvetica-bold" color="text">Go functions </text-chunk>
        <text-chunk color="text">can be made available to templates through the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">helper function map</text-chunk>
        <text-chunk color="text">. The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">subtemplates map </text-chunk>
        <text-chunk color="text">is used to provide additional templates which are automatically available to be accessed both from within the main template and from one another. All the resources added to the maps in the template options can be accessed in templates by their assigned name (map key).</text-chunk>
    </paragraph>

    <division enable-page-wrap="false">
        {{$codeBlock :=
`type TemplateOptions struct {
    FontMap  map[string]*model.PdfFont
    ImageMap map[string]*model.Image
    ColorMap map[string]Color
    ChartMap map[string]render.ChartRenderable

    HelperFuncMap  template.FuncMap
    SubtemplateMap map[string]io.Reader
}`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "CodeLabel" "GO" "Code" $codeBlock}}
    </division>

    <!-- Font map -->
    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Font map"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">font map </text-chunk>
        <text-chunk color="text">is used to pass fonts to templates. The fonts are accessed in templates using their assigned map key.</text-chunk>
    </paragraph>

    <division>
        {{$codeBlock :=
`<paragraph>
    <text-chunk color="#4dd0e1" font-size="12"
        font="deja-vu-sans-mono">Sample text</text-chunk>
</paragraph>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

        {{$codeBlock :=
`// Load DejaVuSansMono font.
dejaVuSansMono, err := model.NewPdfFontFromTTFFile("templates/res/fonts/DejaVuSansMono.ttf")
if err != nil {
    log.Fatal(err)
}

// Create template options.
tplOpts := &creator.TemplateOptions{
    FontMap: map[string]*model.PdfFont{
        "deja-vu-sans-mono": dejaVuSansMono,
    },
}

// Draw template.
c := creator.New()
if err := c.DrawTemplate(tpl, nil, tplOpts); err != nil {
    log.Fatal(err)
}
`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "CodeLabel" "GO" "Code" $codeBlock}}
    </division>

    <!-- Image map -->
    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Image map"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">image map </text-chunk>
        <text-chunk color="text">is used to pass images to templates. The images are accessed in templates using their assigned map key.</text-chunk>
    </paragraph>

    <division>
        {{$codeBlock :=
`<image src="sample"
    width="128" height="72" align="center"></image>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

        {{$codeBlock :=
`// Load sample image.
f, err := os.Open("templates/res/images/sample-image-2.jpg")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

img, err := model.ImageHandling.Read(f)
if err != nil {
    log.Fatal(err)
}

// Create template options.
tplOpts := &creator.TemplateOptions{
    ImageMap: map[string]*model.Image{
        "sample": img,
    },
}

// Draw template.
c := creator.New()
if err := c.DrawTemplate(tpl, nil, tplOpts); err != nil {
    log.Fatal(err)
}`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "CodeLabel" "GO" "Code" $codeBlock}}
    </division>

    <!-- Color map -->
    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Color map"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">color map </text-chunk>
        <text-chunk color="text">is used to pass colors to templates. The colors are accessed in templates using their assigned map key.</text-chunk>
    </paragraph>

    <division>
        {{$codeBlock :=
`<rectangle width="6" height="1.5" position="relative"
    fit-mode="fill-width" border-width="3"
    border-color="secondary" fill-color="primary"></rectangle>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

        {{$codeBlock :=
`// Create template options.
tplOpts := &creator.TemplateOptions{
    ColorMap: map[string]creator.Color{
        "primary":   creator.ColorRGBFromHex("#0772cd"),
        "secondary": creator.ColorRGBFromHex("#f00c27"),
    },
}

// Draw template.
c := creator.New()
if err := c.DrawTemplate(tpl, nil, tplOpts); err != nil {
    log.Fatal(err)
}`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "CodeLabel" "GO" "Code" $codeBlock}}
    </division>

    <!-- Chart map -->
    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Chart map"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">chart map </text-chunk>
        <text-chunk color="text">is used to pass charts to templates. The charts are accessed in templates using their assigned map key.</text-chunk>
    </paragraph>

    <division>
        {{createLineChart "sample-chart" 50 0.0 100.0}}
        {{$codeBlock :=
`<chart src="sample-chart" height="150"></chart>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

        {{$codeBlock :=
`// Create sample chart.
mainSeries := dataset.ContinuousSeries{
    XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(1.0).WithEnd(float64(50)).WithStep(1)}.Values(),
    YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(50).WithMin(0).WithMax(100)}.Values(),
    Style: render.Style{
        FillColor: color.RGBA{R: 38, G: 198, B: 218},
    },
}

linRegSeries := &dataset.LinearRegressionSeries{
    InnerSeries: mainSeries,
    Style: render.Style{
        StrokeColor: color.RGBA{R: 255, G: 167, B: 38},
    },
}

sampleChart := &unichart.Chart{
    Series: []dataset.Series{
        mainSeries,
        linRegSeries,
    },
}

// Create template options.
tplOpts := &creator.TemplateOptions{
    ChartMap: map[string]render.ChartRenderable{
        "sample-chart": sampleChart,
    },
}

// Draw template.
c := creator.New()
if err := c.DrawTemplate(tpl, nil, tplOpts); err != nil {
    log.Fatal(err)
}`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "CodeLabel" "GO" "Code" $codeBlock}}
    </division>

    <!-- Helper function map -->
    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Helper function map"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">helper function map </text-chunk>
        <text-chunk color="text">is used to pass functions to templates. The functions are accessed in templates using their assigned map key.</text-chunk>
    </paragraph>

    <division>
        {{$codeBlock :=
`<paragraph>
    <text-chunk font="helvetica-bold" font-size="12"
        color="#1b5e20">%s</text-chunk>
</paragraph>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" (printf $codeBlock `{{strToUpper "sample text"}}`) "Result" (printf $codeBlock (strToUpper "sample text"))}}

        {{$codeBlock :=
`<paragraph>
    <text-chunk font="helvetica-bold" font-size="12"
        color="#c62828">%s</text-chunk>
</paragraph>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" (printf $codeBlock `{{now.Format "Jan 02 2006"}}`) "Result" (printf $codeBlock (now.Format "Jan 02 2006"))}}

        {{$codeBlock :=
`// Create template options.
tplOpts := &creator.TemplateOptions{
    HelperFuncMap: template.FuncMap{
        "now": time.Now,
        "strToUpper": strings.ToUpper,
    },
}

// Draw template.
c := creator.New()
if err := c.DrawTemplate(tpl, nil, tplOpts); err != nil {
    log.Fatal(err)
}`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "HideResult" true "CodeLabel" "GO" "Code" $codeBlock}}
    </division>

    <paragraph margin="5 0" line-height="1.1">
        <text-chunk color="text">For more information, please see </text-chunk>
        <text-chunk link="url('https://pkg.go.dev/text/template#FuncMap')">text/template#FuncMap</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>

    <!-- Subtemplate map -->
    {{template "paragraph" dict "Margin" "10 0 0 0" "Font" "helvetica-bold" "Text" "Subtemplate map"}}

    <paragraph margin="5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">subtemplate map </text-chunk>
        <text-chunk color="text">is used to include additional templates to the processing pipeline. The added </text-chunk>
        <text-chunk font="helvetica-bold" color="text">subtemplates </text-chunk>
        <text-chunk color="text">can be accessed both by the main template and by the other subtemplates using their assigned map key. Subtemplates defined inside the provided subtemplates are accessible without including their parent templates. All the resources available to the main template are also available to the subtemplates.{{.newline}}{{.newline}}Subtemplates are useful both for splitting long content into more manageable pieces and for defining reusable chunks of content. For example, this documentation is split into chapters using subtemplates.</text-chunk>
    </paragraph>

    <paragraph margin="5 0" line-height="1.1">
        <text-chunk color="text">For more information, please see the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">text/template </text-chunk>
        <text-chunk color="text">package documentation regarding </text-chunk>
        <text-chunk link="url('https://pkg.go.dev/text/template#Associated_templates')">associated and nested templates</text-chunk>
        <text-chunk color="text">.</text-chunk>
    </paragraph>
</chapter>
