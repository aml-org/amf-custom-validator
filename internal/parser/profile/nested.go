package profile

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/path"
)

type NestedExpression struct {
	BaseStatement
	path   path.PropertyPath
	parent Variable
	child  Variable
	value  Rule
}

func newNestedExpression(negated bool, parent Variable, path path.PropertyPath, varGenerator VarGenerator) NestedExpression {
	nested := NestedExpression{
		BaseStatement: BaseStatement{
			Negated: negated,
			Name:    "nested",
		},
		parent: parent,
		child:  varGenerator.GenExpressionVar(ForAll, nil),
		path:   path,
	}
	return nested
}

func (exp NestedExpression) Clone() Rule {
	return NestedExpression{
		BaseStatement: BaseStatement{
			Negated: exp.Negated,
			Name:    exp.Name,
		},
		path:   exp.path,
		parent: exp.parent,
		child:  exp.child,
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

	return fmt.Sprintf("%sNested(%s,%s,'%s')", negation, exp.parent, exp.child, exp.path.Source())
}
