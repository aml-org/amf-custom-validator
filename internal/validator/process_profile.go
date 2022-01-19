package validator

import (
	"context"
	"github.com/aml-org/amf-custom-validator/internal/generator"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func ProcessProfile(profileText string, debug bool, eventChan *chan e.Event) (*rego.PreparedEvalQuery, error) {
	// Generate Rego code
	regoUnit, err := GenerateRego(profileText, debug, eventChan)

	if err != nil {
		return nil, err
	}

	// Compile Rego code
	return CompileRego(regoUnit, eventChan)
}

func GenerateRego(profileText string, debug bool, eventChan *chan e.Event) (*generator.RegoUnit, error) {
	// Parse profile
	dispatchEvent(e.NewEvent(e.ProfileParsingStart), eventChan)
	parsed, err := parser.Parse(profileText)
	dispatchEvent(e.NewEvent(e.ProfileParsingDone), eventChan)

	if err != nil {
		return nil, err
	}

	if debug {
		println("Logic translation")
		println("-------------------------------")
		println(parsed.String())
	}

	// Generate Rego code
	dispatchEvent(e.NewEvent(e.RegoGenerationStart), eventChan)
	module := generator.Generate(*parsed)
	dispatchEvent(e.NewEvent(e.RegoGenerationDone), eventChan)

	if debug {
		println("Generated profile")
		println("-------------------------------")
		println(module.Code)
	}

	return &module, err
}

func CompileRego(regoUnit *generator.RegoUnit, eventChan *chan e.Event) (*rego.PreparedEvalQuery, error) {
	dispatchEvent(e.NewEvent(e.RegoCompilationStart), eventChan)
	query := rego.Query("data." + regoUnit.Name + "." + regoUnit.Entrypoint)
	module := rego.Module(regoUnit.Name+".rego", regoUnit.Code)
	preparedEvalQuery, err := rego.New(query, module).PrepareForEval(context.Background())
	dispatchEvent(e.NewEvent(e.RegoCompilationDone), eventChan)
	return &preparedEvalQuery, err
}
