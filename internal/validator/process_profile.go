package validator

import (
	"context"
	"github.com/aml-org/amf-custom-validator/internal/generator"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func ProcessProfileWASM(profileText string, debug bool, eventChan *chan e.Event) ([]byte, error) {
	// Generate Rego code
	regoUnit, err := GenerateRego(profileText, debug, eventChan)

	if err != nil {
		return nil, err
	}

	result, err := CompileRegoWASM(regoUnit, eventChan)
	// Compile Rego code
	return result, err
}

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

	// Generate Rego code
	dispatchEvent(e.NewEvent(e.RegoGenerationStart), eventChan)
	module := generator.Generate(*parsed)
	dispatchEvent(e.NewEvent(e.RegoGenerationDone), eventChan)

	return &module, err
}

// unsafeBuiltinsMap When updating to 0.35 ast.NetLookupIPAddr will be available and needs to be added and blocked too
var unsafeBuiltinsMap = map[string]struct{}{
	ast.HTTPSend.Name:        {},
	ast.WalkBuiltin.Name:     {},
	ast.OPARuntime.Name:      {},
	ast.RegoParseModule.Name: {},
}

func CompileRego(regoUnit *generator.RegoUnit, eventChan *chan e.Event) (*rego.PreparedEvalQuery, error) {
	dispatchEvent(e.NewEvent(e.RegoCompilationStart), eventChan)
	query := rego.Query("data." + regoUnit.Name + "." + regoUnit.Entrypoint)
	module := rego.Module(regoUnit.Name+".rego", regoUnit.Code)
	unsafeBuiltins := rego.UnsafeBuiltins(unsafeBuiltinsMap)
	preparedEvalQuery, err := rego.New(query, module, unsafeBuiltins).PrepareForEval(context.Background())
	dispatchEvent(e.NewEvent(e.RegoCompilationDone), eventChan)
	return &preparedEvalQuery, err
}

func CompileRegoWASM(regoUnit *generator.RegoUnit, eventChan *chan e.Event) ([]byte, error) {
	dispatchEvent(e.NewEvent(e.RegoCompilationStart), eventChan)

	// create Rego
	query := rego.Query("data." + regoUnit.Name + "." + regoUnit.Entrypoint)
	module := rego.Module(regoUnit.Name+".rego", regoUnit.Code)
	unsafeBuiltins := rego.UnsafeBuiltins(unsafeBuiltinsMap)
	validator := rego.New(query, module, unsafeBuiltins)

	// create wasm
	ctx := context.Background()
	compileResult, err := validator.Compile(ctx)
	dispatchEvent(e.NewEvent(e.RegoCompilationDone), eventChan)
	result := compileResult.Bytes
	return result, err
}

//func CompileRegoWASM2(regoUnit *generator.RegoUnit) ([]byte, error) {
//	buf := bytes.NewBuffer(nil)
//
//	fs := memoryfs.New()
//	fs.WriteFile("profile.rego", []byte(regoUnit.Code), 0o700)
//
//	capabilities := &ast.Capabilities{}
//	for key, _ := range unsafeBuiltinsMap {
//		f := &ast.Builtin{
//			Name: key,
//		}
//		capabilities.Builtins = append(capabilities.Builtins, f)
//	}
//
//	compiler := compile.New().
//		WithPaths("profile.rego").
//		WithTarget(compile.TargetWasm).
//		WithOutput(buf).
//		WithEntrypoints(regoUnit.Name + "/" + regoUnit.Entrypoint).
//		WithFS(fs).
//		WithRegoAnnotationEntrypoints(true)
//	//WithCapabilities(capabilities)
//
//	ctx := context.Background()
//	err := compiler.Build(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	result := compiler.Bundle().WasmModules[0].Raw
//	return result, nil
//}
