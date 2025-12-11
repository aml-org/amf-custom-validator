package notifications

import e "github.com/aml-org/amf-custom-validator/pkg/events"

func DispatchEvent(event e.Event, eventChan *chan e.Event) {
	if eventChan != nil {
		*eventChan <- event
	}
}

func CloseEventChan(eventChan *chan e.Event) {
	if eventChan != nil {
		close(*eventChan)
	}
}
