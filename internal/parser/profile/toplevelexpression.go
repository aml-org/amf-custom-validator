package profile

import (
	"fmt"
)

type TopLevelExpression struct {
	Expression
	Message        string
	Level          string
	ClassGenerator string
}

func newTopLevelExpression(negated bool, name string, message string, level string, targetClass string, varGenerator *VarGenerator) TopLevelExpression {
	exp := TopLevelExpression{
		Expression:     newExpression(negated, name, varGenerator),
		Message:        message,
		Level:          level,
		ClassGenerator: targetClass,
	}
	return exp
}

func (exp TopLevelExpression) Negate() Rule {
	negated := exp
	negated.Negated = !exp.Negated
	return negated
}

func (exp TopLevelExpression) String() string {
	return fmt.Sprintf("%s[%s] :=  %s[Class(%s)] : \n%s", exp.Name, exp.Level, exp.Variable.String(), exp.ClassGenerator, Indent(exp.Value.String()))
}
