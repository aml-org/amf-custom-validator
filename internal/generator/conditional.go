package generator

import "github.com/aml-org/amf-custom-validator/internal/parser/profile"

func GenerateConditional(conditional profile.ConditionalRule) []BranchRegoResult {
	return GenerateOr(conditional.MaterialImplication())
}
