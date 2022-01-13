package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
)

type CountQualifier int

const (
	Min CountQualifier = iota
	Max
	Exact
)

type TargetValue int

const (
	StringLength TargetValue = iota
	ItemsInArray
)

type CountRule struct {
	AtomicStatement
	Qualifier CountQualifier
	Target TargetValue
	Argument  int
}


func (r CountRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
}

func (r CountRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s%s(%s,'%s',%d)", negation, r.Name, r.Variable.Name, r.Path.Source(), r.Argument)
}

func newCount(negated bool, qualifier CountQualifier, target TargetValue, variable Variable, path path.PropertyPath, argument int) CountRule {
	return CountRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
			},
			Variable: variable,
			Path:     path,
		},
		Qualifier: qualifier,
		Argument:  argument,
		Target: target,
	}
}

func newMinCount(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Min, ItemsInArray, variable, path, argument)
	c.Name = "minCount"
	return c
}

func newMaxCount(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Max, ItemsInArray, variable, path, argument)
	c.Name = "maxCount"
	return c
}

func newExactCount(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Exact, ItemsInArray, variable, path, argument)
	c.Name = "exactCount"
	return c
}

func newMinLength(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Min, StringLength, variable, path, argument)
	c.Name = "minLength"
	return c
}

func newMaxLength(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Max, StringLength, variable, path, argument)
	c.Name = "maxLength"
	return c
}

func newExactLength(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Exact, StringLength, variable, path, argument)
	c.Name = "exactLength"
	return c
}
