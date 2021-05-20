package profile

import (
	"fmt"
	"sort"
	"strings"
)

type AndRule ComplexStatement

type OrRule ComplexStatement

func (r AndRule) String() string {
	negation := ""
	if r.Negated {
		negation = "¬"
	}

	body := make([]string, len(r.Body))
	for i, v := range r.Body {
		body[i] = v.String()
	}
	sort.Strings(body)
	if len(body) > 1 {
		return fmt.Sprintf("%s(%s)", negation, strings.Join(body, " ∧ "))
	} else {
		return fmt.Sprintf("%s%s", negation, body[0])
	}
}

func (r OrRule) String() string {
	negation := ""
	if r.Negated {
		negation = "¬"
	}

	body := make([]string, len(r.Body))
	for i, v := range r.Body {
		body[i] = v.String()
	}
	sort.Strings(body)
	if len(r.Body) > 1 {
		return fmt.Sprintf("%s(%s)", negation, strings.Join(body, " ∨ "))
	} else {
		return fmt.Sprintf("%s%s", negation, body[0])
	}

}

func (r AndRule) Clone() Rule {
	return NewAnd(r.Negated, r.Body)
}

func (r OrRule) Clone() Rule {
	return NewOr(r.Negated, r.Body)
}

func (r AndRule) Negate() Rule {
	negatedBody := make([]Rule, len(r.Body))
	for i, br := range r.Body {
		negatedBody[i] = br.Negate()
	}

	return NewOr(!r.Negated, negatedBody)
}

func (r OrRule) Negate() Rule {
	negatedBody := make([]Rule, len(r.Body))
	for i, br := range r.Body {
		negatedBody[i] = br.Negate()
	}

	return NewAnd(!r.Negated, negatedBody)
}

func NewAnd(negated bool, body []Rule) AndRule {
	return AndRule{
		BaseStatement: BaseStatement{
			Negated: negated,
			Name:    "and",
		},
		Body: body,
	}
}

func NewOr(negated bool, body []Rule) OrRule {
	return OrRule{
		BaseStatement: BaseStatement{
			Negated: negated,
			Name:    "or",
		},
		Body: body,
	}
}
