package validator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
	"github.com/aml-org/amfopa/test"
	"strings"
	"testing"
)

const debug = false

func TestProduction(t *testing.T) {
	for _, fixture := range test.ProductionFixtures("../../test/data/production", nil) {
		profile := fixture.Profile()
		for _, example := range fixture.Examples() {
			filter := ""
			if strings.Index(example.File, filter) > -1 {
				//println(example.File)
				report, err := Validate(profile, example.Text, debug)

				if err != nil {
					t.Errorf("Validation failed %v", err)
				}
				if conforms(report) != example.Positive {
					t.Errorf(fmt.Sprintf("%s, %s expected conforms: %t got conforms %t\n\n%s\n", string(fixture), example.File, example.Positive, conforms(report), report))
				}
				//test.ForceWrite(example.Reportfile(), report)
				expected := example.ReadReport()
				actual := report
				if expected != actual {
					t.Errorf(fmt.Sprintf("failed report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", example.File, expected, report))
				}
			}
		}
	}
}

func TestValidate(t *testing.T) {
	//filter := "profile12"
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", nil) {
		prof := fixture.ReadProfile()
		profile.GenReset()
		report, err := Validate(prof, fixture.ReadFixturePositiveData(), debug)
		if err != nil {
			t.Errorf("positive validation failed %v", err)
		}
		if !conforms(report) {
			t.Errorf("positive case failed")
		}
		expected := strings.TrimSpace(fixture.ReadFixturePositiveReport())

		//test.ForceWrite(string(fixture)+"/positive.report.jsonld", strings.TrimSpace(report))
		if strings.TrimSpace(report) != expected {
			t.Errorf(fmt.Sprintf("failed positive report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", fixture, expected, report))
		}

		report, err = Validate(prof, fixture.ReadFixtureNegativeData(), debug)
		if err != nil {
			t.Errorf("positive validation failed %v", err)
		}
		if conforms(report) {
			t.Errorf("negative case failed")
		}
		expected = strings.TrimSpace(fixture.ReadFixtureNegativeReport())
		//test.ForceWrite(string(fixture)+"/negative.report.jsonld", strings.TrimSpace(report))
		if strings.TrimSpace(report) != expected {
			t.Errorf(fmt.Sprintf("failed positive report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", fixture, expected, report))
		}

	}
}

func conforms(report string) bool {
	return strings.Index(report, "\"http://www.w3.org/ns/shacl#conforms\": true") > -1
}
