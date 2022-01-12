package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func Dispatch(r profile.Rule) []GeneratedRegoResult {
	switch e := r.(type) {
	case profile.NestedExpression:
		return GenerateNestedExpression(e)
	case profile.CountRule:
		return simpleAsGeneratedRegoResult(GenerateCount(e))
	case profile.ScalarSetRule:
		if e.SetCriteria == profile.SubSet {
			return simpleAsGeneratedRegoResult(GenerateScalarSubSetRule(e))
		} else if e.SetCriteria == profile.InsersectSet {
			return simpleAsGeneratedRegoResult(GenerateScalarIntersectSetRule(e))
		} else if e.SetCriteria == profile.EqualsSet {
			return simpleAsGeneratedRegoResult(GenerateScalarIsSetRule(e))
		} else {
			return simpleAsGeneratedRegoResult(GenerateScalarSuperSetRule(e))
		}
	case profile.PatternRule:
		return simpleAsGeneratedRegoResult(GeneratePattern(e))
	case profile.PropertyComparisonRule:
		return simpleAsGeneratedRegoResult(GeneratePropertyComparison(e))
	case profile.AndRule:
		return branchAsGeneratedRegoResult(GenerateAnd(e))
	case profile.OrRule:
		return branchAsGeneratedRegoResult(GenerateOr(e))
	case profile.ConditionalRule:
		return branchAsGeneratedRegoResult(GenerateConditional(e))
	case profile.RegoRule:
		return simpleAsGeneratedRegoResult(GenerateRegoRule(e))
	case profile.NumericRule:
		return simpleAsGeneratedRegoResult(GenerateNumericComparison(e))
	case profile.DatatypeRule:
		return simpleAsGeneratedRegoResult(GenerateDatatype(e))
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
