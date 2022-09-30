package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GeneratePropertyComparison(comparison profile.PropertyComparisonRule, iriExpander *misc.IriExpander) []SimpleRegoResult {
	var rego []string
	path := comparison.Path
	altPath := comparison.Argument

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertySet(path, comparison.Variable.Name, iriExpander)
	propVariable := fmt.Sprintf("%sA", pathResult.rule)
	rego = append(rego, fmt.Sprintf("%ss = %s with data.sourceNode as %s", propVariable, pathResult.rule, comparison.Variable.Name))
	rego = append(rego, "#  querying path: "+altPath.Source())
	altPathResult := GeneratePropertySet(altPath, comparison.Variable.Name, iriExpander)
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

	tracePath, err := comparison.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: comparison.Name,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult, altPathResult},
		Path:       tracePath,
		Variable:   comparison.Variable.Name,
		TraceNode:  comparison.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t, \"condition\":\"%s\",\"expected\":%s, \"actual\":%s, \"altPath\": \"%s\"", comparison.Negated, comparison.Operator.String(), altPropVariable, propVariable, altPath.Source())),
	}
	return []SimpleRegoResult{r}
}
