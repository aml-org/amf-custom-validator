package validator

import (
	"context"

	"github.com/aml-org/amf-custom-validator/internal/generator"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func ProcessProfile(profileText string, debug bool, eventChan *chan e.Event, version ast.RegoVersion) (*rego.PreparedEvalQuery, error) {
	regoUnit, err := GenerateRego(profileText, debug, eventChan, version)

	if err != nil {
		return nil, err
	}

	// Compile Rego code
	return CompileRego(regoUnit, eventChan, version)
}

func GenerateRego(profileText string, debug bool, eventChan *chan e.Event, version ast.RegoVersion) (*generator.RegoUnit, error) {
	// Parse profile
	dispatchEvent(e.NewEvent(e.ProfileParsingStart), eventChan)
	parsed, err := parser.Parse(profileText)
	dispatchEvent(e.NewEvent(e.ProfileParsingDone), eventChan)

	if err != nil {
		return nil, err
	}

	// Generate Rego code
	dispatchEvent(e.NewEvent(e.RegoGenerationStart), eventChan)
	module := generator.Generate(*parsed, version)
	dispatchEvent(e.NewEvent(e.RegoGenerationDone), eventChan)

	return &module, err
}
func GenerateRegoV0(profileText string, debug bool, eventChan *chan e.Event) (*generator.RegoUnit, error) {
	return GenerateRego(profileText, debug, eventChan, ast.RegoV0)
}

func GenerateRegoV1(profileText string, debug bool, eventChan *chan e.Event) (*generator.RegoUnit, error) {
	return GenerateRego(profileText, debug, eventChan, ast.RegoV1)
}

// unsafeBuiltinsMap When updating to 0.35 ast.NetLookupIPAddr will be available and needs to be added and blocked too
var unsafeBuiltinsMap = map[string]struct{}{
	ast.HTTPSend.Name:        {},
	ast.WalkBuiltin.Name:     {},
	ast.OPARuntime.Name:      {},
	ast.RegoParseModule.Name: {},
}

func CompileRego(regoUnit *generator.RegoUnit, eventChan *chan e.Event, regoVersion ast.RegoVersion) (*rego.PreparedEvalQuery, error) {
	dispatchEvent(e.NewEvent(e.RegoCompilationStart), eventChan)
	query := rego.Query("data." + regoUnit.Name + "." + regoUnit.Entrypoint)
	module := rego.Module(regoUnit.Name+".rego", regoUnit.Code)
	unsafeBuiltins := rego.UnsafeBuiltins(unsafeBuiltinsMap)
	preparedEvalQuery, err := rego.New(
		query,
		module,
		unsafeBuiltins,
		rego.SetRegoVersion(regoVersion)).
		PrepareForEval(context.Background())
	dispatchEvent(e.NewEvent(e.RegoCompilationDone), eventChan)
	return &preparedEvalQuery, err
}
