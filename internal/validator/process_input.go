package validator

import (
	"bytes"
	"encoding/json"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
)

func ProcessInput(jsonldText string, debug bool, receiver *chan e.Event) (any, error) {
	dispatchEvent(e.NewEvent(e.InputDataParsingStart), receiver)
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonldText)))
	decoder.UseNumber()

	var input any
	if err := decoder.Decode(&input); err != nil {
		return "", nil
	}
	dispatchEvent(e.NewEvent(e.InputDataParsingDone), receiver)

	dispatchEvent(e.NewEvent(e.InputDataNormalizationStart), receiver)
	normalizedInput := Index(Normalize(input))
	dispatchEvent(e.NewEvent(e.InputDataNormalizationDone), receiver)

	return normalizedInput, nil
}
