package path

import (
	"testing"
)

func TestPEGProperty(t *testing.T) {
	parsed, err := Parse("property", []byte("core.name"))
	if err != nil {
		t.Errorf("Exception parsing %s", err)
	}
	switch v := parsed.(type) {
	case IRI:
		if v.Value != "core:name" {
			t.Errorf("expected IRI 'core:name' got %s", v.Value)
		}
		if v.Transitive {
			t.Errorf("expected non-transitive IRI")
		}
		if v.Inverse {
			t.Errorf("expected non-inverse IRI")
		}
	default:
		t.Errorf("expected IRI got %v", v)
	}
}

func TestPEGPropertyInverse(t *testing.T) {
	parsed, err := Parse("property", []byte("core.name ^"))
	if err != nil {
		t.Errorf("Exception parsing %s", err)
	}
	switch v := parsed.(type) {
	case IRI:
		if v.Value != "core:name" {
			t.Errorf("expected IRI 'core:name' got %s", v.Value)
		}
		if v.Transitive {
			t.Errorf("expected non-transitive IRI")
		}
		if !v.Inverse {
			t.Errorf("expected noninverse IRI")
		}
	default:
		t.Errorf("expected IRI got %v", v)
	}
}

func TestPEGPropertyInverse2(t *testing.T) {
	parsed, err := Parse("property", []byte("core.name^"))
	if err != nil {
		t.Errorf("Exception parsing %s", err)
	}
	switch v := parsed.(type) {
	case IRI:
		if v.Value != "core:name" {
			t.Errorf("expected IRI 'core:name' got %s", v.Value)
		}
		if v.Transitive {
			t.Errorf("expected non-transitive IRI")
		}
		if !v.Inverse {
			t.Errorf("expected noninverse IRI")
		}
	default:
		t.Errorf("expected IRI got %v", v)
	}
}

func TestPEGPropertyAND(t *testing.T) {
	parsed, err := Parse("property", []byte("shacl.schema / core.name"))
	if err != nil {
		t.Errorf("Exception parsing %s", err)
	}
	switch v := parsed.(type) {
	case AND:
		if len(v.body) != 2 {
			t.Errorf("expected 2 arguments in AND, got %d", len(v.body))
		}
		if v.body[0].(IRI).Value != "shacl:schema" {
			t.Errorf("first AND component must be shacl:schema")
		}
		if v.body[1].(IRI).Value != "core:name" {
			t.Errorf("first AND component must be core:name")
		}
	default:
		t.Errorf("expected AND got %v", v)
	}
}

func TestPEGPropertyOR(t *testing.T) {
	parsed, err := Parse("property", []byte("shacl.schema | core.name"))
	if err != nil {
		t.Errorf("Exception parsing %s", err)
	}
	switch v := parsed.(type) {
	case OR:
		if len(v.body) != 2 {
			t.Errorf("expected 2 arguments in OR, got %d", len(v.body))
		}
		if v.body[0].(IRI).Value != "shacl:schema" {
			t.Errorf("first OR component must be shacl:schema")
		}
		if v.body[1].(IRI).Value != "core:name" {
			t.Errorf("first OR component must be core:name")
		}
	default:
		t.Errorf("expected OR got %v", v)
	}
}

func TestPEGPropertyComplex(t *testing.T) {
	parsed, err := Parse("property", []byte("apiContract.expects / (apiContract.parameter / shapes.schema) | (apiContract.payload / shapes.schema) / shacl.name"))
	if err != nil {
		t.Errorf("Exception parsing %s", err)
	}
	switch v := parsed.(type) {
	case AND:
		// correct
	default:
		t.Errorf("expected OR got %v", v)
	}
}
