package commands

import (
	"fmt"
	"os"

	h "github.com/aml-org/amf-custom-validator/cmd/commands/helpers"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"github.com/open-policy-agent/opa/ast"
)

func Validate() {

	h.ValidateNsArgs([]int{4, 5}, "[REGO_VERSION=<v0|v1>] acv validate <PROFILE> <DATA> [OUTPUT_FILE]")

	profile := h.ReadOrPanic(os.Args[2])
	data := h.ReadOrPanic(os.Args[3])

	regoVersion := ast.RegoV1 // default to v1
	if v := os.Getenv("REGO_VERSION"); v == "v0" {
		regoVersion = ast.RegoV0
	}

	res, err := validator.Validate(string(profile), string(data), false, regoVersion, nil)
	h.CheckError(err)

	if len(os.Args) == 4 {
		// 4 args means output to standard output
		fmt.Println(res)
	}
	if len(os.Args) == 5 {
		// 5 args means output to file path
		path := os.Args[4]
		file := h.OpenOrCreateFile(path)
		h.WriteString(file, res)
	}

	os.Exit(0)
}
