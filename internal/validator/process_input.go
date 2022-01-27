package validator

import (
	"bytes"
	"encoding/json"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
)

func ProcessInput(jsonldText string, debug bool, receiver *chan e.Event) (interface{}, error) {
	dispatchEvent(e.NewEvent(e.InputDataParsingStart), receiver)
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonldText)))
	decoder.UseNumber()

	var input interface{}
	if err := decoder.Decode(&input); err != nil {
		return "", nil
	}
	dispatchEvent(e.NewEvent(e.InputDataParsingDone), receiver)

	dispatchEvent(e.NewEvent(e.InputDataNormalizationStart), receiver)
	normalizedInput := Index(Normalize(input))
	dispatchEvent(e.NewEvent(e.InputDataNormalizationDone), receiver)

	if debug {
		println("Input data")
		println("-------------------------------")
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "  ")
		err := enc.Encode(normalizedInput)
		if err != nil {
			panic(err)
		}
		println(b.String())
	}
	return normalizedInput, nil
}
