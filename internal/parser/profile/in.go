package profile

import (
	"fmt"
	"github.com/aml-org/amfopa/internal"
	"github.com/aml-org/amfopa/internal/parser/path"
	"strings"
)

type InRule struct {
	AtomicStatement
	Argument []string
}

func (r InRule) Clone() Rule {
	return InRule{
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

func (r InRule) Negate() Rule {
	cloned := r.Clone()
	switch c := cloned.(type) {
	case InRule:
		c.Negated = !r.Negated
		return c
	}
	return cloned
}

func (r InRule) ValueHash() string {
	v := fmt.Sprintf("%s%s", r.Name, strings.Join(r.Argument, "-"))
	return internal.HashString(v)
}

func (r InRule) JSONValues() string {
	var acc []string
	for _, v := range r.Argument {
		acc = append(acc, fmt.Sprintf("\\\"%s\\\"", v))
	}

	return fmt.Sprintf("[%s]", strings.Join(acc, ","))
}

func (r InRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}
	var acc []string
	for _, a := range r.Argument {
		acc = append(acc, fmt.Sprintf("%v", a))
	}
	return fmt.Sprintf("%s%s(%s,'%s',%s)", negation, r.Name, r.Variable.Name, r.Path.Source(), strings.Join(acc, ","))
}

func newIn(negated bool, variable Variable, path path.PropertyPath, argument []string) InRule {
	return InRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    "in",
			},
			Variable: variable,
			Path:     path,
		},
		Argument: argument,
	}
}
