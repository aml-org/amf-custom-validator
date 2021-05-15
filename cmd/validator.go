package main

import (
	"github.com/aml-org/amfopa/internal/parser"
	"os"
)

func main() {
	arg := os.Args[1]
	prof, err := parser.Parse(arg)
	if err != nil {
		panic(err)
	}

	println(prof.String())
}
