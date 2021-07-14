package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateRegoRule(rule profile.RegoRule) []SimpleRegoResult {
	path := rule.Path
	var rego []string

	// By default we use the expression variable, this is the right one if the rego rule is top-level
	checkVariable := rule.Variable.Name

	// let's try generate the path rule for the constraint.
	// This can be a null path if is an inline-rego block at the top-level of a validation
	pathResult := GeneratePropertyArray(path, rule.Variable.Name, "code_"+rule.ValueHash())

	// If this is not a top-level rego rule (the path generates code), we use the bind the check variable for the path computation result
	if len(pathResult.rego) > 0 {
		checkVariable = profile.Genvar(pathResult.rule + "_node")
		rego = append(rego, "#  querying path: "+path.Source())
		rego = append(rego, fmt.Sprintf("%s_array = %s with data.sourceNode as %s", checkVariable, pathResult.rule, rule.Variable.Name))
		rego = append(rego, fmt.Sprintf("%s = %s_array", checkVariable, checkVariable))
	}

	// This is the value where we will store the result of the custom rego code.
	// Must be unique within the wrapping rego rule
	resultVariable := profile.Genvar("rego_result")

	// we first need to replace the variables in the rego template to match the right check and result variables
	// we have generated before
	text := rule.Argument
	text = strings.ReplaceAll(text, "$result", resultVariable)
	text = strings.ReplaceAll(text, "$node", checkVariable)
	// let's add all custom rego code to the code to be generated
	for _, l := range strings.Split(text, "\n") {
		rego = append(rego, l)
	}
	// now we can negate or not the resultVariable, we are checking that the result is erroneous
	if rule.Negated {
		rego = append(rego, fmt.Sprintf("%s == true", resultVariable))
	} else {
		rego = append(rego, fmt.Sprintf("%s != true", resultVariable))
	}

	r := SimpleRegoResult{
		Constraint: "rego",
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult}, // this can be an empty path result
		Path:       rule.Path.Source(),
		Variable:   checkVariable,
		TraceNode:  rule.Variable.Name,
		TraceValue: fmt.Sprintf("{\"negated\":%t}", rule.Negated),
	}
	return []SimpleRegoResult{r}
}
