package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

func GenerateNumericComparison(num profile.NumericRule) []SimpleRegoResult {
	switch num.Operation {
	case profile.GTEQ:
		return generateNumericRule(num, "minimumInclusive", ">=")
	case profile.GT:
		return generateNumericRule(num, "minimumExclusive", ">")
	case profile.LT:
		return generateNumericRule(num, "maximumExclusive", "<")
	case profile.LTEQ:
		return generateNumericRule(num, "maximumInclusive", "<=")
	default:
		panic(fmt.Sprintf("cannot generate unknown numeric constraint: %v", num))
	}
}

func generateNumericRule(num profile.NumericRule, rule string, op string) []SimpleRegoResult {
	path := num.Path
	var rego []string

	// Let's get the path computed and stored in the inValuesVariable
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GenerateNodeArray(path, num.Variable.Name, num.ValueHash())
	valueVariable := profile.Genvar("numeric_comparison")
	rego = append(rego, fmt.Sprintf("%s = %s with data.sourceNode as %s", valueVariable, pathResult.rule, num.Variable.Name))

	// Add the validation
	if num.Negated {
		i, errI := num.IntArgument()
		f, errF := num.FloatArgument()
		if errI == nil {
			rego = append(rego, fmt.Sprintf("%s %s %d", valueVariable, op, i))
		}

		if errF == nil {
			rego = append(rego, fmt.Sprintf("%s %s %f", valueVariable, op, f))
		}

	} else {
		i, errI := num.IntArgument()
		f, errF := num.FloatArgument()
		if errI == nil {
			rego = append(rego, fmt.Sprintf("not %s %s %d", valueVariable, op, i))
		}

		if errF == nil {
			rego = append(rego, fmt.Sprintf("not %s %s %f", valueVariable, op, f))
		}
	}

	r := SimpleRegoResult{
		Constraint: rule,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       num.Path.Source(),
		Value:      fmt.Sprintf("%s", valueVariable),
		Variable:   valueVariable,
		Trace:      fmt.Sprintf("value not matching rule %v", num.Argument),
	}
	return []SimpleRegoResult{r}
}
