package validator

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aml-org/amfopa/internal/generator"
	"github.com/aml-org/amfopa/internal/parser"
	"github.com/open-policy-agent/opa/rego"
)

func Validate(profileText string, jsonldText string, debug bool) (string, error) {
	parsed, err := parser.Parse(profileText)
	if err != nil {
		panic(err)
	}
	module := generator.Generate(*parsed)

	if debug {
		println("Generated profile")
		println("-------------------------------")
		println(module.Code)
	}

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonldText)))
	decoder.UseNumber()

	var input interface{}
	if err := decoder.Decode(&input); err != nil {
		panic(err)
	}

	normalizedInput := Normalize(input)

	if debug {
		println("Input data")
		println("-------------------------------")
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "  ")
		err := enc.Encode(normalizedInput)
		if err != nil {
			panic(err)
		}
		println(b.String())
	}
	validator := rego.New(
		rego.Query("data."+module.Name+"."+module.Entrypoint),
		rego.Module(module.Name+".rego", module.Code),
		rego.Input(normalizedInput),
	)

	ctx := context.Background()
	result, err := validator.Eval(ctx)
	if err != nil {
		return "", err
	} else {
		return BuildReport(result)
	}
}
