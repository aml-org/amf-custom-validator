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

func TestShaclTckPositiveData(t *testing.T) {
	filter := ""
	for _, fixture := range test.ShaclTckFixtures("../../test/data/shacl-tck", &filter) {
		if !fixture.IsIgnored() {
			prof := fixture.ReadProfile()
			profile.GenReset()
			eventsChan := make(chan events.Event)
			milestonesChan := make(chan milestones.Milestone)
			go milestones.GenerateMilestonesFromEvents(&eventsChan, &milestonesChan)
			go printMilestones(fixture, "Positive", &milestonesChan)
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
}

func TestShaclTckNegativeData(t *testing.T) {
	filter := ""
	for _, fixture := range test.ShaclTckFixtures("../../test/data/shacl-tck", &filter) {
		if !fixture.IsIgnored() {
			prof := fixture.ReadProfile()
			profile.GenReset()
			eventsChan := make(chan events.Event)
			milestonesChan := make(chan milestones.Milestone)
			go milestones.GenerateMilestonesFromEvents(&eventsChan, &milestonesChan)
			go printMilestones(fixture, "Negative", &milestonesChan)
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
}
