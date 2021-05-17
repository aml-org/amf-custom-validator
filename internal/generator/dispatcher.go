package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/constraints"
	"github.com/aml-org/amfopa/internal/parser/profile/expression"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
)

func Dispatch(r statements.Rule) []GeneratedRegoResult {
	switch e := r.(type) {
	//case expression.Expression:
	//	return GenerateExpression(e)
	case constraints.CountRule:
		return simpleAsGeneratedRegoResult(GenerateCount(e))
	case expression.AndRule:
		return branchAsGeneratedRegoResult(GenerateAnd(e))
	case expression.OrRule:
		return branchAsGeneratedRegoResult(GenerateOr(e))
	default:
		panic(errors.New(fmt.Sprintf("unknown rule type %v", r)))
	}
}

func simpleAsGeneratedRegoResult(simple []SimpleRegoResult) []GeneratedRegoResult {
	acc := make([]GeneratedRegoResult,len(simple))
	for i,s := range simple {
		acc[i] = GeneratedRegoResult(s)
	}
	return acc
}

func branchAsGeneratedRegoResult(simple []BranchRegoResult) []GeneratedRegoResult {
	acc := make([]GeneratedRegoResult,len(simple))
	for i,s := range simple {
		acc[i] = GeneratedRegoResult(s)
	}
	return acc
}