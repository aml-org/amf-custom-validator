package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
	"strings"
)

func GenerateIn(in profile.InRule) []SimpleRegoResult {
	path := in.Path
	var rego []string

	// Let's get the path computed and stored in the inValuesVariable
	inValuesVariable := profile.Genvar("inValues")
	inValuesTestVariable := fmt.Sprintf("%s_check", in.Variable.Name)

	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, in.Variable.Name, "in_"+in.ValueHash())
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

	r := SimpleRegoResult{
		Constraint: "in",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       in.Path.Source(),
		Value:      inValuesTestVariable,
		Variable:   inValuesTestVariable,
		Trace:      fmt.Sprintf("Error with value %s and enumeration ['%s']", inValuesVariable, strings.Join(in.Argument, "','")),
	}
	return []SimpleRegoResult{r}
}
