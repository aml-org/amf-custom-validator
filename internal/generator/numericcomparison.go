package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateNumericComparison(num profile.NumericRule, iriExpander *misc.IriExpander) []SimpleRegoResult {
	switch num.Operation {
	case profile.GTEQ:
		return generateNumericRule(num, "minimumInclusive", ">=", iriExpander)
	case profile.GT:
		return generateNumericRule(num, "minimumExclusive", ">", iriExpander)
	case profile.LT:
		return generateNumericRule(num, "maximumExclusive", "<", iriExpander)
	case profile.LTEQ:
		return generateNumericRule(num, "maximumInclusive", "<=", iriExpander)
	default:
		panic(fmt.Sprintf("cannot generate unknown numeric constraint: %v", num))
	}
}

func generateNumericRule(num profile.NumericRule, rule string, op string, iriExpander *misc.IriExpander) []SimpleRegoResult {
	path := num.Path
	var rego []string

	// Let's get the path computed and stored in the inValuesVariable
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, num.Variable.Name, iriExpander)
	valueVariable := profile.Genvar("numeric_comparison")
	rego = append(rego, fmt.Sprintf("%s_elem = %s with data.sourceNode as %s", valueVariable, pathResult.rule, num.Variable.Name))
	rego = append(rego, fmt.Sprintf("%s = %s_elem[_]", valueVariable, valueVariable))
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

	tracePath, err := num.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: rule,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		Variable:   valueVariable,
		TraceNode:  num.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"condition\":\"%s\",\"expected\":%s,\"actual\":%s", num.Negated, op, num.StringArgument(), valueVariable)),
	}
	return []SimpleRegoResult{r}
}
