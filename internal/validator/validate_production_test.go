package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	"github.com/aml-org/amf-custom-validator/test"
	"strings"
	"testing"
)

func TestProduction(t *testing.T) {
	filter := ""
	for _, fixture := range test.ProductionFixtures("../../test/data/production", &filter) {
		profile := fixture.Profile()
		for _, example := range fixture.Examples() {
			filter := "" // put the number of the text to filter here
			if strings.Index(example.File, filter) > -1 {
				report, err := ValidateWithConfiguration(profile, example.Text, config.Debug, nil, TestValidationConfiguration{})
				if config.Debug {
					generateDebugOutput(profile, example)
				}
				if err != nil {
					t.Errorf("Validation failed %v", err)
				}

				if report == "" {
					t.Error("Empty report")
				} else {
					conformant := isReportConformant(report)
					if conformant == example.Positive {
						if config.Override {
							test.ForceWrite(example.Reportfile(), report)
						} else {
							expected := example.ReadReport()
							actual := report
							if expected != actual {
								t.Errorf(fmt.Sprintf("Not matching report for %s", example.File))
							}
						}
					} else {
						t.Errorf("Generated report (positive==%v) not matching report (conformat==%v) for example %s", example.Positive, conformant, example.File)
					}
				}
			}
		}
	}
}

func generateDebugOutput(profile string, example test.ProductionExample) {
	parsedProfile, err := parser.Parse(profile)
	if err == nil {
		test.ForceWrite(example.Folfile(), parsedProfile.String())
	}
	compiledRego, err := GenerateRego(profile, config.Debug, nil)
	if err == nil {
		test.ForceWrite(example.Regofile(), compiledRego.Code)
	}

	normalizedInput, err := ProcessInput(example.Text, config.Debug, nil)
	if err == nil {
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "  ")
		err := enc.Encode(normalizedInput)
		if err == nil {
			test.ForceWrite(example.Normalizedinputfile(), b.String())
		}
	}
}

func isReportConformant(report string) bool {
	var output []any
	json.Unmarshal([]byte(report), &output)
	doc := output[0].(map[string]any)
	encoded := doc["doc:encodes"].([]any)
	encodedDoc := encoded[0].(map[string]any)
	return encodedDoc["conforms"].(bool)
}
