package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateScalarSuperSetRule(in profile.ScalarSetRule, iriExpander *misc.IriExpander) []SimpleRegoResult {

	path := in.Path
	var rego []string

	// Let's get the path computed and stored in the inValuesVariable
	inValuesVariable := profile.Genvar("inValues")
	inValuesTestVariable := profile.Genvar(fmt.Sprintf("%s_check", in.Variable.Name))

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertySet(path, in.Variable.Name, iriExpander)
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", inValuesTestVariable, pathResult.rule, in.Variable.Name))
	rego = append(rego, fmt.Sprintf("%s_scalar = %s_array[_]", inValuesTestVariable, inValuesTestVariable))
	rego = append(rego, fmt.Sprintf("%s = as_string(%s_scalar)", inValuesTestVariable, inValuesTestVariable))
	rego = append(rego, fmt.Sprintf("%s = { \"%s\"}", inValuesVariable, strings.Join(in.Argument, "\",\"")))
	// Add the validation
	if in.Negated {
		rego = append(rego, fmt.Sprintf("%s[%s]", inValuesVariable, inValuesTestVariable))
	} else {
		rego = append(rego, fmt.Sprintf("not %s[%s]", inValuesVariable, inValuesTestVariable))
	}

	tracePath, err := in.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: in.Name,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		TraceNode:  in.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"actual\": %s,\"expected\": \"%s\"", in.Negated, strings.ReplaceAll(inValuesTestVariable, "\"", "'"), in.JSONValues()),
		),
		Variable: inValuesTestVariable,
	}
	return []SimpleRegoResult{r}
}
