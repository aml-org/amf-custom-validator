package constraints

import (
	"errors"
	"github.com/aml-org/amfopa/internal/parser/path"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"github.com/aml-org/amfopa/internal/parser/yaml"
	y "github.com/kylelemons/go-gypsy/yaml"
)

func Parse(path path.PropertyPath, variable statements.Variable, constraint y.Map) (statements.Rule, error) {
	min,err := yaml.GetInt(constraint,"minCount")
	if err == nil {
		return newMinCount(false, variable, path, min), nil
	}

	max,err := yaml.GetInt(constraint,"maxCount")
	if err == nil {
		return newMaxCount(false, variable, path, max), nil
	}

	return nil, errors.New("unknown constraint")
}