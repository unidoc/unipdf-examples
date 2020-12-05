#!/bin/bash

mkdir -p bin

echo "Building to bin/ folder"

#CGO required to build example relying on crypto11 dependency.
find . -name "*.go" ! -name "*pkcs11*.go" -print0 | CGO_ENABLED=0 xargs -0 -n1 go build
#Temporarily disabled due to module issues associated with crypto11 dependency.
#TODO: Update crypto11 to newest version.
#find . -name "*pkcs11*.go" -print0 | CGO_ENABLED=1 xargs -0 -n1 go build

mv pdf_* fdf_* jbig2_* bin/

