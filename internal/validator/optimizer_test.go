package validator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/internal/parser"
	"github.com/aml-org/amf-custom-validator/test"
	"testing"
)

func testOptimize(directory string, t *testing.T) {
	profilePath := fmt.Sprintf("%s/%s", directory, "profile.yaml")
	profileText := read(profilePath)
	profile, _ := parser.Parse(profileText)

	inputPath := fmt.Sprintf("%s/%s", directory, "api.lexical.jsonld")
	inputText := read(inputPath)
	input := ParseJson(inputText)

	optimizer := NewOptimizer(profile)
	optimizer.Optimize(&input)
	optimizedInputText := encodeJson(input)

	goldenPath := fmt.Sprintf("%s/%s", directory, "api.optimized.jsonld")
	if config.Override {
		test.ForceWrite(goldenPath, optimizedInputText)
	} else {
		goldenText := read(goldenPath)
		if goldenText != optimizedInputText {
			t.Errorf("Optimization test > %s > Actual did not match expected", directory)
		}
	}
}

// Actual tests --------------------------------------------------------------------------------------------------------

func TestPathToScalar(t *testing.T) {
	testOptimize("../../test/data/optimization/path-to-scalar", t)
}
