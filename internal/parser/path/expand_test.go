package path

import (
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/validator/contexts"
	"testing"
)

func testExpandedPath(source, expected string, t *testing.T) {
	parsed, err := ParsePath(source)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	expander := misc.IriExpander{Context: contexts.DefaultAMFContext}
	actual, err := parsed.Expanded(&expander)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if actual != expected {
		t.Errorf("Actual did not match expected\n\tExpected: %s\n\tActual: %s", expected, actual)
	}
}

func TestExpanded(t *testing.T) {
	source := "core.name"
	expected := "http://a.ml/vocabularies/core#name"
	testExpandedPath(source, expected, t)
}

func TestInverseExpanded(t *testing.T) {
	source := "core.name ^"
	expected := "http://a.ml/vocabularies/core#name"
	testExpandedPath(source, expected, t)
}

func TestInverse2Expanded(t *testing.T) {
	source := "core.name^"
	expected := "http://a.ml/vocabularies/core#name"
	testExpandedPath(source, expected, t)
}

func TestANDExpanded(t *testing.T) {
	source := "shacl.schema / core.name"
	expected := "http://www.w3.org/ns/shacl#schema / http://a.ml/vocabularies/core#name"
	testExpandedPath(source, expected, t)
}

func TestManagementExpanded(t *testing.T) {
	source := "management.labels"
	expected := "http://anypoint.com/vocabs/management#labels"
	testExpandedPath(source, expected, t)
}


func TestDOTExpanded(t *testing.T) {
	source := "management.mulesoft.com"
	expected := "http://anypoint.com/vocabs/management#mulesoft.com"
	testExpandedPath(source, expected, t)
}

func TestSlashExpanded(t *testing.T) {
	source := "management.mulesoft.com\\/apiinstance-id"
	expected := "http://anypoint.com/vocabs/management#mulesoft.com/apiinstance-id"
	testExpandedPath(source, expected, t)
}

func TestORExpanded(t *testing.T) {
	source := "shacl.schema | core.name"
	expected := "http://www.w3.org/ns/shacl#schema | http://a.ml/vocabularies/core#name"
	testExpandedPath(source, expected, t)
}

func TestComplexExpanded(t *testing.T) {
	source := "apiContract.expects / (apiContract.parameter / shapes.schema) | (apiContract.payload / shapes.schema) / shacl.name"
	expected := "http://a.ml/vocabularies/apiContract#expects / ((http://a.ml/vocabularies/apiContract#parameter / http://a.ml/vocabularies/shapes#schema) | (http://a.ml/vocabularies/apiContract#payload / http://a.ml/vocabularies/shapes#schema)) / http://www.w3.org/ns/shacl#name"
	testExpandedPath(source, expected, t)
}

func TestMustNotFailWithNilExpander(t *testing.T) {
	source := "core.name"
	parsed, err := ParsePath(source)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	actual, _ := parsed.Expanded(nil)
	if actual != source {
		t.Errorf("Actual did not match expected\n\tExpected: %s\n\tActual: %s", source, actual)
	}
}
