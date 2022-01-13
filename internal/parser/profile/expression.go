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


func (exp Expression) Negate() Rule {
	negated := exp
	negated.Negated = !exp.Negated
	return negated
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
