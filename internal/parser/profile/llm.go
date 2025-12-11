package profile

import (
	"fmt"
	"sync/atomic"

	"github.com/aml-org/amf-custom-validator/internal/parser/path"
)

type LlmMetric struct {
	AtomicStatement
	MetricName  string
	Threshold   float64
	TargetClass string
	Prompt      string
	Id          string
}

var globalCounter int64 // package-level global variable

// NextID atomically increments and returns the next counter value
func NextID() int64 {
	return atomic.AddInt64(&globalCounter, 1)
}
func GenerateLlmMetricIdentifier(metric LlmMetric) LlmMetric {
	c := NextID()
	id := fmt.Sprintf("llm-metric-%d", c)

	metric.Id = id

	return metric
}

func (l LlmMetric) Negate() Rule {
	negated := l
	negated.Negated = !l.Negated
	return negated
}

func (l LlmMetric) String() string {
	var negation = ""
	if l.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s%s-%s(%s,'%s',%s,'%s')", negation, l.Name, l.MetricName, l.Variable.Name, l.Path.Source(), l.Threshold, l.Prompt)
}

func newLlmMetric(negated bool, variable Variable, path path.PropertyPath, metricName string, prompt string, threshold float64) LlmMetric {
	return LlmMetric{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Name:    "llmMetric",
				Negated: negated,
			},
			Variable: variable,
			Path:     path,
		},
		MetricName: metricName,
		Threshold:  threshold,
		Prompt:     prompt,
	}
}
