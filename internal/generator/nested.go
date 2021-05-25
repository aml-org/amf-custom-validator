package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

func GenerateNested(nested profile.NestedExpression) SimpleRegoResult {
	path := nested.Path
	var rego []string

	pathResult := GenerateNodeArray(nested.Path, nested.Parent.Name, "nested_")
	pluralName := fmt.Sprintf("%ss", nested.Child.Name)

	rego = append(rego, "#  querying path: "+path.Source())
	rego = append(rego, fmt.Sprintf("%s = %s with data.sourceNode as %s", pluralName, pathResult.rule, nested.Parent.Name))
	return SimpleRegoResult{
		Constraint: "nested",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       nested.Path.Source(),
		TraceValue: fmt.Sprintf("{}"),
		TraceNode:  nested.Parent.Name,
		Variable:   pluralName,
	}
}
