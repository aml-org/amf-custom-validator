#!/usr/bin/env bash
set -e

go get -u github.com/mna/pigeon
pigeon ./third_party/propertyparser.peg > ./internal/parser/path/peg.go