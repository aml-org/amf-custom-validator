package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	"strings"
)

type SetCriteria int

const (
	SuperSet SetCriteria = iota // in
	SubSet                      // containsAll
	InsersectSet                // containsSome
)

type ScalarSetRule struct {
	AtomicStatement
	Argument []string
	SetCriteria
}


func (r ScalarSetRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
}

func (r ScalarSetRule) JSONValues() string {
	var acc []string
	for _, v := range r.Argument {
		acc = append(acc, fmt.Sprintf("\\\"%s\\\"", v))
	}

	return fmt.Sprintf("[%s]", strings.Join(acc, ","))
}

func (r ScalarSetRule) String() string {
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

func newContainsAll(negated bool, variable Variable, path path.PropertyPath, argument []string) ScalarSetRule {
	return ScalarSetRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    "containsAll",
			},
			Variable: variable,
			Path:     path,
		},
		Argument: argument,
		SetCriteria: SubSet,
	}
}

func newContainsSome(negated bool, variable Variable, path path.PropertyPath, argument []string) ScalarSetRule {
	return ScalarSetRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    "containsSome",
			},
			Variable: variable,
			Path:     path,
		},
		Argument: argument,
		SetCriteria: InsersectSet,
	}
}

func newIn(negated bool, variable Variable, path path.PropertyPath, argument []string) ScalarSetRule {
	return ScalarSetRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    "in",
			},
			Variable: variable,
			Path:     path,
		},
		Argument: argument,
		SetCriteria: SuperSet,
	}
}
