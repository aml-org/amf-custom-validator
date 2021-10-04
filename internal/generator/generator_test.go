package generator

import (
	"context"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/test"
	"github.com/open-policy-agent/opa/rego"
	"strings"
	"testing"
)

func TestGenerated(t *testing.T) {

	for _, fix := range test.Fixtures("../../test/data/basic") {
		profile.GenReset()
		prof, err := parser.Parse(fix.ReadProfile())
		if err != nil {
			panic(err)
		}
		profile.GenReset()
		generated := Generate(*prof)
		success, err := validateRegoUnit(generated)
		if !success {
			t.Error(err)
		}
		if config.Override {
			test.ForceWrite(fix.Generated, generated.Code)
		} else {
			actual := strings.TrimSpace(generated.Code)
			expected := strings.TrimSpace(fix.ReadGenerated())

			if actual != expected {
				t.Errorf("%s > Actual did not match expected", fix.Profile)
			}
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
