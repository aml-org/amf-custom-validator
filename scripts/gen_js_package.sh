#!/usr/bin/env bash
GOOS=js GOARCH=wasm $GO build -o wrappers/js/lib/main.wasm -ldflags "-s -w" pkg/validator.go