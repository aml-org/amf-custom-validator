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
	targetClass,error := yaml.GetString(data, "targetClass")
	if error != nil {
		return nil,error
	}
	message,error := yaml.GetString(data, "message")
	if error != nil {
		message = "Validation error"
	}

	exp := newTopLevelExpression(false, name, message, level, targetClass)
	v := exp.GenVar(statements.ForAll,nil)

	value, error := parseExpressionValue(v, data)
	if error != nil {
		return nil, error
	}
	exp.Value = value

	return exp, nil
}


func parseExpressionValue(variable statements.Variable, data y.Map) (statements.Rule, error) {
	v,error := yaml.GetMap(data, "propertyConstraints")
	if error == nil {
		return parseImplicitAnd(v, variable)
	}

	and,error := yaml.GetList(data, "and")
	if error == nil {
		return parseAnd(and, variable)
	}
	return nil, errors.New("Unknown expression node, cannot find properties to parse")

}

func parseAnd(and y.List, variable statements.Variable) (statements.Rule, error) {
	values := make([]statements.Rule, len(and))
	for _,n := range and {
		switch p := n.(type) {
		case y.Map:
			c, error := parseExpressionValue(variable, p)
			if error != nil {
				return nil, error
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
		propertyPath, error := path.Parse(pathString)
		if error != nil {
			return nil, error
		}

		switch p := constraint.(type) {
		case y.Map:
			c, error := constraints.Parse(propertyPath, variable, p)
			if error != nil {
				return nil, error
			}
			values = append(values, c)
		default:
			return nil, errors.New("PropertyConstraint must be a map")
		}
	}
	return NewAnd(false, values), nil
}
