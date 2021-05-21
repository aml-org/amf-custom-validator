package validator

import (
	"fmt"
	"github.com/aml-org/amfopa/test"
	"strings"
	"testing"
)

const debug = false

//const filter = "profile6"

func TestValidate(t *testing.T) {

	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", nil) {
		profile := fixture.ReadProfile()

		report, err := Validate(profile, fixture.ReadFixturePositiveData(), debug)
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

		report, err = Validate(profile, fixture.ReadFixtureNegativeData(), debug)
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
