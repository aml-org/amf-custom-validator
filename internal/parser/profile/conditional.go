package profile

import "fmt"

type ConditionalRule ComplexStatement

func (r ConditionalRule) Head() Rule {
	return r.Body[0]
}

func (r ConditionalRule) Tail() Rule {
	return r.Body[1]
}

func (r ConditionalRule) String() string {
	negation := ""
	if r.Negated {
		negation = "¬"
	}

	body := make([]string, len(r.Body))
	for i, v := range r.Body {
		body[i] = v.String()
	}
	if len(body) == 2 {
		return fmt.Sprintf("%s(\n%s\n→\n%s\n)", negation, Indent(body[0]), Indent(body[1]))
	} else {
		panic("Conditional expression must have head and tail")
	}
}

func (r ConditionalRule) Negate() Rule {
	return NewConditional(!r.Negated, r.Head(), r.Tail())
}

func (r ConditionalRule) MaterialImplication() OrRule {
	return NewOr(r.Negated, []Rule{
		r.Head().Negate(),
		r.Tail(),
	})
}

func NewConditional(negated bool, head Rule, tail Rule) ConditionalRule {
	return ConditionalRule{
		BaseStatement: BaseStatement{
			Negated: negated,
			Name:    "conditional",
		},
		Body: []Rule{head, tail},
	}
}
