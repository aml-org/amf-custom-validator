package path

import (
	"testing"
)

func TestNullPath(t *testing.T) {
	p, err := ParsePath("")
	if err != nil {
		t.Errorf("got error parsing null path %v", err)
	}
	switch v := p.(type) {
	case NullPath:
		// correct
		break
	default:
		t.Errorf("expected NullPath, got %v", v)
	}
}
