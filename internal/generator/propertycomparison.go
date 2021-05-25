package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

func GeneratePropertyComparison(comparison profile.PropertyComparisonRule) []SimpleRegoResult {
	var rego []string
	path := comparison.Path
	altPath := comparison.Argument

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, comparison.Variable.Name, "prop_"+comparison.ValueHash())
	propVariable := fmt.Sprintf("%sA", pathResult.rule)
	rego = append(rego, fmt.Sprintf("%ss = %s with data.sourceNode as %s", propVariable, pathResult.rule, comparison.Variable.Name))
	rego = append(rego, "#  querying path: "+altPath.Source())
	altPathResult := GeneratePropertyArray(altPath, comparison.Variable.Name, "alt_prop_"+comparison.ValueHash())
	altPropVariable := fmt.Sprintf("%sB", altPathResult.rule)
	rego = append(rego, fmt.Sprintf("%ss = %s with data.sourceNode as %s", altPropVariable, altPathResult.rule, comparison.Variable.Name))

	// this will compute [a_i,b_j]
	rego = append(rego, fmt.Sprintf("%s = %ss[_]", propVariable, propVariable))
	rego = append(rego, fmt.Sprintf("%s = %ss[_]", altPropVariable, altPropVariable))

	if comparison.Negated {
		rego = append(rego, fmt.Sprintf("%s %s %s", propVariable, comparison.Operator.String(), altPropVariable))
	} else {
		rego = append(rego, fmt.Sprintf("not %s %s %s", propVariable, comparison.Operator.String(), altPropVariable))
	}

	r := SimpleRegoResult{
		Constraint: comparison.Name,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult, altPathResult},
		Path:       comparison.Path.Source(),
		Variable:   comparison.Variable.Name,
		TraceNode:  comparison.Variable.Name,
		TraceValue: fmt.Sprintf("{\"negated\":%t, \"condition\":\"%s\",\"expected\":%s, \"actual\":%s}", comparison.Negated, comparison.Operator.String(), propVariable, altPropVariable),
	}
	return []SimpleRegoResult{r}
}
