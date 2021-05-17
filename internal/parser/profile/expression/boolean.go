package expression

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"strings"
)

type AndRule statements.ComplexStatement

type OrRule statements.ComplexStatement

func (r AndRule) String() string {
	negation := ""
	if r.Negated {
		negation = "¬"
	}

	body := make([]string, len(r.Body))
	for i,v := range r.Body {
		body[i] = v.String()
	}

	return fmt.Sprintf("%s(%s)",negation,strings.Join(body, " ∧ "))
}

func (r OrRule) String() string {
	negation := ""
	if r.Negated {
		negation = "¬"
	}

	body := make([]string, len(r.Body))
	for i,v := range r.Body {
		body[i] = v.String()
	}

	return fmt.Sprintf("%s(%s)",negation,strings.Join(body, " ∨ "))

}

func (r AndRule) Clone() statements.Rule {
	return NewAnd(r.Negated, r.Body)
}

func (r OrRule) Clone() statements.Rule {
	return NewOr(r.Negated, r.Body)
}


func (r AndRule) Negate() statements.Rule {
	negatedBody := make([]statements.Rule, len(r.Body))
	for i,br := range r.Body {
		negatedBody[i] = br.Negate()
	}

	return NewOr(!r.Negated, negatedBody)
}

func (r OrRule) Negate() statements.Rule {
	negatedBody := make([]statements.Rule, len(r.Body))
	for i,br := range r.Body {
		negatedBody[i] = br.Negate()
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

func NewOr(negated bool, body []statements.Rule) OrRule {
	return OrRule{
		BaseStatement: statements.BaseStatement{
			Negated: negated,
			Name: "or",
		},
		Body: body,
	}
}