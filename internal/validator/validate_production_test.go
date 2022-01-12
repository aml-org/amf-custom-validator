package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/test"
	"strings"
	"testing"
)

type ValidationDump struct {
	inputJson string
	inputRego string
}

func Dump(profileText string, jsonldText string)  (ValidationDump, error) {
	err, module := Generate(profileText, false, nil)
	if err != nil {
		return ValidationDump{}, err
	}
	normalizedInput := Normalize_(jsonldText, module, false, nil)
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	err = enc.Encode(normalizedInput)
	if err != nil {
		return ValidationDump{}, err
	}
	return ValidationDump{
		inputJson: b.String(),
		inputRego: module.Code,
	}, nil
}


func TestProduction(t *testing.T) {
	filter := "json-api"
	for _, fixture := range test.ProductionFixtures("../../test/data/production", &filter) {
		profile := fixture.Profile()
		for _, example := range fixture.Examples() {
			filter := "positive1.yaml" // put the number of the text to filter here
			if strings.Index(example.File, filter) > -1 {
				report, err := Validate(profile, example.Text, config.Debug, nil)
				if err != nil {
					t.Errorf("Validation failed at %s: %v", example.File, err)
				}
				if conforms(report) != example.Positive {
					t.Errorf(fmt.Sprintf("%s, %s expected conforms: %t got conforms %t\n", string(fixture), example.File, example.Positive, conforms(report)))
					if config.Dump {
						dump, _ := Dump(profile, example.Text)
						test.ForceWrite(fmt.Sprintf("%s.rego", example.File), dump.inputRego)
						test.ForceWrite(fmt.Sprintf("%s.json", example.File), dump.inputJson)
					}
				}

				if config.Override {
					test.ForceWrite(example.Reportfile(), report)
				} else {
					expected := example.ReadReport()
					actual := report
					if expected != actual {
						t.Errorf(fmt.Sprintf("failed report for %s", example.File))
					}
				}
			}
		}
	}
}
