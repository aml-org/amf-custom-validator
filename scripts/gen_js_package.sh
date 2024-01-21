#!/usr/bin/env bash
GOOS=js GOARCH=wasm go build -o wrappers/js/assets/main.wasm -ldflags "-s -w" js/validator.go
gzip -9 -v -c wrappers/js/assets/main.wasm > wrappers/js/assets/main.wasm.gz
rm wrappers/js/assets/main.wasm