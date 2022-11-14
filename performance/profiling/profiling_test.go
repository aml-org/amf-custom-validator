package profiling

import (
	"github.com/aml-org/amf-custom-validator/internal/config"
	p "github.com/aml-org/amf-custom-validator/internal/parser/profile"
	v "github.com/aml-org/amf-custom-validator/internal/validator"
	"io/ioutil"
	"testing"
)

func TestCase(t *testing.T) {
	profile := "../../test/data/integration/profile1/profile.yaml"
	data := "../../test/data/integration/profile1/positive.data.lexical.jsonld"
	report, err := validate(profile, data)
	if err != nil {
		t.Errorf("Found error: '%s'", err)
	} else if len(report) <= 0 {
		t.Errorf("Report is empty")
	}
}

func read(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func validate(profile string, data string) (string, error) {
	p.GenReset()
	profileText := read(profile)
	dataText := read(data)
	return v.Validate(profileText, dataText, config.Debug, nil)
}
