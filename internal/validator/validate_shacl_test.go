package validator

import (
	"fmt"
	"testing"
)

const shaclTckPath relativePath = "../../test/data/shacl"
const override = false

func testFn(directory relativePath, t *testing.T) {
	resolvedDirectory := fmt.Sprintf("%s/%s", shaclTckPath, directory)
	profile := fmt.Sprintf("%s/profile.yaml", resolvedDirectory)
	data := fmt.Sprintf("%s/data.jsonld", resolvedDirectory)
	actualText := validate(profile, data)
	expected := fmt.Sprintf("%s/report.jsonld", resolvedDirectory)
	if override {
		write(actualText, expected)
	} else {
		if !compare(actualText, expected) {
			t.Errorf("Failed %s. Actual did not match expexted", directory)
		}
	}
}

func ignore(directory relativePath, t *testing.T) {
	t.Skipf("Ignored %s", directory)
}

func TestMiscSeverity001(t *testing.T) {
	testFn("misc/severity-001", t)
}

func TestNodeAnd001(t *testing.T) {
	testFn("node/and-001", t)
}

func TestNodeAnd002(t *testing.T) {
	testFn("node/and-002", t)
}

func TestNodeNot001(t *testing.T) {
	testFn("node/not-001", t)
}

func TestNodeOr001(t *testing.T) {
	ignore("node/or-001", t)
}

func TestNodeXone001(t *testing.T) {
	ignore("node/xone-001", t)
}

func TestPathAlternative001(t *testing.T) {
	testFn("path/path-alternative-001", t)
}

func TestPathInverse001(t *testing.T) {
	ignore("path/path-inverse-001", t)
}

func TestPathSequence001(t *testing.T) {
	testFn("path/path-sequence-001", t)
}

func TestPathSequence002(t *testing.T) {
	testFn("path/path-sequence-002", t)
}

func TestPropertyDatatype001(t *testing.T) {
	ignore("property/datatype-001", t)
}

func TestPropertyDatatype002(t *testing.T) {
	testFn("property/datatype-002", t)
}

func TestPropertyDisjoint001(t *testing.T) {
	ignore("property/disjoint-001", t)
}

func TestPropertyEquals001(t *testing.T) {
	ignore("property/equals-001", t)
}

func TestPropertyIn001(t *testing.T) {
	testFn("property/in-001", t)
}

func TestPropertyLessThan001(t *testing.T) {
	testFn("property/lessThan-001", t)
}

func TestPropertyLessThan002(t *testing.T) {
	ignore("property/lessThan-002", t)
}

func TestPropertyLessThanOrEquals001(t *testing.T) {
	testFn("property/lessThanOrEquals-001", t)
}

func TestPropertyMaxCount001(t *testing.T) {
	testFn("property/maxCount-001", t)
}

func TestPropertyMaxCount002(t *testing.T) {
	testFn("property/maxCount-002", t)
}

func TestPropertyMaxExclusive001(t *testing.T) {
	testFn("property/maxExclusive-001", t)
}

func TestPropertyMaxInclusive001(t *testing.T) {
	testFn("property/maxInclusive-001", t)
}

func TestPropertyMaxLength001(t *testing.T) {
	ignore("property/maxLength-001", t)
}

func TestPropertyMinCount001(t *testing.T) {
	testFn("property/minCount-001", t)
}

func TestPropertyMinCount002(t *testing.T) {
	testFn("property/minCount-002", t)
}

func TestPropertyMinExclusive001(t *testing.T) {
	ignore("property/minExclusive-001", t)
}

func TestPropertyMinExclusive002(t *testing.T) {
	testFn("property/minExclusive-002", t)
}

func TestPropertyMinLength001(t *testing.T) {
	ignore("property/minLength-001", t)
}

func TestPropertyPattern001(t *testing.T) {
	testFn("property/pattern-001", t)
}

func TestPropertyPattern002(t *testing.T) {
	ignore("property/pattern-002", t)
}

func TestPropertyTargetClass001(t *testing.T) {
	testFn("targets/targetClass-001", t)
}
