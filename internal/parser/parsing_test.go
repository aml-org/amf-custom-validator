package parser

import (
	"github.com/aml-org/amfopa/test"
	"strings"
	"testing"
)

func TestParsed(t *testing.T) {

	for _,fix := range test.Fixtures() {
		profile, err := Parse(fix.Profile)
		if err != nil {
			panic(err)
		}
		actual := strings.TrimSpace(profile.String())
		expected := strings.TrimSpace(fix.ReadParsed())

		if actual != expected {
			t.Errorf("Error in expected profile %s\n\nActual:\n%s\n----\nExpected:\n%s", fix.Profile, actual, expected)
		}
	}
}