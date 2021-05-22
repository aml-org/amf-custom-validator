package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

func GenerateDatatype(datatype profile.DatatypeRule) []SimpleRegoResult {
	path := datatype.Path
	var rego []string

	// Let's get the path computed and stored in the inValuesVariable
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GenerateNodeArray(path, datatype.Variable.Name, datatype.ValueHash())
	valueVariable := profile.Genvar("datatype_check")
	rego = append(rego, fmt.Sprintf("%s = %s with data.sourceNode as %s", valueVariable, pathResult.rule, datatype.Variable.Name))

	if datatype.Negated {
		rego = append(rego, fmt.Sprintf("check_datatype(%s,\"%s\")", valueVariable, datatype.Argument))
	} else {
		rego = append(rego, fmt.Sprintf("not check_datatype(%s,\"%s\")", valueVariable, datatype.Argument))
	}
	r := SimpleRegoResult{
		Constraint: "datatype",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       datatype.Path.Source(),
		Value:      valueVariable,
		Variable:   valueVariable,
		Trace:      fmt.Sprintf("Error with datatype %s", valueVariable),
	}
	return []SimpleRegoResult{r}
}
