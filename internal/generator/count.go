package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

// Generates the Rego code snippet for the rule, supports minCount and maxCount
func GenerateCount(count profile.CountRule) []SimpleRegoResult {
	if count.Qualifier == profile.Min {
		return generateRule(count, "minCount", ">=")
	} else if count.Qualifier == profile.Max {
		return generateRule(count, "maxCount", "<=")
	} else {
		return generateRule(count, "exact", "==")
	}
}

// Generates the rule using the 'count'  property from Rego
func generateRule(count profile.CountRule, rule string, condition string) []SimpleRegoResult {
	path := count.Path
	var rego []string

	// Let's get the path computed and stored in the inValuesVariable
	inValuesVariable := profile.Genvar("propValues")
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, count.Variable.Name, rule+"_"+count.ValueHash())
	rego = append(rego, fmt.Sprintf("%s = %s with data.sourceNode as %s", inValuesVariable, pathResult.rule, count.Variable.Name))

	// Add the validation
	if count.Negated {
		rego = append(rego, fmt.Sprintf("count(%s) %s %d", inValuesVariable, condition, count.Argument))
	} else {
		rego = append(rego, fmt.Sprintf("not count(%s) %s %d", inValuesVariable, condition, count.Argument))
	}

	r := SimpleRegoResult{
		Constraint: rule,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       count.Path.Source(),
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"condition\":\"%s\",\"actual\": count(%s),\"expected\": %d", count.Negated, condition, inValuesVariable, count.Argument),
		),
		TraceNode:  count.Variable.Name,
		Variable:   inValuesVariable,
	}
	return []SimpleRegoResult{r}
}
