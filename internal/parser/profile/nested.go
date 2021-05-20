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

func (exp NestedExpression) Clone() Rule {
	return NestedExpression{
		BaseStatement: BaseStatement{
			Negated: exp.Negated,
			Name:    exp.Name,
		},
		Path:   exp.Path,
		Parent: exp.Parent,
		Child:  exp.Child,
	}
}

func (exp NestedExpression) Negate() Rule {
	cloned := exp.Clone()
	switch expr := cloned.(type) {
	case NestedExpression:
		expr.Negated = !exp.Negated
	}
	return cloned
}

func (exp NestedExpression) String() string {
	var negation = ""
	if exp.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s[Nested(%s,%s)] : %s(%s)", exp.Child.String(), exp.Parent.Name, exp.Path.Source(), negation, exp.Value.String())
}
