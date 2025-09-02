# Examples

This example repository demonstrates many use cases for UniDoc's UniPDF library. Example code should make
it easy for users to get started with UniPDF. Feel free to add to this by submitting
a pull request.

While the majority of examples are fully in pure Go, there are a few examples that demonstrate additional 
functionality that requires CGO and external dependencies. Those examples are clarified by filename suffix "_cgo.go".

## License codes
UniPDF requires license codes to operate, there are two options:
- Metered License API keys: Free ones can be obtained at https://cloud.unidoc.io
- Offline codes: Can be purchased at https://unidoc.io/pricing

Most of the examples demonstrate loading the Metered License API keys through an environment
variable `UNIDOC_LICENSE_API_KEY`.

Examples for Offline License Key loading can be found in the license subdirectory.

### Build all examples

#### Building with go modules:
Simply run the build script which builds all the binaries to subfolder `bin/`

```bash
$ ./build_examples.sh
```

#### Building with GOPATH:
Building with GOPATH requires a slightly different approach due to the `/v4` semantic import portion of the unipdf import paths.  There are two options:

Both options start with:
- `go get github.com/unidoc/unipdf/...` to download the packages

Then one can decide between the two options:

1. Remove the `/v4/` in the unipdf import paths, e.g. use `github.com/unidoc/unipdf/core` instead of `github.com/unidoc/unipdf/v4/core`
2. Alternatively create a symbolic link from the v4 subdirectory of unipdf to the unipdf repository, i.e.
```bash
ln -s $GOPATH/src/github.com/unidoc/unipdf $GOPATH/src/github.com/unidoc/unipdf/v4
```
or move/copy the unipdf folder to unipdf/v4 if symbolic links are not an option.

Once this has been done, then can build using the build script as well:
```bash
$ ./build_examples.sh
```
or build individual example codes as desired.
