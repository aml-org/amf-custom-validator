package profile

import (
	"errors"
	"github.com/aml-org/amfopa/internal/parser/path"
	"github.com/aml-org/amfopa/internal/parser/yaml"
	y "github.com/kylelemons/go-gypsy/yaml"
)

func ParseExpression(name string, data y.Map, level string, varGenerator VarGenerator) (Rule, error) {
	targetClass, err := yaml.GetString(data, "targetClass")
	if err != nil {
		return nil, err
	}
	message, err := yaml.GetString(data, "message")
	if err != nil {
		message = "Validation err"
	}

	exp := newTopLevelExpression(false, name, message, level, targetClass, varGenerator)

	value, err := parseExpressionValue(*exp.Variable, data, varGenerator)
	if err != nil {
		return nil, err
	}
	exp.Value = value

	return exp, nil
}

func parseNestedExpression(data y.Map, negated bool, variable Variable, path path.PropertyPath, varGenerator VarGenerator) (Rule, error) {
	nested := newNestedExpression(negated, variable, path, varGenerator)
	value, err := parseExpressionValue(nested.child, data, varGenerator)
	if err != nil {
		return nil, err
	}
	nested.value = value

	return nested, nil
}

func parseExpressionValue(variable Variable, data y.Map, varGenerator VarGenerator) (Rule, error) {
	v, err := yaml.GetMap(data, "propertyConstraints")
	if err == nil {
		return parseImplicitAnd(v, variable, varGenerator)
	}

	and, err := yaml.GetList(data, "and")
	if err == nil {
		return parseAnd(and, variable, varGenerator)
	}
	return nil, errors.New("unknown expression node, cannot find properties to parse")

}

func parseAnd(and y.List, variable Variable, varGenerator VarGenerator) (Rule, error) {
	values := make([]Rule, len(and))
	for _, n := range and {
		switch p := n.(type) {
		case y.Map:
			c, err := parseExpressionValue(variable, p, varGenerator)
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

func parseImplicitAnd(propertyConstraints y.Map, variable Variable, varGenerator VarGenerator) (Rule, error) {
	var values []Rule
	for pathString, constraint := range propertyConstraints {
		propertyPath, err := path.ParsePath(pathString)
		if err != nil {
			return nil, err
		}

		switch p := constraint.(type) {
		case y.Map:
			cs, err := ParseConstraint(propertyPath, variable, p, varGenerator)
			if err != nil {
				return nil, err
			}
			values = append(values, cs...)
		default:
			return nil, errors.New("PropertyConstraint must be a map")
		}
	}
	return NewAnd(false, values), nil
}
