package parser

import (
	"github.com/aml-org/amf-custom-validator/test"
	"strings"
	"testing"
)

func TestParsed(t *testing.T) {

	for _, fix := range test.Fixtures("../../test/data/basic") {
		profile, err := Parse(fix.ReadProfile())
		if err != nil {
			panic(err)
		}
		actual := strings.TrimSpace(profile.String())
		expected := strings.TrimSpace(fix.ReadParsed())

		//test.ForceWrite(fix.Parsed, actual)

		if actual != expected {
			t.Errorf("Error in expected profile %s\n\nActual:\n%s\n----\nExpected:\n%s", fix.Profile, actual, expected)
		}
	}
}
