package generator

import "github.com/aml-org/amfopa/internal/parser/profile"

func GenerateConditional(conditional profile.ConditionalRule) []BranchRegoResult {
	return GenerateOr(conditional.MaterialImplication())
}
