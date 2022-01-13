package profile

import "fmt"

type ConditionalRule ComplexStatement

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
	} else if len(body) == 3 { // defines else rule
		return fmt.Sprintf("%s(\n(\n%s\n→\n%s\n)\n∧\n(\n¬%s\n→\n%s\n)\n)", negation, Indent(body[0]), Indent(body[1]),Indent(body[0]), Indent(body[2]))
	} else {
		panic("Conditional expression must have if/then rules")
	}
}

func (r ConditionalRule) Negate() Rule {
	return NewConditional(!r.Negated, r.IfRule(), r.ThenRule())
}


// ifRule -> thenRule <==> ¬ifRule ∨ thenRule
func (r ConditionalRule) ThenMaterialImplication() OrRule {
	return NewOr(r.Negated, []Rule{
		r.IfRule().Negate(),
		r.ThenRule(),
	})
}

// ¬ifRule -> elseRule <==> ifRule ∨ thenRule
func (r ConditionalRule) ElseMaterialImplication() OrRule {
	return NewOr(r.Negated, []Rule{
		r.IfRule(),
		r.ElseRule(),
	})
}

func (r ConditionalRule) IfRule() Rule {
	return r.Body[0]
}

func (r ConditionalRule) ThenRule() Rule {
	return r.Body[1]
}

func (r ConditionalRule) ElseRule() Rule {
	return r.Body[2]
}

func (r ConditionalRule) ElseIsDefined() bool {
	return r.Body.Len() > 2
}

func newConditional(negated bool, rules *[]Rule) ConditionalRule {
	return ConditionalRule{
		BaseStatement: BaseStatement{
			Negated: negated,
			Name:    "conditional",
		},
		Body: *rules,
	}
}

func NewIfThenElseConditional(negated bool, ifRule Rule, thenRule Rule, elseRule Rule) ConditionalRule {
	return newConditional(negated, &[]Rule{ifRule, thenRule, elseRule})
}

func NewConditional(negated bool, ifRule Rule, thenRule Rule) ConditionalRule {
	return newConditional(negated, &[]Rule{ifRule, thenRule})
}
