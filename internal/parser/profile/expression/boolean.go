package expression

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"strings"
)

type AndRule statements.ComplexStatement

type OrRule statements.ComplexStatement

func (a AndRule) String() string {
	negation := ""
	if a.Negated {
		negation = "¬"
	}

	body := make([]string, len(a.Body))
	for i,v := range a.Body {
		body[i] = v.String()
	}

	return fmt.Sprintf("%s(%s)",negation,strings.Join(body, " ∧ "))
}

func (a OrRule) String() string {
	negation := ""
	if a.Negated {
		negation = "¬"
	}

	body := make([]string, len(a.Body))
	for i,v := range a.Body {
		body[i] = v.String()
	}

	return fmt.Sprintf("%s(%s)",negation,strings.Join(body, " ∨ "))

}

func (a AndRule) Clone() statements.Rule {
	return NewAnd(a.Negated, a.Body)
}

func (a OrRule) Clone() statements.Rule {
	return NewOr(a.Negated, a.Body)
}


func (r AndRule) Negate() statements.Rule {
	negatedBody := make([]statements.Rule, len(r.Body))
	for _,br := range r.Body {
		negatedBody = append(negatedBody, br.Negate())
	}

	return NewOr(!r.Negated, negatedBody)
}

func (r OrRule) Negate() statements.Rule {
	negatedBody := make([]statements.Rule, len(r.Body))
	for _,br := range r.Body {
		negatedBody = append(negatedBody, br.Negate())
	}

	return NewAnd(!r.Negated, negatedBody)
}

func NewAnd(negated bool, body []statements.Rule) AndRule {
	return AndRule{
		BaseStatement: statements.BaseStatement{
			Negated: negated,
			Name: "and",
		},
		Body: body,
	}
}

func NewOr(negated bool, body []statements.Rule) AndRule {
	return AndRule{
		BaseStatement: statements.BaseStatement{
			Negated: negated,
			Name: "or",
		},
		Body: body,
	}
}