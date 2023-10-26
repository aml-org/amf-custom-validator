package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
)

func GenerateTopLevelExpression(exp profile.Rule, iriExpander *misc.IriExpander) string {
	switch e := exp.(type) {
	case profile.TopLevelExpression:
		return generateTopLevel(e, iriExpander)
	case profile.NestedExpression:
		panic(errors.New("nested expressions cannot be generated as a top level expression"))
	default:
		panic(errors.New(fmt.Sprintf("expected expression or top-level expression, got %v", exp)))
	}
}

func GenerateNestedExpression(exp profile.Rule, iriExpander *misc.IriExpander) []GeneratedRegoResult {
	switch e := exp.(type) {
	case profile.TopLevelExpression:
		panic(errors.New("nested expressions not supported yet"))
	case profile.NestedExpression:
		return generateNested(e, iriExpander)
	default:
		panic(errors.New(fmt.Sprintf("expected expression or top-level expression, got %v", exp)))
	}
}

func generateTopLevel(exp profile.TopLevelExpression, iriExpander *misc.IriExpander) string {
	v := exp.Value
	if exp.Negated {
		v = exp.Value.Negate()
	}
	results := Dispatch(v, iriExpander)
	return wrapTopLevelRegoResult(exp, results, iriExpander)
}

func generateNested(exp profile.NestedExpression, iriExpander *misc.IriExpander) []GeneratedRegoResult {
	v := exp.Value

	var rego []string

	// we extract the same node for all the checks in the related branches
	sharedGeneratorRego := GenerateNested(exp, iriExpander)
	rego = append(rego, sharedGeneratorRego.Rego...)
	var variable = sharedGeneratorRego.Variable

	// accumulate path rules
	var pathRules []RegoPathResult
	for _, pr := range sharedGeneratorRego.PathRules {
		pathRules = append(pathRules, pr)
	}

	results := Dispatch(v, iriExpander)
	var branchErrors []string
	var errorAcc = fmt.Sprintf("%s_errorAcc", exp.Child.Name)
	rego = append(rego, fmt.Sprintf("%s0 = []", errorAcc))
	for i, b := range results {
		switch r := b.(type) {
		case BranchRegoResult:
			rego = append(rego, wrapNestedRegoResult(variable, i, errorAcc, exp.Child.Name, sharedGeneratorRego, r, &branchErrors, iriExpander)...)
			// accumulate paths
			for _, r := range r.Branch {
				for _, pr := range r.PathRules {
					pathRules = append(pathRules, pr)
				}
			}
		case SimpleRegoResult:
			var wrappedBranch = BranchRegoResult{Branch: []SimpleRegoResult{r}}
			rego = append(rego, wrapNestedRegoResult(variable, i, errorAcc, exp.Child.Name, sharedGeneratorRego, wrappedBranch, &branchErrors, iriExpander)...)
			// accumulate paths
			for _, pr := range r.PathRules {
				pathRules = append(pathRules, pr)
			}
		}
	}
	rego = append(rego, fmt.Sprintf("%s = %s%d", errorAcc, errorAcc, len(results)))

	rego = append(rego, "# let's accumulate results")
	errorNodeVariable := fmt.Sprintf("%s_error_node_variables_agg", variable)
	aggVar := strings.Join(branchErrors, " | ") // TODO: change this for AND/OR  cases?
	rego = append(rego, fmt.Sprintf("%s = %s", errorNodeVariable, aggVar))

	// to generate the right information in the trace, default values for nested, negated
	var ruleName = ""
	var trace = ""
	// fmt.Sprintf("{\"negated\":%t, \"expected\":%s, \"condition\":\"%s\", \"actual\":count(%s), \"subResult\": %s}", exp.Negated, expected, condition, errorNodeVariable, errorAcc)

	// regular nested we just look for one error
	if exp.Child.Quantification == profile.ForAll {
		if exp.Negated {
			rego = append(rego, fmt.Sprintf("count(%s) == 0", errorNodeVariable))
		} else {
			rego = append(rego, fmt.Sprintf("count(%s) > 0", errorNodeVariable))
		}
		ruleName = "nested"
		trace = fmt.Sprintf("\"negated\":%t, \"failedNodes\":count(%s), \"successfulNodes\":(count(%s)-count(%s)),\"subResult\": %s", exp.Negated, errorNodeVariable, variable, errorNodeVariable, errorAcc)
	} else {
		// quantified nested, we need to check for a particular number of failed nodes
		if exp.Negated {
			rego = append(rego, fmt.Sprintf("count(%s) - count(%s) %s", variable, errorNodeVariable, exp.Child.Cardinality.String()))
		} else {
			rego = append(rego, fmt.Sprintf("not count(%s) - count(%s) %s", variable, errorNodeVariable, exp.Child.Cardinality.String()))
		}
		ruleName = exp.Child.Cardinality.RuleName()
		trace = fmt.Sprintf("\"negated\":%t, \"failedNodes\":count(%s), \"successfulNodes\":(count(%s)-count(%s)), \"cardinality\":%d, \"subResult\": %s", exp.Negated, errorNodeVariable, variable, errorNodeVariable, exp.Child.Cardinality.Value, errorAcc)
	}

	// build result
	return []GeneratedRegoResult{
		BranchRegoResult{
			Constraint: ruleName,
			Branch: []SimpleRegoResult{
				{
					Constraint: ruleName,
					Rego:       rego,
					PathRules:  pathRules,
					Path:       sharedGeneratorRego.Path,
					TraceNode:  exp.Parent.Name,
					TraceValue: BuildTraceValueNode(trace),
					Variable:   fmt.Sprintf("%s", errorNodeVariable),
				},
			},
		},
	}
}

func wrapNestedRegoResult(variable string, i int, errorAcc string, nodeVariable string, sharedGeneratorRego SimpleRegoResult, branchRego BranchRegoResult, branchErrors *[]string, iriExpander *misc.IriExpander) []string {
	branchResultVariable := fmt.Sprintf("%s_br_%d", variable, i)
	branchResultError := fmt.Sprintf("%s_br_%d_errors", variable, i)
	*branchErrors = append(*branchErrors, branchResultError)

	var rego []string
	errorVariable := fmt.Sprintf("%s_error", branchResultVariable) // all failing errors for each nested node will be collected here
	// this is an intermediate variable collecting tuples [inner_node, error] so we can catch the case where the inner
	// constraint validates multiple value of a property inside a single nested node
	innerErrorVariable := fmt.Sprintf("%s_inner_error", branchResultVariable)
	// let's create a comprehension for the conditions over the nested nodes
	rego = append(rego, fmt.Sprintf("%s = [ %s|", branchResultVariable, errorVariable))          // report errors using this variable
	rego = append(rego, fmt.Sprintf("  %s = %s[_]", nodeVariable, sharedGeneratorRego.Variable)) // the underlying rules expect the quantified variable that was passed on profile parsing
	// note we are passing innerErrorVariable here
	message := profile.Message{Expression: fmt.Sprintf("error in nested nodes under %s", sharedGeneratorRego.Path), Variables: make([]string, 0)}
	for _, l := range wrapBranch("nested", message, branchRego, innerErrorVariable, nodeVariable, iriExpander) {
		rego = append(rego, l)
	}
	rego = append(rego, fmt.Sprintf("  %s = [%s[\"@id\"],%s]", errorVariable, nodeVariable, innerErrorVariable))
	rego = append(rego, "]")

	// let's now split the collected tuples into set of failing nodes for cardinality checks
	rego = append(rego, fmt.Sprintf("%s = { nodeId | n = %s[_]; nodeId = n[0] }", branchResultError, branchResultVariable))
	rego = append(rego, fmt.Sprintf("%s_errors = [ node | n = %s[_]; node = n[1] ]", branchResultError, branchResultVariable))
	rego = append(rego, fmt.Sprintf("%s%d = array.concat(%s%d,%s_errors)", errorAcc, i+1, errorAcc, i, branchResultError))

	return rego
}

func wrapTopLevelRegoResult(e profile.TopLevelExpression, results []GeneratedRegoResult, iriExpander *misc.IriExpander) string {
	classTargetResult := GenerateClassTarget(e.Variable.Name, e.ClassGenerator, iriExpander)
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
		for _, l := range wrapBranch(e.Name, e.Message, branch, "matches", classTargetVariable, iriExpander) {
			acc = append(acc, l)
		}
		acc = append(acc, "}")
		branchesAcc = append(branchesAcc, strings.Join(acc, "\n"))
	}

	branchesValidation := strings.Join(branchesAcc, "\n\n")
	total := []string{"# Path rules", pathValidations, "# Constraint rules", branchesValidation}
	return strings.Join(total, "\n\n")
}

func wrapBranch(name string, message profile.Message, branch BranchRegoResult, matchesVariable string, mappingVariable string, iriExpander *misc.IriExpander) []string {
	var acc []string
	resultBindings := make([]string, 0)
	customMessage := false

	for i, r := range branch.Branch {
		bindingResult := fmt.Sprintf("_result_%d", i)
		resultBindings = append(resultBindings, bindingResult)
		traceResultPath := r.Path
		if iriExpander != nil {
			traceResultPath, _ = iriExpander.Expand(r.Path)
		}
		matchesLine := fmt.Sprintf("  %s := trace(\"%s\",\"%s\",%s,%s)", bindingResult, r.ConstraintId(), traceResultPath, r.TraceNode, r.TraceValue)
		for _, l := range r.Rego {
			if strings.Contains(l, "$message") {
				customMessage = true
				l = strings.ReplaceAll(l, "$message", "message")
			}
			acc = append(acc, "  "+l)
		}
		acc = append(acc, matchesLine)
	}

	// if not defined in custom rego, generate message from the message key in the profile yaml
	if !customMessage {
		vars := make([]string, 0)
		for varIdx, compactPath := range message.Variables {
			expandedPath, _ := iriExpander.Expand(compactPath)

			// set variable
			varName := fmt.Sprintf("msg_var_%d", varIdx)
			acc = append(acc, fmt.Sprintf("  %s := object.get(%s, \"%s\", \"null\")", varName, mappingVariable, expandedPath))
			vars = append(vars, varName)
		}

		if len(vars) > 0 {
			acc = append(acc, fmt.Sprintf("  message_vars := [%s]", strings.Join(vars, ",")))
			acc = append(acc, fmt.Sprintf("  message := sprintf(\"%s\", message_vars)", sanitizedMessage(message.Expression)))
		} else {
			acc = append(acc, fmt.Sprintf("  message := \"%s\"", sanitizedMessage(message.Expression)))
		}
	}

	acc = append(acc, fmt.Sprintf("  %s := error(\"%s\",%s, message ,[%s])", matchesVariable, name, mappingVariable, strings.Join(resultBindings, ",")))
	return acc
}

func sanitizedMessage(s string) string {
	result := strings.ReplaceAll(s, "\n", "\\n")
	return strings.ReplaceAll(result, "\"", "'")
}
