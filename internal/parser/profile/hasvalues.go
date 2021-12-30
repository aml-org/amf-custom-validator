package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	"strings"
)

type HasValueRule struct {
	AtomicStatement
	Argument []string
}

func (r HasValueRule) Clone() Rule {
	return HasValueRule{
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

func (r HasValueRule) Negate() Rule {
	cloned := r.Clone()
	switch c := cloned.(type) {
	case HasValueRule:
		c.Negated = !r.Negated
		return c
	}
	return cloned
}

func (r HasValueRule) ValueHash() string {
	v := fmt.Sprintf("%s%s", r.Name, r.Argument)
	return internal.HashString(v)
}

func (r HasValueRule) JSONValues() string {
	var acc []string
	for _, v := range r.Argument {
		acc = append(acc, fmt.Sprintf("\\\"%s\\\"", v))
	}

	return fmt.Sprintf("[%s]", strings.Join(acc, ","))
}

func (r HasValueRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}
	return fmt.Sprintf("%s%s(%s,'%s',%s)", negation, r.Name, r.Variable.Name, r.Path.Source(), r.Argument)
}

func newHasValue(negated bool, variable Variable, path path.PropertyPath, argument string) HasValueRule {
	return HasValueRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    "hasValue",
			},
			Variable: variable,
			Path:     path,
		},
		Argument: []string{argument},
	}
}

func newHasValues(negated bool, variable Variable, path path.PropertyPath, argument []string) HasValueRule {
	return HasValueRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    "hasValues",
			},
			Variable: variable,
			Path:     path,
		},
		Argument: argument,
	}
}
