#!/bin/bash

mkdir -p bin

echo "Building to bin/ folder"

#CGO required to build example relying on crypto11 dependency.
find . -name "*.go" ! -name "*pkcs11*.go" -print0 | CGO_ENABLED=0 xargs -0 -n1 go build
find . -name "*pkcs11*.go" -print0 | CGO_ENABLED=1 xargs -0 -n1 go build

mv pdf_* fdf_* bin/

