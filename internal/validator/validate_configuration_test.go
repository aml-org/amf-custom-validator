package validator

import (
	"encoding/json"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"testing"
)

func TestNoDateCreated(t *testing.T) {
	reportConfig := ReportConfiguration{
		IncludeReportCreationTime: false,
	}

	profile := read("../../test/data/integration/profile1/profile.yaml")
	data := read("../../test/data/integration/profile1/negative.data.jsonld")

	report, err := ValidateWithConfiguration(profile, data, config.Debug, nil, TestValidationConfiguration{}, reportConfig)

	if err != nil {
		t.Errorf("Error during validation\n")
	}

	var output []any
	err = json.Unmarshal([]byte(report), &output)
	if err != nil {
		t.Errorf("Error during report JSON unmarshling\n")
	}
	doc := output[0].(map[string]any)
	encoded := doc["doc:encodes"].([]any)
	encodedDoc := encoded[0].(map[string]any)

	_, hasKey := encodedDoc["dateCreated"]

	if hasKey {
		t.Errorf("Report contains 'dateCreated' despite being disabled\n")
	}
}
