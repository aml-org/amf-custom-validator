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

func newTopLevelExpression(negated bool, name string, message string, level string, targetClass string, varGenerator VarGenerator) TopLevelExpression {
	exp := TopLevelExpression{
		Expression:     newExpression(negated, name, varGenerator),
		Message:        message,
		Level:          level,
		ClassGenerator: targetClass,
	}
	return exp
}

func (exp TopLevelExpression) Clone() Rule {
	cloned := TopLevelExpression{
		Expression: Expression{
			BaseStatement: BaseStatement{
				Negated: exp.Negated,
				Name:    exp.Name,
			},
			Variable: exp.Variable,
		},
		Message:        exp.Message,
		Level:          exp.Level,
		ClassGenerator: exp.ClassGenerator,
	}
	return cloned
}

func (exp TopLevelExpression) Negate() Rule {
	cloned := exp.Clone()
	switch tl := cloned.(type) {
	case TopLevelExpression:
		tl.Negated = !exp.Negated
	}
	return cloned
}

func (exp TopLevelExpression) String() string {
	return fmt.Sprintf("%s[%s] :=  %s : %s", exp.Name, exp.Level, exp.Variable.String(), exp.Value.String())
}
