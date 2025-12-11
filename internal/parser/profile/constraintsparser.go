package profile

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	pathParser "github.com/aml-org/amf-custom-validator/internal/parser/path"
	y "github.com/aml-org/amf-custom-validator/internal/parser/yaml"
)

func ParseConstraint(path pathParser.PropertyPath, variable Variable, constraint *y.Yaml, varGenerator *VarGenerator) ([]Rule, error) {
	var acc []Rule

	minCount, err := constraint.Get("minCount").Int()
	if err == nil {
		acc = append(acc, newMinCount(false, variable, path, minCount))
	}

	maxCount, err := constraint.Get("maxCount").Int()
	if err == nil {
		acc = append(acc, newMaxCount(false, variable, path, maxCount))
	}

	exactCount, err := constraint.Get("exactCount").Int()
	if err == nil {
		acc = append(acc, newExactCount(false, variable, path, exactCount))
	}

	minLength, err := constraint.Get("minLength").Int()
	if err == nil {
		acc = append(acc, newMinLength(false, variable, path, minLength))
	}

	maxLength, err := constraint.Get("maxLength").Int()
	if err == nil {
		acc = append(acc, newMaxLength(false, variable, path, maxLength))
	}

	exactLength, err := constraint.Get("exactLength").Int()
	if err == nil {
		acc = append(acc, newExactLength(false, variable, path, exactLength))
	}

	pattern, err := constraint.Get("pattern").String()
	if err == nil {
		acc = append(acc, newPattern(false, variable, path, pattern))
	}

	in, err := constraint.Get("in").Array()
	if err == nil {
		l, err := scalarList(in)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newIn(false, variable, path, l))
	}

	unique, err := constraint.Get("uniqueValues").Bool()
	if err == nil {
		acc = append(acc, newUniqueValues(false, variable, path, unique))
	}

	containsAll, err := constraint.Get("containsAll").Array()
	if err == nil {
		l, err := scalarList(containsAll)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newContainsAll(false, variable, path, l))
	}

	containsSome, err := constraint.Get("containsSome").Array()
	if err == nil {
		l, err := scalarList(containsSome)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newContainsSome(false, variable, path, l))
	}

	otherProp, err := constraint.Get("lessThanProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newLessThan(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("lessThanOrEqualsToProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newLessThanOrEquals(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("equalsToProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newEquals(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("disjointWithProperty").String()
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newDisjoint(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("moreThanProperty").String() // TODO: not defined in dialect or tested
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newMoreThan(false, variable, path, compPath))
	}

	otherProp, err = constraint.Get("moreThanOrEqualsToProperty").String() // TODO: not defined in dialect or tested
	if err == nil {
		compPath, err := pathParser.ParsePath(otherProp)
		if err != nil {
			return nil, err
		}
		acc = append(acc, newMoreThanOrEquals(false, variable, path, compPath))
	}

	atLeast := constraint.Get("atLeast")
	if atLeast.IsFound() {
		rule, err := parseQualifiedNestedExpression(atLeast, false, variable, path, varGenerator, GTEQ)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	atMost := constraint.Get("atMost")
	if atMost.IsFound() {
		rule, err := parseQualifiedNestedExpression(atMost, false, variable, path, varGenerator, LTEQ)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	exactly := constraint.Get("exactly") // TODO: not defined in dialect or tested
	if exactly.IsFound() {
		rule, err := parseQualifiedNestedExpression(exactly, false, variable, path, varGenerator, EQ)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num := constraint.Get("minInclusive")
	if num.IsFound() {
		rule, err := parseMinInclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num = constraint.Get("minExclusive")
	if num.IsFound() {
		rule, err := parseMinExclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num = constraint.Get("maxInclusive")
	if num.IsFound() {
		rule, err := parseMaxInclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	num = constraint.Get("maxExclusive")
	if num.IsFound() {
		rule, err := parseMaxExclusive(false, variable, path, num)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	dt := constraint.Get("datatype")
	if dt.IsFound() {
		rule, err := parseDatatype(false, variable, path, dt)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	nested := constraint.Get("nested")
	if nested.IsFound() && nested.IsMap() {
		rule, err := parseNestedExpression(nested, false, variable, path, varGenerator)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	code := constraint.Get("rego")
	if code.IsFound() {
		rule, err := ParseRego(code, false, variable, path)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	codeModule := constraint.Get("regoModule")
	if codeModule.IsFound() {
		rule, err := ParseRego(codeModule, false, variable, path)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	llmMetric := constraint.Get("llmCoherence")
	if llmMetric.IsFound() {
		rule, err := ParseLlmMetric("coherence", llmMetric, false, variable, path)
		if err != nil {
			return nil, err
		}
		acc = append(acc, rule)
	}

	return acc, nil

}

// metricDescriptions maps metric names to their detailed descriptions for LLM evaluation
var metricDescriptions = map[string]string{
	"coherence":    "Coherence. A generated answer is coherent if, when read by itself, it’s easy to understand and free of grammar errors.\\\\nA generated answer is not coherent if it’s difficult to understand what the answer is trying to say.\\\\nGenerally, it’s more important that the response is understandable than it being free of grammar errors.",
	"concis":       "Assess how clear and understandable the content is. Consider if the language is precise, if concepts are well-explained, and if the content is easy to comprehend.",
	"completeness": "Determine how complete and comprehensive the content is. Consider if all necessary information is present, if there are missing details, and if the content fully addresses its intended purpose.",
	"accuracy":     "Evaluate the factual correctness and precision of the content. Consider if the information is accurate, if claims are supported, and if there are any factual errors.",
	"relevance":    "Assess how relevant and appropriate the content is for its intended context. Consider if the information is pertinent, if it addresses the right topics, and if it's suitable for the target audience.",
	"consistency":  "Evaluate how consistent the content is in terms of style, terminology, and approach. Consider if naming conventions are uniform, if the tone is consistent, and if similar concepts are treated similarly.",
}

// llmPromptTemplate defines the structure for LLM evaluation prompts
const llmPromptTemplate = `You are an expert evaluator tasked with assessing content quality.

Metric: %s
Description: %s

Please evaluate the following value based on the metric above:
Value: {{value}}

Provide your evaluation as a single number between 0.0 and 1.0, where:
- 0.0 = Poor/Completely fails the metric
- 0.5 = Average/Partially meets the metric  
- 1.0 = Excellent/Fully meets the metric

Return only the numeric score (e.g., 0.7) without any additional text or explanation.`

// generateLlmPrompt creates a standardized prompt for LLM evaluation
func generateLlmPrompt(metricName string) string {
	description, exists := metricDescriptions[metricName]
	if !exists {
		// Fallback description for unknown metrics
		description = fmt.Sprintf("Evaluate the content based on the '%s' metric. Consider the quality, appropriateness, and effectiveness of the content in relation to this metric.", metricName)
	}

	return fmt.Sprintf(llmPromptTemplate, strings.Title(metricName), description)
}

func ParseLlmMetric(metricName string, llmMetric *y.Yaml, negated bool, variable Variable, path pathParser.PropertyPath) (Rule, error) {
	if llmMetric.IsMap() {
		threshold, err := llmMetric.Get("threshold").Float()
		if err != nil {
			return nil, err
		}

		// Generate standardized prompt or use custom prompt if provided
		var prompt string
		if llmMetric.Get("prompt").IsFound() {
			// Use custom prompt if provided
			prompt, err = llmMetric.Get("prompt").String()
			if err != nil {
				return nil, err
			}
		} else {
			// Generate standardized prompt based on metric name
			prompt = generateLlmPrompt(metricName)
		}

		return newLlmMetric(negated, variable, path, metricName, prompt, threshold), nil
	} else {
		return nil, errors.New("llm metric rule must be a map")
	}
}

func ParseRego(code *y.Yaml, negated bool, variable Variable, path pathParser.PropertyPath) (Rule, error) {
	var regoCode string
	message := "Violation in native Rego constraint"
	s, err := code.String()
	if err == nil {
		// code fragment directly embedded in the constraint
		regoCode = s
	}

	if code.IsMap() {
		// message and code version
		s, err := code.Get("code").String()
		if err != nil {
			return nil, err
		}
		regoCode = s
		m, err := code.Get("message").String()
		if err == nil {
			message = m
		}
	}

	return newRego(negated, variable, path, regoCode, message), nil
}

func parseQualifiedNestedExpression(qNested *y.Yaml, negated bool, variable Variable, propertyPath pathParser.PropertyPath, generator *VarGenerator, op CardinalityOperation) (Rule, error) {
	count, err := qNested.Get("count").Int()
	if err != nil {
		return nil, err
	}
	cardinality := VariableCardinality{
		Operator: op,
		Value:    count,
	}

	data := qNested.Get("validation")
	if !data.IsFound() || (data.IsFound() && !data.IsMap()) {
		return nil, errors.New("map containing a validation is required in qualified atLeast/atMost constraint")
	}

	nested, err := parseNestedExpression(data, negated, variable, propertyPath, generator)
	if err != nil {
		return nil, err
	}

	switch n := nested.(type) {
	case NestedExpression:
		n.Child.Cardinality = &cardinality
		n.Child.Quantification = Exists
		return n, nil
	}

	return nil, errors.New(fmt.Sprintf("expected nested expression to build quantified nested expression but got %v", nested))
}

// We are always stringifying the value to be able to compare it easily in Rego
func stringifyNode(node *y.Yaml) (string, error) {
	s, nok := node.String()
	if nok == nil {
		return s, nil
	}

	i, nok := node.Int()
	if nok == nil {
		return strconv.Itoa(i), nil
	}

	f, nok := node.Float()
	if nok == nil {
		return strconv.FormatFloat(f, 'f', 6, 64), nil
	}

	b, nok := node.Bool()
	if nok == nil {
		return strconv.FormatBool(b), nil
	}

	l, c := node.Pos()
	return "", errors.New(fmt.Sprintf("expected scalars, found %v at [%d,%d]", node, l, c))
}

func scalarList(in []*y.Yaml) ([]string, error) {
	var acc []string
	for _, e := range in {
		s, nok := stringifyNode(e)
		if nok == nil {
			acc = append(acc, s)
			continue
		} else {
			return nil, nok
		}

	}

	return acc, nil
}
