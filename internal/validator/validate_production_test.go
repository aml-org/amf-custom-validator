package validator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/test"
	"strings"
	"testing"
)

func TestProduction(t *testing.T) {
	filter := "json-api"
	for _, fixture := range test.ProductionFixtures("../../test/data/production", &filter) {
		profile := fixture.Profile()
		for _, example := range fixture.Examples() {
			filter := "" // put the number of the text to filter here
			if strings.Index(example.File, filter) > -1 {
				report, err := Validate(profile, example.Text, config.Debug, nil)
				if err != nil {
					t.Errorf("Validation failed %v", err)
				}
				if conforms(report) != example.Positive {
					t.Errorf(fmt.Sprintf("%s, %s expected conforms: %t got conforms %t\n", string(fixture), example.File, example.Positive, conforms(report)))
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
