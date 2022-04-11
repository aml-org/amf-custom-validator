package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func Dispatch(r profile.Rule, iriExpander *misc.IriExpander) []GeneratedRegoResult {
	switch e := r.(type) {
	case profile.NestedExpression:
		return GenerateNestedExpression(e, iriExpander)
	case profile.CountRule:
		return simpleAsGeneratedRegoResult(GenerateCount(e, iriExpander))
	case profile.ScalarSetRule:
		if e.SetCriteria == profile.SubSet {
			return simpleAsGeneratedRegoResult(GenerateScalarSubSetRule(e, iriExpander))
		} else if e.SetCriteria == profile.InsersectSet {
			return simpleAsGeneratedRegoResult(GenerateScalarIntersectSetRule(e, iriExpander))
		} else {
			return simpleAsGeneratedRegoResult(GenerateScalarSuperSetRule(e, iriExpander))
		}
	case profile.PatternRule:
		return simpleAsGeneratedRegoResult(GeneratePattern(e, iriExpander))
	case profile.UniqueValuesRule:
		return simpleAsGeneratedRegoResult(GenerateUniqueValues(e, iriExpander))
	case profile.PropertyComparisonRule:
		return simpleAsGeneratedRegoResult(GeneratePropertyComparison(e, iriExpander))
	case profile.AndRule:
		return branchAsGeneratedRegoResult(GenerateAnd(e, iriExpander))
	case profile.OrRule:
		return branchAsGeneratedRegoResult(GenerateOr(e, iriExpander))
	case profile.ConditionalRule:
		return branchAsGeneratedRegoResult(GenerateConditional(e, iriExpander))
	case profile.RegoRule:
		return simpleAsGeneratedRegoResult(GenerateRegoRule(e, iriExpander))
	case profile.NumericRule:
		return simpleAsGeneratedRegoResult(GenerateNumericComparison(e, iriExpander))
	case profile.DatatypeRule:
		return simpleAsGeneratedRegoResult(GenerateDatatype(e, iriExpander))
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
