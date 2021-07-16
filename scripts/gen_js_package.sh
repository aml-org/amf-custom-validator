#!/usr/bin/env bash
GOOS=js GOARCH=wasm go build -o wrappers/js/lib/main.wasm -ldflags "-s -w" js/validator.go