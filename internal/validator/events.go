package validator

import e "github.com/aml-org/amf-custom-validator/pkg/events"

func dispatchEvent(event e.Event, eventChan *chan e.Event) {
	if eventChan != nil {
		*eventChan <- event
	}
}

func CloseEventChan(eventChan *chan e.Event) {
	if eventChan != nil {
		close(*eventChan)
	}
}
