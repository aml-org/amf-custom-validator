package generator

import (
	"errors"
	"fmt"

	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/open-policy-agent/opa/ast"
)

func Dispatch(r profile.Rule, iriExpander *misc.IriExpander, regoVersion ast.RegoVersion) []GeneratedRegoResult {
	switch e := r.(type) {
	case profile.NestedExpression:
		return GenerateNestedExpression(e, iriExpander, regoVersion)
	case profile.CountRule:
		return simpleAsGeneratedRegoResult(GenerateCount(e, iriExpander, regoVersion))
	case profile.ScalarSetRule:
		if e.SetCriteria == profile.SubSet {
			return simpleAsGeneratedRegoResult(GenerateScalarSubSetRule(e, iriExpander, regoVersion))
		} else if e.SetCriteria == profile.InsersectSet {
			return simpleAsGeneratedRegoResult(GenerateScalarIntersectSetRule(e, iriExpander, regoVersion))
		} else {
			return simpleAsGeneratedRegoResult(GenerateScalarSuperSetRule(e, iriExpander, regoVersion))
		}
	case profile.PatternRule:
		return simpleAsGeneratedRegoResult(GeneratePattern(e, iriExpander, regoVersion))
	case profile.UniqueValuesRule:
		return simpleAsGeneratedRegoResult(GenerateUniqueValues(e, iriExpander, regoVersion))
	case profile.PropertyComparisonRule:
		return simpleAsGeneratedRegoResult(GeneratePropertyComparison(e, iriExpander, regoVersion))
	case profile.AndRule:
		return branchAsGeneratedRegoResult(GenerateAnd(e, iriExpander, regoVersion))
	case profile.OrRule:
		return branchAsGeneratedRegoResult(GenerateOr(e, iriExpander, regoVersion))
	case profile.ConditionalRule:
		return branchAsGeneratedRegoResult(GenerateConditional(e, iriExpander, regoVersion))
	case profile.RegoRule:
		return simpleAsGeneratedRegoResult(GenerateRegoRule(e, iriExpander, regoVersion))
	case profile.NumericRule:
		return simpleAsGeneratedRegoResult(GenerateNumericComparison(e, iriExpander, regoVersion))
	case profile.DatatypeRule:
		return simpleAsGeneratedRegoResult(GenerateDatatype(e, iriExpander, regoVersion))
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
