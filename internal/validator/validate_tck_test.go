package validator

import (
	"fmt"
	"testing"
)

const tckPath relativePath = "../../test/data/tck"

func testTckCase(subDirectory relativePath, t *testing.T) {
	directory := fmt.Sprintf("%s/%s", tckPath, subDirectory)
	validateAndCompareDirectory(directory, t)
}

// TODO none of these are working properly, we do not support validation of coerced types. Modify these after defining the expected behavior

func TestDataTypesInteger(t *testing.T) {
	testTckCase("data-types/integer", t)
}

func TestDataTypesDouble(t *testing.T) {
	testTckCase("data-types/double", t)
}

func TestDataTypesFloat(t *testing.T) {
	testTckCase("data-types/float", t)
}

func TestDataTypesBoolean(t *testing.T) {
	testTckCase("data-types/boolean", t)
}

func TestDataTypesString(t *testing.T) {
	testTckCase("data-types/string", t)
}

func TestDataTypesDuration(t *testing.T) {
	testTckCase("data-types/duration", t)
}

func TestDataTypesCoercion(t *testing.T) {
	testTckCase("data-types/coercion", t)
}
