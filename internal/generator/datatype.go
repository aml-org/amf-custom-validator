package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateDatatype(datatype profile.DatatypeRule, iriExpander *misc.IriExpander) []SimpleRegoResult {
	path := datatype.Path
	var rego []string

	// Let's get the path computed and stored in the inValuesVariable
	rego = append(rego, "#  querying path: "+path.Source())
	pathResult := GeneratePropertyArray(path, datatype.Variable.Name, iriExpander)
	valueVariable := profile.Genvar("datatype_check")
	rego = append(rego, fmt.Sprintf("%s_elem = %s with data.sourceNode as %s", valueVariable, pathResult.rule, datatype.Variable.Name))
	rego = append(rego, fmt.Sprintf("%s = %s_elem[_]", valueVariable, valueVariable))
	var datatypeIri string
	if iriExpander != nil {
		datatypeIriExpanded, err := iriExpander.Expand(datatype.Argument)
		if err != nil {
			panic(err)
		}
		datatypeIri = datatypeIriExpanded
	} else {
		datatypeIri = datatype.Argument
	}
	if datatype.Negated {
		rego = append(rego, fmt.Sprintf("check_datatype(%s,\"%s\")", valueVariable, datatypeIri))
	} else {
		rego = append(rego, fmt.Sprintf("not check_datatype(%s,\"%s\")", valueVariable, datatypeIri))
	}
	tracePath, err := datatype.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: "datatype",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       tracePath,
		TraceNode:  datatype.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"actual\": %s,\"expected\": \"%s\"", datatype.Negated, valueVariable, datatypeIri)),
		Variable: valueVariable,
	}
	return []SimpleRegoResult{r}
}
