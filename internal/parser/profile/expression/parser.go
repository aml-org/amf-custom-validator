package expression

import (
	"errors"
	"github.com/aml-org/amfopa/internal/parser/path"
	"github.com/aml-org/amfopa/internal/parser/profile/constraints"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"github.com/aml-org/amfopa/internal/parser/yaml"
	y "github.com/kylelemons/go-gypsy/yaml"
)

func Parse(name string, data y.Map, level string) (statements.Rule, error) {
	targetClass, err := yaml.GetString(data, "targetClass")
	if err != nil {
		return nil, err
	}
	message, err := yaml.GetString(data, "message")
	if err != nil {
		message = "Validation err"
	}

	exp := newTopLevelExpression(false, name, message, level, targetClass)
	v := exp.GenVar(statements.ForAll,nil)

	value, err := parseExpressionValue(v, data)
	if err != nil {
		return nil, err
	}
	exp.Value = value

	return exp, nil
}


func parseExpressionValue(variable statements.Variable, data y.Map) (statements.Rule, error) {
	v, err := yaml.GetMap(data, "propertyConstraints")
	if err == nil {
		return parseImplicitAnd(v, variable)
	}

	and, err := yaml.GetList(data, "and")
	if err == nil {
		return parseAnd(and, variable)
	}
	return nil, errors.New("unknown expression node, cannot find properties to parse")

}

func parseAnd(and y.List, variable statements.Variable) (statements.Rule, error) {
	values := make([]statements.Rule, len(and))
	for _,n := range and {
		switch p := n.(type) {
		case y.Map:
			c, err := parseExpressionValue(variable, p)
			if err != nil {
				return nil, err
			}
			values = append(values, c)
		default:
			return nil, errors.New("AND must be a list")
		}
	}
	return NewAnd(false, values), nil
}

func parseImplicitAnd(propertyConstraints y.Map, variable statements.Variable) (statements.Rule, error) {
	values := make([]statements.Rule, 0)
	for pathString,constraint := range propertyConstraints {
		propertyPath, err := path.ParsePath(pathString)
		if err != nil {
			return nil, err
		}

		switch p := constraint.(type) {
		case y.Map:
			c, err := constraints.Parse(propertyPath, variable, p)
			if err != nil {
				return nil, err
			}
			values = append(values, c)
		default:
			return nil, errors.New("PropertyConstraint must be a map")
		}
	}
	return NewAnd(false, values), nil
}
