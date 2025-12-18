package validator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeBaseUri(t *testing.T) {
	testDataPath := filepath.Join("..", "..", "test", "data", "test-base-uri.jsonld")
	fileData, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}

	originalBaseUri := "urn:ms:d2270dec-7e27-4494-806c-23228cf84dce:asset::d2270dec-7e27-4494-806c-23228cf84dce/required-examples-w-17395811/1.0.0"

	// Check if original JSON contains the base URI
	if !strings.Contains(string(fileData), originalBaseUri) {
		t.Errorf("Original JSON should contain base URI: %s", originalBaseUri)
	}

	// Parse the JSON-LD string into a JSON object
	var jsonData any
	err = json.Unmarshal(fileData, &jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON-LD: %v", err)
	}

	// Now pass the parsed JSON object to Normalize
	normalizedJson := Normalize(jsonData)

	// Convert normalized JSON to pretty-printed string
	normalizedBytes, err := json.MarshalIndent(normalizedJson, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal normalized JSON: %v", err)
	}

	normalizedStr := string(normalizedBytes)

	// Check if normalized JSON does not contain the base URI
	if strings.Contains(normalizedStr, originalBaseUri) {
		t.Errorf("Normalized JSON should not contain base URI: %s", originalBaseUri)
	}
}
