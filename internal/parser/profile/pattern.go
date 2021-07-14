package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
)

type PatternRule struct {
	AtomicStatement
	Argument string
}

func (r PatternRule) Clone() Rule {
	return PatternRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: r.Negated,
				Name:    r.Name,
			},
			Variable: r.Variable,
			Path:     r.Path,
		},
		Argument: r.Argument,
	}
}

func (r PatternRule) Negate() Rule {
	cloned := r.Clone()
	switch c := cloned.(type) {
	case PatternRule:
		c.Negated = !r.Negated
		return c
	}
	return cloned
}

func (r PatternRule) ValueHash() string {
	v := fmt.Sprintf("%s%s", r.Name, r.Argument)
	return internal.HashString(v)
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
