package llm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aml-org/amf-custom-validator/internal/data"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/internal/types"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
)

func EnrichDataWithLlmMetrics(metrics []profile.LlmMetric, jsonldText string, debug bool, eventChan *chan e.Event, iriExpander *misc.IriExpander) map[string]bool {
	normalizedInput, _ := data.ProcessInput(jsonldText, debug, eventChan)
	results := make(map[string]bool)
	
	// Initialize LLM client
	llmClient, err := NewOpenAILLM()
	if err != nil {
		if debug {
			fmt.Printf("Failed to initialize LLM client: %v\n", err)
		}
		return results
	}
	
	for _, metric := range metrics {
		targetValues := findPropertyValues(normalizedInput, metric.TargetClass, metric.Path, iriExpander)
		
		// Process each target value found
		for _, targetValue := range targetValues {
			passed := evaluateMetricWithLLM(llmClient, metric, targetValue, debug)
			
			// Store result using metric ID as key
			results[metric.Id] = passed
		}
	}
	
	return results
}

func findPropertyValues(input interface{}, class string, path path.PropertyPath, iriExpander *misc.IriExpander) []interface{} {
	classes := input.(types.ObjectMap)["@types"].(types.ObjectMap)
	nodes := input.(types.ObjectMap)["@ids"].(types.ObjectMap)

	classUri, _ := iriExpander.Expand(class)
	idsForClasses := classes[classUri].([]string)

	var acc []interface{}
	for _, nodeId := range idsForClasses {
		node, _ := nodes[nodeId]
		expandedPath := expandPath(path, iriExpander)
		data := fetchData(node, expandedPath)
		acc = append(acc, data...)
	}

	return acc
}

func fetchData(node interface{}, expandedPath []string) []interface{} {
	acc := []interface{}{node}
	for _, p := range expandedPath {
		var innAcc []interface{}
		for _, n := range acc {
			nn := n.(types.ObjectMap)[p]
			switch nn.(type) {
			case []interface{}:
				innAcc = append(innAcc, nn.([]interface{})...)
			default:
				innAcc = append(innAcc, nn)
			}
		}
		acc = innAcc
	}
	return acc
}

func expandPath(sourcePath path.PropertyPath, expander *misc.IriExpander) []string {
	acc := []string{}

	switch p := sourcePath.(type) {
	case path.Property:
		var ep string
		if p.IsCustom(expander) {
			ep, _ = p.CustomName(expander)
		} else {
			ep, _ = p.Expanded(expander)
		}

		acc = append(acc, ep)

	default:
		break
	}

	return acc
}

// evaluateMetricWithLLM performs LLM evaluation for a single metric and target value
func evaluateMetricWithLLM(llmClient *OpenAILLM, metric profile.LlmMetric, targetValue interface{}, debug bool) bool {
	ctx := context.Background()
	
	// Build the prompt by substituting the target value into the metric prompt
	prompt := buildPrompt(metric.Prompt, targetValue)
	
	// Call LLM with the constructed prompt
	response, err := llmClient.Call(ctx, prompt)
	if err != nil {
		if debug {
			fmt.Printf("LLM call failed for metric %s: %v\n", metric.MetricName, err)
		}
		return false // Default to false on error
	}
	
	// Evaluate the LLM response against the threshold
	passed := evaluateResponse(response, metric.Threshold)
	
	if debug {
		fmt.Printf("Metric %s: prompt='%s', response='%s', threshold=%.2f, passed=%t\n", 
			metric.MetricName, prompt, response, metric.Threshold, passed)
	}
	
	return passed
}

// buildPrompt constructs the LLM prompt by substituting the target value
func buildPrompt(promptTemplate string, targetValue interface{}) string {
	// Convert target value to string representation
	valueStr := fmt.Sprintf("%v", targetValue)
	
	// Simple placeholder replacement - you can enhance this with more sophisticated templating
	// For now, we'll replace {{value}} or {value} placeholders with the actual value
	prompt := strings.ReplaceAll(promptTemplate, "{{value}}", valueStr)
	prompt = strings.ReplaceAll(prompt, "{value}", valueStr)
	
	return prompt
}

// evaluateResponse processes the LLM response and compares it against the threshold
func evaluateResponse(response string, threshold float64) bool {
	// Clean up the response
	response = strings.TrimSpace(response)
	
	// Try to parse the response as a number (score)
	if score, err := strconv.ParseFloat(response, 64); err == nil {
		return score >= threshold
	}
	
	// If not a number, look for common positive/negative indicators
	responseLower := strings.ToLower(response)
	
	// Check for explicit boolean responses
	if strings.Contains(responseLower, "true") || strings.Contains(responseLower, "yes") || 
	   strings.Contains(responseLower, "pass") || strings.Contains(responseLower, "valid") {
		return threshold <= 1.0 // Pass if threshold is 1.0 or less
	}
	
	if strings.Contains(responseLower, "false") || strings.Contains(responseLower, "no") || 
	   strings.Contains(responseLower, "fail") || strings.Contains(responseLower, "invalid") {
		return threshold <= 0.0 // Only pass if threshold is 0.0 or less
	}
	
	// Try to extract a numeric score from the response using regex-like approach
	words := strings.Fields(response)
	for _, word := range words {
		if score, err := strconv.ParseFloat(word, 64); err == nil {
			return score >= threshold
		}
	}
	
	// Default: if we can't determine the score, assume it doesn't meet the threshold
	return false
}
