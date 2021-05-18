package validator

import(
	"bytes"
	"context"
	"encoding/json"
	"github.com/aml-org/amfopa/internal/generator"
	"github.com/aml-org/amfopa/internal/parser"
	"github.com/open-policy-agent/opa/rego"
	"io/ioutil"
)


func Validate(profilePath string, jsonldPath string, debug bool) (string, error) {
	parsed,err := parser.Parse(profilePath)
	if err != nil {
		panic(err)
	}
	module := generator.Generate(*parsed)

	if debug {
		println("Generated profile")
		println("-------------------------------")
		println(module.Code)
	}

	dataBytes, err := ioutil.ReadFile(jsonldPath)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(bytes.NewBuffer(dataBytes))
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
		enc.Encode(normalizedInput)
		println(b.String())
	}
	validator := rego.New(
		rego.Query("data." + module.Name + "." + module.Entrypoint),
		rego.Module(module.Name + ".rego", module.Code),
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