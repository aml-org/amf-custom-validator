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
	testTckCase("conditionals/if-then-else", t)
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

func TestConditionalsIfNotThen(t *testing.T) {
	testTckCase("conditionals/if-not-then", t)
}

func TestConditionalsIfNotThenNot(t *testing.T) {
	testTckCase("conditionals/if-not-then-not", t)
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

func TestAndNot(t *testing.T) {
	testTckCase("and/and-not", t)
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

func TestOrNot(t *testing.T) {
	testTckCase("or/or-not", t)
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

func TestNestedNot(t *testing.T) {
	testTckCase("nested/nested-not", t)
}

// property constraints (are not natively defined in shacl)

func TestPropertyExactLength(t *testing.T) {
	testTckCase("property/exactLength", t)
}

func TestPropertyExactCount(t *testing.T) {
	testTckCase("property/exactCount", t)
}

func TestPropertyContainsAll(t *testing.T) {
	testTckCase("property/containsAll", t)
}

func TestPropertyContainsSome(t *testing.T) {
	testTckCase("property/containsSome", t)
}

func TestPropertyOnlyValue(t *testing.T) {
	testTckCase("property/is", t)
}