package generator

import (
	"context"
	"github.com/aml-org/amfopa/internal/parser"
	"github.com/aml-org/amfopa/test"
	"github.com/open-policy-agent/opa/rego"
	"strings"
	"testing"
)

func TestGenerated(t *testing.T) {

	for _, fix := range test.Fixtures("../../test/data") {
		profile, err := parser.Parse(fix.ReadProfile())
		if err != nil {
			panic(err)
		}
		generated := Generate(*profile)

		succes, err := validateRegoUnit(generated)
		if !succes {
			t.Error(err)
		}
		//test.ForceWrite(fix.Generated, generated.Code)
		actual := strings.TrimSpace(generated.Code)
		expected := strings.TrimSpace(fix.ReadGenerated())

		if actual != expected {
			t.Errorf("Error in expected profile %s\n\nActual:\n%s\n----\nExpected:\n%s", fix.Profile, actual, expected)
		}
	}
}

func validateRegoUnit(module RegoUnit) (bool, error) {
	validator := rego.New(
		rego.Query("data."+module.Name+"."+module.Entrypoint),
		rego.Module(module.Name+".rego", module.Code),
	)
	ctx := context.Background()
	_, err := validator.Eval(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
