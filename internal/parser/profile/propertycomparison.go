package profile

import (
	"fmt"
	"github.com/aml-org/amfopa/internal"
	"github.com/aml-org/amfopa/internal/parser/path"
)

type PropertyComparisonRule struct {
	AtomicStatement
	Operator CardinalityOperation
	Argument path.PropertyPath
}

func (r PropertyComparisonRule) Clone() Rule {
	return PropertyComparisonRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: r.Negated,
				Name:    r.Name,
			},
			Variable: r.Variable,
			Path:     r.Path,
		},
		Operator: r.Operator,
		Argument: r.Argument,
	}
}

func (r PropertyComparisonRule) Negate() Rule {
	cloned := r.Clone()
	switch c := cloned.(type) {
	case PropertyComparisonRule:
		c.Negated = !r.Negated
		return c
	}
	return cloned
}

func (r PropertyComparisonRule) ValueHash() string {
	v := fmt.Sprintf("%s%s%d", r.Name, r.Operator, r.Argument)
	return internal.HashString(v)
}

func (r PropertyComparisonRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s(Property(%s,'%s') %s Property(%s,'%s'))", negation, r.Variable.Name, r.Path.Source(), r.Operator.String(), r.Variable.Name, r.Argument.Source())
}

func newPropertyComparison(negated bool, operator CardinalityOperation, variable Variable, path path.PropertyPath, argument path.PropertyPath) PropertyComparisonRule {
	return PropertyComparisonRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
			},
			Variable: variable,
			Path:     path,
		},
		Operator: operator,
		Argument: argument,
	}
}

func newLessThan(negated bool, variable Variable, path path.PropertyPath, argument path.PropertyPath) PropertyComparisonRule {
	r := newPropertyComparison(negated, LT, variable, path, argument)
	r.Name = "lessThan"
	return r
}

func newLessThanOrEquals(negated bool, variable Variable, path path.PropertyPath, argument path.PropertyPath) PropertyComparisonRule {
	r := newPropertyComparison(negated, LTEQ, variable, path, argument)
	r.Name = "lessThanOrEqualsTo"
	return r
}

func newEquals(negated bool, variable Variable, path path.PropertyPath, argument path.PropertyPath) PropertyComparisonRule {
	r := newPropertyComparison(negated, EQ, variable, path, argument)
	r.Name = "equalsTo"
	return r
}

func newDisjoint(negated bool, variable Variable, path path.PropertyPath, argument path.PropertyPath) PropertyComparisonRule {
	r := newPropertyComparison(negated, NEQ, variable, path, argument)
	r.Name = "disjointWith"
	return r
}

func newMoreThan(negated bool, variable Variable, path path.PropertyPath, argument path.PropertyPath) PropertyComparisonRule {
	r := newPropertyComparison(negated, GT, variable, path, argument)
	r.Name = "moreThan"
	return r
}

func newMoreThanOrEquals(negated bool, variable Variable, path path.PropertyPath, argument path.PropertyPath) PropertyComparisonRule {
	r := newPropertyComparison(negated, GTEQ, variable, path, argument)
	r.Name = "moreThanOrEqualsTo"
	return r
}
