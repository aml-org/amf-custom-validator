package profile

import (
	"fmt"
	"github.com/aml-org/amfopa/internal"
	"github.com/aml-org/amfopa/internal/parser/path"
)

type RegoRule struct {
	AtomicStatement
	Message  string
	Argument string
}

func (r RegoRule) Clone() Rule {
	return RegoRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: r.Negated,
				Name:    r.Name,
			},
			Variable: r.Variable,
			Path:     r.Path,
		},
		Message:  r.Message,
		Argument: r.Argument,
	}
}

func (r RegoRule) Negate() Rule {
	cloned := r.Clone()
	switch c := cloned.(type) {
	case RegoRule:
		c.Negated = !r.Negated
		return c
	}
	return cloned
}

func (r RegoRule) ValueHash() string {
	v := fmt.Sprintf("%s", r.Name)
	return internal.HashString(v)
}

func (r RegoRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s%s(%s,'%s','%s')", negation, r.Name, r.Variable.Name, r.Path.Source(), r.Message)
}

func newRego(negated bool, variable Variable, path path.PropertyPath, code string, message string) RegoRule {
	return RegoRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Name:    "rego",
				Negated: negated,
			},
			Variable: variable,
			Path:     path,
		},
		Message:  message,
		Argument: code,
	}
}