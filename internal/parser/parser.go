package parser

import (
	"github.com/aml-org/amfopa/internal/parser/profile"
	y "github.com/smallfish/simpleyaml"
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
