package generator

import (
	"context"
	"strings"
	"testing"

	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/test"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func regoVersionFromBool(isV0 bool) ast.RegoVersion {
	if isV0 {
		return ast.RegoV0
	}
	return ast.RegoV1
}

func TestGenerated(t *testing.T) {

	for _, fix := range test.Fixtures("../../test/data/basic") {
		profile.GenReset()
		prof, err := parser.Parse(fix.ReadProfile())
		if err != nil {
			panic(err)
		}
		testProfile(prof, t, fix, false)
		testProfile(prof, t, fix, true)
	}
}

func testProfile(prof *profile.Profile, t *testing.T, fix test.Fixture, isV0 bool) {
	profile.GenReset()
	regoVersion := regoVersionFromBool(isV0)
	generated := Generate(*prof, regoVersion)
	success, err := validateRegoUnit(generated, isV0)
	if !success {
		println(generated.Code)
		t.Error(err)
	}
	if config.Override {
		if isV0 {
			test.ForceWrite(fix.Generated, generated.Code)
		} else {
			test.ForceWrite(fix.GeneratedV1, generated.Code)
		}
	} else {
		actual := strings.TrimSpace(generated.Code)
		var expected = ""
		if isV0 {
			expected = strings.TrimSpace(fix.ReadGenerated())
		} else {
			expected = strings.TrimSpace(fix.ReadGeneratedV1())
		}

		if actual != expected {
			t.Errorf("%s > Actual did not match expected (is V0? %t)", fix.Profile, isV0)
		}
	}
}

func validateRegoUnit(module RegoUnit, isV0 bool) (bool, error) {
	var regoVersion ast.RegoVersion
	if isV0 {
		regoVersion = ast.RegoV0
	} else {
		regoVersion = ast.RegoV1
	}

	validator := rego.New(
		rego.Query("data."+module.Name+"."+module.Entrypoint),
		rego.Module(module.Name+".rego", module.Code),
		rego.SetRegoVersion(regoVersion),
	)
	ctx := context.Background()
	_, err := validator.Eval(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
