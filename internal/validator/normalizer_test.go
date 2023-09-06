package validator

import (
	"bytes"
	"encoding/json"
	"github.com/aml-org/amf-custom-validator/test"
	"testing"
)

func TestNormalize(t *testing.T) {
	for _, fixture := range test.IntegrationFixtures("../../test/data/integration", nil) {
		data := fixture.ReadFixturePositiveData()
		jsonData := decode(data)
		normalized := Normalize(jsonData)
		indexed := Index(normalized)
		encodeJson(indexed)
		//println(res)
	}
}

func decode(text string) interface{} {
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(text)))
	decoder.UseNumber()

	var input interface{}
	if err := decoder.Decode(&input); err != nil {
		panic(err)
	}

	return input
}

func encodeJson(data interface{}) string {
	s := ""
	bs := bytes.NewBufferString(s)
	encoder := json.NewEncoder(bs)
	_ = encoder.Encode(data)
	return bs.String()
}
