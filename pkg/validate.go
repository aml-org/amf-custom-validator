package pkg

import (
	internal "github.com/aml-org/amf-custom-validator/internal/validator"
	config2 "github.com/aml-org/amf-custom-validator/pkg/config"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func Validate(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return internal.Validate(profileText, jsonldText, debug, eventChan)
}

func ValidateCompiled(compiledRegoPtr *rego.PreparedEvalQuery, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return internal.ValidateCompiled(compiledRegoPtr, jsonldText, debug, eventChan)
}

func ValidateWithConfiguration(profileText string, jsonldText string, debug bool, eventChan *chan e.Event, validationConfig config2.ValidationConfiguration, reportConfig config2.ReportConfiguration) (string, error) {
	return internal.ValidateWithConfiguration(profileText, jsonldText, debug, eventChan, validationConfig, reportConfig)
}

func ValidateCompiledWithConfiguration(compiledRegoPtr *rego.PreparedEvalQuery, jsonldText string, debug bool, eventChan *chan e.Event, validationConfig config2.ValidationConfiguration, reportConfig config2.ReportConfiguration) (string, error) {
	return internal.ValidateCompiledWithConfiguration(compiledRegoPtr, jsonldText, debug, eventChan, validationConfig, reportConfig)
}
