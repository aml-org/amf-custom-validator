package generator

import (
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateConditional(conditional profile.ConditionalRule, iriExpander *misc.IriExpander, llmMetrics *[]profile.LlmMetric) []BranchRegoResult {
	thenMaterialImplication := conditional.ThenMaterialImplication()
	var results = GenerateOr(thenMaterialImplication, iriExpander, llmMetrics)
	if conditional.ElseIsDefined() {
		elseMaterialImplication := conditional.ElseMaterialImplication()
		results = append(results, GenerateOr(elseMaterialImplication, iriExpander, llmMetrics)...)
	}
	return results
}
