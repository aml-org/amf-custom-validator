package constraints

import (
	"errors"
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/path"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"github.com/aml-org/amfopa/internal/parser/yaml"
	y "github.com/kylelemons/go-gypsy/yaml"
)

func Parse(path path.PropertyPath, variable statements.Variable, constraint y.Map) ([]statements.Rule, error) {
	var acc []statements.Rule


	min,err := yaml.GetInt(constraint,"minCount")
	if err == nil {
		acc = append(acc, newMinCount(false, variable, path, min))
	}

	max,err := yaml.GetInt(constraint,"maxCount")
	if err == nil {
		acc = append(acc, newMaxCount(false, variable, path, max))
	}

	in,err := yaml.GetList(constraint, "in")
	if err == nil {
		l, err := scalarList(in)
		if err != nil {
			return nil,err
		}
		acc = append(acc, newIn(false, variable, path, l))
	}

	return acc, nil
}

func scalarList(in y.List) ([]string, error) {
	var acc []string
	for _,e := range in {
		switch pe := e.(type) {
		case y.Scalar:
			acc = append(acc, pe.String())
		default:
			return nil, errors.New(fmt.Sprintf("expected scalars in 'in' constraint, found %v", e))
		}
	}

	return acc,nil
}