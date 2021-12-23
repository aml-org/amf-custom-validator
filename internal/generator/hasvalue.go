package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateHasValue(hasValue profile.HasValueRule) []SimpleRegoResult {

	path := hasValue.Path
	var rego []string

	// Let's get the path computed and stored hasValue the inValuesVariable
	actualValuesVariable := profile.Genvar(fmt.Sprintf("%s_check", hasValue.Variable.Name))

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, hasValue.Variable.Name, "in_"+hasValue.ValueHash())
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", actualValuesVariable, pathResult.rule, hasValue.Variable.Name))
	rego_convert_to_string_set := "%s_string_set = {\n" +
								  "  mapped |\n" +
		                          "    original := %s_array[_]\n" +
		                          "    mapped := as_string(original)\n}" // cast value to string for matching with argument value
	rego = append(rego, fmt.Sprintf(rego_convert_to_string_set, actualValuesVariable, actualValuesVariable))

	// used for actual value in trace
	rego = append(rego, fmt.Sprintf("%s_string = concat(\"\", [\"[ \", concat(\", \",%s_string_set), \" ]\"])", actualValuesVariable, actualValuesVariable))

	// assert that hasValue string value is contained within queried values
	if hasValue.Negated {
		rego = append(rego, fmt.Sprintf("%s_string_set[\"%s\"]", actualValuesVariable, hasValue.Argument))
	} else {
		rego = append(rego, fmt.Sprintf("not %s_string_set[\"%s\"]", actualValuesVariable, hasValue.Argument))
	}


	r := SimpleRegoResult{
		Constraint: "hasValue",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       hasValue.Path.Source(),
		TraceNode:  hasValue.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"actual\": %s,\"expected\": \"%s\"", hasValue.Negated, fmt.Sprintf("%s_string", actualValuesVariable), hasValue.Argument),
		),
		Variable: actualValuesVariable,
	}
	return []SimpleRegoResult{r}
}
