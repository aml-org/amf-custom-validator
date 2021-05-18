package generator

import (
	"fmt"
	"strings"
)

func GenerateClassTarget(variable string, class string) SimpleRegoResult {
	rego := fmt.Sprintf("target_class[%s] with data.class as \"%s\"",variable,strings.ReplaceAll(class,".",":"))
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
