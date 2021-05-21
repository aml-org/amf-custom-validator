package profile

import (
	"errors"
	"github.com/aml-org/amfopa/internal/parser/path"
	y "github.com/smallfish/simpleyaml"
)

func ParseExpression(name string, data *y.Yaml, level string, varGenerator *VarGenerator) (Rule, error) {
	targetClass, err := data.Get("targetClass").String()
	if err != nil {
		return nil, errors.New("missing targetClass in validation definition")
	}
	message, err := data.Get("message").String()
	if err != nil {
		message = "Validation error"
	}
	exp := newTopLevelExpression(false, name, message, level, targetClass, varGenerator)

	value, err := parseExpressionValue(*exp.Variable, data, varGenerator)
	if err != nil {
		return nil, err
	}
	exp.Value = value

	return exp, nil
}

func parseNestedExpression(data *y.Yaml, negated bool, variable Variable, path path.PropertyPath, varGenerator *VarGenerator) (Rule, error) {
	nested := newNestedExpression(negated, variable, path, varGenerator)
	value, err := parseExpressionValue(nested.Child, data, varGenerator)
	if err != nil {
		return nil, err
	}
	nested.Value = value

	return nested, nil
}

func parseExpressionValue(variable Variable, data *y.Yaml, varGenerator *VarGenerator) (Rule, error) {
	v := data.Get("propertyConstraints")
	if v.IsFound() {
		return parseImplicitAnd(v, variable, varGenerator)
	}

	and := data.Get("and")
	if and.IsFound() {
		if and.IsArray() {
			return parseAnd(and, variable, varGenerator)
		} else {
			return nil, errors.New("and constraint must be a list")
		}

	}

	or := data.Get("or")
	if or.IsFound() {
		if or.IsArray() {
			return parseOr(or, variable, varGenerator)
		} else {
			return nil, errors.New("or constraint must be a list")
		}

	}

	return nil, errors.New("unknown expression node, cannot find properties to parse")

}

func parseAnd(and *y.Yaml, variable Variable, varGenerator *VarGenerator) (Rule, error) {
	var values []Rule
	size, _ := and.GetArraySize()
	for i := 0; i < size; i++ {
		n := and.GetIndex(i)
		if !n.IsMap() {
			return nil, errors.New("not found expected map for and constraint element")
		}
		c, err := parseExpressionValue(variable, n, varGenerator)
		if err != nil {
			return nil, err
		}
		values = append(values, c)
	}
	return NewAnd(false, values), nil
}

func parseOr(or *y.Yaml, variable Variable, varGenerator *VarGenerator) (Rule, error) {
	var values []Rule
	size, _ := or.GetArraySize()
	for i := 0; i < size; i++ {
		n := or.GetIndex(i)
		if !n.IsMap() {
			return nil, errors.New("not found expected map for or constraint element")
		}
		c, err := parseExpressionValue(variable, n, varGenerator)
		if err != nil {
			return nil, err
		}
		values = append(values, c)
	}
	return NewOr(false, values), nil
}

func parseImplicitAnd(data *y.Yaml, variable Variable, varGenerator *VarGenerator) (Rule, error) {
	var values []Rule
	propertyConstraints, _ := data.GetMapKeys()

	for _, pathString := range propertyConstraints {
		propertyPath, err := path.ParsePath(pathString)
		if err != nil {
			return nil, err
		}

		constraint := data.Get(pathString)
		if !constraint.IsMap() {
			return nil, errors.New("PropertyConstraint must be a map")
		}
		cs, err := ParseConstraint(propertyPath, variable, constraint, varGenerator)
		if err != nil {
			return nil, err
		}
		values = append(values, cs...)
	}

	return NewAnd(false, values), nil
}
