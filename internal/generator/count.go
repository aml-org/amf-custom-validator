package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

// Generates the Rego code snippet for the rule, supports minCount/maxCount and minLength/maxLength
func GenerateCount(count profile.CountRule) []SimpleRegoResult {
	return generateCountRule(count, obtainCondition(count.Qualifier))
}

func obtainCondition(qualifier profile.CountQualifier) string {
	switch qualifier {
	case profile.Min: return ">="
	case profile.Max: return "<="
	default: return "=="
	}
}

// Generates the rule using the 'count'  property from Rego
func generateCountRule(count profile.CountRule, condition string) []SimpleRegoResult {
	path := count.Path
	rule := count.Name
	var rego []string

	// Let's get the path computed and stored in the arrayVariable
	arrayVariable := profile.Genvar("propValues")
	singleValueVariable := fmt.Sprintf("%s_elem", arrayVariable)

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, count.Variable.Name, rule+"_"+count.ValueHash())
	rego = append(rego, fmt.Sprintf("%s = %s with data.sourceNode as %s", arrayVariable, pathResult.rule, count.Variable.Name))

	var targetValueVariable string
	if count.Target == profile.ItemsInArray {
		targetValueVariable = arrayVariable
	} else { // profile.StringLength
		targetValueVariable = singleValueVariable
		rego = append(rego, fmt.Sprintf("%s = %s[_]", singleValueVariable, arrayVariable))
	}


	// Add the validation
	if count.Negated {
		rego = append(rego, fmt.Sprintf("count(%s) %s %d", targetValueVariable, condition, count.Argument))
	} else {
		rego = append(rego, fmt.Sprintf("not count(%s) %s %d", targetValueVariable, condition, count.Argument))
	}

	r := SimpleRegoResult{
		Constraint: rule,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       count.Path.Source(),
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"condition\":\"%s\",\"actual\": count(%s),\"expected\": %d", count.Negated, condition, targetValueVariable, count.Argument),
		),
		TraceNode: count.Variable.Name,
		Variable:  targetValueVariable,
	}
	return []SimpleRegoResult{r}
}
