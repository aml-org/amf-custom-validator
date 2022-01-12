package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateScalarSubSetRule(hasValue profile.ScalarSetRule) []SimpleRegoResult {

	path := hasValue.Path
	var rego []string

	// Let's get the path computed and stored hasValue the inValuesVariable
	actualValuesVariable := profile.Genvar(fmt.Sprintf("%s_check", hasValue.Variable.Name))
	containsAllVariable := profile.Genvar("containsAll")


	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, hasValue.Variable.Name, "in_"+hasValue.ValueHash())
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", actualValuesVariable, pathResult.rule, hasValue.Variable.Name))
	rego = append(rego, fmt.Sprintf("count(%s_array) != 0 # validation applies if property was defined", actualValuesVariable))
	rego_convert_to_string_set := "%s_string_set = { mapped |\n" +
		                          "    original := %s_array[_]\n" +
		                          "    mapped := as_string(original)\n}\n" // cast value to string for matching with argument value
	rego = append(rego, fmt.Sprintf(rego_convert_to_string_set, actualValuesVariable, actualValuesVariable))

	rego = append(rego, fmt.Sprintf("%s = { \"%s\"}", containsAllVariable, strings.Join(hasValue.Argument, "\",\"")))

	// assert that all containsAll are contained in actualValues
	if hasValue.Negated {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) == 0", containsAllVariable, actualValuesVariable))
	} else {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) != 0", containsAllVariable, actualValuesVariable))
	}

	// used for actual value in trace
	rego = append(rego, fmt.Sprintf("%s_quoted = [concat(\"\", [\"\\\"\", res, \"\\\"\"]) |  res := %s_string_set[_]]", actualValuesVariable, actualValuesVariable))
	rego = append(rego, fmt.Sprintf("%s_string = concat(\"\", [\"[\", concat(\", \",%s_quoted), \"]\"])", actualValuesVariable, actualValuesVariable))

	r := SimpleRegoResult{
		Constraint: hasValue.Name,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       hasValue.Path.Source(),
		TraceNode:  hasValue.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"actual\": %s,\"expected\": \"%s\"", hasValue.Negated, fmt.Sprintf("%s_string", actualValuesVariable), hasValue.JSONValues()),
		),
		Variable: actualValuesVariable,
	}
	return []SimpleRegoResult{r}
}
