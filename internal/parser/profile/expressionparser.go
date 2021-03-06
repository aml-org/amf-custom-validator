package profile

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	y "github.com/aml-org/amf-custom-validator/internal/parser/yaml"
)

func ParseExpression(name string, data *y.Yaml, level string, varGenerator *VarGenerator) (Rule, error) {
	targetClass, err := data.Get("targetClass").String()
	if err != nil {
		l, c := data.Pos()
		return nil, errors.New(fmt.Sprintf("missing targetClass in validation definition at [%d,%d]", l, c))
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

	code := data.Get("rego")
	if code.IsFound() {
		return parseImplicitRego(code, variable)
	}

	codeModule := data.Get("regoModule")
	if codeModule.IsFound() {
		return parseImplicitRego(codeModule, variable)
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

	not := data.Get("not")
	if not.IsFound() {
		if not.IsMap() {
			return parseNot(not, variable, varGenerator)
		} else {
			return nil, errors.New("not constraint must be a mpa")
		}
	}

	ifContent := data.Get("if")
	if ifContent.IsFound() {
		thenContent := data.Get("then")
		elseContent := data.Get("else")
		if thenContent.IsFound() {
			return parseConditional(ifContent, thenContent, elseContent, variable, varGenerator)
		} else {
			l, c := ifContent.Pos()
			return nil, errors.New(fmt.Sprintf("Found if clause without then statement at [%d,%d]", l, c))
		}
	}

	l, c := data.Pos()
	return nil, errors.New(fmt.Sprintf("unknown expression node, cannot find properties to parse at [%d,%d]", l, c))

}

func parseImplicitRego(code *y.Yaml, variable Variable) (Rule, error) {
	return ParseRego(code, false, variable, path.NullPath{})
}

func parseNot(not *y.Yaml, variable Variable, varGenerator *VarGenerator) (Rule, error) {
	parsed, err := parseExpressionValue(variable, not, varGenerator)
	if err != nil {
		return nil, err
	}

	return parsed.Negate(), nil
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

func parseConditional(ifContent *y.Yaml, thenContent *y.Yaml, optElseContent *y.Yaml, variable Variable, generator *VarGenerator) (Rule, error) {
	ifRule, err := parseExpressionValue(variable, ifContent, generator)
	if err != nil {
		return nil, err
	}
	thenRule, err := parseExpressionValue(variable, thenContent, generator)
	if err != nil {
		return nil, err
	}
	if optElseContent.IsFound() {
		elseRule, err := parseExpressionValue(variable, optElseContent, generator)
		if err != nil {
			return nil, err
		}
		return NewIfThenElseConditional(false, ifRule, thenRule, elseRule), nil
	}
	return NewConditional(false, ifRule, thenRule), nil
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
