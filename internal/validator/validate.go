package validator

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aml-org/amf-custom-validator/internal/generator"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func Validate(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	err, module := Generate(profileText, debug, eventChan)
	normalizedInput := Normalize_(jsonldText, debug, eventChan)

	dispatch(e.NewEvent(e.OpaValidationStart), eventChan)
	validator := rego.New(
		rego.Query("data."+module.Name+"."+module.Entrypoint),
		rego.Module(module.Name+".rego", module.Code),
		rego.Input(normalizedInput),
	)
	ctx := context.Background()
	result, err := validator.Eval(ctx)
	dispatch(e.NewEvent(e.OpaValidationDone), eventChan)

	if err != nil {
		closeIfNotNil(eventChan)
		return "", err
	} else {
		dispatch(e.NewEvent(e.BuildReportStart), eventChan)
		report, err := BuildReport(result)
		dispatch(e.NewEvent(e.BuildReportDone), eventChan)
		closeIfNotNil(eventChan)
		return report, err
	}
}

func Normalize_(jsonldText string, debug bool, receiver *chan e.Event) interface{} {
	dispatch(e.NewEvent(e.InputDataParsingStart), receiver)
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonldText)))
	decoder.UseNumber()

	var input interface{}
	if err := decoder.Decode(&input); err != nil {
		panic(err)
	}
	dispatch(e.NewEvent(e.InputDataParsingDone), receiver)

	dispatch(e.NewEvent(e.InputDataNormalizationStart), receiver)
	normalizedInput := Index(Normalize(input))
	dispatch(e.NewEvent(e.InputDataNormalizationDone), receiver)

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
	return normalizedInput
}

func Generate(profileText string, debug bool, eventChan *chan e.Event) (error, generator.RegoUnit) {
	dispatch(e.NewEvent(e.ProfileParsingStart), eventChan)
	parsed, err := parser.Parse(profileText)
	dispatch(e.NewEvent(e.ProfileParsingDone), eventChan)
	if err != nil {
		panic(err)
	}

	if debug {
		println("Logic translation")
		println("-------------------------------")
		println(parsed.String())
	}
	dispatch(e.NewEvent(e.RegoGenerationStart), eventChan)
	module := generator.Generate(*parsed)
	dispatch(e.NewEvent(e.RegoGenerationDone), eventChan)
	if debug {
		println("Generated profile")
		println("-------------------------------")
		println(module.Code)
	}
	return err, module
}

func dispatch(event e.Event, eventChan *chan e.Event) {
	if eventChan != nil {
		*eventChan <- event
	}
}

func closeIfNotNil(eventChan *chan e.Event) {
	if eventChan != nil {
		close(*eventChan)
	}
}
