package main

import (
	"github.com/aml-org/amf-custom-validator/pkg"
	"io/ioutil"
	"os"
)

const (
	ITERATIONS int = 100
	PROFILE        = "test/data/production/best-practices/profile.yaml"
	DATA           = "test/data/production/best-practices/negative1.raml.jsonld"
)

func main() {
	if len(os.Args) < 2 {
		noPrecompilation()
	} else if os.Args[1] == "--pre-compiled" {
		precompiled()
	} else {
		panic("Usage: 'go run performance/main/performance.go' with optional --pre-compiled flag")
	}
}

func noPrecompilation() {
	profile := read(PROFILE)
	data := read(DATA)

	for i := 0; i < ITERATIONS; i += 1 {
		_, _ = pkg.Validate(profile, data, false, nil)
	}
}

func precompiled() {
	profile := read(PROFILE)
	data := read(DATA)
	preparedEvalQuery, _ := pkg.CompileProfile(profile, false, nil)

	for i := 0; i < ITERATIONS; i += 1 {
		_, _ = pkg.ValidateCompiled(preparedEvalQuery, data, false, nil)
	}
}

func read(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
