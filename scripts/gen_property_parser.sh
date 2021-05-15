#!/usr/bin/env bash
go get -u github.com/mna/pigeon
pigeon propertyparser.peg > ../internal/parser/path/peg.go