package profile

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	y "github.com/aml-org/amf-custom-validator/internal/parser/yaml"
)

type NumericRule struct {
	AtomicStatement
	Operation CardinalityOperation
	Argument  *y.Yaml
}

func (r NumericRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
}

func (r NumericRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	i, err := r.Argument.Int()
	if err == nil {
		return fmt.Sprintf("%s%s(%s,'%s',%d)", negation, r.Name, r.Variable.Name, r.Path.Source(), i)
	}
	f, _ := r.Argument.Float()
	return fmt.Sprintf("%s%s(%s,'%s',%f)", negation, r.Name, r.Variable.Name, r.Path.Source(), f)
}

func (r NumericRule) IntArgument() (int, error) {
	return r.Argument.Int()
}

func (r NumericRule) FloatArgument() (float64, error) {
	return r.Argument.Float()
}

func (r NumericRule) StringArgument() string {
	i, err := r.Argument.Int()
	if err == nil {
		return fmt.Sprintf("%d", i)
	}
	f, _ := r.Argument.Float()
	return fmt.Sprintf("%f", f)
}

func newNumericComparison(negated bool, name string, operation CardinalityOperation, variable Variable, path path.PropertyPath, argument *y.Yaml) (NumericRule, error) {
	n := NumericRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    name,
			},
			Variable: variable,
			Path:     path,
		},
		Operation: operation,
		Argument:  argument,
	}

	_, err1 := n.IntArgument()
	_, err2 := n.FloatArgument()
	if err1 != nil && err2 != nil {
		return n, errors.New(fmt.Sprintf("expected float or int argument for numeric comparison, found %v", argument))
	}

	return n, nil
}

func parseMinInclusive(negated bool, variable Variable, path path.PropertyPath, argument *y.Yaml) (NumericRule, error) {
	return newNumericComparison(negated, "minInclusive", GTEQ, variable, path, argument)
}

func parseMaxInclusive(negated bool, variable Variable, path path.PropertyPath, argument *y.Yaml) (NumericRule, error) {
	return newNumericComparison(negated, "maxInclusive", LTEQ, variable, path, argument)
}

func parseMinExclusive(negated bool, variable Variable, path path.PropertyPath, argument *y.Yaml) (NumericRule, error) {
	return newNumericComparison(negated, "minExclusive", GT, variable, path, argument)
}

func parseMaxExclusive(negated bool, variable Variable, path path.PropertyPath, argument *y.Yaml) (NumericRule, error) {
	return newNumericComparison(negated, "maxExclusive", LT, variable, path, argument)
}
