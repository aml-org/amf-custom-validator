package validator

import (
	"fmt"
	"testing"
)

const shaclTckPath relativePath = "../../test/data/shacl"

func testShaclCase(subDirectory relativePath, t *testing.T) {
	directory := fmt.Sprintf("%s/%s", shaclTckPath, subDirectory)
	validateAndCompareDirectory(directory, t)
}

func ignoreShaclCase(subDirectory relativePath, t *testing.T) {
	directory := fmt.Sprintf("%s/%s", shaclTckPath, subDirectory)
	ignoreDirectory(directory, t)
}

func TestMiscSeverity001(t *testing.T) {
	testShaclCase("misc/severity-001", t)
}

func TestNodeAnd001(t *testing.T) {
	testShaclCase("node/and-001", t)
}

func TestNodeAnd002(t *testing.T) {
	testShaclCase("node/and-002", t)
}

func TestNodeNot001(t *testing.T) {
	testShaclCase("node/not-001", t)
}

func TestNodeOr001(t *testing.T) {
	ignoreShaclCase("node/or-001", t)
}

func TestNodeXone001(t *testing.T) {
	ignoreShaclCase("node/xone-001", t)
}

func TestPathAlternative001(t *testing.T) {
	testShaclCase("path/path-alternative-001", t)
}

func TestPathInverse001(t *testing.T) {
	ignoreShaclCase("path/path-inverse-001", t)
}

func TestPathSequence001(t *testing.T) {
	testShaclCase("path/path-sequence-001", t)
}

func TestPathSequence002(t *testing.T) {
	testShaclCase("path/path-sequence-002", t)
}

func TestPropertyDatatype001(t *testing.T) {
	ignoreShaclCase("property/datatype-001", t)
}

func TestPropertyDatatype002(t *testing.T) {
	testShaclCase("property/datatype-002", t)
}

func TestPropertyDisjoint001(t *testing.T) {
	ignoreShaclCase("property/disjoint-001", t)
}

func TestPropertyEquals001(t *testing.T) {
	ignoreShaclCase("property/equals-001", t)
}

func TestPropertyIn001(t *testing.T) {
	testShaclCase("property/in-001", t)
}

func TestPropertyLessThan001(t *testing.T) {
	testShaclCase("property/lessThan-001", t)
}

func TestPropertyLessThan002(t *testing.T) {
	ignoreShaclCase("property/lessThan-002", t)
}

func TestPropertyLessThanOrEquals001(t *testing.T) {
	testShaclCase("property/lessThanOrEquals-001", t)
}

func TestPropertyMaxCount001(t *testing.T) {
	testShaclCase("property/maxCount-001", t)
}

func TestPropertyMaxCount002(t *testing.T) {
	testShaclCase("property/maxCount-002", t)
}

func TestPropertyMaxExclusive001(t *testing.T) {
	testShaclCase("property/maxExclusive-001", t)
}

func TestPropertyMaxInclusive001(t *testing.T) {
	testShaclCase("property/maxInclusive-001", t)
}

func TestPropertyMaxLength001(t *testing.T) {
	ignoreShaclCase("property/maxLength-001", t)
}

func TestPropertyMinCount001(t *testing.T) {
	testShaclCase("property/minCount-001", t)
}

func TestPropertyMinCount002(t *testing.T) {
	testShaclCase("property/minCount-002", t)
}

func TestPropertyMinExclusive001(t *testing.T) {
	ignoreShaclCase("property/minExclusive-001", t)
}

func TestPropertyMinExclusive002(t *testing.T) {
	testShaclCase("property/minExclusive-002", t)
}

func TestPropertyMinLength001(t *testing.T) {
	ignoreShaclCase("property/minLength-001", t)
}

func TestPropertyPattern001(t *testing.T) {
	testShaclCase("property/pattern-001", t)
}

func TestPropertyPattern002(t *testing.T) {
	ignoreShaclCase("property/pattern-002", t)
}

func TestPropertyTargetClass001(t *testing.T) {
	testShaclCase("targets/targetClass-001", t)
}
