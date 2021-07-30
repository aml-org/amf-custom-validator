package milestones

import (
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"time"
)

type Operation string

const (
	ProfileParsing         Operation = "ProfileParsing"
	InputDataParsing       Operation = "InputDataParsing"
	InputDataNormalization Operation = "InputDataNormalization"
	RegoGeneration         Operation = "RegoGeneration"
	OpaValidation          Operation = "OpaValidation"
	BuildReport            Operation = "BuildReport"
)

type Milestone struct {
	Operation Operation
	Start     time.Time
	Duration  time.Duration
}

func GenerateMilestonesFromEvents(eventChan *chan e.Event, milestoneChan *chan Milestone) {
	startEvents := make(map[e.EventType]e.Event)
	for event := range *eventChan {
		switch eventType := event.EventType; eventType {
		case e.ProfileParsingStart, e.InputDataParsingStart, e.InputDataNormalizationStart, e.RegoGenerationStart, e.OpaValidationStart, e.BuildReportStart:
			startEvents[eventType] = event
		case e.ProfileParsingDone:
			start := startEvents[e.ProfileParsingStart]
			end := event
			*milestoneChan <- generateMilestone(ProfileParsing, start, end)
		case e.InputDataParsingDone:
			start := startEvents[e.InputDataParsingStart]
			end := event
			*milestoneChan <- generateMilestone(InputDataParsing, start, end)
		case e.InputDataNormalizationDone:
			start := startEvents[e.InputDataNormalizationStart]
			end := event
			*milestoneChan <- generateMilestone(InputDataNormalization, start, end)
		case e.RegoGenerationDone:
			start := startEvents[e.RegoGenerationStart]
			end := event
			*milestoneChan <- generateMilestone(RegoGeneration, start, end)
		case e.OpaValidationDone:
			start := startEvents[e.OpaValidationStart]
			end := event
			*milestoneChan <- generateMilestone(OpaValidation, start, end)
		case e.BuildReportDone:
			start := startEvents[e.BuildReportStart]
			end := event
			*milestoneChan <- generateMilestone(BuildReport, start, end)
		}
	}
	close(*milestoneChan)
}

func generateMilestone(operation Operation, start, end e.Event) Milestone {
	return Milestone{
		Operation: operation,
		Start:     start.Time,
		Duration:  end.Time.Sub(start.Time),
	}
}
