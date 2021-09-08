package validator

import (
	"fmt"
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
				report, err := Validate(profile, example.Text, debug, nil)

				if err != nil {
					t.Errorf("Validation failed %v", err)
				}
				if conforms(report) != example.Positive {
					t.Errorf(fmt.Sprintf("%s, %s expected conforms: %t got conforms %t\n\n%s\n", string(fixture), example.File, example.Positive, conforms(report), report))
				}

				// test.ForceWrite(example.Reportfile(), report)
				expected := example.ReadReport()
				actual := report
				if expected != actual {
					t.Errorf(fmt.Sprintf("failed report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", example.File, expected, report))
				}
			}
		}
	}
}
