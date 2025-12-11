package data

import (
	"bytes"
	"encoding/json"

	"github.com/aml-org/amf-custom-validator/internal/notifications"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
)

func ProcessInput(jsonldText string, debug bool, receiver *chan e.Event) (any, error) {
	notifications.DispatchEvent(e.NewEvent(e.InputDataParsingStart), receiver)
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonldText)))
	decoder.UseNumber()

	var input any
	if err := decoder.Decode(&input); err != nil {
		return "", nil
	}
	notifications.DispatchEvent(e.NewEvent(e.InputDataParsingDone), receiver)

	notifications.DispatchEvent(e.NewEvent(e.InputDataNormalizationStart), receiver)
	normalizedInput := Index(Normalize(input))
	notifications.DispatchEvent(e.NewEvent(e.InputDataNormalizationDone), receiver)

	return normalizedInput, nil
}
