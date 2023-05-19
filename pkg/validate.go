package pkg

import (
	internal "github.com/aml-org/amf-custom-validator/internal/validator"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func Validate(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return internal.Validate(profileText, jsonldText, debug, eventChan)
}

func ValidateCompiled(compiledRegoPtr *rego.PreparedEvalQuery, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return internal.ValidateCompiled(compiledRegoPtr, jsonldText, debug, eventChan)
}
