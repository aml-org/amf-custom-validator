package main

import internal "github.com/aml-org/amf-custom-validator/internal/validator"

func Validate(profileText string, jsonldText string, debug bool) (string, error) {
	return internal.Validate(profileText, jsonldText, debug)
}