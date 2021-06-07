package profile

import (
	"fmt"
	"github.com/aml-org/amfopa/internal"
	"github.com/aml-org/amfopa/internal/parser/path"
)

type CountQualifier int

const (
	Min CountQualifier = iota
	Max
)

type CountRule struct {
	AtomicStatement
	Qualifier CountQualifier
	Argument  int
}

func (r CountRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
}

func (r CountRule) ValueHash() string {
	v := fmt.Sprintf("%s%d", r.Name, r.Argument)
	return internal.HashString(v)
}

func (r CountRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s%s(%s,'%s',%d)", negation, r.Name, r.Variable.Name, r.Path.Source(), r.Argument)
}

func newCount(negated bool, qualifier CountQualifier, variable Variable, path path.PropertyPath, argument int) CountRule {
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
	}
}

func newMinCount(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Min, variable, path, argument)
	c.Name = "minCount"
	return c
}

func newMaxCount(negated bool, variable Variable, path path.PropertyPath, argument int) CountRule {
	c := newCount(negated, Max, variable, path, argument)
	c.Name = "maxCount"
	return c
}
