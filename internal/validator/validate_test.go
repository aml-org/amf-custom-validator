package validator

import (
	"fmt"
	"github.com/aml-org/amfopa/test"
	"strings"
	"testing"
)

const debug = false

func TestValidate(t *testing.T) {
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration") {
		profile := fixture.ReadProfile()

		report, err := Validate(profile, fixture.ReadFixturePositiveData(), debug)
		if err != nil {
			t.Errorf("positive validation failed %v", err)
		}
		//test.ForceWrite(string(fixture) + "/negative.report.jsonld", strings.TrimSpace(report))
		expected := strings.TrimSpace(fixture.ReadFixturePositiveReport())
		if strings.TrimSpace(report) != expected {
			t.Errorf(fmt.Sprintf("failed positive report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", fixture, expected, report))
		}

		report, err = Validate(profile, fixture.ReadFixtureNegativeData(), debug)
		if err != nil {
			t.Errorf("positive validation failed %v", err)
		}
		expected = strings.TrimSpace(fixture.ReadFixtureNegativeReport())
		//test.ForceWrite(string(fixture) + "/negative.report.jsonld", strings.TrimSpace(report))
		if strings.TrimSpace(report) != expected {
			t.Errorf(fmt.Sprintf("failed positive report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", fixture, expected, report))
		}
	}
}
