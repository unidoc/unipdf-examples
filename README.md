# Examples

This contains demonstration of the many features UniDoc is capable of. Example code should make
it easy for users to know how to use all aspects of UniDoc. Feel free to add to this by submitting
a pull request.

The examples are also available on UniDoc's website: [https://unidoc.io/examples/](https://unidoc.io/examples/). 

### Build all examples

```bash
find . -name "*.go" -print0 | CGO_ENABLED=1 xargs -0 -n1 go build
```
