package profile

import (
	"errors"
	"fmt"
	pathParser "github.com/aml-org/amfopa/internal/parser/path"
	y "github.com/aml-org/amfopa/internal/parser/yaml"
	"strconv"
)

func ParseConstraint(path pathParser.PropertyPath, variable Variable, constraint *y.Yaml, varGenerator *VarGenerator) ([]Rule, error) {
	var acc []Rule

	min, err := constraint.Get("minCount").Int()
	if err == nil {
		acc = append(acc, newMinCount(false, variable, path, min))
	}

	max, err := constraint.Get("maxCount").Int()
	if err == nil {
		acc = append(acc, newMaxCount(false, variable, path, max))
	}

	pattern, err := constraint.Get("pattern").String()
	if err == nil {
		acc = append(acc, newPattern(false, variable, path, pattern))
	}

	in, err := constraint.Get("in").Array()
	if err == nil {
		l, err := scalarList(in)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newIn(false, variable, path, l))
	}

	otherProp, err := constraint.Get("lessThanProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newLessThan(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("lessThanOrEqualsToProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newLessThanOrEquals(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("equalsToProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newEquals(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("disjointWithProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newDisjoint(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("moreThanProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newMoreThan(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("moreThanOrEqualsToProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newMoreThanOrEquals(false, variable, path, compPath))
	}

	atLeast := constraint.Get("atLeast")
	if atLeast.IsFound() {
		rule, err := parseQualifiedNestedExpression(atLeast, false, variable, path, varGenerator, GTEQ)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	atMost := constraint.Get("atMost")
	if atMost.IsFound() {
		rule, err := parseQualifiedNestedExpression(atMost, false, variable, path, varGenerator, LTEQ)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num := constraint.Get("minInclusive")
	if num.IsFound() {
		rule, err := parseMinInclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num = constraint.Get("minExclusive")
	if num.IsFound() {
		rule, err := parseMinExclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num = constraint.Get("maxInclusive")
	if num.IsFound() {
		rule, err := parseMaxInclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num = constraint.Get("maxExclusive")
	if num.IsFound() {
		rule, err := parseMaxExclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	dt := constraint.Get("datatype")
	if dt.IsFound() {
		rule, err := parseDatatype(false, variable, path, dt)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	nested := constraint.Get("nested")
	if nested.IsFound() && nested.IsMap() {
		rule, err := parseNestedExpression(nested, false, variable, path, varGenerator)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	code := constraint.Get("rego")
	if code.IsFound() {
		rule, err := ParseRego(code, false, variable, path)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	codeModule := constraint.Get("regoModule")
	if codeModule.IsFound() {
		rule, err := ParseRego(codeModule, false, variable, path)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}
	return acc, nil

}

func ParseRego(code *y.Yaml, negated bool, variable Variable, path pathParser.PropertyPath) (Rule, error) {
	var regoCode string
	message := "Violation in native Rego constraint"
	s, err := code.String()
	if err == nil {
		// code fragment directly embedded in the constraint
		regoCode = s
	}

	if code.IsMap() {
		// message and code version
		s, err := code.Get("code").String()
		if err != nil {
			return nil, err
		}
		regoCode = s
		m, err := code.Get("message").String()
		if err == nil {
			message = m
		}
	}

	return newRego(negated, variable, path, regoCode, message), nil
}

func parseQualifiedNestedExpression(qNested *y.Yaml, negated bool, variable Variable, propertyPath pathParser.PropertyPath, generator *VarGenerator, op CardinalityOperation) (Rule, error) {
	count, err := qNested.Get("count").Int()
	if err != nil {
		return nil, err
	}
	cardinality := VariableCardinality{
		Operator: op,
		Value:    count,
	}

	data := qNested.Get("validation")
	if !data.IsFound() || (data.IsFound() && !data.IsMap()) {
		return nil, errors.New("map containing a validation is required in qualified atLeast/atMost constraint")
	}

	nested, err := parseNestedExpression(data, negated, variable, propertyPath, generator)
	if err != nil {
		return nil, err
	}

	switch n := nested.(type) {
	case NestedExpression:
		n.Child.Cardinality = &cardinality
		n.Child.Quantification = Exists
		return n, nil
	}

	return nil, errors.New(fmt.Sprintf("expected nested expression to build quantified nested expression but got %v", nested))
}

// We are always stringifying the value to be able to compare it easily in Rego
func scalarList(in []*y.Yaml) ([]string, error) {
	var acc []string
	for _, e := range in {
		s, nok := e.String()
		if nok == nil {
			acc = append(acc, s)
			continue
		}

		i, nok := e.Int()
		if nok == nil {
			acc = append(acc, strconv.Itoa(i))
			continue
		}

		f, nok := e.Float()
		if nok == nil {
			acc = append(acc, strconv.FormatFloat(f, 'f', 6, 64))
			continue
		}

		b, nok := e.Bool()
		if nok == nil {
			acc = append(acc, strconv.FormatBool(b))
			continue
		}

		l, c := e.Pos()
		return nil, errors.New(fmt.Sprintf("expected scalars in 'in' constraint, found %v at [%d,%d]", e, l, c))

	}

	return acc, nil
}
