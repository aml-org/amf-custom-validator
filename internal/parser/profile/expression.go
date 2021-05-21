package profile

import (
	"fmt"
)

type Expression struct {
	BaseStatement
	Variable *Variable
	Value    Rule
}

func newExpression(negated bool, name string, varGenerator *VarGenerator) Expression {
	newVar := varGenerator.GenExpressionVar(ForAll, nil)
	return Expression{
		BaseStatement: BaseStatement{
			Negated: negated,
			Name:    name,
		},
		Variable: &newVar,
	}
}

func (exp Expression) Clone() Rule {
	return Expression{
		BaseStatement: BaseStatement{
			Negated: exp.Negated,
			Name:    exp.Name,
		},
		Variable: exp.Variable,
	}
}

func (exp Expression) Negate() Rule {
	cloned := exp.Clone()
	switch expr := cloned.(type) {
	case Expression:
		expr.Negated = !exp.Negated
	}
	return cloned
}

func (exp Expression) String() string {
	var negation = ""
	if exp.Negated {
		negation = "¬"
	}
	varsText := ""
	if exp.Variable != nil {
		varsText = negation + exp.Variable.String()
	}

	return fmt.Sprintf("%s :\n%s", varsText, Indent(exp.Value.String()))
}
