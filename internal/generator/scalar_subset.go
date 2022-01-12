package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateScalarSubSetRule(containsAll profile.ScalarSetRule, iriExpander *misc.IriExpander) []SimpleRegoResult {

	path := containsAll.Path
	var rego []string

	// Let's get the path computed and stored containsAll the inValuesVariable
	actualValuesVariable := profile.Genvar(fmt.Sprintf("%s_check", containsAll.Variable.Name))
	containsAllVariable := profile.Genvar("containsAll")


	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, containsAll.Variable.Name, iriExpander)
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", actualValuesVariable, pathResult.rule, containsAll.Variable.Name))
	rego = append(rego, fmt.Sprintf("count(%s_array) != 0 # validation applies if property was defined", actualValuesVariable))
	rego_convert_to_string_set := "%s_string_set = { mapped |\n" +
		                          "    original := %s_array[_]\n" +
		                          "    mapped := as_string(original)\n}\n" // cast value to string for matching with argument value
	rego = append(rego, fmt.Sprintf(rego_convert_to_string_set, actualValuesVariable, actualValuesVariable))

	rego = append(rego, fmt.Sprintf("%s = { \"%s\"}", containsAllVariable, strings.Join(containsAll.Argument, "\",\"")))

	// assert that all containsAll are contained in actualValues
	if containsAll.Negated {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) == 0", containsAllVariable, actualValuesVariable))
	} else {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) != 0", containsAllVariable, actualValuesVariable))
	}

	// used for actual value in trace
	rego = append(rego, fmt.Sprintf("%s_quoted = [concat(\"\", [\"\\\"\", res, \"\\\"\"]) |  res := %s_string_set[_]]", actualValuesVariable, actualValuesVariable))
	rego = append(rego, fmt.Sprintf("%s_string = concat(\"\", [\"[\", concat(\", \",%s_quoted), \"]\"])", actualValuesVariable, actualValuesVariable))

	tracePath, err := containsAll.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: containsAll.Name,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		TraceNode:  containsAll.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"actual\": %s,\"expected\": \"%s\"", containsAll.Negated, fmt.Sprintf("%s_string", actualValuesVariable), containsAll.JSONValues()),
		),
		Variable: actualValuesVariable,
	}
	return []SimpleRegoResult{r}
}
