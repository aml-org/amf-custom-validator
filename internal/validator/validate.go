package validator

import (
	"context"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

// Note: Only exported functions should be able to close channels. Otherwise, we would try to close an already closed channel
// These should not even be closed here, but in the `pkg` module. Closing here because we have internal function calls
// directly to `internal.Validate` rather than `pkg.Validate`
func Validate(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return validateKernel(profileText, jsonldText, debug, eventChan, false)
}

func ValidateWithOptimization(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return validateKernel(profileText, jsonldText, debug, eventChan, true)
}

func validateKernel(profileText string, jsonldText string, debug bool, eventChan *chan e.Event, optimize bool) (string, error) {
	// Generate and compile Rego code
	compiledRego, err := ProcessProfile(profileText, debug, eventChan)

	if err != nil {
		CloseEventChan(eventChan)
		return "", err
	}

	// TODO add events for optimization
	// TODO refactor code to parse profile once only
	var optimizerPtr *Optimizer = nil
	if optimize {
		parsed, _ := parser.Parse(profileText)
		optimizer := NewOptimizer(parsed)
		optimizerPtr = &optimizer
	}

	return ValidateCompiled(compiledRego, jsonldText, debug, eventChan, optimizerPtr)
}

func ValidateCompiled(compiledRegoPtr *rego.PreparedEvalQuery, jsonldText string, debug bool, eventChan *chan e.Event, optimizer *Optimizer) (string, error) {
	compiledRego := *compiledRegoPtr

	// Normalize input
	normalizedInput, err := ProcessInput(jsonldText, debug, eventChan, optimizer)

	if err != nil {
		CloseEventChan(eventChan)
		return "", err
	}

	// Execute validation
	validationResult, err := executeValidation(eventChan, err, compiledRego, normalizedInput)

	if err != nil {
		CloseEventChan(eventChan)
		return "", err
	}

	// Build report
	report, err := processResult(validationResult, eventChan)

	if err != nil {
		CloseEventChan(eventChan)
		return "", err
	}

	CloseEventChan(eventChan)
	return report, err
}

func executeValidation(eventChan *chan e.Event, err error, compiledRego rego.PreparedEvalQuery, normalizedInput interface{}) (*rego.ResultSet, error) {
	dispatchEvent(e.NewEvent(e.OpaValidationStart), eventChan)
	validationResult, err := compiledRego.Eval(context.Background(), rego.EvalInput(normalizedInput))
	dispatchEvent(e.NewEvent(e.OpaValidationDone), eventChan)
	return &validationResult, err
}
