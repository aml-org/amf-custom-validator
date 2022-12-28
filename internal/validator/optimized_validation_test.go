package validator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/config"
	"github.com/aml-org/amf-custom-validator/test"
	"testing"
)

func testValidateOptimized(directory string, t *testing.T) {
	profilePath := fmt.Sprintf("%s/%s", directory, "profile.yaml")
	profileText := read(profilePath)

	inputPath := fmt.Sprintf("%s/%s", directory, "api.lexical.jsonld")
	inputText := read(inputPath)

	report, err := ValidateWithOptimization(profileText, inputText, false, nil)
	if err != nil {
		return
	}

	goldenPath := fmt.Sprintf("%s/%s", directory, "report.jsonld")
	if config.Override {
		test.ForceWrite(goldenPath, report)
	} else {
		goldenText := read(goldenPath)
		if goldenText != report {
			t.Errorf("Optimized validation test > %s > Actual did not match expected", directory)
		}
	}
}

// Actual tests --------------------------------------------------------------------------------------------------------

func TestPathToScalarValidation(t *testing.T) {
	testValidateOptimized("../../test/data/optimization/path-to-scalar", t)
}
