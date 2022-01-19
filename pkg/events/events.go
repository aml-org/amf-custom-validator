package events

import "time"

type EventType int

const (
	ProfileParsingStart EventType = iota
	ProfileParsingDone
	InputDataParsingStart
	InputDataParsingDone
	InputDataNormalizationStart
	InputDataNormalizationDone
	RegoGenerationStart
	RegoGenerationDone
	RegoCompilationStart
	RegoCompilationDone
	OpaValidationStart
	OpaValidationDone
	BuildReportStart
	BuildReportDone
)

type Event struct {
	EventType EventType
	Time      time.Time
}

func NewEvent(eventType EventType) Event {
	return Event{eventType, time.Now()}
}
