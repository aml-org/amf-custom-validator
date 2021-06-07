package profile

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/path"
)

type NestedExpression struct {
	BaseStatement
	Path   path.PropertyPath
	Parent Variable
	Child  Variable
	Value  Rule
}

func newNestedExpression(negated bool, parent Variable, path path.PropertyPath, varGenerator *VarGenerator) NestedExpression {
	nested := NestedExpression{
		BaseStatement: BaseStatement{
			Negated: negated,
			Name:    "nested",
		},
		Parent: parent,
		Child:  varGenerator.GenExpressionVar(ForAll, nil),
		Path:   path,
	}
	return nested
}

func (exp NestedExpression) Negate() Rule {
	negated := exp
	negated.Negated = !exp.Negated
	return negated
}

func (exp NestedExpression) String() string {
	var negation = ""
	if exp.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s[Nested(%s,%s,%s)] : %s\n  (\n  %s\n  )", exp.Child.String(), exp.Parent.Name, exp.Child.String(), exp.Path.Source(), negation, Indent(exp.Value.String()))
}
