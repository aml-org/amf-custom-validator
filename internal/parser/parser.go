package parser

import (
	"github.com/aml-org/amfopa/internal/parser/profile"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	yamlparser "github.com/aml-org/amfopa/internal/parser/yaml"
)

func Parse(path string) (*statements.Profile, error) {

	node,err := yamlparser.Parse(path)
	if err != nil {
		return nil,err
	}

	prof,err := profile.Parse(node)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}