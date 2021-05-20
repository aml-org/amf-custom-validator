package parser

import (
	"github.com/aml-org/amfopa/internal/parser/profile"
	yamlparser "github.com/aml-org/amfopa/internal/parser/yaml"
)

func Parse(profileText string) (*profile.Profile, error) {

	node, err := yamlparser.Parse(profileText)
	if err != nil {
		return nil, err
	}

	prof, err := profile.Parse(node)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}
