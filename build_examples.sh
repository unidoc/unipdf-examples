#!/bin/bash

mkdir -p bin

echo "Building to bin/ folder"

# CGO required to build example relying on crypto11 and imagick dependency.
find . -name "*.go" ! -name "*_cgo.go" -print0 | CGO_ENABLED=0 xargs -0 -n1 go build -o bin
find . -name "*pkcs11*_cgo.go" -print0 | CGO_ENABLED=1 xargs -0 -n1 go build -o bin

# Custom encoder example need another CGO flags to build imagick lib.
find . -name "*custom_encoder*_cgo.go" -print0 | CGO_ENABLED=1 CGO_CFLAGS_ALLOW='-Xpreprocessor' xargs -0 -n1 go build -o bin
