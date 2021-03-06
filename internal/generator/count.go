package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

// GenerateCount Generates the Rego code snippet for the rule, supports minCount/maxCount and minLength/maxLength
func GenerateCount(count profile.CountRule, iriExpander *misc.IriExpander) []SimpleRegoResult {
	return generateCountRule(count, obtainCondition(count.Qualifier), iriExpander)
}

func obtainCondition(qualifier profile.CountQualifier) string {
	switch qualifier {
	case profile.Min:
		return ">="
	case profile.Max:
		return "<="
	default:
		return "=="
	}
}

// Generates the rule using the 'count'  property from Rego
func generateCountRule(count profile.CountRule, condition string, iriExpander *misc.IriExpander) []SimpleRegoResult {
	path := count.Path
	rule := count.Name
	var rego []string

	// Let's get the path computed and stored in the arrayVariable
	arrayVariable := profile.Genvar("propValues")
	singleValueVariable := fmt.Sprintf("%s_elem", arrayVariable)

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertySet(path, count.Variable.Name, iriExpander)
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

	tracePath, err := count.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: rule,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"condition\":\"%s\",\"actual\": count(%s),\"expected\": %d", count.Negated, condition, targetValueVariable, count.Argument),
		),
		TraceNode: count.Variable.Name,
		Variable:  targetValueVariable,
	}
	return []SimpleRegoResult{r}
}
