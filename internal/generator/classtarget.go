package generator

import (
	"fmt"

	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/open-policy-agent/opa/ast"
)

func GenerateClassTarget(variable string, class string, iriExpander *misc.IriExpander, regoVersion ast.RegoVersion) SimpleRegoResult {
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

	var rego = ""
	if IsV0(regoVersion) {
		rego = fmt.Sprintf("target_class[%s] with data.class as \"%s\"", variable, classIri)
	} else {
		rego = fmt.Sprintf("%s = target_class(\"%s\")[_]", variable, classIri)
	}

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
