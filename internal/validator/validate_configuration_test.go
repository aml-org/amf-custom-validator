package validator

import (
	"encoding/json"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/internal/validator/contexts"
	config3 "github.com/aml-org/amf-custom-validator/pkg/config"
	"testing"
)

func TestNoDateCreated(t *testing.T) {
	reportConfig := config3.ReportConfiguration{
		IncludeReportCreationTime: false,
	}

	profile := read("../../test/data/integration/profile1/profile.yaml")
	data := read("../../test/data/integration/profile1/negative.data.jsonld")

	report, err := ValidateWithConfiguration(profile, data, config.Debug, nil, config3.TestValidationConfiguration{}, reportConfig)

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

func TestAlternativeSchemas(t *testing.T) {
	reportSchemaIri := "http://a.ml/report"
	lexicalSchemaIri := "http://a.ml/lexical"

	reportConfig := config3.ReportConfiguration{
		ReportSchemaIri:  reportSchemaIri,
		LexicalSchemaIri: lexicalSchemaIri,
	}

	profile := read("../../test/data/integration/profile1/profile.yaml")
	data := read("../../test/data/integration/profile1/negative.data.jsonld")

	report, err := ValidateWithConfiguration(profile, data, config.Debug, nil, config3.TestValidationConfiguration{}, reportConfig)

	if err != nil {
		t.Errorf("Error during validation\n")
	}

	var output []any
	err = json.Unmarshal([]byte(report), &output)
	if err != nil {
		t.Errorf("Error during report JSON unmarshling\n")
	}
	doc := output[0].(map[string]any)
	context := doc["@context"].(map[string]any)

	actual := context["reportSchema"].(string)
	expected := contexts.DeclarationsFrom(reportSchemaIri)

	if actual != expected {
		t.Errorf("Actual '%s' does not match expected '%s'", actual, expected)
	}

	actual = context["lexicalSchema"].(string)
	expected = contexts.DeclarationsFrom(lexicalSchemaIri)

	if actual != expected {
		t.Errorf("Actual '%s' does not match expected '%s'", actual, expected)
	}
}
