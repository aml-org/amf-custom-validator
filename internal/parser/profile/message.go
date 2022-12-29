package profile

import (
	"regexp"
	"strings"
)

type Message struct {
	Expression string
	Variables  []string
}

func (a Message) Compare(b Message) int {
	exprComp := strings.Compare(a.Expression, b.Expression)
	if exprComp == 0 {
		for i, variable := range a.Variables {
			varComp := strings.Compare(variable, b.Variables[i])
			if varComp != 0 {
				return varComp
			}
		}
		return 0
	} else {
		return exprComp
	}

}

func ParseMessageExpression(rawExpression string) Message {
	// result
	expression := rawExpression
	var variables []string

	// find variables
	re := regexp.MustCompile(`\{\{\s*([\w-]+\.[\w-]+)\s*}}`)
	for _, v := range re.FindAllStringSubmatch(rawExpression, -1) {
		expression = strings.ReplaceAll(expression, v[0], "%v") // replace variable by template string variable
		variables = append(variables, v[1])
	}

	return Message{
		expression,
		variables,
	}
}
