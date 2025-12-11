package validator

import (
	"context"

	"github.com/aml-org/amf-custom-validator/internal/generator"
	"github.com/aml-org/amf-custom-validator/internal/notifications"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func ProcessProfile(profileText string, debug bool, eventChan *chan e.Event) (*rego.PreparedEvalQuery, []profile.LlmMetric, error) {
	// Generate Rego code
	regoUnit, err := GenerateRego(profileText, debug, eventChan)

	if err != nil {
		return nil, nil, err
	}

	// Compile Rego code
	return CompileRego(regoUnit, eventChan)
}

func GenerateRego(profileText string, debug bool, eventChan *chan e.Event) (*generator.RegoUnit, error) {
	// Parse profile
	notifications.DispatchEvent(e.NewEvent(e.ProfileParsingStart), eventChan)
	parsed, err := parser.Parse(profileText)
	notifications.DispatchEvent(e.NewEvent(e.ProfileParsingDone), eventChan)

	if err != nil {
		return nil, err
	}

	// Generate Rego code
	notifications.DispatchEvent(e.NewEvent(e.RegoGenerationStart), eventChan)
	module := generator.Generate(*parsed) // TODO: do something with the LLM metrics
	notifications.DispatchEvent(e.NewEvent(e.RegoGenerationDone), eventChan)

	return &module, err
}

// unsafeBuiltinsMap When updating to 0.35 ast.NetLookupIPAddr will be available and needs to be added and blocked too
var unsafeBuiltinsMap = map[string]struct{}{
	ast.HTTPSend.Name:        {},
	ast.WalkBuiltin.Name:     {},
	ast.OPARuntime.Name:      {},
	ast.RegoParseModule.Name: {},
}

func CompileRego(regoUnit *generator.RegoUnit, eventChan *chan e.Event) (*rego.PreparedEvalQuery, []profile.LlmMetric, error) {
	notifications.DispatchEvent(e.NewEvent(e.RegoCompilationStart), eventChan)
	query := rego.Query("data." + regoUnit.Name + "." + regoUnit.Entrypoint)
	module := rego.Module(regoUnit.Name+".rego", regoUnit.Code)
	unsafeBuiltins := rego.UnsafeBuiltins(unsafeBuiltinsMap)
	preparedEvalQuery, err := rego.New(query, module, unsafeBuiltins).PrepareForEval(context.Background())
	notifications.DispatchEvent(e.NewEvent(e.RegoCompilationDone), eventChan)
	return &preparedEvalQuery, regoUnit.LlmMetrics, err
}
