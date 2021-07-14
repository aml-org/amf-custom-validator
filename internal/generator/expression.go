package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateTopLevelExpression(exp profile.Rule) string {
	switch e := exp.(type) {
	case profile.TopLevelExpression:
		return generateTopLevel(e)
	case profile.NestedExpression:
		panic(errors.New("nested expressions cannot be generated as a top level expression"))
	default:
		panic(errors.New(fmt.Sprintf("expected expression or top-level expression, got %v", exp)))
	}
}

func GenerateNestedExpression(exp profile.Rule) []GeneratedRegoResult {
	switch e := exp.(type) {
	case profile.TopLevelExpression:
		panic(errors.New("nested expressions not supported yet"))
	case profile.NestedExpression:
		return generateNested(e)
	default:
		panic(errors.New(fmt.Sprintf("expected expression or top-level expression, got %v", exp)))
	}
}

func generateTopLevel(exp profile.TopLevelExpression) string {
	v := exp.Value
	if exp.Negated {
		v = exp.Value.Negate()
	}
	results := Dispatch(v)
	return wrapTopLevelRegoResult(exp, results)
}

func generateNested(exp profile.NestedExpression) []GeneratedRegoResult {
	v := exp.Value
	sharedGeneratorRego := GenerateNested(exp)
	results := Dispatch(v)

	var acc []GeneratedRegoResult
	for _, rego := range results {
		switch r := rego.(type) {
		case BranchRegoResult:
			acc = append(acc, wrapNestedRegoResult(exp, sharedGeneratorRego, r))
		}
	}
	return acc
}

func wrapNestedRegoResult(exp profile.NestedExpression, sharedGeneratorRego SimpleRegoResult, branchRego BranchRegoResult) BranchRegoResult {
	nestedVariable := exp.Child.Name
	var rego []string
	// we append the shared rego generator. all the conditions of the branch will be check over the nested nodes through
	// the nested property path.
	rego = append(rego, sharedGeneratorRego.Rego...)
	// this variable holds each of the nested nodes
	variable := sharedGeneratorRego.Variable
	errorVariable := fmt.Sprintf("%s_error", variable)
	// let's create a comprehension for the conditions over the nested nodes
	rego = append(rego, fmt.Sprintf("%ss = [ %s|", errorVariable, errorVariable)) // report errors using this variable
	rego = append(rego, fmt.Sprintf("  %s = %s[_]", nestedVariable, variable))    // the underlying rules expect the quantified variable that was passed on profile parsing
	for _, l := range wrapBranch("nested", fmt.Sprintf("error in nested nodes under %s", sharedGeneratorRego.Path), branchRego, errorVariable, nestedVariable) {
		rego = append(rego, l)
	}
	rego = append(rego, "]")

	// Now Let's check if there was an error counting the failed nodes
	if exp.Child.Cardinality != nil {
		// quantified nested, we need to check for a particular number of failed nodes
		if exp.Negated {
			rego = append(rego, fmt.Sprintf("count(%s) - count(%ss) %s", variable, errorVariable, exp.Child.Cardinality.String()))
		} else {
			rego = append(rego, fmt.Sprintf("not count(%s) - count(%ss) %s", variable, errorVariable, exp.Child.Cardinality.String()))
		}
	} else {
		// regular nested we just loo for one error
		if exp.Negated {
			rego = append(rego, fmt.Sprintf("count(%ss) == 0", errorVariable))
		} else {
			rego = append(rego, fmt.Sprintf("count(%ss) > 0", errorVariable))
		}
	}

	// accumulate path rules
	var pathRules []RegoPathResult
	for _, pr := range sharedGeneratorRego.PathRules {
		pathRules = append(pathRules, pr)
	}
	for _, r := range branchRego.Branch {
		for _, pr := range r.PathRules {
			pathRules = append(pathRules, pr)
		}
	}

	// build result
	return BranchRegoResult{
		Constraint: "nested",
		Branch: []SimpleRegoResult{
			{
				Constraint: "nested",
				Rego:       rego,
				PathRules:  pathRules,
				Path:       sharedGeneratorRego.Path,
				TraceNode:  exp.Parent.Name,
				TraceValue: fmt.Sprintf("{\"negated\":%t, \"expected\":0, \"actual\":count(%ss)}", exp.Negated, errorVariable),
				Variable:   fmt.Sprintf("%ss", errorVariable),
			},
		},
	}

}

func wrapTopLevelRegoResult(e profile.TopLevelExpression, results []GeneratedRegoResult) string {
	classTargetResult := GenerateClassTarget(e.Variable.Name, e.ClassGenerator)
	classTargetVariable := classTargetResult.Variable

	var branches []BranchRegoResult
	for _, r := range results {
		switch rr := r.(type) {
		case SimpleRegoResult:
			branches = append(branches, BranchRegoResult{
				Constraint: rr.Constraint,
				Branch:     []SimpleRegoResult{rr},
			})
		case BranchRegoResult:
			branches = append(branches, rr)
		}
	}
	pathValidationsAcc := make([]string, 0)
	for _, branch := range branches {
		for _, r := range branch.Branch {
			for _, pr := range r.PathRules {
				pathValidationsAcc = append(pathValidationsAcc, strings.Join(pr.rego, "\n"))
			}
		}
	}
	pathValidations := strings.Join(pathValidationsAcc, "\n\n")

	branchesAcc := make([]string, 0)
	for _, branch := range branches {
		var acc []string
		acc = append(acc, fmt.Sprintf("%s[matches] {", strings.ToLower(e.Level)))
		for _, r := range classTargetResult.Rego {
			acc = append(acc, "  "+r)
		}
		for _, l := range wrapBranch(e.Name, e.Message, branch, "matches", classTargetVariable) {
			acc = append(acc, l)
		}
		acc = append(acc, "}")
		branchesAcc = append(branchesAcc, strings.Join(acc, "\n"))
	}

	branchesValidation := strings.Join(branchesAcc, "\n\n")
	total := []string{"# Path rules", pathValidations, "# Constraint rules", branchesValidation}
	return strings.Join(total, "\n\n")
}

func wrapBranch(name string, message string, branch BranchRegoResult, matchesVariable string, mappingVariable string) []string {
	var acc []string
	resultBindings := make([]string, 0)

	for i, r := range branch.Branch {
		bindingResult := fmt.Sprintf("_result_%d", i)
		resultBindings = append(resultBindings, bindingResult)
		matchesLine := fmt.Sprintf("  %s := trace(\"%s\",\"%s\",%s,%s)", bindingResult, r.ConstraintId(), r.Path, r.TraceNode, r.TraceValue)
		for _, l := range r.Rego {
			acc = append(acc, "  "+l)
		}
		acc = append(acc, matchesLine)
	}

	acc = append(acc, fmt.Sprintf("  %s := error(\"%s\",%s,\"%s\",[%s])", matchesVariable, name, mappingVariable, sanitizedMessage(message), strings.Join(resultBindings, ",")))
	return acc
}

func sanitizedMessage(s string) string {
	result := strings.ReplaceAll(s, "\n", "\\n")
	return strings.ReplaceAll(result, "\"", "'")
}
