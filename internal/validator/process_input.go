package validator

import (
	e "github.com/aml-org/amf-custom-validator/pkg/events"
)

func ProcessInput(jsonldText string, debug bool, receiver *chan e.Event, optimizer *Optimizer) (interface{}, error) {
	dispatchEvent(e.NewEvent(e.InputDataParsingStart), receiver)
	input := ParseJson(jsonldText)
	dispatchEvent(e.NewEvent(e.InputDataParsingDone), receiver)

	// TODO add events here
	if optimizer != nil {
		optimizer.Optimize(&input)
	}

	dispatchEvent(e.NewEvent(e.InputDataNormalizationStart), receiver)
	normalizedInput := Index(Normalize(input))
	dispatchEvent(e.NewEvent(e.InputDataNormalizationDone), receiver)

	return normalizedInput, nil
}
