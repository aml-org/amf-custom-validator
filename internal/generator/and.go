package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

// Generates code snippets for an AND rule. The result is a set of branches,
// since we just need one of the predicates to produce an error for the validation
// to fail, we branch for each clause
func GenerateAnd(and profile.AndRule) []BranchRegoResult {
	if and.Negated {
		negated := and.Negate()
		switch orRule := negated.(type) {
		case profile.OrRule:
			return GenerateOr(orRule)
		default:
			panic(fmt.Sprintf("only OR can be the result of negating and AND, got %v", orRule))
		}
	} else {
		branches := make([]BranchRegoResult, 0)
		for _, r := range and.Body {
			results := Dispatch(r)
			for _, rr := range results {
				switch tr := rr.(type) {
				case SimpleRegoResult:
					branches = append(branches, BranchRegoResult{
						Constraint: tr.Constraint,
						Branch:     []SimpleRegoResult{tr},
					})
				case BranchRegoResult:
					branches = append(branches, tr)
				}
			}
		}
		return branches
	}
}
