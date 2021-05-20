package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
)

func Dispatch(r profile.Rule) []GeneratedRegoResult {
	switch e := r.(type) {
	case profile.NestedExpression:
		return GenerateNestedExpression(e)
	case profile.CountRule:
		return simpleAsGeneratedRegoResult(GenerateCount(e))
	case profile.InRule:
		return simpleAsGeneratedRegoResult(GenerateIn(e))
	case profile.PatternRule:
		return simpleAsGeneratedRegoResult(GeneratePattern(e))
	case profile.AndRule:
		return branchAsGeneratedRegoResult(GenerateAnd(e))
	case profile.OrRule:
		return branchAsGeneratedRegoResult(GenerateOr(e))
	default:
		panic(errors.New(fmt.Sprintf("unknown rule type %v", r)))
	}
}

func simpleAsGeneratedRegoResult(simple []SimpleRegoResult) []GeneratedRegoResult {
	acc := make([]GeneratedRegoResult, len(simple))
	for i, s := range simple {
		acc[i] = GeneratedRegoResult(s)
	}
	return acc
}

func branchAsGeneratedRegoResult(simple []BranchRegoResult) []GeneratedRegoResult {
	acc := make([]GeneratedRegoResult, len(simple))
	for i, s := range simple {
		acc[i] = GeneratedRegoResult(s)
	}
	return acc
}
