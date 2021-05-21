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
func expandBranches(regoResults []SimpleRegoResult, branchsets [][]BranchRegoResult) []BranchRegoResult {
	// we start with a single branch-set holding only one branch for the single results of the OR
	// otherwise this is an branch-set with an empty 0 branches inside
	acc := [][]SimpleRegoResult{regoResults}

	// now let's expand the branch-sets
	for _, branches := range branchsets {
		// we merge the branches in this branch-set
		// with each of the previous computed branch-sets
		// Now there should be n-branches(acc) * m-branches(this-branch-set) branch-sets
		var newAcc [][]SimpleRegoResult
		for _, branch := range branches {
			for _, sourceBranchArray := range acc {
				var sourceExpanded []SimpleRegoResult
				sourceExpanded = append(sourceExpanded, sourceBranchArray...)
				for _, r := range branch.Branch {
					sourceExpanded = append(sourceBranchArray, r)
				}
				newAcc = append(newAcc, sourceExpanded)
			}
		}
		acc = newAcc
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
	var acc [][]BranchRegoResult
	for _, r := range rego {
		var branches []BranchRegoResult
		for _, rr := range r {
			switch branch := rr.(type) {
			case BranchRegoResult:
				branches = append(branches, branch)
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
