package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
)

type RegoRule struct {
	AtomicStatement
	Message  string
	Argument string
}

func (r RegoRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
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
