package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

func GeneratePattern(pattern profile.PatternRule) []SimpleRegoResult {
	path := pattern.Path
	var rego []string
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, pattern.Variable.Name, "pattern_"+pattern.ValueHash())
	checkVariable := profile.Genvar(fmt.Sprintf("%s_node", pathResult.rule))
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", checkVariable, pathResult.rule, pattern.Variable.Name))
	rego = append(rego, fmt.Sprintf("%s = %s_array[_]", checkVariable, checkVariable))
	// Add the validation
	if pattern.Negated {
		rego = append(rego, fmt.Sprintf("regex.match(\"%s\",%s)", pattern.Argument, checkVariable))
	} else {
		rego = append(rego, fmt.Sprintf("not regex.match(\"%s\",%s)", pattern.Argument, checkVariable))
	}

	r := SimpleRegoResult{
		Constraint: "pattern",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       pattern.Path.Source(),
		TraceNode:  pattern.Variable.Name,
		TraceValue: fmt.Sprintf("{\"negated\":%t,\"argument\": %s}", pattern.Negated, checkVariable),
		Variable:   checkVariable,
	}
	return []SimpleRegoResult{r}
}
