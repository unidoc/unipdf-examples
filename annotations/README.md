Annotations are different from normal PDF page contents and are stored separately.  They are intended to mark up
document contents, without changing the original document contents.

Support for creating annotations in UniDoc is through the unidoc/pdf/annotator package.
The annotator package creates the Annotation object and also an appearance stream, which is required to make
the annotation look the same in all viewers.

UniDoc's model package has support for all types of PDF annotations, whereas the annotator currently supports Square,
Circle, Line annotations.  Support for more annotation types will be added over time.  If you need support for an
unsupported annotation, please file a new issue or contact support.

The examples in this folder illustrate a few capabilities for creating ellipse, lines, rectangles.
