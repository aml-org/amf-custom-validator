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

func ignoreTckCase(subDirectory relativePath, t *testing.T) {
	directory := fmt.Sprintf("%s/%s", tckPath, subDirectory)
	ignoreDirectory(directory, t)
}

// Data types

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

// Conditionals

func TestConditionalsIfThen(t *testing.T) {
	testTckCase("conditionals/if-then", t)
}

func TestConditionalsIfThenElse(t *testing.T) {
	ignoreTckCase("conditionals/if-then-else", t) // else not implemented
}


func TestConditionalsIfAndThenAnd(t *testing.T) {
	testTckCase("conditionals/if-and-then-and", t)
}

func TestConditionalsIfOrThenOr(t *testing.T) {
	testTckCase("conditionals/if-or-then-or", t)
}

func TestConditionalsIfNestedThenNested(t *testing.T) {
	testTckCase("conditionals/if-nested-then-nested", t)
}

// And ('simple and' tested in SHACL TCK)

func TestAndAnd(t *testing.T) {
	testTckCase("and/and-and", t)
}

func TestAndIfThen(t *testing.T) {
	testTckCase("and/and-if-then", t)
}

func TestAndNested(t *testing.T) {
	testTckCase("and/and-nested", t)
}

func TestAndOr(t *testing.T) {
	testTckCase("and/and-or", t)
}

// Or ('simple or' tested in SHACL TCK)

func TestOrAnd(t *testing.T) {
	testTckCase("or/or-and", t)
}

func TestOrIfThen(t *testing.T) {
	testTckCase("or/or-if-then", t)
}

func TestOrNested(t *testing.T) {
	testTckCase("or/or-nested", t)
}

func TestOrOr(t *testing.T) {
	testTckCase("or/or-or", t)
}

// Severity ('violation' & 'warning' severities tested in SHACL TCK)

func TestSeverityInfo(t *testing.T) {
	testTckCase("severity/info", t)
}

// Nested

func TestNestedAnd(t *testing.T) {
	testTckCase("nested/nested-and", t)
}

func TestNestedIfThen(t *testing.T) {
	testTckCase("nested/nested-if-then", t)
}

func TestNestedNested(t *testing.T) {
	testTckCase("nested/nested-nested", t)
}

func TestNestedOr(t *testing.T) {
	testTckCase("nested/nested-or", t)
}