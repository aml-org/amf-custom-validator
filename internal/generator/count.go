package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/constraints"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
)

// Generates the Rego code snippet for the rule, supports minCount and maxCount
func GenerateCount(count constraints.CountRule) []SimpleRegoResult {
	if count.Qualifier == constraints.Min {
		return generateRule(count, "minCount", ">=")
	} else {
		return generateRule(count, "maxCount", "<=")
	}
}

// Generates the rule using the 'count'  property from Rego
func generateRule(count constraints.CountRule, rule string, condition string) []SimpleRegoResult {
	path := count.Path
	rego := make([]string, 0)

	// Let's get the path computed and stored in the inValuesVariable
	inValuesVariable := statements.Genvar("propValues")
	rego = append(rego, "#  storing path in "+inValuesVariable+" : "+path.Source())
	pathResult := GeneratePropertyArray(path, count.Variable.Name, rule+"_"+count.ValueHash())
	rego = append(rego, fmt.Sprintf("%s = %s with data.sourceNode as %s", inValuesVariable, pathResult.rule, count.Variable.Name))

	// Add the validation
	if count.Negated {
		rego = append(rego, fmt.Sprintf("count(%s) %s %d", inValuesVariable, condition, count.Argument))
	} else {
		rego = append(rego, fmt.Sprintf("not count(%s) %s %d", inValuesVariable, condition, count.Argument))
	}

	r := SimpleRegoResult{
		Constraint: rule,
		Rego:       rego,
		PathRules:  []RegoPathResult{pathResult},
		Path:       count.Path.Source(),
		Value:      fmt.Sprintf("count(%s)", inValuesVariable),
		Variable:   inValuesVariable,
		Trace:      fmt.Sprintf("value not matching rule %d", count.Argument),
	}
	return []SimpleRegoResult{r}
}
