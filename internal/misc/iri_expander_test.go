package misc

import (
	"github.com/aml-org/amf-custom-validator/internal/types"
	"testing"
)

func TestExpandCompactIriWithDot(t *testing.T) {
	expander := IriExpander{
		Context: map[string]types.Object{
			"aml": "https://a.ml#",
		},
	}
	actual, err := expander.Expand("aml.vocabulary")
	if err != nil {
		t.Errorf("Not expecting error: %s", err)
	}

	expected := "https://a.ml#vocabulary"

	if actual != expected {
		t.Errorf("Actual did not match expeted\n\tExpected: %s\n\tActual: %s", expected, actual)
	}
}

func TestExpandCompactIriWithColon(t *testing.T) {
	expander := IriExpander{
		Context: map[string]types.Object{
			"aml": "https://a.ml#",
		},
	}
	actual, err := expander.Expand("aml:vocabulary")
	if err == nil {
		t.Errorf("Error expected: %s is not in compact form", actual)
	}
}

func TestExpandNonCompactIri(t *testing.T) {
	expander := IriExpander{
		Context: map[string]types.Object{
			"aml": "https://a.ml#",
		},
	}
	actual, err := expander.Expand("https://example.org#vocabulary")
	if err == nil {
		t.Errorf("Error expected: %s is not in compact form", actual)
	}

}

func TestExpandCompactIriWithUnknownTerm(t *testing.T) {
	expander := IriExpander{
		Context: map[string]types.Object{
			"aml": "https://a.ml#",
		},
	}
	_, err := expander.Expand("example.vocabulary")
	if err == nil {
		t.Errorf("Expecting error when expanding unknown term")
	}
}

func TestMustNotExpandNorFailWithReservedKeyword(t *testing.T) {
	expander := IriExpander{
		Context: map[string]types.Object{
			"aml": "https://a.ml#",
		},
	}
	actual, err := expander.Expand("@type")
	if err != nil {
		t.Errorf("Not expecting error: %s", err)
	}

	expected := "@type"

	if actual != expected {
		t.Errorf("Actual did not match expeted\n\tExpected: %s\n\tActual: %s", expected, actual)
	}
}
