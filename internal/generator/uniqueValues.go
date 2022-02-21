package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateUniqueValues(uniqueValues profile.UniqueValuesRule, iriExpander *misc.IriExpander) []SimpleRegoResult {
	// Fetch variable
	path := uniqueValues.Path
	var rego []string
	rego = append(rego, "#  querying path: "+path.Source())
	arrayRule := GeneratePropertyArray(path, uniqueValues.Variable.Name, iriExpander)

	// Set variable values
	arrayVariable := profile.Genvar("array_values")
	duplicatesVariable := profile.Genvar("duplicates")
	rego = append(rego, fmt.Sprintf("%s = %s with data.sourceNode as %s", arrayVariable, arrayRule.rule, uniqueValues.Variable.Name))

	findDuplicatesTemplate := `
  %s = { duplicate |
    array_value = %s[_]
    indices_for_value := [ idx | array_value == %s[idx]]
    count(indices_for_value) > 1
    duplicate = array_value
  }
`
	rego = append(rego, fmt.Sprintf(findDuplicatesTemplate, duplicatesVariable, arrayVariable, arrayVariable))

	// Add the validation
	if uniqueValues.Argument != uniqueValues.Negated { // Argument xor Negated
		rego = append(rego, fmt.Sprintf("count(%s) > 0", duplicatesVariable))
	} else {
		// When setting the Argument as false this rule will get activated if no duplicate values are found. This might
		// not be the expected behavior
		rego = append(rego, fmt.Sprintf("not count(%s) > 0", duplicatesVariable))
	}

	// Build result
	tracePath, err := uniqueValues.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}
	r := SimpleRegoResult{
		Constraint: "uniqueValues",
		Rego:       rego,
		PathRules:  []RegoPathResult{arrayRule},
		Path:       tracePath,
		TraceNode:  uniqueValues.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t", uniqueValues.Negated)),
		Variable: arrayVariable,
	}
	return []SimpleRegoResult{r}
}
