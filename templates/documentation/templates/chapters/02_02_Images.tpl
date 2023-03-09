<chapter margin="10 0">
    {{template "chapter-title" dict "Text" "Images"}}

    <paragraph margin="10 0 5 0" line-height="1.1" enable-word-wrap="true">
        <text-chunk color="text">Images can be rendered using the </text-chunk>
        <text-chunk font="helvetica-bold" color="text">image </text-chunk>
        <text-chunk color="text">component. The component is represented in templates using the </text-chunk>
        <text-chunk color="secondary">{{xmlEscape "<paragraph>"}} </text-chunk>
        <text-chunk color="text">tag.</text-chunk>
        <text-chunk color="text">{{.newline}}{{.newline}}</text-chunk>
    </paragraph>

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "0 0 10 0" "Text" "Basic syntax of images:"}}

    {{$codeBlock :=
`<image src="path('templates/res/images/sample-image.jpg')"
    fit-mode="fill-width"></image>`
    }}
    {{template "code-block" dict "Margin" "0 0 5 0" "Code" $codeBlock}}

    {{template "paragraph" dict "Font" "helvetica-bold" "Margin" "15 0 10 0" "Text" "Supported attributes:"}}

    <!-- src -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-4.jpg')"
    fit-mode="fill-width"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "src"
        "AttrDescr" "The `src` attribute is used to specify the source of the image. The source image can be specified by path in the template, in which case it will be loaded, or it can be specified by name if the image is loaded in the image map of the options used to draw the template."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.2. Images"
        "AttrLink" "page(48, 0, 150)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "empty string"
    }}

    <!-- fit-mode -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-5.jpg')"
    fit-mode="fill-width"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "fit-mode"
        "AttrDescr" "The `fit-mode` attribute controls the sizing of the component relative to the available space. When the attribute value is set to `fill-width`, the component is scaled so that it occupies the entire available width, preserving the original aspect ratio."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "none" "fill-width")
        "AttrDefault" "none"
    }}

    <!-- align -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-2.jpg')"
    width="160" height="90"></image>`
`<image src="path('templates/res/images/sample-image-2.jpg')"
    width="160" height="90" align="center"></image>`
`<image src="path('templates/res/images/sample-image-2.jpg')"
    width="160" height="90" align="right"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "align"
        "AttrDescr" "The `align` attribute aligns the image in the available space, based on the specified option."
        "CodeBlocks" $codeBlocks
        "AttrValues" (array "left" "right" "center")
        "AttrDefault" "left"
    }}

    <!-- width -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-3.jpg')"
    width="160" height="90"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "width"
        "AttrDescr" "The `width` attribute sets the width of the rendered image."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "original image width"
    }}

    <!-- height -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image.jpg')"
    width="160" height="90"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "height"
        "AttrDescr" "The `height` attribute sets the height of the rendered image."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "original image height"
    }}

    <!-- opacity -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-2.jpg')"
    width="160" height="90"></image>`
`<image src="path('templates/res/images/sample-image-2.jpg')"
    width="160" height="90" opacity="0.5"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "opacity"
        "AttrDescr" "The `opacity` attribute allows setting the transparency of the image."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "1.0"
        "AttrValues" (array "any value between 0 and 1")
    }}

    <!-- margin -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-3.jpg')"
    width="160" height="90"></image>`
`<image src="path('templates/res/images/sample-image-3.jpg')"
    width="160" height="90" margin="10 0 0 50"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "margin"
        "AttrDescr" "The `margin` attribute allows setting a configurable amount of space around the component."
        "AttrLinkDescr" "For more information see chapter "
        "AttrLinkText" "4.4. Margin and Padding"
        "AttrLink" "page(52, 0, 230)"
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}

    <!-- x -->
    {{template "attr-showcase" dict "AttrName" "x"
        "AttrDescr" "The `x` attribute sets the X coordinate of the top left corner of the component. The attribute is used only when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- y -->
    {{template "attr-showcase" dict "AttrName" "y"
        "AttrDescr" "The `y` attribute sets the Y coordinate of the top left corner of the component. The attribute is used only when manually positioning the component."
        "CodeBlocks" ""
        "AttrDefault" "0"
    }}

    <!-- angle -->
    {{$codeBlocks := array
`<image src="path('templates/res/images/sample-image-5.jpg')"
    width="160" height="90" angle="180"></image>`
    }}
    {{template "attr-showcase" dict "AttrName" "angle"
        "AttrDescr" "The `angle` attribute represents an angle specified in degrees to rotate the image by. The rotation is applied anti-clockwise."
        "CodeBlocks" $codeBlocks
        "AttrDefault" "0"
    }}
</chapter>
