package pkg

import (
	internal "github.com/aml-org/amf-custom-validator/internal/validator"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func CompileProfile(profileText string, debug bool, eventChan *chan e.Event) (*rego.PreparedEvalQuery, error) {
	compiledProfile, err := internal.ProcessProfile(profileText, debug, eventChan)
	if err != nil {
		internal.CloseEventChan(eventChan)
		return nil, err
	}
	return compiledProfile, err
}
