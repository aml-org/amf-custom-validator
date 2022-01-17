package generator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
)

func GenerateClassTarget(variable string, class string, iriExpander *misc.IriExpander) SimpleRegoResult {
	var classIri string
	if iriExpander != nil {
		var err error
		classIri, err = iriExpander.Expand(class)
		if err != nil {
			panic(err)
		}
	} else {
		classIri = class
	}

	rego := fmt.Sprintf("target_class[%s] with data.class as \"%s\"", variable, classIri)
	return SimpleRegoResult{
		Constraint: "classTarget",
		Rego:       []string{rego},
		Path:       "",
		Variable:   variable,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"classTarget\":\"%s\"", class)),
		TraceNode: variable,
		PathRules: []RegoPathResult{},
	}
}
