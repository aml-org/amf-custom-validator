package generator

import (
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/open-policy-agent/opa/ast"
)

func GenerateConditional(conditional profile.ConditionalRule, iriExpander *misc.IriExpander, regoVersion ast.RegoVersion) []BranchRegoResult {
	thenMaterialImplication := conditional.ThenMaterialImplication()
	var results = GenerateOr(thenMaterialImplication, iriExpander, regoVersion)
	if conditional.ElseIsDefined() {
		elseMaterialImplication := conditional.ElseMaterialImplication()
		results = append(results, GenerateOr(elseMaterialImplication, iriExpander, regoVersion)...)
	}
	return results
}
