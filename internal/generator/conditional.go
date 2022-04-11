package generator

import (
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateConditional(conditional profile.ConditionalRule, iriExpander *misc.IriExpander) []BranchRegoResult {
	thenMaterialImplication := conditional.ThenMaterialImplication()
	var results = GenerateOr(thenMaterialImplication, iriExpander)
	if conditional.ElseIsDefined() {
		elseMaterialImplication := conditional.ElseMaterialImplication()
		results = append(results, GenerateOr(elseMaterialImplication, iriExpander)...)
	}
	return results
}
