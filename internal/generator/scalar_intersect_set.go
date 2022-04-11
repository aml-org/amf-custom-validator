package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateScalarIntersectSetRule(containsSome profile.ScalarSetRule, iriExpander *misc.IriExpander) []SimpleRegoResult {

	path := containsSome.Path
	var rego []string

	// Let's get the path computed and stored containsSome the inValuesVariable
	actualValuesVariable := profile.Genvar(fmt.Sprintf("%s_check", containsSome.Variable.Name))
	containsSomeVariable := profile.Genvar("containsSome")

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertySet(path, containsSome.Variable.Name, iriExpander)
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", actualValuesVariable, pathResult.rule, containsSome.Variable.Name))
	rego = append(rego, fmt.Sprintf("count(%s_array) != 0 # validation applies if property was defined", actualValuesVariable))
	rego_convert_to_string_set := "%s_string_set = { mapped |\n" +
		"    original := %s_array[_]\n" +
		"    mapped := as_string(original)\n}\n" // cast value to string for matching with argument value
	rego = append(rego, fmt.Sprintf(rego_convert_to_string_set, actualValuesVariable, actualValuesVariable))

	rego = append(rego, fmt.Sprintf("%s = { \"%s\"}", containsSomeVariable, strings.Join(containsSome.Argument, "\",\"")))

	// assert that the difference between containsSome and actualValues is different from all containsSome
	if containsSome.Negated {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) != count(%s)", containsSomeVariable, actualValuesVariable, containsSomeVariable))
	} else {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) == count(%s)", containsSomeVariable, actualValuesVariable, containsSomeVariable))
	}

	// used for actual value in trace
	rego = append(rego, fmt.Sprintf("%s_quoted = [concat(\"\", [\"\\\"\", res, \"\\\"\"]) |  res := %s_string_set[_]]", actualValuesVariable, actualValuesVariable))
	rego = append(rego, fmt.Sprintf("%s_string = concat(\"\", [\"[\", concat(\", \",%s_quoted), \"]\"])", actualValuesVariable, actualValuesVariable))

	tracePath, err := containsSome.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: containsSome.Name,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		TraceNode:  containsSome.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"actual\": %s,\"expected\": \"%s\"", containsSome.Negated, fmt.Sprintf("%s_string", actualValuesVariable), containsSome.JSONValues()),
		),
		Variable: actualValuesVariable,
	}
	return []SimpleRegoResult{r}
}
