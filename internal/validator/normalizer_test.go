package validator

import (
	"bytes"
	"encoding/json"
	"testing"

	data2 "github.com/aml-org/amf-custom-validator/internal/data"
	"github.com/aml-org/amf-custom-validator/test"
)

func TestNormalize(t *testing.T) {
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", nil) {
		data := fixture.ReadFixturePositiveData()
		jsonData := decode(data)
		normalized := data2.Normalize(jsonData)
		indexed := data2.Index(normalized)
		encodeJson(indexed)
		//println(res)
	}
}

func decode(text string) any {
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(text)))
	decoder.UseNumber()

	var input any
	if err := decoder.Decode(&input); err != nil {
		panic(err)
	}

	return input
}

func encodeJson(data any) string {
	s := ""
	bs := bytes.NewBufferString(s)
	encoder := json.NewEncoder(bs)
	_ = encoder.Encode(data)
	return bs.String()
}
