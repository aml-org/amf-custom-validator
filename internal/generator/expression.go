package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/expression"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"strings"
)

func GenerateTopLevelExpression(exp statements.Rule) string {
	switch e := exp.(type) {
	case expression.TopLevelExpression:
		return generateTopLevel(e)
	case expression.Expression:
		panic(errors.New("nested expressions cannot be generated as a top level expression"))
	default:
		panic(errors.New(fmt.Sprintf("expected expression or top-level expression, got %v", exp)))
	}
}

/*
func GenerateNestedExpression(exp statements.Rule) string {
	switch e := exp.(type) {
	case expression.TopLevelExpression:
		panic(errors.New("nested expressions not supported yet"))
	case expression.Expression:
		return generateTopLevel(e)
	default:
		panic(errors.New(fmt.Sprintf("expected expression or top-level expression, got %v", exp)))
	}
}
*/
func generateTopLevel(e expression.TopLevelExpression) string {
	if e.Negated {
		return generateBodyNegated(e)
	} else {
		return generateBody(e)
	}
}

func generateBodyNegated(e expression.TopLevelExpression) string {
	negated := e.Value.Negate()
	results := Dispatch(negated)
	return wrapRegoResult(e, results)
}

func generateBody(e expression.TopLevelExpression) string {
	results := Dispatch(e.Value)
	return wrapRegoResult(e, results)
}

func wrapRegoResult(e expression.TopLevelExpression, results []GeneratedRegoResult) string {
	classTargetResult := GenerateClassTarget(e.Variables[0].Name, e.ClassGenerator)
	classTargetVariable := classTargetResult.Variable

	branches := make([]BranchRegoResult, 0)
	for _,r := range results {
		switch rr := r.(type) {
		case SimpleRegoResult:
			branches = append(branches,BranchRegoResult{
			  Constraint: rr.Constraint,
			  Branch: []SimpleRegoResult{rr},
			})
		case BranchRegoResult:
			branches = append(branches,rr)
		}
	}
	pathValidationsAcc := make([]string,0)
	for _,branch := range branches {
		for _,r := range branch.Branch {
			for _,pr := range r.PathRules {
				pathValidationsAcc = append(pathValidationsAcc, strings.Join(pr.rego, "\n"))
			}
		}
	}
	pathValidations := strings.Join(pathValidationsAcc, "\n\n")

	branchesAcc := make([]string, 0)
	for _,branch := range branches {
		acc := make([]string,0)
		acc = append(acc,fmt.Sprintf("%s[matches] {", strings.ToLower(e.Level)))
		for _,r := range classTargetResult.Rego {
			acc = append(acc, "  " + r)
		}
		for _,l := range wrapBranch(e,branch,"matches",classTargetVariable) {
			acc = append(acc, l)
		}
		acc = append(acc, "}")
		branchesAcc = append(branchesAcc, strings.Join(acc,"\n"))
	}

	branchesValidation := strings.Join(branchesAcc,"\n\n")
	total := []string{"# Path rules", pathValidations, "# Constraint rules", branchesValidation}
	return strings.Join(total, "\n\n")
}


func wrapBranch(e expression.TopLevelExpression, branch BranchRegoResult, matchesVariable string, mappingVariable string) []string {
	acc := make([]string,0)
	resultBindings := make([]string,0)

	for i,r := range branch.Branch {
		bindingResult := fmt.Sprintf("_result_%d",i)
		resultBindings = append(resultBindings, bindingResult)
		matchesLine := fmt.Sprintf("  %s := trace(\"%s\",\"%s\",%s,\"%s\")",bindingResult,r.ConstraintId(),r.Path, r.Value,r.Trace)
		for _,l := range r.Rego {
			acc = append(acc, "  " + l)
		}
		acc = append(acc, matchesLine)
	}

	acc = append(acc, fmt.Sprintf("  %s := error(\"%s\",%s,\"%s\",[%s])", matchesVariable, e.Name, mappingVariable,e.Message,strings.Join(resultBindings,",")))
	return acc
}
