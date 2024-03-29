package generator

import (
	"encoding/json"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GeneratePattern(pattern profile.PatternRule, iriExpander *misc.IriExpander) []SimpleRegoResult {
	path := pattern.Path
	var rego []string
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertySet(path, pattern.Variable.Name, iriExpander)
	checkVariable := profile.Genvar(fmt.Sprintf("%s_node", pathResult.rule))
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", checkVariable, pathResult.rule, pattern.Variable.Name))
	rego = append(rego, fmt.Sprintf("%s = %s_array[_]", checkVariable, checkVariable))
	// Add the validation
	if pattern.Negated {
		rego = append(rego, fmt.Sprintf("regex.match(`%s`,%s)", pattern.Argument, checkVariable))
	} else {
		rego = append(rego, fmt.Sprintf("not regex.match(`%s`,%s)", pattern.Argument, checkVariable))
	}

	tracePath, err := pattern.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}

	escapedArgumentStringByes, err := json.Marshal(pattern.Argument)
	if err != nil {
		escapedArgumentStringByes = []byte{}
	}
	escapedArgumentString := string(escapedArgumentStringByes)

	r := SimpleRegoResult{
		Constraint: "pattern",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		TraceNode:  pattern.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"expected\": %s,\"actual\": %s", pattern.Negated, escapedArgumentString, checkVariable)),
		Variable: checkVariable,
	}
	return []SimpleRegoResult{r}
}
