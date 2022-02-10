package validator

import (
	"context"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

// Note: Only exported functions should be able to close channels. Otherwise, we would try to close an already closed channel
// These should not even be closed here, but in the `pkg` module. Closing here because we have internal function calls
// directly to `internal.Validate` rather than `pkg.Validate`

func Validate(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	// Generate and compile Rego code
	compiledRego, err := ProcessProfile(profileText, debug, eventChan)

	if err != nil {
		CloseEventChan(eventChan)
		return "", err
	}

	return ValidateCompiled(compiledRego, jsonldText, debug, eventChan)
}

func ValidateCompiled(compiledRegoPtr *rego.PreparedEvalQuery, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	// Normalize input
	normalizedInput, err := ProcessInput(jsonldText, debug, eventChan)

	if err != nil {
		CloseEventChan(eventChan)
		return "", err
	}

	// Execute validation
	validationResult, err := executeValidation(eventChan, err, compiledRegoPtr, normalizedInput)

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

func executeValidation(eventChan *chan e.Event, err error, compiledRego *rego.PreparedEvalQuery, normalizedInput interface{}) (*rego.ResultSet, error) {
	dispatchEvent(e.NewEvent(e.OpaValidationStart), eventChan)
	//fmt.Println("hello")
	validationResult, err := compiledRego.Eval(context.Background(), rego.EvalInput(normalizedInput))
	dispatchEvent(e.NewEvent(e.OpaValidationDone), eventChan)
	return &validationResult, err
}
