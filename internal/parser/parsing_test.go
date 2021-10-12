package parser

import (
	"github.com/aml-org/amf-custom-validator/internal/config"
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

		if config.Override {
			test.ForceWrite(fix.Parsed, actual)
		}

		if actual != expected {
			t.Errorf("%s> Actual did not match expected", fix.Profile)
		}
	}
}
