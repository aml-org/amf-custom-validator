package parser

import (
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	y "github.com/aml-org/amf-custom-validator/internal/parser/yaml"
)

func Parse(profileText string) (*profile.Profile, error) {

	node, err := y.NewYaml([]byte(profileText))
	if err != nil {
		return nil, err
	}

	prof, err := profile.Parse(node)

	if err != nil {
		return nil, err
	}
	return &prof, nil
}
