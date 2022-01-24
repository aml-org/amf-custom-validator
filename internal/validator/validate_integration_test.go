package validator

import (
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/test"
	"strings"
	"testing"
)

func TestIntegrationPositiveData(t *testing.T) {
	filter := ""
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", &filter) {
		prof := fixture.ReadProfile()
		profile.GenReset()
		report, err := Validate(prof, fixture.ReadFixturePositiveData(), config.Debug, nil)

		if err != nil {
			t.Errorf("%s > Positive case > Failed with error %v", fixture, err)
		}

		if config.Override {
			test.ForceWrite(string(fixture)+"/positive.report.jsonld", strings.TrimSpace(report))
		} else {
			expected := strings.TrimSpace(fixture.ReadFixturePositiveReport())
			if strings.TrimSpace(report) != expected {
				t.Errorf("%s > Positive case > Actual did not match expected", fixture)
			}
		}
	}
}

func TestIntegrationNegativeData(t *testing.T) {
	filter := ""
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", &filter) {
		prof := fixture.ReadProfile()
		profile.GenReset()
		report, err := Validate(prof, fixture.ReadFixtureNegativeData(), config.Debug, nil)
		if err != nil {
			t.Errorf("%s > Negative case > Failed with error %v", fixture, err)
		}

		if config.Override {
			test.ForceWrite(string(fixture)+"/negative.report.jsonld", strings.TrimSpace(report))
		} else {
			expected := strings.TrimSpace(fixture.ReadFixtureNegativeReport())
			if strings.TrimSpace(report) != expected {
				t.Errorf("%s > Negative case > Actual did not match expected", fixture)
			}
		}
	}
}

func TestIntegrationNegativeDataWithLexical(t *testing.T) {
	filter := ""
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", &filter) {
		prof := fixture.ReadProfile()
		profile.GenReset()

		lexicalFixture, fixtureError := fixture.ReadFixtureNegativeDataWithLexical()
		if fixtureError == nil {
			report, err := Validate(prof, lexicalFixture, config.Debug, nil)
			if err != nil {
				t.Errorf("%s > Negative case with lexical > Failed with error %v", fixture, err)
			}

			if config.Override {
				test.ForceWrite(string(fixture)+"/negative.report.lexical.jsonld", strings.TrimSpace(report))
			} else {
				expected := strings.TrimSpace(fixture.ReadFixtureNegativeReportWithLexical())
				if strings.TrimSpace(report) != expected {
					t.Errorf("%s > Negative case with lexical > Actual did not match expected", fixture)
				}
			}
		}
	}
}
