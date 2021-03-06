package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
)

type PatternRule struct {
	AtomicStatement
	Argument string
}

func (r PatternRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
}

func (r PatternRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s%s(%s,'%s','%s')", negation, r.Name, r.Variable.Name, r.Path.Source(), r.Argument)
}

func newPattern(negated bool, variable Variable, path path.PropertyPath, argument string) PatternRule {
	return PatternRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Name:    "pattern",
				Negated: negated,
			},
			Variable: variable,
			Path:     path,
		},
		Argument: argument,
	}
}
