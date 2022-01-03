package generator

import "github.com/aml-org/amf-custom-validator/internal/parser/profile"

func GenerateConditional(conditional profile.ConditionalRule) []BranchRegoResult {
	thenMaterialImplication := conditional.ThenMaterialImplication()
	var results = GenerateOr(thenMaterialImplication)
	if conditional.ElseIsDefined() {
		elseMaterialImplication := conditional.ElseMaterialImplication()
		results = append(results, GenerateOr(elseMaterialImplication)...)
	}
	return results
}
