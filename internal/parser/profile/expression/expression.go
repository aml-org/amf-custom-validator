package expression

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"strings"
)

type Expression struct {
	statements.BaseStatement
	Variables    []statements.Variable
	Value        statements.Rule
	varGenerator statements.VarGenerator
}

func newExpression(negated bool, name string) Expression {
	return Expression{
		BaseStatement: statements.BaseStatement{
			Negated: negated,
			Name: name,
		},
		Variables:    make([]statements.Variable,0),
		varGenerator: statements.NewVarGenerator(),
	}
}

func (exp Expression) Clone() statements.Rule {
	return Expression{
		BaseStatement: statements.BaseStatement{
			Negated: exp.Negated,
			Name: exp.Name,
		},
		Variables: exp.Variables,
		varGenerator: exp.varGenerator,
	}
}

func (exp Expression) Negate() statements.Rule {
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
	varsText := make([]string, len(exp.Variables))

	for i,v := range exp.Variables {
		varsText[i] = negation +  " " + v.String()
	}

	return fmt.Sprintf("%s :%s", strings.Join(varsText, ","), exp.Value.String())
}

