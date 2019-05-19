# Examples

This example repository demonstrates many use cases for UniDoc's UniPDF library. Example code should make
it easy for users to get started with UniPDF. Feel free to add to this by submitting
a pull request.

The examples are also available on UniDoc's website: [https://unidoc.io/examples/](https://unidoc.io/examples/). 

### Build all examples

#### Building with go modules:
Simply run the build script which builds all the binaries to subfolder `bin/`

```bash
$ ./build_examples.sh
```

#### Building with GOPATH:
Building with GOPATH requires a slightly different approach due to the `/v3` semantic import portion of the unipdf import paths.  There are two options:

Both options start with:
- `go get github.com/unidoc/unipdf/...` to download the packages

Then one can decide between the two options:

1. Remove the `/v3/` in the unipdf import paths, e.g. use `github.com/unidoc/unipdf/core` instead of `github.com/unidoc/unipdf/v3/core`
2. Alternatively create a symbolic link from the v3 subdirectory of unipdf to the unipdf repository, i.e.
```bash
ln -s $GOPATH/src/github.com/unidoc/unipdf $GOPATH/src/github.com/unidoc/unipdf/v3
```
or move/copy the unipdf folder to unipdf/v3 if symbolic links are not an option.

Once this has been done, then can build using the build script as well:
```bash
$ ./build_examples.sh
```
or build individual example codes as desired.
