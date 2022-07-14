#!/bin/bash

# Stop script execution on error
set -euo pipefail

mkdir -p bin

echo "Building to bin/ folder"

# CGO required to build example relying on crypto11 and imagick dependency.
find . -name "*.go" ! -name "*_cgo.go" ! -name "lib_*" -print0 | CGO_ENABLED=0 xargs -0 -n1 -I% bash -c 'go build -o bin % || exit 255'
find . -name "*_cgo.go" -print0 | CGO_ENABLED=1 CGO_CFLAGS_ALLOW='-Xpreprocessor' xargs -0 -n1 -I% bash -c 'go build -o bin % || exit 255'
