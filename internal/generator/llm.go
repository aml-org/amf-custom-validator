package generator

import (
	"fmt"

	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

func GenerateLlmMetric(llmMetric profile.LlmMetric, iriExpander *misc.IriExpander, llmMetrics *[]profile.LlmMetric) []SimpleRegoResult {
	llmMetric = profile.GenerateLlmMetricIdentifier(llmMetric)
	// Collect the LLM metric if a collector is provided
	if llmMetrics != nil {
		*llmMetrics = append(*llmMetrics, llmMetric)
	}

	var rego []string

	rego = append(rego, "#  inserting LLM Metric check for : "+llmMetric.Id)
	if llmMetric.Negated {
		rego = append(rego, fmt.Sprintf("data.LLM_VALUES[%s] == false", llmMetric.Id))
	} else {
		rego = append(rego, fmt.Sprintf("data.LLM_VALUES[%s] == true", llmMetric.Id))
	}

	tracePath, err := llmMetric.Path.Trace(iriExpander)
	if err != nil {
		panic(err)
	}

	r := SimpleRegoResult{
		Constraint: "llmMetric",
		Rego:       rego,
		PathRules:  []RegoPathResult{},
		Path:       tracePath,
		TraceNode:  llmMetric.Variable.Name,
		TraceValue: BuildTraceValueNode(
			fmt.Sprintf("\"negated\":%t,\"llmdData\":\"data.LLM_VALUES[%s]\"", llmMetric.Negated, llmMetric.Id)),
		Variable: llmMetric.Id,
	}
	return []SimpleRegoResult{r}
}
