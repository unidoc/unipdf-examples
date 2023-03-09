<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Lists"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">The </text-chunk>
        <text-chunk font="helvetica-bold" color="text">list </text-chunk>
        <text-chunk color="text">component is used to create a list of items and is similar to a </text-chunk>
        <text-chunk font="helvetica-bold" color="text">HTML list</text-chunk>
        <text-chunk color="text">. Currently, the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">paragraph </text-chunk>
        <text-chunk color="text">and </text-chunk>
        <text-chunk font="helvetica-bold" color="text">list </text-chunk>
        <text-chunk color="text">components are supported as list items. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<list>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of lists:"}}

    {{$codeBlock :=
`<list>
    <list-item>
        <paragraph><text-chunk>Coffee</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>Tea</text-chunk></paragraph>
    </list-item>
    <list-item>
        <list-marker> </list-marker>
        <list>
            <list-marker>» </list-marker>
            <list-item>
                <paragraph><text-chunk>Rooibos</text-chunk></paragraph>
            </list-item>
            <list-item>
                <paragraph><text-chunk>Oolong</text-chunk></paragraph>
            </list-item>
        </list>
    </list-item>
    <list-item>
        <paragraph><text-chunk>Milk</text-chunk></paragraph>
    </list-item>
</list>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "15 0 10 0" "Text" "Supported attributes:"}}

    <!-- margin -->
    {{$codeBlocks := array
`<list margin="5 0 5 20">
    <list-item>
        <paragraph><text-chunk>Leonardo</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>Donatello</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>Raphael</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>Michelangelo</text-chunk></paragraph>
    </list-item>
</list>`
    }}
    {{template "attr-showcase" dict "AttrName" "margin"
        "AttrDescr" "The `margin` attribute allows setting a configurable amount of space around the component."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- indent -->
    {{$codeBlocks := array
`<list indent="5">
    <list-item>
        <paragraph><text-chunk>Blue shades</text-chunk></paragraph>
    </list-item>
    <list-item>
        <list-marker> </list-marker>
        <list indent="10">
            <list-item>
                <paragraph>
                    <text-chunk color="#48aaad">Teal</text-chunk>
                </paragraph>
            </list-item>
            <list-item>
                <paragraph>
                    <text-chunk color="#241571">Berry</text-chunk>
                </paragraph>
            </list-item>
        </list>
    </list-item>
</list>`
    }}
    {{template "attr-showcase" dict "AttrName" "indent"
        "AttrDescr" "The `indent` attribute represents an horizontal offset applied on the left side of the list items. By default, for list items of type list, the `indent` attribute is 15. Otherwise, it defaults to 0."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <chapter>
        {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "List Items" "AlternateText" "List » Items"}}

        <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
            <text-chunk font="helvetica-bold" color="text">List items </text-chunk>
            <text-chunk color="text">are the main building blocks of the </text-chunk>
            <text-chunk font="helvetica-bold" color="text">list </text-chunk>
            <text-chunk color="text">component. Currently, list items can be either </text-chunk>
            <text-chunk font="helvetica-bold" color="text">paragraphs </text-chunk>
            <text-chunk color="text">or </text-chunk>
            <text-chunk font="helvetica-bold" color="text">lists </text-chunk>
            <text-chunk color="text">. The component is represented in templates using the </text-chunk>
            <text-chunk color="secondary">{{xmlEscape "<list-item>"}} </text-chunk>
            <text-chunk color="text">tag.</text-chunk>
        </paragraph>

        {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "10 0 0 0" "Text" "Supported attributes:"}}
        {{template "paragraph" dict "Margin" "5 0 15 0" "Text" "List items have no attributes."}}
    </chapter>

    <chapter>
        {{template "chapter-title" dict "FillColor" "primary-light-bg-gradient" "Text" "List Markers" "AlternateText" "List » Markers"}}

        <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
            <text-chunk font="helvetica-bold" color="text">List markers </text-chunk>
            <text-chunk color="text">are sequences of symbols which precede each </text-chunk>
            <text-chunk font="helvetica-bold" color="text">list item</text-chunk>
            <text-chunk color="text">. By default, the list marker is the </text-chunk>
            <text-chunk font="helvetica-bold" color="text">bullet (`•`) </text-chunk>
            <text-chunk color="text">symbol. The component is represented in templates using the </text-chunk>
            <text-chunk color="secondary">{{xmlEscape "<list-marker>"}} </text-chunk>
            <text-chunk color="text">tag. The tag is supported by both the </text-chunk>
            <text-chunk font="helvetica-bold" color="text">list </text-chunk>
            <text-chunk color="text">and </text-chunk>
            <text-chunk font="helvetica-bold" color="text">list item</text-chunk>
            <text-chunk color="text">components. All list items inherit the properties of the list marker defined by the list. However, individual list items can override all or some of the properties of the list marker.</text-chunk>
            <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
        </paragraph>

        {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of list markers:"}}

        {{$codeBlock :=
`<list>
    <list-item>
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </list-item>
</list>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

        {{$codeBlock :=
`<list>
    <list-marker font="zapf-dingbats" text-rise="4">☞ </list-marker>
    <list-item>
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </list-item>
    <list-item>
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </list-item>
</list>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

        {{$codeBlock :=
`<list>
    <list-item>
        <list-marker color="#c62828">1. </list-marker>
        <paragraph><text-chunk>A</text-chunk></paragraph>
    </list-item>
    <list-item>
        <list-marker color="#1b5e20">2. </list-marker>
        <paragraph><text-chunk>B</text-chunk></paragraph>
    </list-item>
    <list-item>
        <list-marker color="#0d47a1">3. </list-marker>
        <paragraph><text-chunk>C</text-chunk></paragraph>
    </list-item>
</list>`
        }}
        {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

        {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "15 0 10 0" "Text" "Supported attributes:"}}

        <paragraph>
            <text-chunk font="helvetica-bold">List markers </text-chunk>
            <text-chunk>are basically </text-chunk>
            <text-chunk font="helvetica-bold">text chunks </text-chunk>
            <text-chunk>so they share the same attributes. Please see chapter </text-chunk>
            <text-chunk link="page(7, 0, 50)">2.1.1. Text Chunks</text-chunk>
            <text-chunk> in order to see the full list of supported attributes.</text-chunk>
        </paragraph>
    </chapter>
</chapter>
