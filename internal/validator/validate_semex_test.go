package validator

import (
	"fmt"
	"testing"
)

const semexPath relativePath = "../../test/data/semex"

func testSemexCase(subDirectory relativePath, dataFileName string, t *testing.T) {
	profile := fmt.Sprintf("%s/%s/profile.yaml", semexPath, subDirectory)
	dataPath := fmt.Sprintf("%s/%s/%s", semexPath, subDirectory, dataFileName)
	expectedPath := fmt.Sprintf("%s/%s/%s.report.jsonld", semexPath, subDirectory, dataFileName)
	validateAndCompare(profile, dataPath, expectedPath, t)
}

func TestPaginationAsync20(t *testing.T) {
	testSemexCase("semantic-extension-pagination", "api.async.yaml.jsonld", t)
}
func TestPaginationOas20(t *testing.T) {
	testSemexCase("semantic-extension-pagination", "api.oas20.yaml.jsonld", t)
}
func TestPaginationOas30(t *testing.T) {
	testSemexCase("semantic-extension-pagination", "api.oas30.yaml.jsonld", t)
}
func TestPaginationRaml10(t *testing.T) {
	testSemexCase("semantic-extension-pagination", "api.raml.jsonld", t)
}
