package pkg

import (
	internal "github.com/aml-org/amf-custom-validator/internal/validator"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
)

func Validate(profileText string, jsonldText string, debug bool, eventChan *chan e.Event) (string, error) {
	return internal.Validate(profileText, jsonldText, debug, eventChan)
}