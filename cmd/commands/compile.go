package commands

import (
	"fmt"
	"os"

	"github.com/aml-org/amf-custom-validator/cmd/commands/helpers"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"github.com/open-policy-agent/opa/ast"
)

func Compile() {
	helpers.ValidateNArgs(3, "[REGO_VERSION=<v0|v1>] acv compile <PROFILE>")

	profile := helpers.ReadOrPanic(os.Args[2])

	regoVersion := ast.RegoV1 // default to v1
	if v := os.Getenv("REGO_VERSION"); v == "v0" {
		regoVersion = ast.RegoV0
	}

	_, err := validator.ProcessProfile(string(profile), true, nil, regoVersion)

	if err != nil {
		panic(err)
	}
	fmt.Println("Compile Success!")
	os.Exit(0)
}
