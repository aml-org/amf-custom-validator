#!/usr/bin/env bash
GOOS=js GOARCH=wasm go build -o wrappers/js/lib/main.wasm -ldflags "-s -w" js/validator.go
gzip -9 -v -c wrappers/js/lib/main.wasm > wrappers/js/lib/main.wasm.gz
rm wrappers/js/lib/main.wasm