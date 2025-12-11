package validator

import (
	"context"

	"github.com/aml-org/amf-custom-validator/internal/data"
	"github.com/aml-org/amf-custom-validator/internal/llm"
	"github.com/aml-org/amf-custom-validator/internal/notifications"
	c "github.com/aml-org/amf-custom-validator/pkg/config"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

// Note: Only exported functions should be able to close channels. Otherwise, we would try to close an already closed channel
// These should not even be closed here, but in the `pkg` module. Closing here because we have internal function calls
// directly to `internal.Validate` rather than `pkg.Validate`

func Validate(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return ValidateWithConfiguration(profileText, jsonldText, debug, eventChan, c.DefaultValidationConfiguration{}, c.DefaultReportConfiguration())
}

func ValidateCompiled(compiledRegoPtr *rego.PreparedEvalQuery, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return ValidateCompiledWithConfiguration(compiledRegoPtr, jsonldText, debug, eventChan, c.DefaultValidationConfiguration{}, c.DefaultReportConfiguration())
}

func ValidateWithConfiguration(profileText string, jsonldText string, debug bool, eventChan *chan e.Event, validationConfig c.ValidationConfiguration, reportConfig c.ReportConfiguration) (string, error) {
	// Generate and compile Rego code
	compiledRego, llmMetrics, err := ProcessProfile(profileText, debug, eventChan)

	if err != nil {
		notifications.CloseEventChan(eventChan)
		return "", err
	}

	var llmEvaluationResults map[string]bool
	if len(llmMetrics) > 0 {
		// This will evaluate the metrics and add them to the input JSONLD
		llmEvaluationResults = llm.EnrichDataWithLlmMetrics(llmMetrics, jsonldText, debug, eventChan)
	} else {
		llmEvaluationResults = make(map[string]bool)
	}

	return ValidateCompiledWithConfiguration(compiledRego, jsonldText, debug, eventChan, validationConfig, reportConfig)
}

func ValidateCompiledWithConfiguration(compiledRegoPtr *rego.PreparedEvalQuery, jsonldText string, debug bool, eventChan *chan e.Event, validationConfig c.ValidationConfiguration, reportConfig c.ReportConfiguration) (string, error) {
	compiledRego := *compiledRegoPtr

	// Normalize input
	normalizedInput, err := data.ProcessInput(jsonldText, debug, eventChan)

	if err != nil {
		notifications.CloseEventChan(eventChan)
		return "", err
	}

	// Execute validation
	validationResult, err := executeValidation(eventChan, err, compiledRego, normalizedInput)

	if err != nil {
		notifications.CloseEventChan(eventChan)
		return "", err
	}

	// Build report
	report, err := processResult(validationResult, eventChan, validationConfig, reportConfig)

	if err != nil {
		notifications.CloseEventChan(eventChan)
		return "", err
	}

	notifications.CloseEventChan(eventChan)
	return report, err
}

func executeValidation(eventChan *chan e.Event, err error, compiledRego rego.PreparedEvalQuery, normalizedInput any) (*rego.ResultSet, error) {
	notifications.DispatchEvent(e.NewEvent(e.OpaValidationStart), eventChan)
	validationResult, err := compiledRego.Eval(context.Background(), rego.EvalInput(normalizedInput))
	notifications.DispatchEvent(e.NewEvent(e.OpaValidationDone), eventChan)
	return &validationResult, err
}
