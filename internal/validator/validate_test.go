package validator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/test"
	"strings"
	"testing"
)

const debug = false

func TestProduction(t *testing.T) {
	filter := ""
	for _, fixture := range test.ProductionFixtures("../../test/data/production", &filter) {
		profile := fixture.Profile()
		for _, example := range fixture.Examples() {
			filter := "" // put the number of the text to filter here
			if strings.Index(example.File, filter) > -1 {
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
	filter := ""
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", &filter) {
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
			t.Errorf("negative validation failed %v", err)
		}
		if conforms(report) {
			t.Errorf("negative case failed")
		}
		expected = strings.TrimSpace(fixture.ReadFixtureNegativeReport())
		//test.ForceWrite(string(fixture)+"/negative.report.jsonld", strings.TrimSpace(report))
		if strings.TrimSpace(report) != expected {
			t.Errorf(fmt.Sprintf("failed negative report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", fixture, expected, report))
		}

		lexicalFixture, fixtureError := fixture.ReadFixtureNegativeDataWithLexical()
		if fixtureError == nil {
			report, err = Validate(prof, lexicalFixture, debug)
			if err != nil {
				t.Errorf("negative validation failed %v", err)
			}
			if conforms(report) {
				t.Errorf("negative case failed")
			}
			expected = strings.TrimSpace(fixture.ReadFixtureNegativeReportWithLexical())
			if strings.TrimSpace(report) != expected {
				t.Errorf(fmt.Sprintf("failed negative lexical report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", lexicalFixture, expected, report))
			}
			//test.ForceWrite(string(fixture)+"/negative.report.lexical.jsonld", strings.TrimSpace(report))
		}
	}
}

func conforms(report string) bool {
	return strings.Index(report, "\"http://www.w3.org/ns/shacl#conforms\": true") > -1
}
