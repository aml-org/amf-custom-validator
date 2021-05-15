package expression

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"strings"
)

type TopLevelExpression struct {
	Expression
	Message string
	Level string
	ClassGenerator string
}

func newTopLevelExpression(negated bool, name string, message string, level string, targetClass string) TopLevelExpression {
	exp := TopLevelExpression{
		Expression:     newExpression(negated,name),
		Message:        message,
		Level:          level,
		ClassGenerator: targetClass,
	}
	return exp
}

func (exp TopLevelExpression) SubExpression(negated bool) Expression {
	subExp := newExpression(negated,exp.Name + "_sub")
	subExp.varGenerator = exp.varGenerator
	return subExp
}

func (exp *Expression) GenVar(quantification statements.Quantification, cardinality *statements.VariableCardinality) statements.Variable {
	variable := exp.varGenerator.GenExpressionVar(quantification, cardinality)
	exp.Variables = append(exp.Variables, variable)
	return variable
}

func (exp TopLevelExpression) Clone() statements.Rule {
	tl := newExpression(exp.Negated, exp.Name)
	tl.Variables = exp.Variables
	tl.varGenerator = exp.varGenerator
	cloned := TopLevelExpression{
		Expression: tl,
		Message: exp.Message,
		Level: exp.Level,
		ClassGenerator: exp.ClassGenerator,
	}
	return cloned
}


func (exp TopLevelExpression) Negate() statements.Rule {
	cloned := exp.Clone()
	switch tl := cloned.(type) {
	case TopLevelExpression:
		tl.Negated = !exp.Negated
	}
	return cloned
}

func (exp TopLevelExpression) String() string {
	negation := ""
	if exp.Negated {
		negation = "¬"
	}
	varsText := make([]string, len(exp.Variables))

	for i,v := range exp.Variables {
		varsText[i] = negation +  " " + v.String()
	}
	return fmt.Sprintf("%s[%s] := %s : %s", exp.Name, exp.Level, strings.Join(varsText, ","), exp.Value.String())
}

