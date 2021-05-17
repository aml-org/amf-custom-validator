package generator

import "fmt"

func GenerateClassTarget(variable string, class string) SimpleRegoResult {
	rego := fmt.Sprintf("target_class[%s] with data.class as \"%s\"",variable,class)
	return SimpleRegoResult{
		Constraint: "classTarget",
		Rego: []string{rego},
		Path: "",
		Value: variable,
		Variable: variable,
		Trace: "",
		PathRules: []RegoPathResult{},
	}
}
