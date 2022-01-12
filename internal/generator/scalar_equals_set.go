package generator

import (
"fmt"
"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateScalarIsSetRule(isRule profile.ScalarSetRule) []SimpleRegoResult {
	path := isRule.Path
	var rego []string
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, isRule.Variable.Name, "is_"+isRule.ValueHash())
	checkVariable := profile.Genvar(fmt.Sprintf("%s_node", pathResult.rule))
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", checkVariable, pathResult.rule, isRule.Variable.Name))
	rego = append(rego, fmt.Sprintf("%s = %s_array[_]", checkVariable, checkVariable))
	// Add the validation
	if isRule.Negated {
		rego = append(rego, fmt.Sprintf("\"%s\" = as_string(%s)", isRule.Argument[0], checkVariable))
	} else {
		rego = append(rego, fmt.Sprintf("not \"%s\" = as_string(%s)", isRule.Argument[0], checkVariable))
	}

	r := SimpleRegoResult{
		Constraint: "isRule",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       isRule.Path.Source(),
		TraceNode:  isRule.Variable.Name,
		TraceValue: BuildTraceValueNode(fmt.Sprintf("\"negated\":%t,\"argument\": %s", isRule.Negated, checkVariable)),
		Variable:   checkVariable,
	}
	return []SimpleRegoResult{r}
}

