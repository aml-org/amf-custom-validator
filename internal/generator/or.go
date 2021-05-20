package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

func GenerateOr(or profile.OrRule) []BranchRegoResult {
	if or.Negated {
		negated := or.Negate()
		switch andRule := negated.(type) {
		case profile.AndRule:
			return GenerateAnd(andRule)
		default:
			panic(fmt.Sprintf("ony AND can be the result of a naged OR, got %v", andRule))
		}

	} else {
		rego := collectAllResults(or)
		regoBranches := filterBranchesResult(rego)
		regoResults := filterSimpleResults(rego)
		return expandBranches(regoResults, regoBranches)
	}
}

// Flatten all the branches as a list of branches containing only simple rego results.
// Starts with th set of simple results in the OR and multiples it for each additional nested
// branch obtaining I * B1 * B2 * ... * Bn branches.
func expandBranches(regoResults []SimpleRegoResult, regoBranches [][]BranchRegoResult) []BranchRegoResult {
	// we start with a single branch holding only
	// the simple rego in the or if there was simple statements
	// in the OR, otherwise this is an array with an empty array inside
	acc := [][]SimpleRegoResult{regoResults}

	// now let's expand the branches
	for _, branches := range regoBranches {
		// we merge the branch and variables in the branch
		// with each of the previous computed branches
		// Now there should be n(acc) * m(branch) branches
		for _, branch := range branches {
			newAcc := make([][]SimpleRegoResult, 0)
			for _, source := range acc {
				for _, r := range branch.Branch {
					source = append(source, r)
				}
				newAcc = append(newAcc, source)
			}
			acc = newAcc
		}
	}

	orAcc := make([]BranchRegoResult, len(acc))
	for i, rs := range acc {
		orAcc[i] = BranchRegoResult{
			Constraint: "or",
			Branch:     rs,
		}
	}
	return orAcc
}

func filterSimpleResults(rego [][]GeneratedRegoResult) []SimpleRegoResult {
	acc := make([]SimpleRegoResult, 0)
	for _, r := range rego {
		for _, rr := range r {
			switch simple := rr.(type) {
			case SimpleRegoResult:
				acc = append(acc, simple)
			default:
				continue
			}
		}
	}
	return acc
}

func filterBranchesResult(rego [][]GeneratedRegoResult) [][]BranchRegoResult {
	acc := make([][]BranchRegoResult, 0)
	for _, r := range rego {
		branches := make([]BranchRegoResult, 0)
		for i, rr := range r {
			switch branch := rr.(type) {
			case BranchRegoResult:
				branches[i] = branch
			default:
				continue
			}
		}
		if len(branches) > 0 {
			acc = append(acc, branches)
		}
	}
	return acc
}

func collectAllResults(or profile.OrRule) [][]GeneratedRegoResult {
	acc := make([][]GeneratedRegoResult, len(or.Body))
	for i, r := range or.Body {
		acc[i] = Dispatch(r)
	}
	return acc
}
