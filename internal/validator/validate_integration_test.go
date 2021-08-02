package validator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/aml-org/amf-custom-validator/pkg/milestones"
	"github.com/aml-org/amf-custom-validator/test"
	"strings"
	"testing"
)

func TestIntegrationPositiveData(t *testing.T) {
	filter := ""
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", &filter) {
		prof := fixture.ReadProfile()
		profile.GenReset()
		eventsChan := make(chan events.Event)
		milestonesChan := make(chan milestones.Milestone)
		go milestones.GenerateMilestonesFromEvents(&eventsChan, &milestonesChan)
		go printMilestones(fixture,"Positive", &milestonesChan)
		report, err := Validate(prof, fixture.ReadFixturePositiveData(), debug, &eventsChan)

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
	}
}

func TestIntegrationNegativeData(t *testing.T) {
	filter := ""
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", &filter) {
		prof := fixture.ReadProfile()
		profile.GenReset()
		eventsChan := make(chan events.Event)
		milestonesChan := make(chan milestones.Milestone)
		go milestones.GenerateMilestonesFromEvents(&eventsChan, &milestonesChan)
		go printMilestones(fixture,"Negative", &milestonesChan)
		report, err := Validate(prof, fixture.ReadFixtureNegativeData(), debug, &eventsChan)
		if err != nil {
			t.Errorf("negative validation failed %v", err)
		}
		if conforms(report) {
			t.Errorf("negative case failed")
		}
		expected := strings.TrimSpace(fixture.ReadFixtureNegativeReport())
		//test.ForceWrite(string(fixture)+"/negative.report.jsonld", strings.TrimSpace(report))
		if strings.TrimSpace(report) != expected {
			t.Errorf(fmt.Sprintf("failed negative report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", fixture, expected, report))
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
			eventsChan := make(chan events.Event)
			milestonesChan := make(chan milestones.Milestone)
			go milestones.GenerateMilestonesFromEvents(&eventsChan, &milestonesChan)
			go printMilestones(fixture, "Negative+Lexical", &milestonesChan)
			report, err := Validate(prof, lexicalFixture, debug, &eventsChan)
			if err != nil {
				t.Errorf("negative validation failed %v", err)
			}
			if conforms(report) {
				t.Errorf("negative case failed")
			}
			expected := strings.TrimSpace(fixture.ReadFixtureNegativeReportWithLexical())
			if strings.TrimSpace(report) != expected {
				t.Errorf(fmt.Sprintf("failed negative lexical report for %s\n-------------Expected:\n%s\n-------------Actual:\n%s\n", lexicalFixture, expected, report))
			}
			//test.ForceWrite(string(fixture)+"/negative.report.lexical.jsonld", strings.TrimSpace(report))
		}
	}
}
