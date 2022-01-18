package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateScalarSubSetRule(hasValue profile.ScalarSetRule, iriExpander *misc.IriExpander) []SimpleRegoResult {

	path := hasValue.Path
	var rego []string

	// Let's get the path computed and stored hasValue the inValuesVariable
	actualValuesVariable := profile.Genvar(fmt.Sprintf("%s_check", hasValue.Variable.Name))
	hasValuesVariable := profile.Genvar("hasValues")


	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, hasValue.Variable.Name, iriExpander)
	rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", actualValuesVariable, pathResult.rule, hasValue.Variable.Name))
	rego = append(rego, fmt.Sprintf("count(%s_array) != 0 # validation applies if property was defined", actualValuesVariable))
	rego_convert_to_string_set := "%s_string_set = {\n" +
								  "  mapped |\n" +
		                          "    original := %s_array[_]\n" +
		                          "    mapped := as_string(original)\n}" // cast value to string for matching with argument value
	rego = append(rego, fmt.Sprintf(rego_convert_to_string_set, actualValuesVariable, actualValuesVariable))

	rego = append(rego, fmt.Sprintf("%s = { \"%s\"}", hasValuesVariable, strings.Join(hasValue.Argument, "\",\"")))

	// assert that all hasValues are contained in actualValues
	if hasValue.Negated {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) == 0", hasValuesVariable, actualValuesVariable))
	} else {
		rego = append(rego, fmt.Sprintf("count(%s - %s_string_set) != 0", hasValuesVariable, actualValuesVariable))
	}

	// used for actual value in trace
	rego = append(rego, fmt.Sprintf("%s_quoted = [concat(\"\", [\"\\\"\", res, \"\\\"\"]) |  res := %s_string_set[_]]", actualValuesVariable, actualValuesVariable))
	rego = append(rego, fmt.Sprintf("%s_string = concat(\"\", [\"[\", concat(\", \",%s_quoted), \"]\"])", actualValuesVariable, actualValuesVariable))

	tracePath, err := hasValue.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: hasValue.Name,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		TraceNode:  hasValue.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"actual\": %s,\"expected\": \"%s\"", hasValue.Negated, fmt.Sprintf("%s_string", actualValuesVariable), hasValue.JSONValues()),
		),
		Variable: actualValuesVariable,
	}
	return []SimpleRegoResult{r}
}
